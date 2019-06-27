package graph

import (
	"github.com/n3wscott/knap/pkg/knative"
	"k8s.io/client-go/dynamic"
)

func ForTriggers(client dynamic.Interface, ns string) string {
	g := New(ns)

	c := knative.New(client)

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
	return g.String()
}

func ForSubscriptions(client dynamic.Interface, ns string) string {
	g := New(ns)

	c := knative.New(client)

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

	for _, subscription := range c.Subscriptions(ns) {
		g.AddSubscription(subscription)
	}

	return g.String()
}
