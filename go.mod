module github.com/n3wscott/graph

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.1.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/tmc/dot v0.0.0-20180926222610-6d252d5ff882
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/eventing v0.15.2
	knative.dev/pkg v0.0.0-20200625173728-dfb81cf04a7c // release-0.15
	knative.dev/serving v0.15.2
)

// Pin k8s deps to 1.16.5
replace (
	k8s.io/api => k8s.io/api v0.16.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5
	k8s.io/client-go => k8s.io/client-go v0.16.5
	k8s.io/code-generator => k8s.io/code-generator v0.16.5
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190327210449-e17681d19d3a
)
