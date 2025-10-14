package testutil

import m "github.com/james/tasks-planner/internal/model"

// StubTasksFile returns a minimal TasksFile suitable for tests.
func StubTasksFile() *m.TasksFile {
	tf := &m.TasksFile{}
	tf.Meta.Version = "v8"
	tf.Meta.MinConfidence = 0.7
	tf.Tasks = []m.Task{{
		ID:        "T001",
		FeatureID: "F001",
		Title:     "Do thing",
		Duration: m.DurationPERT{
			Optimistic:  1,
			MostLikely:  2,
			Pessimistic: 3,
		},
		DurationUnit: "hours",
		AcceptanceChecks: []m.AcceptanceCheck{{
			Type: "command",
			Cmd:  "echo ok",
		}},
	}}
	return tf
}
