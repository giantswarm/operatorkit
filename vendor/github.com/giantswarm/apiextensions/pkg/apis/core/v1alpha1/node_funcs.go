package v1alpha1

import "time"

const (
	NodeConfigStatusStatusTrue  = "True"
	NodeConfigStatusTypeDrained = "Drained"
)

func (s NodeConfigStatus) HasFinalCondition() bool {
	for _, c := range s.Conditions {
		if c.Type == NodeConfigStatusTypeDrained && c.Status == NodeConfigStatusStatusTrue {
			return true
		}
	}

	return false
}

func (s NodeConfigStatus) NewFinalCondition() NodeConfigStatusCondition {
	return NodeConfigStatusCondition{
		LastTransitionTime: DeepCopyTime{time.Now()},
		Status:             NodeConfigStatusStatusTrue,
		Type:               NodeConfigStatusTypeDrained,
	}
}
