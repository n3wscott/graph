package knative

import (
	"github.com/knative/eventing/pkg/apis/duck/v1alpha1"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	eventingv1alpha1 "knative.dev/eventing/pkg/apis/eventing/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func (c *Client) Sources(namespace string, yv *[]YamlView) []duckv1.Source {
	gvrs := crdsToGVR(c.SourceCRDs())
	all := make([]duckv1.Source, 0)

	for _, gvr := range gvrs {
		objs := c.Source(namespace, gvr, yv)
		all = append(all, objs...)
	}
	return all
}

func (c *Client) Triggers(namespace string, yv *[]YamlView) []eventingv1alpha1.Trigger {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "triggers",
	}
	like := eventingv1alpha1.Trigger{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Triggers, %v", err)
		return nil
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

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}

func (c *Client) Brokers(namespace string, yv *[]YamlView) []eventingv1alpha1.Broker {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "brokers",
	}
	like := eventingv1alpha1.Broker{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Brokers, %v", err)
		return nil
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

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}

func (c *Client) Channels(namespace string, yv *[]YamlView) []v1alpha1.Channelable {
	gvrs := crdsToGVR(c.ChannelCRDs())
	all := make([]v1alpha1.Channelable, 0)

	for _, gvr := range gvrs {
		objs := c.Channelable(namespace, gvr, yv)
		all = append(all, objs...)
	}
	return all
}

func (c *Client) EventTypes(namespace string, yv *[]YamlView) []eventingv1alpha1.EventType {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1alpha1",
		Resource: "eventtypes",
	}
	like := eventingv1alpha1.EventType{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List EventTypes, %v", err)
		return nil
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

		// Yaml View
		AddToYamlView(item, yv)
	}
	return all
}
