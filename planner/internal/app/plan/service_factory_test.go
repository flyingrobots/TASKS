package plan

import "testing"

func TestNewDefaultServiceProvidesAdapters(t *testing.T) {
	svc := NewDefaultService()
	if svc.BuildTasks == nil {
		t.Error("BuildTasks adapter is nil")
	}
	if svc.AnalyzeRepo == nil {
		t.Error("AnalyzeRepo adapter is nil")
	}
	if svc.ResolveDeps == nil {
		t.Error("ResolveDeps adapter is nil")
	}
	if svc.BuildDAG == nil {
		t.Error("BuildDAG adapter is nil")
	}
	if svc.BuildCoordinator == nil {
		t.Error("BuildCoordinator adapter is nil")
	}
	if svc.ValidateTasks == nil {
		t.Error("ValidateTasks adapter is nil")
	}
	if svc.ValidateDAG == nil {
		t.Error("ValidateDAG adapter is nil")
	}
	if svc.BuildWaves == nil {
		t.Error("BuildWaves adapter is nil")
	}
	if svc.WriteArtifacts == nil {
		t.Error("WriteArtifacts adapter is nil")
	}
	if svc.NewValidatorRunner == nil {
		t.Error("NewValidatorRunner adapter is nil")
	}
}
