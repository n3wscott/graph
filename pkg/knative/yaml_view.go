package knative

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type YamlView struct {
	ID    string
	Title string
	Yaml  string
}

func ToYamlViewURL(name, kind, apiVersion string) string {
	return fmt.Sprintf("#%s-%s-%s", name, kind, apiVersion)
}

func AddToYamlView(item unstructured.Unstructured, yv *[]YamlView) {
	if yv != nil {
		if yb, err := yaml.Marshal(item.Object); err == nil {
			*yv = append(*yv, YamlView{
				ID:    fmt.Sprintf("%s-%s-%s", item.GetName(), item.GetKind(), item.GetAPIVersion()),
				Title: fmt.Sprintf("%s %s.%s", item.GetName(), item.GetKind(), item.GetAPIVersion()),
				Yaml:  string(yb),
			})
		}
	}
}
