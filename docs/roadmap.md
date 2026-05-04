# 開発計画とロードマップ

Status: 初期ロードマップ
Date: 2026-04-25

## First Implementation Bias

最初は、プロダクト価値が一番早く証明できる細い道に集中する。

```text
Go service + GitHub Actions + Datadog + Notion + Slack + hado.yaml + critical-api standard + gobce
```

この workflow が、明確な理由と次の行動を伴う ready / blocked 判定を返せるなら、HADO は observability、operations、security、compliance、regulated-system readiness へ広げるだけの重力を持つ。

## Initial Readiness Expansion Order

初期の readiness 対応は、実際の新サービスリリースでよく起きる順序に合わせる。

```text
QAする -> 運用ルールを決める -> 監視を用意する
```

HADO の gate 展開もこの順序に寄せる。

1. Coverage Readiness
   - `test.c0_coverage`
   - `test.c1_coverage`
   - QA と test confidence を最初に確認する。

2. Operation Readiness
   - `operations.owner_exists`
   - `operations.runbook_exists`
   - 本番運用の責任者と障害対応の入口を確認する。

3. Observability Readiness
   - `observability.slo_exists`
   - `observability.monitor_exists`
   - `observability.dashboard_exists`
   - 本番投入後に異常を検知し、利用者影響を追える状態かを確認する。

## MVP Scope

最初の実用ターゲットはこれ。

```text
GitHub Actions 上で、Go service の最低限の Production Readiness を評価できる。
```

### このリポジトリの現状（実装済みの範囲）

**実装の詳細は [implementation-status.md](implementation-status.md)（手保守）。** コードを変えたら Cursor の `.cursor/rules/hado-implementation-docs.mdc` と Skill `hado-doc-sync` でドキュメントを同時更新する。

現時点の要約:

- HADO Manifest loader（`internal/manifest`）
- Readiness Standard loader（`internal/standard`）
- coverage adapter 層と C0 / C1 coverage gate（`internal/coverage`, `internal/gate`）
- operations / observability / infra / release（rollback と `release.automation_declared`）の existence gate（manifest 由来）
- `hado target`（manifest に `service` / `standard` と **standard に沿った evidence プレースホルダー**; `cmd/hado`）
- `hado evaluate` の text / JSON 出力と終了コード 0 / 1 / 2（`cmd/hado`）

**設計上の CLI 体系:** `hado target` に続く `hado charge` / `hado fire` は [概要](overview.md) と [アーキテクチャ](architecture.md) に記載。実装は当面 `evaluate` が判定を担う。

ロードマップ上は、module runner、閾値ベースの infra gate、Markdown / PR 連携などはこのあとである。

MVP に含めるもの:

- HADO Manifest loader
- Readiness Standard loader
- module runner
- coverage adapter layer
- C0 coverage gate
- C1 gate
- SLO existence gate
- monitor existence gate
- runbook existence gate
- owner existence gate
- rollback plan existence gate
- JSON report
- Markdown report
- CI exit code
- GitHub PR comment integration

MVP に含めないもの:

- Web UI
- AI analysis
- Slack bot
- long-lived audit database
- advanced Datadog API analysis
- Backstage integration
- every language runtime

## Phase 0: Concept and documentation

成果物:

- project overview
- architecture document
- Production Readiness evaluation model
- roadmap
- `gobce` concept document
- open design decisions
- README
- example `hado.yaml`
- example Readiness Standard
- example report

## Phase 1: C1 coverage producers

目的:

```text
C1 coverage を HADO 最初の差別化シグナルにする。
```

重要な前提:

```text
HADO core は C1 coverage の producer に依存しない。
gobce は最初の producer 候補だが、gobco などの実測 producer も同じ metric として扱う。
```

Deliverables:

- `go test -coverprofile` parse
- Go AST parse
- supported branch candidates の抽出
- branch coverage 推定
- standalone CLI
- JSON output
- HADO module result output

Success criteria:

- 実在する Go service に対して30秒以内に動く
- 人間と AI が理解しやすい uncovered branch finding を出す
- HADO gate 経由で CI を fail できる

## Phase 2: HADO Core

Deliverables:

- manifest loader
- standard loader
- policy / gate evaluator
- module interface
- scoring
- exception evaluation の最小実装
- JSON / Markdown report
- **CLI:** `hado target`（manifest の `service` / `standard` と **evidence プレースホルダー**）は実装済み。設計上の `hado charge` / `hado fire` は [概要](overview.md)）。移行期間は `hado evaluate` で判定を一括し、内部でフェーズに分解できるようにする

## Phase 3: First Gate Set

初期 gate:

```text
test.c0_coverage
test.c1_coverage
operations.owner_exists
operations.runbook_exists
observability.slo_exists
observability.monitor_exists
observability.dashboard_exists
release.rollback_plan_exists
```

## Phase 4: Initial Integrations

Deliverables:

- GitHub Action
- GitHub Checks / PR comment integration
- Datadog evidence collector
- Datadog service catalog / ownership / runbook extraction
- Notion evidence collector
- Slack reporter
- PR comment reporter
- artifact upload
- stable exit-code behavior

## Phase 5: Additional Language Coverage

Deliverables:

- Java / Kotlin coverage analyzer integration
- JavaScript / TypeScript coverage analyzer integration
- language-specific coverage metric normalization
- cross-language test readiness gates

## Phase 6: Regulated and high-reliability standards

Deliverables:

- regulated-system standard
- exchange-grade standard
- audit trail checks
- retention checks
- RTO/RPO checks
- separation-of-duties checks
