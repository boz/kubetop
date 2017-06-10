package service

import (
	"fmt"

	"github.com/boz/kubetop/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

type Adapter interface {
	FromResource(metav1.Object) (Service, error)
	FromResourceList([]metav1.Object) ([]Service, error)
	ToResource(Service) (*v1.Service, error)
}

type adapter struct {
	env util.Env
}

func newAdapter(env util.Env) *adapter {
	return &adapter{env}
}

func (a adapter) ToResource(p Service) (*v1.Service, error) {
	return p.Resource(), nil
}

func (a adapter) FromResource(obj metav1.Object) (Service, error) {
	switch obj := obj.(type) {
	case *v1.Service:
		return newService(a.env, obj), nil
	default:
		return nil, fmt.Errorf("invalid type: %T", obj)
	}
}

func (a adapter) FromResourceList(objs []metav1.Object) ([]Service, error) {
	services := make([]Service, 0, len(objs))
	for _, obj := range objs {
		service, err := a.FromResource(obj)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}
