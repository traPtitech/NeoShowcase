package domain

import "context"

type CDService interface {
	Run()
	RegisterBuild(appID string)
	StartBuildLocal()
	SyncDeploymentsLocal()
	Stop(ctx context.Context) error
}
