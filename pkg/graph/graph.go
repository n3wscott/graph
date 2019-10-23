package graph

import (
	"fmt"
	"knative.dev/pkg/apis"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/tmc/dot"
	eventingv1alpha1 "knative.dev/eventing/pkg/apis/eventing/v1alpha1"
	messagingv1alpha1 "knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	duckv1alpha1 "github.com/n3wscott/graph/pkg/apis/duck/v1alpha1"
	"github.com/n3wscott/graph/pkg/knative"
)

type Graph struct {
	*dot.Graph
	nodes     map[string]*dot.Node
	subgraphs map[string]*dot.SubGraph
	dnsToKey  map[string]string // maps domain name to node key

	edgeCount   int
	rainbowEdge bool
}

func New(ns string) *Graph {
	g := dot.NewGraph("G")
	_ = g.Set("tooltip", "Graph View")
	_ = g.Set("shape", "box")
	_ = g.Set("label", "Triggers in "+ns)
	_ = g.Set("rankdir", "LR")

	//_ = g.Set("compound", "true")

	graph := &Graph{
		Graph:       g,
		nodes:       make(map[string]*dot.Node),
		subgraphs:   make(map[string]*dot.SubGraph),
		dnsToKey:    make(map[string]string),
		rainbowEdge: true,
	}

	return graph
}

func (g *Graph) newEdge(src, dst *dot.Node) *dot.Edge {
	e := dot.NewEdge(src, dst)
	if g.rainbowEdge {
		color := colors[g.edgeCount%len(colors)]
		_ = e.Set("color", color)
		g.edgeCount++
	}
	return e
}

//func (g *Graph) AddChannel(channel eventingv1alpha1.Channel) {
//	ck := channelKey(channel.Name)
//	uri := channel.Status.Address.GetURL()
//	dns := strings.TrimSuffix((&uri).String(), "/")
//	cn := dot.NewNode("Channel " + channel.Name)
//
//	_ = cn.Set("URL", knative.ToYamlViewURL(channel.Name, channel.Kind, channel.APIVersion))
//	setNodeShapeForKind(cn, channel.Kind, channel.APIVersion)
//	setNodeColorForStatus(cn, channel.Status.Status)
//
//	_ = cn.Set("shape", "oval") // TODO move to setNodeShapeForKind
//	_ = cn.Set("label", "Ingress")
//
//	g.nodes[ck] = cn
//	g.dnsToKey[dns] = ck
//
//	cg := dot.NewSubgraph(fmt.Sprintf("cluster_%d", len(g.subgraphs)))
//	_ = cg.Set("label", fmt.Sprintf("Channel %s\n%s", channel.Name, dns))
//	g.subgraphs[ck] = cg
//	cg.AddNode(cn)
//	g.AddSubgraph(cg)
//}

// TODO: add channel ducktype.

func (g *Graph) AddInMemoryChannel(channel messagingv1alpha1.InMemoryChannel) {
	ck := inMemoryChannelKey(channel.Name)
	uri := channel.Status.Address.GetURL()
	dns := strings.TrimSuffix((&uri).String(), "/")
	cn := dot.NewNode("InMemoryChannel " + channel.Name)

	_ = cn.Set("URL", knative.ToYamlViewURL(channel.Name, channel.Kind, channel.APIVersion))
	setNodeShapeForKind(cn, channel.Kind, channel.APIVersion)
	setNodeColorForStatusV1Beta1(cn, channel.Status.Status)

	_ = cn.Set("shape", "oval") // TODO move to setNodeShapeForKind
	_ = cn.Set("label", "Ingress")

	g.nodes[ck] = cn
	g.dnsToKey[dns] = ck

	cg := dot.NewSubgraph(fmt.Sprintf("cluster_%d", len(g.subgraphs)))
	_ = cg.Set("label", fmt.Sprintf("InMemoryChannel %s\n%s", channel.Name, dns))
	g.subgraphs[ck] = cg
	cg.AddNode(cn)
	g.AddSubgraph(cg)
}

func (g *Graph) AddSubscription(subscription messagingv1alpha1.Subscription) {
	sk := subscriptionKey(subscription.Name)
	sn := dot.NewNode("Subscription " + subscription.Name)
	_ = sn.Set("URL", knative.ToYamlViewURL(subscription.Name, subscription.Kind, subscription.APIVersion))
	setNodeColorForStatusV1Beta1(sn, subscription.Status.Status)

	ck := gvkKey(subscription.Spec.Channel.GroupVersionKind(), subscription.Spec.Channel.Name)
	if cg, ok := g.subgraphs[ck]; !ok {
		g.AddNode(sn)
	} else {
		cg.AddNode(sn)
	}
	g.nodes[sk] = sn

	if sub := g.getOrCreateSubscriber(subscription.Spec.Subscriber); sub != nil {
		e := dot.NewEdge(sn, sub)
		_ = e.Set("dir", "both")
		setEdgeColorForStatusV1Beta1(e, subscription.Status.Status)
		g.AddEdge(e)
	}

	if rep := g.getOrCreateReply(subscription.Spec.Reply); rep != nil {
		e := g.newEdge(sn, rep)
		_ = e.Set("dir", "forward")
		g.AddEdge(e)
	}
}

func (g *Graph) AddBroker(broker eventingv1alpha1.Broker) {
	key := brokerKey(broker.Name)
	uri := broker.Status.Address.GetURL()
	dns := strings.TrimSuffix((&uri).String(), "/")
	bn := dot.NewNode("Broker " + dns)
	_ = bn.Set("shape", "oval")
	_ = bn.Set("label", "Ingress")
	_ = bn.Set("URL", knative.ToYamlViewURL(broker.Name, broker.Kind, broker.APIVersion))
	setNodeColorForStatusV1Beta1(bn, broker.Status.Status)

	g.nodes[key] = bn
	g.dnsToKey[dns] = key

	bg := dot.NewSubgraph(fmt.Sprintf("cluster_%d", len(g.subgraphs)))
	_ = bg.Set("label", fmt.Sprintf("Broker %s\n%s", broker.Name, dns))
	g.subgraphs[key] = bg
	bg.AddNode(bn)
	g.AddSubgraph(bg)
}

func (g *Graph) AddSource(source duckv1alpha1.SourceType) {
	key := gvkKey(source.GroupVersionKind(), source.Name)
	sn := dot.NewNode(fmt.Sprintf("Source %s\n%s\n%s", source.Name, source.Kind, source.GroupVersionKind().Group))
	_ = sn.Set("shape", "box")

	setNodeColorForStatusV1Beta1(sn, source.Status.Status)
	_ = sn.Set("URL", knative.ToYamlViewURL(source.Name, source.Kind, source.APIVersion))
	g.AddNode(sn)
	g.nodes[key] = sn

	fmt.Println("source ", source.Name, sn.String())

	sink := sinkDNS(source)

	if sink != "" {
		var bn *dot.Node
		var bk string
		var ok bool
		if bk, ok = g.dnsToKey[sink]; !ok {
			// TODO: unknown sink.
			bn = dot.NewNode("UnknownSink " + sink)
			g.AddNode(bn)
		} else {
			if bn, ok = g.nodes[bk]; !ok {
				// TODO: unknown broker.
				bn = dot.NewNode("UnknownSink " + sink)
				g.AddNode(bn)
			}
		}

		e := dot.NewEdge(sn, bn)
		setEdgeColorForStatusV1Beta1(e, source.Status.Status)
		if sg, ok := g.subgraphs[bk]; ok {
			// This is not working.
			_ = e.Set("lhead", sg.Name())
		}
		g.AddEdge(e)
	}
}

func (g *Graph) AddTrigger(trigger eventingv1alpha1.Trigger) {
	broker := trigger.Spec.Broker
	bk := brokerKey(broker)
	bn, ok := g.nodes[bk]
	if !ok {
		bn = dot.NewNode("UnknownBroker " + broker)
		g.AddNode(bn)
		g.nodes[bk] = bn
	}

	tn := dot.NewNode("Trigger " + trigger.Name)
	_ = tn.Set("shape", "box")
	_ = tn.Set("URL", knative.ToYamlViewURL(trigger.Name, trigger.Kind, trigger.APIVersion))
	setNodeColorForStatusV1Beta1(tn, trigger.Status.Status)

	if sg, ok := g.subgraphs[bk]; ok {
		sg.AddNode(tn)
	} else {
		g.AddNode(tn)
	}
	g.nodes[triggerKey(trigger.Name)] = tn

	if trigger.Spec.Filter != nil && trigger.Spec.Filter.Attributes != nil {
		filter := ""
		for k, v := range *trigger.Spec.Filter.Attributes {
			filter = fmt.Sprintf("%s\n%s=%s", filter, k, v)
		}
		_ = tn.Set("label", fmt.Sprintf("%s%s", tn.Name(), filter))
	}

	if sub := g.getOrCreateSubscriber(trigger.Spec.Subscriber); sub != nil {
		e := dot.NewEdge(tn, sub)
		_ = e.Set("dir", "both")
		setEdgeColorForStatusV1Beta1(e, trigger.Status.Status)
		fmt.Println("sub", sub, e)
		g.AddEdge(e)
	}
}

func (g *Graph) LoadKnService(service servingv1.Service) {
	key := servingKey(service.Kind, service.Name)

	var svc *dot.Node
	var ok bool
	label := ""
	if svc, ok = g.nodes[key]; !ok {
		label = fmt.Sprintf("%s\n%s\n%s",
			service.Name,
			service.Kind,
			service.GroupVersionKind().Group,
		)
		svc = dot.NewNode(label)

		_ = svc.Set("URL", knative.ToYamlViewURL(service.Name, service.Kind, service.APIVersion))
		setNodeShapeForKind(svc, service.Kind, service.APIVersion)
		setNodeColorForStatus(svc, service.Status.Status)

		//_ = svc.Set("shape", "septagon")

		g.nodes[key] = svc
		g.AddNode(svc)

		if service.Status.Address != nil && service.Status.Address.URL != nil {
			dns := service.Status.Address.URL.String()
			g.dnsToKey[dns] = key
			fmt.Println(key, "-->", dns)
		}
	}
}

func (g *Graph) AddKnService(service servingv1.Service) {
	config := service.Spec.ConfigurationSpec
	key := servingKey(service.Kind, service.Name)

	var svc *dot.Node
	var ok bool
	label := ""
	if svc, ok = g.nodes[key]; !ok {
		label = fmt.Sprintf("%s\n%s\n%s",
			service.Name,
			service.Kind,
			service.GroupVersionKind().Group,
		)
		svc = dot.NewNode(label)
		_ = svc.Set("URL", knative.ToYamlViewURL(service.Name, service.Kind, service.APIVersion))
		setNodeShapeForKind(svc, service.Kind, service.APIVersion)
		setNodeColorForStatus(svc, service.Status.Status)

		//_ = svc.Set("shape", "septagon")

		g.nodes[key] = svc
		g.AddNode(svc)
	}

	//	fmt.Println(service, "kn svc:", svc)

	for _, env := range config.Template.Spec.Containers[0].Env {
		switch env.Name {
		case "SINK":
			fallthrough
		case "TARGET":
			// Assume full dns name.
			target := g.getOrCreateSink(env.Value)
			e := dot.NewEdge(svc, target)
			setEdgeColorForStatus(e, service.Status.Status)
			g.AddEdge(e)
		}
	}
}

func (g *Graph) AddSequence(seq messagingv1alpha1.Sequence) {

	key := sequenceKey(seq.Name)

	uri := seq.Status.Address.GetURL()
	dns := strings.TrimSuffix((&uri).String(), "/")

	sg := dot.NewSubgraph(fmt.Sprintf("cluster_%d", len(g.subgraphs)))
	_ = sg.Set("label", fmt.Sprintf("Sequence %s\n%s", seq.Name, dns))
	//	_ = sg.Set("rankdir", "BT")

	g.dnsToKey[dns] = key
	sn := dot.NewNode("Sequence " + dns)
	_ = sn.Set("label", "Start")
	_ = sn.Set("URL", knative.ToYamlViewURL(seq.Name, seq.Kind, seq.APIVersion))
	setNodeColorForStatusV1Beta1(sn, seq.Status.Status)

	g.nodes[key] = sn
	sg.AddNode(sn)

	previousNode := sn

	for num, step := range seq.Spec.Steps {
		stepKey := sequenceStepKey(seq.Name, num)
		stepn := dot.NewNode(stepKey)
		_ = stepn.Set("label", fmt.Sprintf("Step %d", num))
		_ = stepn.Set("shape", "box")

		// Add to seq subgraph.
		sg.AddNode(stepn)

		g.nodes[stepKey] = stepn

		if sub := g.getOrCreateSubscriber(&step); sub != nil {
			e := dot.NewEdge(stepn, sub)
			_ = e.Set("dir", "both")
			setEdgeColorForStatusV1Beta1(e, seq.Status.Status)
			g.AddEdge(e)
		}

		e := dot.NewEdge(previousNode, stepn)
		setEdgeColorForStatusV1Beta1(e, seq.Status.Status)
		g.AddEdge(e)
		previousNode = stepn
	}

	if seq.Spec.Reply != nil {
		replyn := dot.NewNode("Reply " + dns)
		_ = replyn.Set("label", "Reply")
		//_ = replyn.Set("rank", "max")
		//g.nodes[] = rn
		sg.AddNode(replyn)

		// TODO where this points.
		e := dot.NewEdge(previousNode, replyn)
		setEdgeColorForStatusV1Beta1(e, seq.Status.Status)
		g.AddEdge(e)

		rk := gvkKey(seq.Spec.Reply.GroupVersionKind(), seq.Spec.Reply.Name)
		if rn, ok := g.nodes[rk]; ok {
			e := dot.NewEdge(replyn, rn)
			setEdgeColorForStatusV1Beta1(e, seq.Status.Status)
			g.AddEdge(e)
		}
	}

	g.subgraphs[key] = sg
	g.AddSubgraph(sg)

}

func setNodeShapeForKind(node *dot.Node, kind, apiVersion string) {
	if strings.HasPrefix(apiVersion, "serving.knative.dev") {
		switch kind {
		case "Service":
			_ = node.Set("shape", "septagon")
		}
	}
}

func getColorMapForStatus(status duckv1.Status) map[string]string {
	cond := status.GetCondition(apis.ConditionReady)
	attrs := make(map[string]string)
	if cond.IsTrue() {
		attrs["color"] = "black"
		attrs["tooltip"] = fmt.Sprintf("Ready as of %s", cond.LastTransitionTime.Inner.String())
	} else if cond.IsUnknown() {
		attrs["color"] = "darkorange2"
		attrs["tooltip"] = fmt.Sprintf("[%s] %s: %s", cond.Status, cond.Reason, cond.Message)
	} else if cond.IsFalse() {
		attrs["color"] = "deeppink"
		attrs["tooltip"] = fmt.Sprintf("[%s] %s: %s", cond.Status, cond.Reason, cond.Message)
	}
	return attrs
}

func getColorMapForStatusV1Beta1(status duckv1beta1.Status) map[string]string {
	cond := status.GetCondition(apis.ConditionReady)
	attrs := make(map[string]string)
	if cond.IsTrue() {
		attrs["color"] = "black"
		attrs["tooltip"] = fmt.Sprintf("Ready as of %s", cond.LastTransitionTime.Inner.String())
	} else if cond.IsUnknown() {
		attrs["color"] = "darkorange2"
		attrs["tooltip"] = fmt.Sprintf("[%s] %s: %s", cond.Status, cond.Reason, cond.Message)
	} else if cond.IsFalse() {
		attrs["color"] = "deeppink"
		attrs["tooltip"] = fmt.Sprintf("[%s] %s: %s", cond.Status, cond.Reason, cond.Message)
	}
	return attrs
}

func setNodeColorForStatus(node *dot.Node, status duckv1.Status) {
	_ = node.Set("fillcolor", "white")
	_ = node.Set("style", "filled")
	for name, value := range getColorMapForStatus(status) {
		_ = node.Set(name, value)
	}
}

func setNodeColorForStatusV1Beta1(node *dot.Node, status duckv1beta1.Status) {
	_ = node.Set("fillcolor", "white")
	_ = node.Set("style", "filled")
	for name, value := range getColorMapForStatusV1Beta1(status) {
		_ = node.Set(name, value)
	}
}

func setEdgeColorForStatus(edge *dot.Edge, status duckv1.Status) {
	for name, value := range getColorMapForStatus(status) {
		_ = edge.Set(name, value)
	}
}

func setEdgeColorForStatusV1Beta1(edge *dot.Edge, status duckv1beta1.Status) {
	for name, value := range getColorMapForStatusV1Beta1(status) {
		_ = edge.Set(name, value)
	}
}

func (g *Graph) getOrCreateSink(uri string) *dot.Node {
	uri = strings.TrimSuffix(uri, "/")

	var node *dot.Node
	var key string
	var ok bool
	if key, ok = g.dnsToKey[uri]; !ok {
		node = dot.NewNode("UnknownSink " + uri)
		g.AddNode(node)
		g.nodes[key] = node
	}
	return g.nodes[key]
}

func (g *Graph) getOrCreateSubscriber(subscriber *messagingv1alpha1.SubscriberSpec) *dot.Node {
	key := "?"
	label := "?"

	if subscriber != nil {
		if subscriber.URI != nil {
			label = *subscriber.URI
			key = uriKey(*subscriber.URI)
		} else if subscriber.Ref != nil {
			label = fmt.Sprintf("%s\n%s\n%s",
				subscriber.Ref.Name,
				subscriber.Ref.Kind,
				subscriber.Ref.GroupVersionKind().Group,
			)
			key = refKey(
				subscriber.Ref.GroupVersionKind().Group,
				subscriber.Ref.Kind,
				subscriber.Ref.Name,
			)
		}
	}
	var sub *dot.Node
	var ok bool
	if sub, ok = g.nodes[key]; !ok {
		sub = dot.NewNode(label)
		if subscriber != nil && subscriber.Ref != nil {
			setNodeShapeForKind(sub, subscriber.Ref.Kind, subscriber.Ref.APIVersion)
		}

		g.nodes[key] = sub
		g.AddNode(sub)
	}
	return sub
}

func (g *Graph) getOrCreateReply(rep *messagingv1alpha1.ReplyStrategy) *dot.Node {
	if rep != nil && rep.Channel != nil {
		ck := gvkKey(rep.Channel.GroupVersionKind(), rep.Channel.Name)
		if cn, ok := g.nodes[ck]; !ok {
			cn = dot.NewNode("Unknown Channel " + rep.Channel.Name)
		} else {
			return cn
		}
	}
	return nil
}

func sinkDNS(source duckv1alpha1.SourceType) string {
	if source.Status.SinkURI != nil {
		return strings.TrimSuffix(*(source.Status.SinkURI), "/")
	}
	return ""
}

func channelKey(name string) string {
	return eventingKey("channel", name)
}

func inMemoryChannelKey(name string) string {
	return messagingKey("inmemorychannel", name)
}

func subscriptionKey(name string) string {
	return eventingKey("subscription", name)
}

func brokerKey(name string) string {
	return eventingKey("broker", name)
}

func triggerKey(name string) string {
	return eventingKey("trigger", name)
}

func sequenceKey(name string) string {
	return messagingKey("sequence", name)
}

func sequenceStepKey(name string, step int) string {
	return messagingKey("sequencestep", fmt.Sprintf("%s-%d", name, step))
}

func gvkKey(gvk schema.GroupVersionKind, name string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s/%s", gvk.Group, gvk.Kind, name))
}

func key(group, kind, name string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s/%s", group, kind, name))
}

func uriKey(uri string) string {
	return strings.ToLower(fmt.Sprintf("uri/%s", uri))
}

func refKey(group, kind, name string) string {
	return strings.ToLower(fmt.Sprintf("%s/%s/%s", group, kind, name))
}

func eventingKey(kind, name string) string {
	return key("eventing.knative.dev", kind, name)
}

func messagingKey(kind, name string) string {
	return key("messaging.knative.dev", kind, name)
}

func servingKey(kind, name string) string {
	return key("serving.knative.dev", kind, name)
}
