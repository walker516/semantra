package ctrl

import "api/internal/flow"

type Ctrls struct {
	Feature *FeatureCtrl
	// Tag    *TagCtrl
}

func NewCtrls(f *flow.Flows) *Ctrls {
	return &Ctrls{
		Feature: NewFeatureCtrl(f.FeatureCmd, f.Feature),
	}
}
