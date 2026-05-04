# 実装状況

**手書きで保守する。** 実装を変えたら Cursor の Rule（`.cursor/rules/hado-implementation-docs.mdc`）と Skill（`hado-doc-sync`）に従い、**同じ PR / セッションで**このファイルを更新する。

## CLI（`cmd/hado`）

| コマンド | 状態 |
| --- | --- |
| 引数なし | 一行ヘルプ |
| `version` / `-v` / `--version` | 実装済み |
| `target` | 実装済み（`--manifest` 必須。TTY / フラグで `service` / `standard`。既定で resolved standard に応じ **evidence のプレースホルダー**をマージ。`manifest.EvidencePlaceholder` と空は未入力扱いで `evaluate` も同様） |
| `evaluate` | 実装済み（設計上は `fire` に相当する判定を、当面は一括で実行） |

**設計（未実装）:** `hado charge`（manifest メタデータから evidence を自動補完）/ `hado fire`（判定のみ）の 2 段階は [overview.md](overview.md) と [architecture.md](architecture.md) に記載。実装時は本表を更新する。

`target` の主なフラグ（`cmd/hado/target.go` の `runTarget`）:

- `--manifest`（必須）
- `--service-name`（任意; 非 TTY では既存 manifest かフラグのどちらかが必要）
- `--service-id`（任意; 空のときは `service-name` と同じにできる）
- `--standard-id`（任意; 非 TTY では既存 manifest かフラグのどちらかが必要）
- `--standards-dir`（任意; プレースホルダー用の standard YAML を探すディレクトリ。既定は manifest と同じ階層の `standards/`）
- `--rewrite-placeholders`（既定 `true`。`false` で service/standard のみ更新し evidence は触らない）

`evaluate` は manifest の文字列が **`PLACEHOLDER`**（`manifest.EvidencePlaceholder`）または空のとき、該当 evidence は **未設定**として existence gate を評価します。

`evaluate` の主なフラグ（`cmd/hado/main.go` の `runEvaluate`）:

- `--standard`（必須）
- `--manifest`（任意）
- `--coverage-input`（繰り返し可; `<adapter>:<path>`。**指定時は manifest の `evidence.coverage.inputs` より優先**）
- `--output`：`text` または `json`（それ以外はエラー）

**Coverage 入力の必須条件:** Readiness Standard が `test.c0_coverage` または `test.c1_coverage` のいずれかを含む場合、`--coverage-input` か manifest の `evidence.coverage.inputs` のどちらかが必要。どちらも無いと `evaluate` はエラー終了（exit 2）。

終了コード: `0` = ready、`1` = blocked（required gate 失敗）、`2` = error（引数・読み込み・未対応 gate など）。

**未実装の例:** `--output markdown`、module runner、score / exception フィールド。

## 実装済みゲート（`internal/gate/evaluate.go` の `switch` 順）

required として宣言されているが、ここに無い gate id は **error**（`unsupported required gate`）になる。optional の未知 gate は無視。

- `test.c0_coverage`（`internal/standard` の `C0CoverageGateID`）
- `test.c1_coverage`
- `operations.owner_exists`
- `operations.runbook_exists`
- `observability.slo_exists`（manifest `evidence.observability.slo` が非空）
- `observability.monitor_exists`（`evidence.observability.monitors` が非空）
- `observability.dashboard_exists`（`evidence.observability.dashboard` が非空）
- `infra.deployment_spec_exists`（`evidence.infra.deployment_spec` が非空; パス・URL・カタログ ID などの参照文字列として扱う）
- `release.rollback_plan_exists`（`evidence.release.rollback_plan` が非空）
- `release.automation_declared`（`evidence.release.automation.workflow_refs` に **実値**（`PLACEHOLDER` 以外）の要素が **1 件以上**。`systems` は任意のメタデータで現ゲートでは未使用）

## Coverage adapter（`internal/coverage/parse.go` の `ParseAdapterInput`）

`--coverage-input` および manifest の `evidence.coverage.inputs[].adapter` に使える文字列（`types.go` の `Format*` 定数と一致）:

- `hado-json`（正規化 JSON の `c0Coverage` / `c1Coverage`）
- `go-coverprofile`（C0 のみ）
- `gobce-json`（C0 / C1; `keyskey/gobce` の JSON）

## Manifest（`internal/manifest`）

- `evidence.operations`（owner, runbook）
- `evidence.observability`（`slo`, `monitors`, `dashboard` …各フィールドが該当 gate の「存在」判定に使われる）
- `evidence.infra`（`deployment_spec`）
- `evidence.release`（`rollback_plan`; `automation.workflow_refs`, 任意で `automation.systems`）
- `evidence.coverage.inputs`（`adapter`, `path`）

## MVP・ロードマップとの差（メモ）

計画全体は [roadmap.md](roadmap.md)。コードにまだ無い例:

- module runner、インフラ向け threshold 型 gate（例: PDB の数値比較）、Markdown レポート、GitHub PR 連携
- `test.uncovered_branch` など gobce findings の評価結果への載せ方
