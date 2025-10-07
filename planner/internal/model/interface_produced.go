package model

// InterfaceProduced describes an interface that a task creates or modifies.
type InterfaceProduced struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
}
