package sensors

import (
	"context"

	"github.com/bacalhau-project/bacalhau/pkg/compute"
	"github.com/bacalhau-project/bacalhau/pkg/compute/store"
	"github.com/bacalhau-project/bacalhau/pkg/models"
)

type RunningExecutionsInfoProviderParams struct {
	Name          string
	BackendBuffer *compute.ExecutorBuffer
}

// RunningExecutionsInfoProvider provides DebugInfo about the currently running executions.
// The info can be used for logging, metric, or to handle /debug API implementation.
type RunningExecutionsInfoProvider struct {
	name          string
	backendBuffer *compute.ExecutorBuffer
}

func NewRunningExecutionsInfoProvider(params RunningExecutionsInfoProviderParams) *RunningExecutionsInfoProvider {
	return &RunningExecutionsInfoProvider{
		name:          params.Name,
		backendBuffer: params.BackendBuffer,
	}
}

func (r RunningExecutionsInfoProvider) GetDebugInfo(ctx context.Context) (models.DebugInfo, error) {
	executions := r.backendBuffer.RunningExecutions()
	summaries := make([]store.ExecutionSummary, 0, len(executions))
	for _, execution := range executions {
		summaries = append(summaries, execution.ToSummary())
	}

	return models.DebugInfo{
		Component: r.name,
		Info:      summaries,
	}, nil
}

// compile-time check that we implement the interface
var _ models.DebugInfoProvider = (*RunningExecutionsInfoProvider)(nil)
