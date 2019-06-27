package knative

import (
	"k8s.io/client-go/dynamic"
)

func New(dc dynamic.Interface) *Client {
	c := &Client{
		dc: dc,
	}
	return c
}

type Client struct {
	dc dynamic.Interface
}
