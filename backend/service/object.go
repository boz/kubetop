package service

import (
	"github.com/boz/kubetop/util"
	"k8s.io/client-go/pkg/api/v1"
)

type Service interface {
	ID() string
	Resource() *v1.Service
	Name() string
}

type service struct {
	resource *v1.Service
	env      util.Env
}

func newService(env util.Env, resource *v1.Service) *service {
	return &service{resource, env}
}

func (p *service) ID() string {
	return p.resource.Namespace + "/" + p.resource.Name
}

func (p *service) Resource() *v1.Service {
	return p.resource
}

func (p *service) Name() string {
	return p.resource.GetName()
}
