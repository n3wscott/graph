package graph

import (
	"k8s.io/client-go/dynamic"

	"github.com/n3wscott/graph/pkg/knative"
)

func ForTriggers(client dynamic.Interface, ns string) (string, []knative.YamlView) {
	g := New(ns)

	c := knative.New(client)

	yv := make([]knative.YamlView, 0)

	// First pre-load the services.
	for _, service := range c.KnServices(ns, nil) {
		g.LoadKnService(service)
	}

	// load the brokers
	for _, broker := range c.Brokers(ns, &yv) {
		g.AddBroker(broker)
	}

	// load the triggers
	for _, trigger := range c.Triggers(ns, &yv) {
		g.AddTrigger(trigger)
	}

	// load the services
	for _, service := range c.KnServices(ns, &yv) {
		g.AddKnService(service)
	}

	// load the sequences
	for _, sequence := range c.Sequences(ns, &yv) {
		g.AddSequence(sequence)
	}

	// Last load the sources.
	for _, source := range c.Sources(ns, &yv) {
		g.AddSource(source)
	}

	return g.String(), yv
}

func ForSubscriptions(client dynamic.Interface, ns string) (string, []knative.YamlView) {
	g := New(ns)

	c := knative.New(client)

	yv := make([]knative.YamlView, 0)

	// First pre-load the services.
	for _, service := range c.KnServices(ns, nil) {
		g.LoadKnService(service)
	}

	// load the brokers
	for _, broker := range c.Brokers(ns, &yv) {
		g.AddBroker(broker)
	}

	// load the sources
	for _, source := range c.Sources(ns, &yv) {
		g.AddSource(source)
	}

	// load the triggers
	for _, trigger := range c.Triggers(ns, &yv) {
		g.AddTrigger(trigger)
	}

	// load the services
	for _, service := range c.KnServices(ns, &yv) {
		g.AddKnService(service)
	}

	for _, channel := range c.Channels(ns, &yv) {
		g.AddChannel(channel)
	}

	for _, subscription := range c.Subscriptions(ns, &yv) {
		g.AddSubscription(subscription)
	}

	return g.String(), yv
}
