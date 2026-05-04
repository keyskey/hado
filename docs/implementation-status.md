# 実装状況

**手書きで保守する。** 実装を変えたら Cursor の Rule（`.cursor/rules/hado-implementation-docs.mdc`）と Skill（`hado-doc-sync`）に従い、**同じ PR / セッションで**このファイルを更新する。

## CLI（`cmd/hado`）

| コマンド | 状態 |
| --- | --- |
| 引数なし | 一行ヘルプ |
| `version` / `-v` / `--version` | 実装済み |
| `target` | 実装済み（`--manifest` 必須。TTY / フラグで `service` / `standard`。既定で resolved standard に応じ **evidence のスキャフォールド**（空文字のキーなど）をマージ） |
| `charge` | 実装済み（`--manifest` 必須。coverage artifact の adapter/path を manifest `evidence.coverage.inputs` に不足分マージ。既存値は置換しない） |
| `fire` | 実装済み（`--manifest` 必須。判定専用。manifest の evidence を gate 評価して READY/BLOCKED/ERROR を返す） |

`evaluate` は廃止し、`target` / `charge` / `fire` に一本化した。

`target` の主なフラグ（`cmd/hado/target/run.go` の `Run`）:

- `--manifest`（必須）
- `--service-name`（任意; 非 TTY では既存 manifest かフラグのどちらかが必要）
- `--service-id`（任意; 空のときは `service-name` と同じにできる）
- `--standard-id`（任意; 非 TTY では既存 manifest かフラグのどちらかが必要）
- `--standards-dir`（任意; スキャフォールド用の standard YAML を探すディレクトリ。既定は manifest と同じ階層の `standards/`）
- `--rewrite-placeholders`（既定 `true`。`false` で service/standard のみ更新し evidence は触らない）

`fire` は manifest の該当フィールドが **空（前後の空白を除いた長さ 0）** のとき、existence gate では **未設定**として評価します。

`charge` の主なフラグ（`cmd/hado/charge/run.go`）:

- `--manifest`（必須）
- `--standard`（任意。未指定時は manifest の `standard.id` を利用）
- `--coverage-input`（繰り返し可; `<adapter>:<path>`。**指定時も既存 manifest 値は置換せず不足分だけマージ**）

`fire` の主なフラグ（`cmd/hado/fire/run.go`）:

- `--manifest`（必須）
- `--standard`（任意。未指定時は manifest の `standard.id` を利用、指定時は上書き）
- `--output`：`text` または `json`（それ以外はエラー）

`--output text` は各 gate の判定行に `severity` を表示し、FAIL 行には「リリース前に必須対応か / リリース後対応可か」の運用ヒントを併記する。総合判定（`HADO: READY/BLOCKED/ERROR`）は一覧の最後に出力する。TTY では ANSI カラーを付与し、`PASS` は緑、`FAIL` は赤/黄（required+critical の FAIL を最強調）で表示する（`NO_COLOR` が設定されている場合は無効）。

**Coverage 入力の必須条件:** Readiness Standard が `test.c0_coverage` または `test.c1_coverage` のいずれかを含む場合、`charge` で `--coverage-input` を渡すか、manifest の `evidence.coverage.inputs` が必要。`fire` 実行時にどちらも無いとエラー終了（exit 2）。

終了コード: `0` = ready、`1` = blocked（`required: true` かつ `severity: critical` の gate が失敗）、`2` = error（引数・読み込み・未対応 gate など）。

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
- `release.automation_declared`（`evidence.release.automation.workflow_refs` に **空白以外**の文字列が **1 件以上**。`systems` は任意のメタデータで現ゲートでは未使用）

### Gate severity（現状）

`severity` は `internal/standard/types.go` の独自型 `Severity` として扱い、次の 3 値へ制限する（未知値は `standard.Load` でエラー）:

- `critical`
- `major`
- `minor`

`required: true` の gate が失敗した場合の扱い:

- `critical`: **blocked**（原則リリース不可）
- `major`: ready のまま（リリースは可能だが、リリース後の早期対応を要求する運用を想定）
- `minor`: ready のまま（リリース可能。リリース後の任意タイミングでの対応を想定）

`severity` 未指定は `minor` と同等に扱う（非ブロック）。

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
