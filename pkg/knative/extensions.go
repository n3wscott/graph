package knative

import (
	"github.com/knative/eventing/pkg/apis/duck/v1alpha1"
	"log"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func (c *Client) SourceCRDs() []apiextensions.CustomResourceDefinition {
	events := c.GetCRDs("eventing.knative.dev/source=true")
	ducks := c.GetCRDs("duck.knative.dev/source=true")

	contains := make(map[string]bool, 0)

	all := make([]apiextensions.CustomResourceDefinition, 0)

	for _, crd := range append(events, ducks...) {
		if _, ok := contains[crd.Name]; ok {
			continue
		}
		all = append(all, crd)
		contains[crd.Name] = true
	}

	return all
}

func (c *Client) ChannelCRDs() []apiextensions.CustomResourceDefinition {
	return c.GetCRDs("messaging.knative.dev/subscribable")
}

func (c *Client) GetCRDs(label string) []apiextensions.CustomResourceDefinition {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1beta1",
		Resource: "customresourcedefinitions",
	}
	like := apiextensions.CustomResourceDefinition{}

	list, err := c.dc.Resource(gvr).List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		log.Printf("Failed to List CRDs for %q, %v", label, err)
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
	for _, crd := range crds {
		for _, v := range crd.Spec.Versions {
			if !v.Served {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    crd.Spec.Group,
				Version:  v.Name,
				Resource: crd.Spec.Names.Plural,
			}
			gvrs = append(gvrs, gvr)
		}
	}
	return gvrs
}

func (c *Client) Addressable(namespace string, gvr schema.GroupVersionResource) []duckv1.AddressableType {
	like := duckv1.AddressableType{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
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

func (c *Client) Source(namespace string, gvr schema.GroupVersionResource, yv *[]YamlView) []duckv1.Source {
	like := duckv1.Source{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Source, [%+v], %v", gvr, err)
		return nil
	}

	all := make([]duckv1.Source, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj

		AddToYamlView(item, yv)
	}
	return all
}

func (c *Client) Channelable(namespace string, gvr schema.GroupVersionResource, yv *[]YamlView) []v1alpha1.Channelable {
	like := v1alpha1.Channelable{}

	list, err := c.dc.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Printf("Failed to List Channeable, [%+v], %v", gvr, err)
		return nil
	}

	all := make([]v1alpha1.Channelable, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			log.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj

		AddToYamlView(item, yv)
	}
	return all
}
