---
name: hado-doc-sync
description: >-
  HADO の実装（gates, CLI, coverage adapters, manifest）を変えたあと、
  docs/implementation-status.md と関連ドキュメントを実装に合わせて更新する。
  cmd/hado、internal/gate、internal/coverage、internal/manifest、internal/standard を触ったときに使う。
---

# HADO 実装とドキュメントの同期

## いつ使うか

- `cmd/hado`（サブコマンド、フラグ、終了コード、出力形式）
- `internal/gate`（評価する gate id、`Evaluate` の分岐）
- `internal/coverage`（adapter 名、`ParseAdapterInput` の分岐）
- `internal/manifest`（読み込む evidence の形）
- `internal/standard`（gate id 定数やバリデーション）

のいずれかを変更したあと、**同じタスク内**でドキュメントを直す。

## 手順

1. **真実のソースはコード**  
   実装を読み、実際にサポートしている gate id・adapter 名・CLI フラグを確定する。

2. **`docs/implementation-status.md` を更新**  
   - 実装済みゲートの箇条書き（`internal/gate/evaluate.go` の `switch` と一致）  
   - coverage adapter 一覧（`internal/coverage/parse.go` の `Format*` 定数と一致）  
   - `hado target` / `charge` / `fire` / `manifest doc` のフラグ・`--output` の取りうる値・終了コードの説明（`cmd/hado` と一致）  
   - MVP / 未実装として追いたい項目があれば表や箇条書きで維持（ロードマップの [docs/roadmap.md](docs/roadmap.md) と矛盾させない）

3. **`docs/hado.manifest.reference.yaml`**  
   - `internal/manifest/types.go` または `field_docs.go` を変えたら **`make gen-manifest-doc`** で再生成し、同じ変更セットに含める。

4. **`docs/roadmap.md`**  
   「このリポジトリの現状」の要約が古ければ 2〜3 文だけ直す。詳細は `implementation-status.md` に任せる。

5. **利用者向け**  
   ルート `README.md` の CLI 例やフラグ説明がずれていたら合わせる。

6. **push / PR 前**  
   `docs/` や `README.md` を触れたら `make lint`（少なくとも `make lint-markdown`）を通す。pre-push を使うなら `make setup-hooks`（[docs/local-development.md](docs/local-development.md)）。

## やらないこと

- `make docstatus` のような別コマンドは**ない**。Manifest の参考 YAML（コメント付き）は **`make gen-manifest-doc`**（`hado manifest doc`）で生成し、それ以外は手書き＋この Skill。
- 実装と無関係な長いロードマップの書き換えは、ユーザーが求めた範囲に留める。

## 参照

- プロジェクトルール: [.cursor/rules/hado-implementation-docs.mdc](.cursor/rules/hado-implementation-docs.mdc)
