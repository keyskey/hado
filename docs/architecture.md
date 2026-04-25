# Project HADO アーキテクチャ

Status: 初期アーキテクチャノート
Date: 2026-04-25

## 論理アーキテクチャ

```text
hado
├── hado-core
│   ├── policy engine
│   ├── gate evaluator
│   ├── score calculator
│   ├── exception evaluator
│   └── report generator
│
├── hado-cli
│   └── hado evaluate
│
├── hado-gate
│   └── deployment and release decision layer
│
├── hado-modules
│   ├── github
│   ├── datadog
│   ├── notion
│   ├── slack
│   ├── sbom
│   ├── openapi
│   ├── terraform
│   └── custom modules
│
├── hado-standards
│   ├── web-service
│   ├── batch-job
│   ├── critical-api
│   ├── data-pipeline
│   ├── regulated-system
│   └── exchange-grade
│
└── hado-integrations
    ├── GitHub Actions
    ├── Datadog
    ├── Notion
    ├── Slack reporter
    └── PR comment reporter
```

ここでの名前は論理コンポーネント名である。実装初期は HADO 本体を単一リポジトリで進め、インターフェースが安定してから必要に応じて分割してよい。

HADO core、CLI、gate evaluator の初期実装言語は Go とする。単一バイナリとして配布しやすく、GitHub Actions やローカル CI で起動しやすく、YAML / JSON 処理と外部 module 実行の実装をシンプルに保ちやすいためである。

ただし、Go の C1 カバレッジ計測ツールである `gobce` は最初から別リポジトリとして開発する。HADO からは外部 analyzer module として利用する。

## 初期ターゲット統合

HADO は OSS として広く使えることを前提にする。ただし MVP では実装範囲を絞るため、初期の integration target は次に集中する。

```text
GitHub
  repository, pull request, checks, comments, artifacts

Datadog
  observability だけでなく、service catalog / ownership / runbook など
  Internal Developer Portal 的な情報源としても使えるようにする

Notion
  PRR、運用手順、リリースノート、監査証跡、チームドキュメント

Slack
  readiness result、release notification、exception / follow-up 通知
```

Backstage への対応は当面見送る。将来的に需要が出た場合は module として追加できる設計に留める。

テストカバレッジ計測については、まず次の言語・ランタイムを対象にする。

```text
Go
Java / Kotlin
JavaScript / TypeScript
```

## 想定リポジトリ構成

HADO 本体リポジトリは、root に `go.mod` を置く単一 Go module として開始し、将来的に次の形へ育てる。

```text
.
├── go.mod
│
├── cmd
│   └── hado
│       └── Go CLI entrypoint
│
├── internal
│   ├── manifest
│   ├── core
│   ├── gate
│   ├── module
│   ├── standard
│   ├── report
│   └── scoring
│
├── modules
│   ├── github
│   ├── datadog
│   ├── notion
│   ├── slack
│   ├── sbom
│   └── openapi
│
├── standards
│   ├── web-service.yaml
│   ├── critical-api.yaml
│   └── regulated-system.yaml
│
├── examples
│   ├── go-service
│   └── github-actions
│
└── docs
```

`gobce` はこのリポジトリに同居させない。別リポジトリで library / CLI / HADO module interface を持たせ、HADO 側はその実行結果を module execution contract 経由で取り込む。

## 実行フロー

```text
developer / CI
    |
    v
hado evaluate --manifest hado.yaml --standard critical-api
    |
    v
HADO Manifest と Readiness Standard を読み込む
    |
    v
必要な module と evidence input を解決する
    |
    v
module を実行して evidence を収集する
    |
    v
evidence を facts / metrics / findings に正規化する
    |
    v
gates と exceptions を評価する
    |
    v
score と release decision を計算する
    |
    v
JSON / Markdown / GitHub PR comment / CI exit code を出力する
```

## ドメインモデル

```text
Service
  評価対象のシステム。name, owner, tier, language, runtime, catalog 参照、
  environment などを持つ。

HADO Manifest
  サービスリポジトリ側に置く宣言ファイル。サービス情報、入力ファイル、
  利用する module、出力形式などを定義する。

Readiness Standard
  再利用可能な readiness 基準。critical-api, web-service, exchange-grade
  など。必須 gate、閾値、severity、必要スコアを定義する。

Module
  evidence を収集・解析・報告する拡張単位。

Gate
  リリース基準。例: test.estimated_c1, observability.slo_exists,
  operations.runbook_exists, release.rollback_plan_exists。

Evidence
  module が収集した生データまたは正規化済みデータ。例: coverage metric,
  Datadog monitor 定義, Datadog service catalog metadata, Notion page, SBOM。

Finding
  module または gate が出す指摘。severity, message, location,
  evidence reference, recommendation を持つ。

Decision
  最終リリース判定。ready, warning, blocked, error。

Exception
  期限付き・監査可能な waiver。特定 gate の失敗を一時的に許容する。

Report
  CI、PR、監査、人間、AI エージェントが読むための出力。
```

## HADO Manifest

`hado.yaml` は、評価対象サービスが自分自身と evidence の場所を宣言するファイルである。単なる設定ではなく、リリース準備状態の入口になる manifest として扱う。

```yaml
version: v1

service:
  name: order-api
  owner: trading-platform
  tier: critical
  language: go
  catalog: datadog:service/order-api

standard:
  id: critical-api

evidence:
  coverage:
    go:
      coverprofile: coverage.out

  observability:
    slo_file: slo.yaml
    datadog_monitors_file: datadog-monitors.yaml
    datadog_service: order-api

  operations:
    runbook_url: https://www.notion.so/example/order-api-runbook
    rollback_plan: docs/rollback.md

modules:
  - id: hado.gobce
  - id: hado.datadog
  - id: hado.notion
  - id: hado.slack

report:
  formats:
    - json
    - markdown
  github_pr_comment: true
```

## Readiness Standard

Readiness Standard は「この種類のサービスは、本番に出る前に最低限これを満たすべき」という再利用可能な基準である。

```yaml
version: v1
id: critical-api
description: Critical API services must satisfy this minimum production readiness standard.

required_score: 85

gates:
  - id: test.statement_coverage
    severity: major
    required: true
    threshold:
      min: 80

  - id: test.estimated_c1
    severity: major
    required: true
    threshold:
      min: 70

  - id: observability.slo_exists
    severity: critical
    required: true

  - id: observability.monitor_exists
    severity: critical
    required: true

  - id: operations.runbook_exists
    severity: critical
    required: true

  - id: release.rollback_plan_exists
    severity: critical
    required: true
```

## Module アーキテクチャ

Module は HADO を language-agnostic / vendor-neutral に保つための拡張点である。

### Module 種別

```text
evidence collector
  外部ソースからデータを読み、正規化された evidence を出す。
  例: Datadog, Notion, GitHub, Kubernetes, Terraform

analyzer
  コードや成果物を解析し、metric と finding を出す。
  例: gobce, SBOM parser, OpenAPI checker

policy pack
  gate や standard fragment、規制業界向けルールを追加する。
  例: fintech, exchange-grade, regulated-system

reporter
  レポートを外部へ出す。
  例: GitHub PR comment, Slack, Markdown file
```

### MVP の実行契約

最初は、単純な外部プロセス実行モデルにする。

```text
HADO core が module executable を起動する
HADO が RunRequest JSON を stdin に渡す
Module が RunResult JSON を stdout に書く
Module がログを stderr に書く
HADO が timeout, exit code, schema validation を管理する
```

これにより、Go、TypeScript、Python、Rust など、どのランタイムでも module を書きやすくなる。将来的に性能や streaming が必要になったら、gRPC や WASM に移行できる。

### Module Manifest

```yaml
apiVersion: hado.dev/v1
kind: Module

metadata:
  id: hado.gobce
  name: gobce
  version: 0.1.0

runtime:
  type: exec
  command: gobce
  args:
    - hado-module

capabilities:
  emits:
    metrics:
      - test.statement_coverage
      - test.estimated_c1
    findings:
      - test.uncovered_branch

evidence:
  required:
    - evidence.coverage.go.coverprofile
  optional:
    - service.language
```

### Run Request

```json
{
  "apiVersion": "hado.dev/v1",
  "runId": "01JZ0000000000000000000000",
  "service": {
    "name": "order-api",
    "owner": "trading-platform",
    "tier": "critical",
    "language": "go"
  },
  "standard": {
    "id": "critical-api"
  },
  "evidence": {
    "coverage": {
      "go": {
        "coverprofile": "coverage.out"
      }
    }
  }
}
```

### Run Result

```json
{
  "apiVersion": "hado.dev/v1",
  "module": {
    "id": "hado.gobce",
    "version": "0.1.0"
  },
  "metrics": [
    {
      "id": "test.statement_coverage",
      "value": 82.1,
      "unit": "percent"
    },
    {
      "id": "test.estimated_c1",
      "value": 68.4,
      "unit": "percent"
    }
  ],
  "findings": [
    {
      "id": "test.uncovered_branch",
      "severity": "major",
      "message": "この分岐の false path がテストされていません。",
      "location": {
        "file": "internal/order/validator.go",
        "line": 42
      },
      "evidence": {
        "kind": "if_false_path"
      },
      "recommendations": [
        "validation が失敗するケースのテストを追加してください。"
      ]
    }
  ]
}
```

### Module 安全性

Module runner は次を強制する。

- 入出力 schema の厳格な検証
- timeout
- 最大出力サイズ
- secret 値そのものではなく secret reference を使うこと
- CI での deterministic な module version
- レポートに secret を出さないこと
- module failure と readiness failure を区別すること

## 開発者体験

MVP では、最小構成で GitHub Actions から使えることを重視する。

```yaml
name: HADO

on:
  pull_request:

jobs:
  hado:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod

      - run: go test ./... -coverprofile coverage.out

      - uses: hado-dev/hado-action@v1
        with:
          manifest: hado.yaml
          standard: critical-api
```

GitHub Action は次を行う。

- HADO CLI を install または実行する
- manifest に定義された modules を実行する
- Markdown PR comment を投稿する
- JSON / Markdown artifact を upload する
- HADO が `blocked` または `error` を返したら job を fail する
