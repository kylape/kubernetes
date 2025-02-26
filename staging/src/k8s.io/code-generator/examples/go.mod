// This is a submodule to isolate k8s.io/code-generator from k8s.io/{api,apimachinery,client-go} dependencies in generated code

module k8s.io/code-generator/examples

go 1.16

require (
	github.com/kcp-dev/logicalcluster/v2 v2.0.0-alpha.1
	k8s.io/api v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v0.0.0
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42
)

replace (
	k8s.io/api => ../../api
	k8s.io/apimachinery => ../../apimachinery
	k8s.io/client-go => ../../client-go
)

replace sigs.k8s.io/json => github.com/liggitt/json v0.0.0-20211020163728-48258682683b
