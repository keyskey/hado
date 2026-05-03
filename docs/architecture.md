# Project HADO アーキテクチャ

Status: 初期アーキテクチャノート
Date: 2026-04-25（CLI 3 段階の設計追記: 2026-05-03）

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
│   ├── hado target   # manifest に service / standard を書く（対話またはフラグ）
│   ├── hado charge   # （設計）証跡収集・manifest の自動補完
│   ├── hado fire     # （設計）gate 判定のみ（デプロイしない）
│   └── hado evaluate # 現行: 判定を一括実行
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
│   ├── manifest      # HADO Manifest の読み込み
│   ├── standard      # Readiness Standard の読み込みと検証
│   ├── coverage      # coverage artifact の adapter と正規化 metric
│   ├── gate          # gate 評価と最終 decision
│   └── integration   # 外部システム連携の契約・基盤（予定: observability / Datadog など）
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

このツリーは、このリポジトリで **いま実在する** `internal` パッケージを反映している。将来、スコア計算やレポート生成、評価オーケストレーションを独立パッケージに切り出す場合は、その時点の責務に合わせて `internal` 配下に追加する。

`gobce` はこのリポジトリに同居させない。別リポジトリで library / CLI / HADO module interface を持たせ、HADO 側はその実行結果を module execution contract 経由で取り込む。

## CLI の 3 段階（target → charge → fire）

HADO の由来である **波動砲のオペレーション**（照準 → 艦のリソースを集めてからの一撃に向けた準備 → 発射の可否判断）を、リリース準備に対応づけた利用モデルである。**実装の正本は常に HADO Manifest** とし、作品固有の用語や話数に依存した説明は [概要](overview.md) のように **責務ベースの語**で書く。

評価を **「対象と基準の確定 → 証跡の収集・manifest への反映 → gate 判定」** の 3 段階に分ける。データの流れの詳細は [概要](overview.md) の該当節を参照する。

**一段ずつの論理フロー（マニフェストを正とする）:**

```text
developer（初回・メンテ）
    |
    v
hado target --manifest hado.yaml
    |
    v
  プロンプトで service / standard などを聞き、回答を hado.yaml（および必要なら隣接ファイル）に書き戻す
    |
    v
developer / CI
    |
    v
hado charge --manifest hado.yaml --standard ...
    |
    v
  manifest のメタデータ（repo URL、Datadog service 等）と target で入れた値を手がかりに、
  未充足の evidence を module / adapter / ローカルコマンドで埋める（主な出力は更新された manifest と参照パス）
    |
    v
hado fire --manifest hado.yaml --standard ...
    |
    v
  充填済み manifest と Standard を照合し gates / exceptions / score / decision を確定
    |
    v
JSON / Markdown / GitHub PR comment / CI exit code を出力（デプロイは行わない）
```

**現行実装（`hado evaluate`）との関係:** いまのコードパスは、上記の **manifest / standard の読み込み + charge に相当する coverage 解決 + fire に相当する gate 評価**を **1 コマンド**にまとめている。将来、フェーズ分割を実装するときの「分解点」は次のイメージで固定する。

- **target 相当:** 対話またはフラグによる **manifest のブートストラップ**（service × standard の確定をファイルに残す）。
- **charge 相当:** manifest に書かれた手がかりに基づく **evidence の自動補完**（module、外部 API、テスト成果物の取り込み、adapter による正規化）。CI で最も時間がかかりうる層。
- **fire 相当:** 更新済み manifest と Standard だけを見た **機械的な gate 判定**とレポート出力。中間 JSON を必須にはしない。

## 実行フロー（単一コマンド・現状）

移行期間およびローカル用途では、次の **一括実行**が引き続き有効である。

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
  リリース基準。例: test.c1_coverage, observability.slo_exists,
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

**実装済みのトップレベル（v1）:** `service`（`id` / `name`）と `standard`（`id`：Readiness Standard の論理 id または標準 YAML へのパス）は `hado target` が書き込める。`evidence` 以下の形はこのリポジトリの [実装状況](implementation-status.md) に従う。

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

インフラ readiness や多数の standard パターンを抱えるときでも **マニフェストのトップレベル構造を増やさない**方針と、本体を薄く保つための **Adapter / Analyzer / Module の分担**については、[Infrastructure Readiness とマニフェスト設計](infrastructure-readiness-and-manifest-design.md) にまとめる。上記 YAML は論理モデルの例であり、長期的には `evidence` を adapter 参照や evidence bundle などの安定した形へ寄せることを想定する。

## Readiness Standard

Readiness Standard は「この種類のサービスは、本番に出る前に最低限これを満たすべき」という再利用可能な基準である。

```yaml
version: v1
id: critical-api
description: Critical API services must satisfy this minimum production readiness standard.

required_score: 85

gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80

  - id: test.c1_coverage
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

Coverage は特に producer ごとの差が大きい。Go の coverprofile、keyskey/gobce
JSON、JaCoCo XML、lcov、Istanbul JSON などを HADO core が直接 gate として
解釈するのではなく、adapter が `test.c0_coverage` / `test.c1_coverage` の
正規化済み metric に変換する。HADO core は言語、ライブラリ、推定値か実測値か
を判定せず、Readiness Standard の threshold と正規化済み metric だけを照合する。

### Module 種別

```text
evidence collector
  外部ソースからデータを読み、正規化された evidence を出す。
  例: Datadog, Notion, GitHub, Kubernetes, Terraform

analyzer
  コードや成果物を解析し、metric と finding を出す。
  例: gobce, gobco, SBOM parser, OpenAPI checker

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

### Coverage adapter contract

Coverage adapter は、各言語・各ツールの coverage artifact を HADO の
正規化 coverage metric へ変換する境界である。HADO core の gate evaluator は
JaCoCo、lcov、Istanbul、gobce、gobco などの producer-specific schema を直接扱わない。

初期 CLI では次の形式で adapter と artifact を指定する（現状は `hado evaluate` 一括。将来の `hado charge` はここに相当する）。

```bash
hado evaluate \
  --standard standards/web-service.yaml \
  --coverage-input hado-json:coverage-metrics.json

hado evaluate \
  --standard standards/web-service.yaml \
  --coverage-input gobce-json:gobce.json
```

正規化後の metric は C0 / C1 の coverage percentage である。

```json
{
  "c0Coverage": 82.1,
  "c1Coverage": 68.4
}
```

今後の Java / TypeScript 対応は、`jacoco-xml`、`lcov`、`istanbul-json`
などの adapter を追加して同じ metric に変換する。標準側の gate ID は
言語や tool に依存せず、`test.c0_coverage` / `test.c1_coverage` のままにする。

Adapter contract の未決事項は [未解決課題](open-design-decisions.md) で管理する。

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
      - test.c0_coverage
      - test.c1_coverage
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
      "id": "test.c0_coverage",
      "value": 82.1,
      "unit": "percent"
    },
    {
      "id": "test.c1_coverage",
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

MVP では、最小構成で GitHub Actions から使えることを重視する。現状は **`hado evaluate` を Action から呼ぶ**形が主である。

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

フェーズ分割 CLI を採用したあとの **目標形**（[概要](overview.md) の 3 段階）は次のとおり。**CI では manifest が既にリポジトリに揃っていることが普通**なので、`target`（対話）は省略し **charge → fire** だけにすることが多い。初回だけローカルで `target` して `hado.yaml` をコミットする運用を想定する。`fire` は **デプロイをしない**（release gate の意思決定とレポートのみ）。GitHub Check / PR コメントへの投稿は `fire` 側（または reporter module）の責務とする。

```yaml
      - run: hado charge --manifest hado.yaml --standard critical-api
      - run: hado fire --manifest hado.yaml --standard critical-api
```

GitHub Action は次を行う。

- HADO CLI を install または実行する
- manifest に定義された modules を実行する
- Markdown PR comment を投稿する
- JSON / Markdown artifact を upload する
- HADO が `blocked` または `error` を返したら job を fail する
