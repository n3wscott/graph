package knative

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	flowsv1 "knative.dev/eventing/pkg/apis/flows/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
)

func (c *Client) Sequences(namespace string, yv *[]YamlView) []flowsv1.Sequence {
	gvr := schema.GroupVersionResource{
		Group:    "messaging.knative.dev",
		Version:  "v1",
		Resource: "sequences",
	}
	like := flowsv1.Sequence{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Sequences, %v", err)
		return nil
	}

	all := make([]flowsv1.Sequence, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}

func (c *Client) InMemoryChannels(namespace string, yv *[]YamlView) []messagingv1.InMemoryChannel {
	gvr := schema.GroupVersionResource{
		Group:    "messaging.knative.dev",
		Version:  "v1",
		Resource: "inmemorychannels",
	}
	like := messagingv1.InMemoryChannel{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List InMemoryChannels, %v", err)
		return nil
	}

	all := make([]messagingv1.InMemoryChannel, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}

func (c *Client) Subscriptions(namespace string, yv *[]YamlView) []messagingv1.Subscription {
	gvr := schema.GroupVersionResource{
		Group:    "messaging.knative.dev",
		Version:  "v1",
		Resource: "subscriptions",
	}
	like := messagingv1.Subscription{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Subscriptions, %v", err)
		return nil
	}

	all := make([]messagingv1.Subscription, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}
