package manifest

import "testing"

func TestReleaseEvidenceAutomationDeclared(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		rel  ReleaseEvidence
		want bool
	}{
		{name: "empty", rel: ReleaseEvidence{}, want: false},
		{name: "nil refs", rel: ReleaseEvidence{Automation: ReleaseAutomationEvidence{}}, want: false},
		{name: "whitespace only", rel: ReleaseEvidence{Automation: ReleaseAutomationEvidence{WorkflowRefs: []string{" ", ""}}}, want: false},
		{name: "one path", rel: ReleaseEvidence{Automation: ReleaseAutomationEvidence{WorkflowRefs: []string{".github/workflows/release.yml"}}}, want: true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.rel.AutomationDeclared(); got != tc.want {
				t.Fatalf("AutomationDeclared() = %v, want %v", got, tc.want)
			}
		})
	}
}
