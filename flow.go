package flow

import (
	"api/internal/flow/command"
	"api/internal/flow/query"
	"api/internal/repo"
)

type Flows struct {
	// Command
	FeatureCmd command.FeatureCommandFlow
	// Query
	Feature query.FeatureQueryFlow
	// Tag    query.TagQueryFlow
}

func NewFlows(r *repo.Repositories) *Flows {
	return &Flows{
		FeatureCmd: command.NewFeatureCommand(r.Feature),
		Feature:    query.NewFeatureQuery(r.Feature),
	}
}
