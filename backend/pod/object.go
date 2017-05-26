package pod

import (
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/pkg/api/v1"
)

type Pod interface {
	Resource() *v1.Pod
	Name() string
}

type pod struct {
	resource *v1.Pod
	env      util.Env
}

func newPod(env util.Env, resource *v1.Pod) *pod {
	return &pod{resource, env}
}

func (p *pod) Resource() *v1.Pod {
	return p.resource
}

func (p *pod) Name() string {
	return p.resource.GetName()
}
