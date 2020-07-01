package knative

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	duckv1 "knative.dev/pkg/apis/duck/v1"
	eventingv1beta1 "knative.dev/eventing/pkg/apis/eventing/v1beta1"
)

func (c *Client) Sources(namespace string, yv *[]YamlView) []duckv1.Source {
	gvrs := crdsToGVR(c.SourceCRDs())
	all := make([]duckv1.Source, 0)

	for _, gvr := range gvrs {
		like := duckv1.Source{}
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

			// Yaml View
			AddToYamlView(item, yv)
		}
	}
	return all
}

func (c *Client) Triggers(namespace string, yv *[]YamlView) []eventingv1beta1.Trigger {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1beta1",
		Resource: "triggers",
	}
	like := eventingv1beta1.Trigger{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Triggers, %v", err)
		return nil
	}

	all := make([]eventingv1beta1.Trigger, len(list.Items))

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

func (c *Client) Brokers(namespace string, yv *[]YamlView) []eventingv1beta1.Broker {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1beta1",
		Resource: "brokers",
	}
	like := eventingv1beta1.Broker{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Brokers, %v", err)
		return nil
	}

	all := make([]eventingv1beta1.Broker, len(list.Items))

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

//
//func (c *Client) Channels(namespace string, yv *[]YamlView) []eventingv1beta1.Channel {
//	gvr := schema.GroupVersionResource{
//		Group:    "eventing.knative.dev",
//		Version:  "v1alpha1",
//		Resource: "channels",
//	}
//	like := eventingv1beta1.Channel{}
//
//	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
//	if err != nil {
//		log.Printf("Failed to List Channels, %v", err)
//		return nil
//	}
//
//	all := make([]eventingv1beta1.Channel, len(list.Items))
//
//	for i, item := range list.Items {
//		obj := like.DeepCopy()
//		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
//			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
//		}
//		obj.ResourceVersion = gvr.Version
//		obj.APIVersion = gvr.GroupVersion().String()
//		all[i] = *obj
//
//		// Yaml View
//		AddToYamlView(item, yv)
//	}
//	return all
//}

func (c *Client) EventTypes(namespace string, yv *[]YamlView) []eventingv1beta1.EventType {
	gvr := schema.GroupVersionResource{
		Group:    "eventing.knative.dev",
		Version:  "v1beta1",
		Resource: "eventtypes",
	}
	like := eventingv1beta1.EventType{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List EventTypes, %v", err)
		return nil
	}

	all := make([]eventingv1beta1.EventType, len(list.Items))

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
