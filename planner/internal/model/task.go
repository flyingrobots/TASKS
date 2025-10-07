package model

// Task represents a single, well-defined unit of work to be executed.
type Task struct {
	ID                 string              `json:"id"`
	FeatureID          string              `json:"feature_id"`
	Title              string              `json:"title"`
	Description        string              `json:"description,omitempty"`
	Category           string              `json:"category,omitempty"`
	Duration           DurationPERT        `json:"duration"`
	DurationUnit       string              `json:"durationUnits"`
	InterfacesProduced []InterfaceProduced `json:"interfaces_produced,omitempty"`
	InterfacesConsumed []InterfaceConsumed `json:"interfaces_consumed,omitempty"`
	AcceptanceChecks   []AcceptanceCheck   `json:"acceptance_checks"`
	Evidence           []Evidence          `json:"source_evidence"`
	Resources          struct {
		Exclusive []string       `json:"exclusive,omitempty"`
		Limited   []ResourceNeed `json:"limited,omitempty"`
	} `json:"resources"`
	ExecutionLogging struct {
		Format         string   `json:"format"`
		RequiredFields []string `json:"required_fields"`
	} `json:"execution_logging"`
	Compensation struct {
		Idempotent  bool   `json:"idempotent"`
		RollbackCmd string `json:"rollback_cmd,omitempty"`
	} `json:"compensation"`
}
