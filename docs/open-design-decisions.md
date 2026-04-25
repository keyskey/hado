# 未解決課題

Status: 初期 Open Design Decisions
Date: 2026-04-25

## Policy language

Readiness Standard の gate 条件をどの表現にするか。

決定:

```text
Policy language は YAML に統一する。
```

理由:

- Production Readiness の基準をできるだけ読みやすくするため
- 実装と運用をシンプルに保つため
- PR review で差分を理解しやすくするため
- 独自言語や複雑な式評価を MVP に持ち込まないため

## Module transport

Module 実行方式をどうするか。

候補:

- JSON-over-stdio
- gRPC
- WASM

初期方針:

```text
MVP は JSON-over-stdio。
実装言語を問わず module を作れることを優先する。
```

## HADO implementation language

HADO core / CLI / gate evaluator をどの言語で実装するか。

決定:

```text
MVP の HADO 本体は Go で実装する。
```

理由:

- 単一バイナリとして配布しやすく、GitHub Actions やローカル CI で扱いやすい
- CLI、設定ファイル読み込み、外部プロセス実行、JSON / YAML 処理との相性がよい
- CI tool として起動が速く、利用者側に大きな runtime 前提を要求しにくい
- `gobce` との親和性が高い一方で、HADO 本体と Go 専用 analyzer の責務は分離できる

補足:

```text
Module は JSON-over-stdio contract を優先し、Go に限定しない。
Datadog、Notion、Slack などの公式 module は、SDK や保守性に応じて Go / TypeScript などを選べる。
```

## Score model

Readiness score をどう計算するか。

候補:

- weighted category score
- required-gates-first model
- severity-based penalty model

初期方針:

```text
MVP は required-gates-first。
score は補助情報として扱い、blocking decision は required gate failure を優先する。
```

## Exception model

Exception / waiver をどう管理するか。

候補:

- local YAML exception
- GitHub approval integration
- external approval workflow
- audit database

初期方針:

```text
MVP では最小限の local exception から始める。
監査要件が強くなった段階で外部承認や audit database を検討する。
```

## Repository strategy

HADO 本体と analyzer module をどう分けるか。

決定済み:

```text
gobce は最初から別リポジトリとして作る。
```

未決:

- HADO 公式 modules を本体リポジトリに含めるか、個別リポジトリに分けるか
- standards を HADO 本体に同梱するか、別パッケージにするか
- GitHub Action を同一リポジトリに置くか、別リポジトリにするか

## Naming

基本方針:

```text
迷ったら、宇宙戦艦ヤマト・波動砲からの連想としてしっくりくる方を選ぶ。
```

現時点の命名:

```text
HADO Manifest
Readiness Standard
Module
Plugin
Gate
Evidence
Finding
Decision
Exception
Report
```

決定:

- 文中・user-facing の表記は原則 `HADO` に統一する。
- 内部的な拡張単位は `Module` と呼ぶ。
- 外部システムから HADO を呼び出しやすくするためのものは `Plugin` と呼ぶ。

例:

```text
Module
  HADO が実行して evidence / metrics / findings を得る内部拡張。
  例: gobce module, Datadog module, Notion module

Plugin
  外部システムに組み込む接続口。
  例: GitHub Action, Slack app, Notion integration
```

未決:

- `Readiness Standard` の最終名をどうするか

## HADO Manifest schema

未決:

- `evidence` と `inputs` のどちらが直感的か
- `modules` を manifest に明示するか、自動 discovery するか
- service metadata を Datadog service catalog から補完するか、manifest に必須で書くか

## Readiness Standard schema

`Readiness Standard` は少し長く、`Standard` や `Policy` だけだとソフトウェア運用の世界で一般語すぎて混乱しやすい。

候補:

```text
HADO Template
  readiness 基準のテンプレートであることが分かりやすい。

HADO Target
  目指すべき readiness target というニュアンスが出る。

Readiness Target
  独自性は弱いが、意味は直感的。
```

初期メモ:

```text
「波動砲を撃てる状態」から考えると、Target は世界観にも合う。
一方、テンプレートとして再利用する実体だと考えると HADO Template も自然。
```

未決:

- `Readiness Standard` を `HADO Template` または `HADO Target` に改名するか
- standard inheritance を許可するか
- standard composition を許可するか
- environment-specific override をどう表現するか
- regulated-system / exchange-grade のような高信頼 standard をどの粒度にするか

## Report format

未決:

- JSON schema の安定化タイミング
- Markdown report の標準レイアウト
- GitHub PR comment の更新方式
- SARIF 互換を持たせるか
- AI agent 向け report format を別に持つか

## gobce integration

決定済み:

```text
gobce は別リポジトリとして開発する。
HADO とは analyzer module として連携する。
```

未決:

- HADO module mode の CLI サブコマンド名
- `estimatedBranchCoverage` と `test.estimated_c1` の正確な metric naming
- coverprofile と source path の解決方法
- generated code の除外ルール
- monorepo 内 Go module の扱い

## Documentation language

未決:

- 英語化のタイミング
- 日本語版を残すか、英語版に完全移行するか
