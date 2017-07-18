package backend

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func ObjectID(obj metav1.Object) string {
	return obj.GetNamespace() + "/" + obj.GetName()
}
