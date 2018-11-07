package v1alpha1

import "time"

func (s DrainerConfigStatus) HasDrainedCondition() bool {
	return hasDrainerConfigCondition(s.Conditions, DrainerConfigStatusStatusTrue, DrainerConfigStatusTypeDrained)
}

func (s DrainerConfigStatus) HasTimeoutCondition() bool {
	return hasDrainerConfigCondition(s.Conditions, DrainerConfigStatusStatusTrue, DrainerConfigStatusTypeTimeout)
}

func (s DrainerConfigStatus) NewDrainedCondition() DrainerConfigStatusCondition {
	return DrainerConfigStatusCondition{
		LastTransitionTime: DeepCopyTime{time.Now()},
		Status:             DrainerConfigStatusStatusTrue,
		Type:               DrainerConfigStatusTypeDrained,
	}
}

func (s DrainerConfigStatus) NewTimeoutCondition() DrainerConfigStatusCondition {
	return DrainerConfigStatusCondition{
		LastTransitionTime: DeepCopyTime{time.Now()},
		Status:             DrainerConfigStatusStatusTrue,
		Type:               DrainerConfigStatusTypeTimeout,
	}
}

func hasDrainerConfigCondition(conditions []DrainerConfigStatusCondition, s string, t string) bool {
	for _, c := range conditions {
		if c.Status == s && c.Type == t {
			return true
		}
	}

	return false
}
