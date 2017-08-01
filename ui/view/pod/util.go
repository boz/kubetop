package pod

import "k8s.io/api/core/v1"

func abbreviatePodConditionType(ct v1.PodConditionType) string {
	switch ct {
	case v1.PodScheduled:
		return "S"
	case v1.PodReady:
		return "R"
	case v1.PodInitialized:
		return "I"
	default:
		return "?"
	}
}
