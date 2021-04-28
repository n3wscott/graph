package knative

import (
	"context"
	"log"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func (c *Client) SourceCRDs() []apiextensions.CustomResourceDefinition {
	// kubectl get crd -l "eventing.knative.dev/source=true"

	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1",
		Resource: "customresourcedefinitions",
	}
	like := apiextensions.CustomResourceDefinition{}

	list, err := c.dc.Resource(gvr).List(context.Background(), metav1.ListOptions{LabelSelector: "eventing.knative.dev/source=true"})
	if err != nil {
		log.Printf("Failed to List Sources, %v", err)
		return nil
	}

	all := make([]apiextensions.CustomResourceDefinition, len(list.Items))

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

func crdsToGVR(crds []apiextensions.CustomResourceDefinition) []schema.GroupVersionResource {
	gvrs := make([]schema.GroupVersionResource, 0)
	log.Println("Source GVRs ----")
	for _, crd := range crds {
		for _, v := range crd.Spec.Versions {
			if !v.Storage {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    crd.Spec.Group,
				Version:  v.Name,
				Resource: crd.Spec.Names.Plural,
			}
			gvrs = append(gvrs, gvr)
			log.Printf(" %v", gvr)
		}
	}
	return gvrs
}

func (c *Client) Addressable(namespace string, gvr schema.GroupVersionResource) []duckv1.AddressableType {
	like := duckv1.AddressableType{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Addressables, [%+v], %v", gvr, err)
		return nil
	}

	all := make([]duckv1.AddressableType, len(list.Items))

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
