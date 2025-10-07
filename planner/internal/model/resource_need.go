package model

// ResourceNeed defines a requirement for a shared, limited resource.
type ResourceNeed struct {
	Name   string `json:"name"`
	Units  int    `json:"units,omitempty"`
	Access string `json:"access,omitempty"`
}
