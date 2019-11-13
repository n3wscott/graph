# graph

Visualize your [Knative Eventing](http://github.com/knative/eventing)
connections.

<img src="./image/preview.png">

## Usage

Visit the root of the graph service in a web browser. This will show you the
graph of the current Knative resources in the namespace the graph resource is
installed.

> Note: Work is required to support installation of `graph` into multiple
> namespaces.

# Deploying

## From Release v0.5.0 >= Knative 0.10.0

> Note: Graph uses Serving v1. To change this, you can edit the release yaml.

To install into the `default` namespace,

```shell
kubectl apply -f https://github.com/n3wscott/graph/releases/download/v0.5.0/release.yaml
```

To install into a `test` namespace,

```shell
export NAMESPACE=test # <-- update test to your target namespace.
curl -L https://github.com/n3wscott/graph/releases/download/v0.5.0/release.yaml \
  | sed "s/default/${NAMESPACE}/" \
  | kubectl apply -n $NAMESPACE --filename -
```

## From Source

To install into the `default` namespace,

```shell
ko apply -f config
```

To install into a `test` namespace,

```shell
export NAMESPACE=test # <-- update test to your target namespace.
ko resolve -f config \
  | sed "s/default/${NAMESPACE}/" \
  | kubectl apply -n $NAMESPACE --filename -
```


### TODO:

- [ ] Get Deployments working when broker is the sink.
- [ ] Work with owner ref graphs.