package knative

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	messagingv1alpha1 "github.com/knative/eventing/pkg/apis/messaging/v1alpha1"
)

func (c *Client) Sequences(namespace string) []messagingv1alpha1.Sequence {
	gvr := schema.GroupVersionResource{
		Group:    "messaging.knative.dev",
		Version:  "v1alpha1",
		Resource: "sequences",
	}
	like := messagingv1alpha1.Sequence{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Sequences, %v", err)
		return nil
	}

	all := make([]messagingv1alpha1.Sequence, len(list.Items))

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

func (c *Client) InMemoryChannels(namespace string) []messagingv1alpha1.InMemoryChannel {
	gvr := schema.GroupVersionResource{
		Group:    "messaging.knative.dev",
		Version:  "v1alpha1",
		Resource: "inmemorychannels",
	}
	like := messagingv1alpha1.InMemoryChannel{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List InMemoryChannels, %v", err)
		return nil
	}

	all := make([]messagingv1alpha1.InMemoryChannel, len(list.Items))

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
