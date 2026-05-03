# hado

HADO は、サービスを本番環境へリリースする前に「本番で生き残れる状態か」を判定するための、オープンソースの Production Readiness プラットフォームです。

## 名前の由来

HADO という名前は「波動砲」から来ています。波動砲は、日本のSFアニメ『宇宙戦艦ヤマト』に登場する、宇宙戦艦ヤマトを象徴する主砲です。

ソフトウェアを本番環境へ送り出すことは、未知の宇宙へ船を出航させることに似ています。どんな障害、サイバー攻撃、運用上の想定外、監査上の問題に遭遇するかは、実際に本番へ出るまで完全には分かりません。

このプロジェクトは、宇宙戦艦ヤマトそのものではなく、ヤマトが未知の宇宙へ出ていく前に「戦える状態か」「波動砲を撃てる状態か」を確かめるための装備、あるいは備えとして位置づけています。

つまり HADO が問うのは、単にテストが通ったか、コードが綺麗かではありません。

```text
このサービスは、本番という未知の宇宙に出しても生き残れるか？
```

その問いを、Production Readiness as Code として扱えるようにすることが HADO の出発点です。

## Docs

- [Project HADO ドキュメント](docs/README.md)
- [実装状況（手保守; Cursor Skill `hado-doc-sync`）](docs/implementation-status.md)
- [ローカル開発コマンド](docs/local-development.md)

## Build and run

ローカルで `hado` CLI をビルドして実行する最小手順です。

```bash
make build
./bin/hado version
./bin/hado
```

## Target manifest（service / standard）

`hado target` は、HADO Manifest に **評価対象の service** と **適用する Readiness Standard** を書き込む。対話モード（TTY）ではプロンプトで聞き、非対話（CI）ではフラグで更新する。既存の `evidence` ブロックは上書きしない。

```bash
./bin/hado target --manifest hado.yaml \
  --service-name order-api \
  --standard-id web-service
```

TTY で実行すると、現在の manifest の値をデフォルトにしながら対話入力できる。

## Evaluate readiness

`hado evaluate` は、Manifest や CLI option で渡された evidence を
Readiness Standard の gate と照合し、required gate を満たしていれば
`READY`、満たしていなければ `BLOCKED` を返します。終了コードは
`0`（ready）、`1`（blocked）、`2`（エラー・未対応の required gate など）です。
`BLOCKED` のときは CI で扱いやすいように 1 で終了します。

**設計上の CLI 体系:** `hado target` で manifest に service / standard を記録したうえで、将来的には `hado charge`（未充足 evidence の自動補完）→ `hado fire`（判定のみ・デプロイはしない）へ分ける想定です。いまは `hado evaluate` が判定を一括で行います。詳細は [docs/overview.md](docs/overview.md) と [docs/architecture.md](docs/architecture.md) を参照してください。

HADO core は、特定の runtime、tool、SaaS、infrastructure provider の
フォーマットに直接依存しません。Coverage、Operation、Observability、
Infrastructure、Application、Security などの readiness domain は、
adapter や module が evidence を正規化し、standard の gate が判定します。

現在の evaluator は、coverage・operations・observability・release（rollback と
自動リリース用 `workflow_refs` の宣言）・infra（deployment 参照）の各 evidence を
Manifest から読み、対応する existence 系 gate を評価できます（詳細は
[docs/implementation-status.md](docs/implementation-status.md)）。
Coverage tool 固有の出力は adapter が `c0Coverage` / `c1Coverage` に正規化します。

```bash
printf '{"c0Coverage": 82.1, "c1Coverage": 72.5}\n' > coverage-metrics.json

cat > hado.yaml <<'YAML'
version: v1
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: coverage-metrics.json
  operations:
    owner: platform-team
    runbook: https://example.com/runbooks/order-api
YAML

./bin/hado evaluate \
  --standard standards/web-service.yaml \
  --manifest hado.yaml
```

Go coverprofile や `keyskey/gobce` の JSON output も adapter 経由で扱えます。

```bash
go test ./... -coverprofile=coverage.out
gobce analyze --coverprofile coverage.out --format json --output gobce.json

./bin/hado evaluate \
  --standard standards/web-service.yaml \
  --manifest hado.yaml \
  --coverage-input go-coverprofile:coverage.out \
  --coverage-input gobce-json:gobce.json
```

`--coverage-input` は従来どおり直接指定にも使え、指定された場合は
manifest の `evidence.coverage.inputs` より優先されます。
