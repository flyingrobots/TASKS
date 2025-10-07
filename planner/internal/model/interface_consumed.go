package model

// InterfaceConsumed describes an interface that a task depends on.
type InterfaceConsumed struct {
	Name               string `json:"name"`
	VersionRequirement string `json:"version_requirement"`
	Type               string `json:"type"`
	Required           bool   `json:"required"`
}
