package manifest

// AutomationDeclared reports whether manifest declares at least one non-empty release workflow reference.
func (r ReleaseEvidence) AutomationDeclared() bool {
	for _, ref := range r.Automation.WorkflowRefs {
		if !EvidenceUnset(ref) {
			return true
		}
	}
	return false
}
