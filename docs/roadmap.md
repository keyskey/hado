# 開発計画とロードマップ

Status: 初期ロードマップ
Date: 2026-04-25

## First Implementation Bias

最初は、プロダクト価値が一番早く証明できる細い道に集中する。

```text
Go service + GitHub Actions + Datadog + Notion + Slack + hado.yaml + critical-api standard + gobce
```

この workflow が、明確な理由と次の行動を伴う ready / blocked 判定を返せるなら、HADO は observability、operations、security、compliance、regulated-system readiness へ広げるだけの重力を持つ。

## MVP Scope

最初の実用ターゲットはこれ。

```text
GitHub Actions 上で、Go service の最低限の Production Readiness を評価できる。
```

MVP に含めるもの:

- HADO Manifest loader
- Readiness Standard loader
- module runner
- `gobce` integration
- statement coverage gate
- estimated C1 gate
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

## Phase 1: gobce

目的:

```text
estimated C1 coverage を HADO 最初の差別化シグナルにする。
```

重要な前提:

```text
gobce は最初から別リポジトリとして開発する。
HADO 本体からは analyzer module として利用する。
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

## Phase 3: First Gate Set

初期 gate:

```text
test.statement_coverage
test.estimated_c1
observability.slo_exists
observability.monitor_exists
observability.dashboard_exists
operations.runbook_exists
operations.owner_exists
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
