package analysis

// FileCensusCounts is a compact summary of file counts for embedding in artifacts.
type FileCensusCounts struct {
    Files   int `json:"files"`
    GoFiles int `json:"go_files"`
}

