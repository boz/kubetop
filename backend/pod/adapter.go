package pod

import (
	"fmt"

	"github.com/boz/kubetop/util"
	"k8s.io/client-go/pkg/api/v1"
)

type Adapter interface {
	FromResource(interface{}) (Pod, error)
	FromResourceList([]interface{}) ([]Pod, error)
	ToResource(Pod) (*v1.Pod, error)
}

type adapter struct {
	env util.Env
}

func newAdapter(env util.Env) *adapter {
	return &adapter{env}
}

func (a adapter) ToResource(p Pod) (*v1.Pod, error) {
	return p.Resource(), nil
}

func (a adapter) FromResource(obj interface{}) (Pod, error) {
	switch obj := obj.(type) {
	case *v1.Pod:
		return newPod(a.env, obj), nil
	default:
		return nil, fmt.Errorf("invalid type: %T", obj)
	}
}

func (a adapter) FromResourceList(objs []interface{}) ([]Pod, error) {
	pods := make([]Pod, 0, len(objs))
	for _, obj := range objs {
		pod, err := a.FromResource(obj)
		if err != nil {
			return nil, err
		}
		pods = append(pods, pod)
	}
	return pods, nil
}