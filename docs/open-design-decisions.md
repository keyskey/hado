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

## Coverage adapter and evidence contract

決定済み:

```text
C0 / C1 coverage は producer-neutral な metric として扱う。
HADO core は gobce / gobco / JaCoCo / lcov / Istanbul などの計測器や、
推定値 / 実測値の違いを判定しない。
producer-specific artifact は adapter が HADO の normalized coverage metrics へ変換する。

HADO Manifest は coverage evidence を `evidence.coverage.inputs` で宣言する。
各 input は `adapter` と `path` を持ち、相対 path は manifest file の directory
から解決する。

`hado evaluate` は `--manifest hado.yaml` から coverage input を読める。
`--coverage-input` が指定された場合は direct CLI override として扱い、
manifest の coverage input より優先する。
```

未決:

- normalized coverage metrics schema に持たせる metadata
  - language
  - producer / tool name
  - producer version
  - metric source が estimated か measured か
  - 対象 package / module / path scope
  - confidence や adapter warning を表現するか
- 複数 coverage evidence が同じ metric を出した場合の merge ルール
  - 後勝ちにするか
  - 明示 priority を持たせるか
  - conflict として evaluation error にするか
  - report に overwritten evidence を残すか
- coverage metric の意味を adapter がどこまで保証するか
  - Java / Kotlin の JaCoCo branch coverage を HADO の C1 とみなす条件
  - JavaScript / TypeScript の lcov / Istanbul branch coverage を HADO の C1 とみなす条件
  - Go coverprofile 由来 C0 と他言語 C0 の比較可能性
- producer-specific adapter の schema versioning
  - `keyskey/gobce` は pre-1.0 で JSON output が変わる可能性がある
  - adapter が producer version を検出するか
  - adapter version と producer version の compatibility をどう表現するか
- HADO Manifest の coverage input 宣言を将来拡張する schema
  - CI matrix / monorepo で複数 coverage artifact を扱う表現
  - language autodetection を許可するか
  - input ごとの include / exclude / scope をどう表現するか
- adapter warning / recommendation を report に含める形式
  - generated code の除外漏れ
  - unsupported coverage feature
  - partial parse
  - source path resolution failure
- report schema の coverage evidence 表現
  - gate result には normalized metric の値を出す
  - 詳細 report には adapter / source artifact / producer metadata を出すか
  - 監査用途で元 artifact への参照を保持するか

Go / gobce 固有で残る未決:

- producer が HADO module として動く場合の CLI サブコマンド名
- coverprofile と source path の解決方法
- generated code の除外ルール
- monorepo 内 Go module の扱い

## Documentation language

未決:

- 英語化のタイミング
- 日本語版を残すか、英語版に完全移行するか
