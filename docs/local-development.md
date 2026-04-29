# ローカル開発コマンド

CI で実行している `lint` / `format` / `test` をローカルでも手軽に実行できるように、`Makefile` を用意しています。

## 事前準備

- Go
- Docker（YAML / Markdown lint 用）
- gobce（HADO 自身の readiness check 用）

依存ツールのインストールは次の 1 コマンドで実行できます。

```bash
make setup
```

`make setup` は Go ツールチェーンの存在確認のみを行います。
YAML / Markdown lint は都度 Docker コンテナで実行するため、ローカルに Python / Node.js の導入は不要です。

HADO 自身の readiness を評価する場合は、事前に `gobce` CLI をインストールします。

```bash
go install github.com/keyskey/gobce/cmd/gobce@latest
```

## 主要コマンド

```bash
make lint       # YAML, Markdown, Go の lint
make fmt        # Go ファイルを gofmt で整形
make fmt-check  # Go の整形漏れチェック (CI 相当)
make test       # Go テスト
make readiness-check # coverage evidence 生成 + HADO 自身の readiness 評価
make ci         # fmt-check + lint + test
```

## 補足

`make ci` は CI のチェック順序に合わせて、`fmt-check` -> `lint` -> `test` を順に実行します。
`make readiness-check` は `coverage.out` と `hado-coverage.json` を生成します。これらはローカル/CI の生成物として Git 管理から除外しています。
