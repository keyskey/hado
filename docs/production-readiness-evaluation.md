# Production Readiness の計測と評価

Status: 初期評価モデルノート
Date: 2026-04-25

## 評価の目的

HADO の評価は、コード品質だけを見るものではない。リリース後にサービスを安全に運用し、障害時に復旧し、監査に説明し、利用者から信頼される状態かを定量的に判断する。

中心となる問いはこれ。

```text
このサービスは、本番に出しても運用・復旧・監査・信頼に耐えられるか？
```

## 評価カテゴリ

HADO は将来的に次のカテゴリを扱う。

```text
Test readiness
  unit test, integration test, statement coverage, estimated C1 branch coverage,
  flaky test rate, critical path test

Code and build readiness
  lint, static analysis, complexity, dependency freshness, SBOM, license,
  reproducible build

Observability readiness
  SLO, SLI, metrics, logs, traces, dashboards, alerts, customer-impact detection

Operability readiness
  runbook, on-call owner, escalation path, rollback plan, feature flag,
  migration plan

Reliability readiness
  RTO/RPO, DR design, redundancy, backup, recovery test, dependency mapping,
  rate limit, graceful degradation

Security readiness
  SAST, dependency vulnerability, secret scan, authn/authz, PII handling,
  encryption

Compliance and audit readiness
  change management, approval history, audit log, evidence retention,
  WORM/retention, separation of duties, release notes

Release readiness
  deployment plan, rollback plan, blast-radius control, progressive delivery,
  post-release validation
```

## Evidence

Evidence は、gate 評価に使う証拠である。

例:

```text
coverage.out
  Go test coverage profile

slo.yaml
  SLO / SLI 定義

datadog-monitors.yaml
Datadog monitor 定義

Datadog service catalog metadata
  owner, lifecycle, system, tier, runbook, oncall など

Notion pages
  runbook, PRR, release note, operation procedure, audit evidence など

SBOM
  依存関係、license、脆弱性評価の入力
```

Evidence は HADO Manifest で宣言し、module が読み取る。

## Gate

Gate は、リリース可否に直接関係する評価基準である。

例:

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

Gate は Readiness Standard 側で required / threshold / severity を定義する。

```yaml
gates:
  - id: test.estimated_c1
    severity: major
    required: true
    threshold:
      min: 70

  - id: operations.runbook_exists
    severity: critical
    required: true
```

## Decision

HADO は module output を Readiness Standard の gate と照合し、最終 decision を返す。

```text
ready
  required gate がすべて通過し、score >= required_score

warning
  required gate は通過したが、non-blocking finding や低確度リスクがある

blocked
  1つ以上の required gate が失敗した

error
  manifest、standard、module、runtime error により評価を完了できなかった
```

## Score

Score は、人間が readiness の全体感を把握するための補助情報である。ただし MVP では、score より required gate の pass/fail を優先する。

初期方針:

```text
required gate failure がある場合:
  decision = blocked

required gate failure がなく、score >= required_score:
  decision = ready

required gate failure はないが warning finding がある場合:
  decision = warning
```

Score model の詳細は未決定であり、[未解決課題](open-design-decisions.md) で管理する。

## Report

HADO の report は、人間と機械の両方が読める必要がある。

JSON report 例:

```json
{
  "status": "blocked",
  "score": 72,
  "requiredScore": 85,
  "failedGates": [
    {
      "id": "test.estimated_c1",
      "severity": "major",
      "message": "Estimated C1 coverage is 68.4%, below the required 70%."
    },
    {
      "id": "operations.runbook_exists",
      "severity": "critical",
      "message": "Critical API must define a runbook."
    }
  ],
  "recommendations": [
    "未カバーの validation branch に対するテストを追加する。",
    "hado.yaml、Datadog service catalog、または Notion に runbook URL を追加する。"
  ]
}
```

Markdown report 例:

```text
HADO: BLOCKED

Score: 72 / 85

Blocking:
- estimated C1 coverage is 68.4%, required 70%
- no runbook URL configured

Warnings:
- dashboard exists but no latency panel was detected

Recommended actions:
- add tests for uncovered validation branches
- link a Notion runbook for order-api
```

## Exception

Exception は、期限付き・監査可能な waiver である。

例:

```text
一時的に estimated C1 が基準未満だが、期限付きで release を許可する。
ただし、理由、承認者、期限、対象 gate、対象 service を残す。
```

Exception は readiness を満たしていることを意味しない。あくまで、リスクを認識した上で一時的に gate failure を許容する仕組みである。

## 初期 MVP Gate

最初は次の gate に絞る。

```text
test.statement_coverage
test.estimated_c1
observability.slo_exists
observability.monitor_exists
operations.runbook_exists
operations.owner_exists
release.rollback_plan_exists
```

この範囲だけでも、SRE が最低限安心してリリースできる状態を定義する価値がある。
