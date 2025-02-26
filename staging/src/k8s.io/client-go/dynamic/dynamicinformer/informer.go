/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dynamicinformer

import (
	"context"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamiclister"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// NewDynamicSharedInformerFactory constructs a new instance of dynamicSharedInformerFactory for all namespaces.
func NewDynamicSharedInformerFactory(client dynamic.Interface, defaultResync time.Duration) DynamicSharedInformerFactory {
	return NewFilteredDynamicSharedInformerFactory(client, defaultResync, metav1.NamespaceAll, nil)
}

// NewFilteredDynamicSharedInformerFactory constructs a new instance of dynamicSharedInformerFactory.
// Listers obtained via this factory will be subject to the same filters as specified here.
func NewFilteredDynamicSharedInformerFactory(client dynamic.Interface, defaultResync time.Duration, namespace string, tweakListOptions TweakListOptionsFunc) DynamicSharedInformerFactory {
	return NewFilteredDynamicSharedInformerFactoryWithOptions(
		client,
		namespace,
		tweakListOptions,
		cache.WithResyncPeriod(defaultResync),
		cache.WithIndexers(cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}),
	)
}

// NewFilteredDynamicSharedInformerFactoryWithOptions constructs a new instance of dynamicSharedInformerFactory.
// Listers obtained via this factory will be subject to the same filters as specified here.
func NewFilteredDynamicSharedInformerFactoryWithOptions(client dynamic.Interface, namespace string, tweakListOptions TweakListOptionsFunc, opts ...cache.SharedInformerOption) DynamicSharedInformerFactory {
	return &dynamicSharedInformerFactory{
		client:                 client,
		namespace:              namespace,
		informers:              map[schema.GroupVersionResource]informers.GenericInformer{},
		startedInformers:       make(map[schema.GroupVersionResource]bool),
		tweakListOptions:       tweakListOptions,
		dynamicInformerOptions: opts,
	}
}

type dynamicSharedInformerFactory struct {
	client        dynamic.Interface
	defaultResync time.Duration
	namespace     string

	lock      sync.Mutex
	informers map[schema.GroupVersionResource]informers.GenericInformer
	// startedInformers is used for tracking which informers have been started.
	// This allows Start() to be called multiple times safely.
	startedInformers map[schema.GroupVersionResource]bool
	tweakListOptions TweakListOptionsFunc

	dynamicInformerOptions []cache.SharedInformerOption
}

var _ DynamicSharedInformerFactory = &dynamicSharedInformerFactory{}

func (f *dynamicSharedInformerFactory) ForResource(gvr schema.GroupVersionResource) informers.GenericInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	key := gvr
	informer, exists := f.informers[key]
	if exists {
		return informer
	}

	informer = NewFilteredDynamicInformerWithOptions(
		f.client,
		gvr,
		f.namespace,
		f.tweakListOptions,
		f.dynamicInformerOptions...,
	)
	f.informers[key] = informer

	return informer
}

// Start initializes all requested informers.
func (f *dynamicSharedInformerFactory) Start(stopCh <-chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	for informerType, informer := range f.informers {
		if !f.startedInformers[informerType] {
			go informer.Informer().Run(stopCh)
			f.startedInformers[informerType] = true
		}
	}
}

// WaitForCacheSync waits for all started informers' cache were synced.
func (f *dynamicSharedInformerFactory) WaitForCacheSync(stopCh <-chan struct{}) map[schema.GroupVersionResource]bool {
	informers := func() map[schema.GroupVersionResource]cache.SharedIndexInformer {
		f.lock.Lock()
		defer f.lock.Unlock()

		informers := map[schema.GroupVersionResource]cache.SharedIndexInformer{}
		for informerType, informer := range f.informers {
			if f.startedInformers[informerType] {
				informers[informerType] = informer.Informer()
			}
		}
		return informers
	}()

	res := map[schema.GroupVersionResource]bool{}
	for informType, informer := range informers {
		res[informType] = cache.WaitForCacheSync(stopCh, informer.HasSynced)
	}
	return res
}

// NewFilteredDynamicInformer constructs a new informer for a dynamic type.
func NewFilteredDynamicInformer(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions TweakListOptionsFunc) informers.GenericInformer {
	return NewFilteredDynamicInformerWithOptions(
		client,
		gvr,
		namespace,
		tweakListOptions,
		cache.WithResyncPeriod(resyncPeriod),
		cache.WithIndexers(indexers),
	)
}

func NewFilteredDynamicInformerWithOptions(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string, tweakListOptions TweakListOptionsFunc, opts ...cache.SharedInformerOption) informers.GenericInformer {
	return &dynamicInformer{
		gvr: gvr,
		informer: cache.NewSharedIndexInformerWithOptions(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					if tweakListOptions != nil {
						tweakListOptions(&options)
					}
					return client.Resource(gvr).Namespace(namespace).List(context.TODO(), options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					if tweakListOptions != nil {
						tweakListOptions(&options)
					}
					return client.Resource(gvr).Namespace(namespace).Watch(context.TODO(), options)
				},
			},
			&unstructured.Unstructured{},
			opts...,
		)}
}

type dynamicInformer struct {
	informer cache.SharedIndexInformer
	gvr      schema.GroupVersionResource
}

var _ informers.GenericInformer = &dynamicInformer{}

func (d *dynamicInformer) Informer() cache.SharedIndexInformer {
	return d.informer
}

func (d *dynamicInformer) Lister() cache.GenericLister {
	return dynamiclister.NewRuntimeObjectShim(dynamiclister.New(d.informer.GetIndexer(), d.gvr))
}
