package datadog

// Monitor is a subset of Datadog monitor API fields needed for discovery and linking.
type Monitor struct {
	ID   int64    `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
