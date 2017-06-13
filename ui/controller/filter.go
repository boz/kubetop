package controller

import (
	"reflect"

	"github.com/boz/kcache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Filter interface {
	kcache.Filter
	Equals(Filter) bool
}

type labelsSelector struct {
	filter labels.Set
}

func LabelsSelector(match map[string]string) Filter {
	return &labelsSelector{labels.Set(match)}
}

func (f *labelsSelector) Accept(obj metav1.Object) bool {
	if len(f.filter) == 0 {
		return false
	}
	return labels.AreLabelsInWhiteList(f.filter, labels.Set(obj.GetLabels()))
}

func (f *labelsSelector) Equals(other Filter) bool {
	if other, ok := other.(*labelsSelector); ok {
		return labels.Equals(f.filter, other.filter)
	}
	return false
}

type nsName struct {
	ns   string
	name string
}

type nsNameSelector map[nsName]bool

func NSNamesSelector(objs ...metav1.Object) Filter {
	set := make(map[nsName]bool)
	for _, obj := range objs {
		key := nsName{obj.GetNamespace(), obj.GetName()}
		set[key] = true
	}
	return nsNameSelector(set)
}

func (f nsNameSelector) Accept(obj metav1.Object) bool {
	key := nsName{obj.GetNamespace(), obj.GetName()}
	_, ok := f[key]
	return ok
}

func (f nsNameSelector) Equals(other Filter) bool {
	return reflect.DeepEqual(f, other)
}
