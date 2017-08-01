package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AbbreviateKind(kind string) string {
	switch kind {
	case "ReplicationController":
		return "RC"
	case "ReplicaSet":
		return "RS"
	case "Service":
		return "S"
	case "Node":
		return "N"
	default:
		return kind
	}
}

func FormatOwnerReference(ref v1.OwnerReference) string {
	if ref.Controller != nil && *ref.Controller {
		return fmt.Sprintf("+%v{%v}", AbbreviateKind(ref.Kind), ref.Name)
	}
	return fmt.Sprintf("%v{%v}", AbbreviateKind(ref.Kind), ref.Name)
}
