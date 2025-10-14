package model

// ArtifactMeta captures shared metadata fields for planner artifacts.
type ArtifactMeta struct {
	Version      string `json:"version"`
	ArtifactHash string `json:"artifact_hash"`
}

// FeaturesArtifact is the serialized features.json contract.
type FeaturesArtifact struct {
	Meta     ArtifactMeta   `json:"meta"`
	Features []FeatureEntry `json:"features"`
}

// FeatureEntry summarizes a single feature in features.json.
type FeatureEntry struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// WavesArtifact models waves.json with versioned metadata.
type WavesArtifact struct {
	Meta  WavesMeta  `json:"meta"`
	Waves [][]string `json:"waves"`
}

// WavesMeta extends artifact metadata with the plan identifier reference.
type WavesMeta struct {
	Version      string `json:"version"`
	PlanID       string `json:"planId"`
	ArtifactHash string `json:"artifact_hash"`
}

// TitlesArtifact carries canonical task titles for DOT rendering.
type TitlesArtifact struct {
	Meta   ArtifactMeta      `json:"meta"`
	Titles map[string]string `json:"titles"`
}
