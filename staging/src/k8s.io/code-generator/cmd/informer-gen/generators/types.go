/*
Copyright 2016 The Kubernetes Authors.

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

package generators

import "k8s.io/gengo/types"

var (
	apiScheme                 = types.Name{Package: "k8s.io/kubernetes/pkg/api/legacyscheme", Name: "Scheme"}
	cacheGenericLister        = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "GenericLister"}
	cacheIndexers             = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "Indexers"}
	cacheSharedInformerOption = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "SharedInformerOption"}
	cacheWithIndexers         = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "WithIndexers"}
	cacheWithKeyFunction      = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "WithKeyFunction"}
	cacheWithResyncPeriod     = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "WithResyncPeriod"}

	cacheListWatch              = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "ListWatch"}
	cacheMetaNamespaceIndexFunc = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "MetaNamespaceIndexFunc"}
	cacheNamespaceIndex         = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "NamespaceIndex"}
	cacheNewGenericLister       = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "NewGenericLister"}
	cacheNewSharedIndexInformer = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "NewSharedIndexInformer"}
	cacheSharedIndexInformer    = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "SharedIndexInformer"}
	cacheKeyFunc                = types.Name{Package: "k8s.io/client-go/tools/cache", Name: "KeyFunc"}
	listOptions                 = types.Name{Package: "k8s.io/kubernetes/pkg/apis/core", Name: "ListOptions"}
	reflectType                 = types.Name{Package: "reflect", Name: "Type"}
	runtimeObject               = types.Name{Package: "k8s.io/apimachinery/pkg/runtime", Name: "Object"}
	schemaGroupResource         = types.Name{Package: "k8s.io/apimachinery/pkg/runtime/schema", Name: "GroupResource"}
	schemaGroupVersionResource  = types.Name{Package: "k8s.io/apimachinery/pkg/runtime/schema", Name: "GroupVersionResource"}
	syncMutex                   = types.Name{Package: "sync", Name: "Mutex"}
	timeDuration                = types.Name{Package: "time", Name: "Duration"}
	v1ListOptions               = types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "ListOptions"}
	metav1NamespaceAll          = types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "NamespaceAll"}
	metav1Object                = types.Name{Package: "k8s.io/apimachinery/pkg/apis/meta/v1", Name: "Object"}
	watchInterface              = types.Name{Package: "k8s.io/apimachinery/pkg/watch", Name: "Interface"}
)
