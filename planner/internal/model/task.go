package model

// Task represents a single, well-defined unit of work to be executed.
type Task struct {
	ID                 string              `json:"id"`
	FeatureID          string              `json:"featureID"`
	Title              string              `json:"title"`
	Description        string              `json:"description,omitempty"`
	Category           string              `json:"category,omitempty"`
	Duration           DurationPERT        `json:"duration"`
	DurationUnit       string              `json:"durationUnits"`
	InterfacesProduced []InterfaceProduced `json:"interfacesProduced,omitempty"`
	InterfacesConsumed []InterfaceConsumed `json:"interfacesConsumed,omitempty"`
	AcceptanceChecks   []AcceptanceCheck   `json:"acceptanceChecks"`
	Evidence           []Evidence          `json:"sourceEvidence"`
	Resources          struct {
		Exclusive []string       `json:"exclusive,omitempty"`
		Limited   []ResourceNeed `json:"limited,omitempty"`
	} `json:"resources"`
	ExecutionLogging struct {
		Format         string   `json:"format"`
		RequiredFields []string `json:"requiredFields"`
	} `json:"executionLogging"`
	Compensation struct {
		Idempotent  bool   `json:"idempotent"`
		RollbackCmd string `json:"rollbackCmd,omitempty"`
	} `json:"compensation"`
}
