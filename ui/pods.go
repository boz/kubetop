package ui

import (
	"k8s.io/client-go/pkg/api/v1"

	"github.com/boz/kcache"
	"github.com/boz/kubetop/backend/pod"
	"github.com/boz/kubetop/ui/elements/table"
	"github.com/boz/kubetop/ui/theme"
	"github.com/boz/kubetop/util"
)

/*

namespace name containers owner ctime phase condition message

- Labels
- Annotations
- CreationTimestamp
- OwnerReferences
- ResourceVersion
- ClusterName

- PodPhase
- Conditions
- Message
- Reason
- HostIP
- PodIP
- StartTime
- ContainerStatuses

- Volumes
- Containers
- RestartPolicy
- TerminationGracePeriodSeconds
- ActiveDeadlineSeconds
- DNSPolicy
- NodeSelector
*/

type podIndexBuilder struct {
	env util.Env
	ds  pod.BaseDatasource
}

func newPodIndexBuilder(env util.Env, ds pod.BaseDatasource) IndexBuilder {
	return &podIndexBuilder{env, ds}
}

func (b *podIndexBuilder) Model() []table.TH {
	header := []table.TH{
		table.NewTH("ns", "Namespace", true, 0),
		table.NewTH("name", "Name", true, 1),
		table.NewTH("version", "Version", true, -1),
		table.NewTH("phase", "Phase", true, -1),
		table.NewTH("conditions", "Conditions", true, -1),
		table.NewTH("message", "Message", true, -1),
	}
	return header
}

func (b *podIndexBuilder) Create(w IndexWidget, donech <-chan struct{}) IndexProvider {
	return newPodIndexProvider(b.env, b.ds, w, donech)
}

type podIndexProvider struct {
	widget IndexWidget
	donech <-chan struct{}
	env    util.Env
	sub    pod.Subscription
}

func newPodIndexProvider(env util.Env, ds pod.BaseDatasource, w IndexWidget, donech <-chan struct{}) IndexProvider {
	p := &podIndexProvider{
		widget: w,
		donech: donech,
		env:    env,
		sub:    ds.Subscribe(),
	}
	go p.run()
	return p
}

func (p *podIndexProvider) Stop() {
	p.sub.Close()
}

func (p *podIndexProvider) run() {
	defer p.sub.Close()

	readych := p.sub.Ready()

	for {
		select {
		case <-p.donech:
		case <-p.sub.Closed():
			return
		case <-readych:
			readych = nil
			p.doInitialize()
		case ev, ok := <-p.sub.Events():
			if !ok {
				return
			}
			p.handleEvent(ev)
		}
	}
}

func (p *podIndexProvider) doInitialize() {
	pods, err := p.sub.List()
	if err != nil {
	}
	rows := make([]table.TR, 0, len(pods))
	for _, pod := range pods {
		rows = append(rows, p.generateRow(pod))
	}
	p.widget.ResetRows(rows)
}

func (p *podIndexProvider) handleEvent(ev pod.Event) {
	obj := ev.Resource()

	switch ev.Type() {
	case kcache.EventTypeDelete:
		p.widget.RemoveRow(obj.ID())
	case kcache.EventTypeCreate:
		p.widget.InsertRow(p.generateRow(obj))
	case kcache.EventTypeUpdate:
		p.widget.UpdateRow(p.generateRow(obj))
	}
}

func (p *podIndexProvider) generateRow(pod pod.Pod) table.TR {

	stat := pod.Resource().Status

	phase := string(stat.Phase)
	message := stat.Message

	conditions := ""
	for _, c := range stat.Conditions {
		conditions += string(c.Type)[0:1]
		switch c.Status {
		case v1.ConditionTrue:
			conditions += "+"
		case v1.ConditionFalse:
			conditions += "-"
		case v1.ConditionUnknown:
			conditions += "?"
		}
	}

	//message += strings.Repeat("x", 50)

	cols := []table.TD{
		table.NewTD("ns", pod.Resource().GetNamespace(), theme.LabelNormal),
		table.NewTD("name", pod.Resource().GetName(), theme.LabelNormal),
		table.NewTD("version", pod.Resource().GetResourceVersion(), theme.LabelNormal),
		table.NewTD("phase", phase, theme.LabelNormal),
		table.NewTD("conditions", conditions, theme.LabelNormal),
		table.NewTD("message", message, theme.LabelNormal),
	}
	return table.NewTR(pod.ID(), cols)
}
