package manifest

import "strings"

// AutomationDeclared reports whether manifest declares at least one non-empty release workflow reference.
func (r ReleaseEvidence) AutomationDeclared() bool {
	for _, ref := range r.Automation.WorkflowRefs {
		if strings.TrimSpace(ref) != "" {
			return true
		}
	}
	return false
}
