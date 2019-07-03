package graph

import (
	"github.com/n3wscott/graph/pkg/knative"
	"k8s.io/client-go/dynamic"
)

func ForTriggers(client dynamic.Interface, ns string) string {
	g := New(ns)

	c := knative.New(client)

	// First pre-load the services.
	for _, service := range c.KnServices(ns) {
		g.LoadKnService(service)
	}

	// load the brokers
	for _, broker := range c.Brokers(ns) {
		g.AddBroker(broker)
	}

	// load the triggers
	for _, trigger := range c.Triggers(ns) {
		g.AddTrigger(trigger)
	}

	// load the services
	for _, service := range c.KnServices(ns) {
		g.AddKnService(service)
	}

	// load the sequences
	for _, sequence := range c.Sequences(ns) {
		g.AddSequence(sequence)
	}

	// Last load the sources.
	for _, source := range c.Sources(ns) {
		g.AddSource(source)
	}

	return g.String()
}

func ForSubscriptions(client dynamic.Interface, ns string) string {
	g := New(ns)

	c := knative.New(client)

	// First pre-load the services.
	for _, service := range c.KnServices(ns) {
		g.LoadKnService(service)
	}

	// load the brokers
	for _, broker := range c.Brokers(ns) {
		g.AddBroker(broker)
	}

	// load the sources
	for _, source := range c.Sources(ns) {
		g.AddSource(source)
	}

	// load the triggers
	for _, trigger := range c.Triggers(ns) {
		g.AddTrigger(trigger)
	}

	// load the services
	for _, service := range c.KnServices(ns) {
		g.AddKnService(service)
	}

	for _, channel := range c.Channels(ns) {
		g.AddChannel(channel)
	}

	for _, channel := range c.InMemoryChannels(ns) {
		g.AddInMemoryChannel(channel)
	}

	for _, subscription := range c.Subscriptions(ns) {
		g.AddSubscription(subscription)
	}

	return g.String()
}
