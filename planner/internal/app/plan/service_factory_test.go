package plan

import "testing"

func TestNewDefaultServiceProvidesAdapters(t *testing.T) {
	svc := NewDefaultService()
	if svc.BuildTasks == nil || svc.AnalyzeRepo == nil || svc.ResolveDeps == nil || svc.BuildDAG == nil || svc.BuildCoordinator == nil || svc.ValidateTasks == nil || svc.ValidateDAG == nil || svc.BuildWaves == nil || svc.WriteArtifacts == nil || svc.NewValidatorRunner == nil {
		t.Fatalf("default service missing adapters: %+v", svc)
	}
}
