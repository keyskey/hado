package manifest

// manifestYAMLDoc maps dotted YAML paths (as produced by manifestYAMLPaths) to human descriptions.
// Keep in sync with struct fields in types.go; refdoc_test fails if keys drift.
var manifestYAMLDoc = map[string]string{
	"version": "Manifest スキーマの版。現行は `v1` を使う。",

	"service":      "評価対象サービスの識別子（ブロック全体は任意）。",
	"service.id":   "サービス ID。未指定時は `target` で `service.name` と同じにできる。",
	"service.name": "サービス名。",

	"standard":    "適用する Readiness Standard への参照（ブロック）。",
	"standard.id": "Standard のファイル名（例: `web-service.yaml`）またはパス。`standards-dir` / manifest 隣の `standards/` から解決される。",

	"evidence": "本番準備の証跡宣言。ゲートごとに必要なブロックだけでよい（各サブブロックは多くが `omitempty`）。",

	"evidence.coverage":                "カバレッジ成果物と adapter（ブロック）。C0/C1 ゲートがある standard で必要。",
	"evidence.coverage.inputs":         "`CoverageInput` の配列。",
	"evidence.coverage.inputs.adapter": "パーサ名。`hado-json` / `go-coverprofile` / `gobce-json` など（実装は `internal/coverage`）。",
	"evidence.coverage.inputs.path":    "リポジトリまたは manifest 相対の成果物パス。",

	"evidence.operations":         "運用責任と障害対応の入口（ブロック）。",
	"evidence.operations.owner":   "オーナー（チーム名・Slack チャンネル等）。`operations.owner_exists` で非空判定。",
	"evidence.operations.runbook": "Runbook の URL またはパス。`operations.runbook_exists` で非空判定。",

	"evidence.observability":           "観測可能性の証跡（ブロック）。SLO / モニター / ダッシュボードは **ベンダー UI 等で辿れる URL** のリストで宣言する（監査・運用オペ向け）。",
	"evidence.observability.slos":       "SLO / SLI への名前付きリンクの配列。`observability.slo_exists` はいずれか 1 件の `url`（trim 後非空）で PASS。",
	"evidence.observability.slos.name":  "人間可読な表示名（任意）。",
	"evidence.observability.slos.url":   "ブラウザで開ける SLO の URL（例: Datadog SLO の管理画面）。",
	"evidence.observability.monitors":   "モニターへの名前付きリンクの配列。`observability.monitor_exists` はいずれか 1 件の `url` で PASS。",
	"evidence.observability.monitors.name": "人間可読な表示名（任意）。",
	"evidence.observability.monitors.url": "モニターの URL（例: Datadog monitor）。",
	"evidence.observability.dashboards":    "ダッシュボードへの名前付きリンクの配列。`observability.dashboard_exists` はいずれか 1 件の `url` で PASS。",
	"evidence.observability.dashboards.name": "人間可読な表示名（任意）。",
	"evidence.observability.dashboards.url":  "ダッシュボードの URL。",

	"evidence.infra":                 "インフラ関連の参照（ブロック）。",
	"evidence.infra.deployment_spec": "デプロイ仕様の参照（パス・URL・カタログ ID）。`infra.deployment_spec_exists`。",

	"evidence.release":                        "リリース・ロールバック（ブロック）。",
	"evidence.release.rollback_plan":          "ロールバック手順の参照。`release.rollback_plan_exists`。",
	"evidence.release.automation":             "自動リリースパイプライン（ブロック）。",
	"evidence.release.automation.workflow_refs": "ワークフロー識別子のリスト（文字列の配列）。1 件以上非空で `release.automation_declared`。",
	"evidence.release.automation.systems":       "任意メタデータ（例: `github_actions`）。現行ゲートでは未使用。",
}
