# ローカル開発コマンド

CI で実行している `lint` / `format` / `test` をローカルでも手軽に実行できるように、`Makefile` を用意しています。

## 事前準備

- **Go** … 次のどちらか
  - **グローバルに入れずに使う:** 下記 **`make bootstrap-go`** で `.tools/go` にだけ公式 Go を置く（`direnv` 等で `export PATH="$PWD/.tools/go/bin:$PATH"` とするとシェルでも同じバイナリになる）
  - **既に PATH に `go` がある:** そのまま利用（`.tools/go` が無ければそちらを使う）
- Docker（YAML / Markdown lint 用）
- gobce（HADO 自身の readiness check 用）

### リポジトリ内だけに Go を置く（Homebrew 不要）

`go.mod` の `go 1.22` に合わせ、`Makefile` の **`GO_BOOTSTRAP_VERSION`**（既定 `1.22.12`）を [go.dev/dl](https://go.dev/dl/) から **`.tools/go`** に展開する（ディレクトリは `.gitignore` 済み）。

```bash
make bootstrap-go
```

初回のみ **ネットワーク必須**。以降の `make test` などは **`ensure-go`** 経由で、`.tools/go` が無い・壊れている場合だけ再取得する。

依存ツールのインストールは次の 1 コマンドで実行できます。

```bash
make setup
```

`make setup` は `ensure-go` のあと、HADO 自身の readiness check で使う `gobce` CLI を `go install` します。
YAML / Markdown の lint は Docker 上で実行します。Markdown は `npx` ではなく、バージョン固定の [davidanson/markdownlint-cli2](https://hub.docker.com/r/davidanson/markdownlint-cli2) イメージを使うため、ローカルに Node.js は不要です。

### pre-push で lint を強制（推奨）

CI の **Lint** ジョブと同じく **`make ci-lint`**（`fmt-check` + `lint`）を **`git push` の前**に走らせるには、リポジトリルートで次を一度実行する。

```bash
make setup-hooks
```

**`go test` は走らない**（Test ワークフローで `readiness-check` が **1 回だけ** `go test` をかけるため、Lint と二重にしない）。push 前にユニットテストもローカルで確かめたいときは **`make pre-pr`**（`ci-lint` + `test`）か **`make test`** を別途実行する。

Go は **PATH 上の `go`** か **`make` が参照する `.tools/go`**（`bootstrap-go` 済み）があればよい。Docker が起動している必要がある（YAML / Markdown lint 用）。

CI の [Test ワークフロー](https://github.com/keyskey/hado/blob/main/.github/workflows/test.yml) は `setup` のあと `readiness-check`（coverage・gobce・fire）を実行する。ローカルで合わせるときは `make readiness-check`（事前に `make setup`）。

無効にする場合: `git config --unset core.hooksPath`（またはリポジトリの config から該当行を削除）。

## 主要コマンド

```bash
make bootstrap-go  # .tools/go にだけ Go を入れる（任意・初回）
make lint       # YAML, Markdown, Go の lint
make fmt        # Go ファイルを gofmt で整形
make fmt-check  # Go の整形漏れチェック (CI 相当)
make test       # Go テスト
make readiness-check # coverage evidence 生成 + HADO 自身の readiness 評価
make ci-lint    # fmt-check + lint（CI Lint ジョブ・pre-push と同じ; go test なし）
make pre-pr     # ci-lint + test（ローカルで PR 前にまとめて叩く用）
make setup-hooks # pre-push を有効化
```

## 補足

`make pre-pr` は `ci-lint` のあと `test` を実行します。**CI 上で `go test` が動くのは Test ワークフローだけ**なので、`pre-push` だけではカバレッジ付きテストや gobce は確認されない。リリース前は `make readiness-check` も忘れずに。

日常的には `make setup-hooks` による pre-push（`ci-lint`）と、PR 作成前に手元で **`make pre-pr`** または **`make test`** を推奨します。

実装を変えたあとのドキュメント更新は、`.cursor/rules/hado-implementation-docs.mdc` と Skill `hado-doc-sync` を参照してください。

`make readiness-check` は `coverage.out` と `hado-coverage.json` を生成します。これらはローカル/CI の生成物として Git 管理から除外しています。
