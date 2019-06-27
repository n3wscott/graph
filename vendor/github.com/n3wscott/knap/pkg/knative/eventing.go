package knative

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	eventingv1alpha1 "github.com/knative/eventing/pkg/apis/eventing/v1alpha1"
	duckv1alpha1 "github.com/n3wscott/knap/pkg/apis/duck/v1alpha1"
)

func (c *Client) Sources(namespace string) []duckv1alpha1.SourceType {
	gvrs := crdsToGVR(c.SourceCRDs())
	all := make([]duckv1alpha1.SourceType, 0)

	for _, gvr := range gvrs {
		like := duckv1alpha1.SourceType{}
		list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to List %s, %v", gvr.String(), err)
			continue
		}

		for _, item := range list.Items {
			obj := like.DeepCopy()
			if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
				log.Fatalf("Error DefaultUnstructuredConverter.FromUnstructured. %v", err)
			}
			obj.ResourceVersion = gvr.Version
			obj.APIVersion = gvr.GroupVersion().String()
			all = append(all, *obj)
		}
	}
	return all
}

func (c *Client) Triggers(namespace string) []eventingv1alpha1.Trigger {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "triggers",
	}
	like := eventingv1alpha1.Trigger{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to List Triggers, %v", err)
	}

	all := make([]eventingv1alpha1.Trigger, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

func (c *Client) Brokers(namespace string) []eventingv1alpha1.Broker {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "brokers",
	}
	like := eventingv1alpha1.Broker{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to List Brokers, %v", err)
	}

	all := make([]eventingv1alpha1.Broker, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

func (c *Client) Channels(namespace string) []eventingv1alpha1.Channel {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "channels",
	}
	like := eventingv1alpha1.Channel{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to List Channels, %v", err)
	}

	all := make([]eventingv1alpha1.Channel, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

func (c *Client) Subscriptions(namespace string) []eventingv1alpha1.Subscription {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}
	like := eventingv1alpha1.Subscription{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to List Subscriptions, %v", err)
	}

	all := make([]eventingv1alpha1.Subscription, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

func (c *Client) EventTypes(namespace string) []eventingv1alpha1.EventType {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "eventtypes",
	}
	like := eventingv1alpha1.EventType{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to List EventTypes, %v", err)
	}

	all := make([]eventingv1alpha1.EventType, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}
