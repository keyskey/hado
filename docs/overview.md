# Project HADO 概要

Status: 初期コンセプトノート
Date: 2026-04-25（CLI 体系の設計追記: 2026-05-03）

## 概要

Project HADO は、サービスを本番環境へリリースする前に「本番で生き残れる状態か」を判定するための、オープンソースの Production Readiness プラットフォームである。

HADO は CI、監視基盤、セキュリティスキャナ、コード品質ツールを置き換えるものではない。それらの上に位置する「信頼性の意思決定レイヤー」である。既存ツールから証拠を集め、リリース準備状態の基準に照らして評価し、リリース可否・理由・次に取るべきアクションを返す。

## 最終的に目指すゴール

HADO が目指すのは、SRE、Platform Engineering、Security、Compliance、Product Engineering が共通して使える、Production Readiness as Code の標準レイヤーである。

最終的には次の状態を目指す。

```text
開発チーム:
  Pull Request やリリース前に、自分たちで HADO を実行して readiness を確認できる。

SRE / Platform:
  サービス種別ごとの Readiness Standard を定義し、チームにセルフサービスのガードレールを提供できる。

Security / Compliance:
  セキュリティ、監査、証跡、変更管理の要件を release gate としてコード化できる。

経営 / 監査:
  重要サービスがどの基準を満たしてリリースされたのかを、後から説明できる。
```

## 中心となる問い

HADO が答える中心的な問いはこれ。

```text
このサービスは、本番に出しても運用・復旧・監査・信頼に耐えられるか？
```

実務的には、各ツールの問いはこう違う。

```text
CI:       テストは通ったか？
Sonar:    コードは綺麗か？
Datadog:  システムは観測できているか、動いているか？
HADO:     このサービスはリリース後に運用・復旧・監査・信頼に耐えられるか？
```

## コンセプト

HADO の世界観は意図的にこう置く。

```text
service      = 宇宙へ出航する船
production   = 未知の宇宙
release      = 出航
incident     = 敵、障害、想定外の宇宙環境
HADO         = 出航前に「撃てる状態」を証明する readiness amplifier
```

このメタファーで重要なのは、HADO が「船そのもの」ではなく、「出航しても戦える状態かを証明する装備・備え」であること。

つまり HADO は、リリースを止める門番そのものではなく、リリース可能状態を作るための増幅器である。

## CLI の 3 段階（target → charge → fire）

プロジェクト名の由来である『宇宙戦艦ヤマト』の **波動砲**は、本番投入のメタファーとして HADO の世界観の核になっている。CLI の 3 段階は、その **オペレーションの流れ**（照準 → 艦のリソースを集めてからの一撃 → 発射可否の判断）を **リリース準備の作業**に読み替えたものである。作品の話数や細かい演出に依存した説明はドキュメントでは使わず、次の **責務**だけで固定する。

リリース可否の判定を、**いきなり一発で走らせず**次の 3 段階に分ける。いずれも **HADO Manifest（例: `hado.yaml`）を正本**とし、別ディレクトリに「第二の設定ファイル」を増やして運用の複雑さを上げないことを原則とする。

```text
1. hado target … 照準（何を、どの基準で評価するかを manifest に書く）
   ターミナル上の **対話（プロンプト）** で、評価対象の service と適用する Readiness Standard を聞き、回答を **manifest に書き戻す**。
   初回セットアップや基準の切り替えで使い、**変更は Git の diff としてレビュー**できる形に残す。

2. hado charge … 充填（証跡を集め、manifest を埋める）
   **target で回収した値**と、manifest に **すでに書かれているサービスメタデータ**（リポジトリ URL、Datadog の APM service 名や service catalog 参照など）を入力にし、Readiness Standard に照らして **まだ埋まっていない evidence を自動で埋める**。
   外部 module の実行、adapter による正規化、CI 内の `go test` 成果物の取り込みなどはここに集約されうる。永続の主たる形は **更新後の manifest**（および manifest が指すパス上の artifact）である。

3. hado fire … 発射判定（gate を評価し、リリース可否だけを返す）
   充填済みの manifest（と参照 artifact）と Standard を照合し、gate を評価する。**デプロイはしない**（波動砲の「射撃」に相当する本番反映は既存の CD に任せる）。HADO は **release gate の意思決定**だけを返す（`hado fire` は「この変更を本番に出してよいか」の判定に相当する）。

コマンド名は **短く・意味がブレない**ことを優先する（毎日・毎 PR で打つ前提）。世界観は overview や README で語り、**サブコマンド名は業務で迷わない語**に留める。

**マニフェスト中心の原則:** `.hado/context.json` のような **評価の正本を二重化する中間ファイルを前提にしない**。ログや人間向けサマリが必要なら **任意の副産物**（標準出力、CI artifact、一時ファイル）とし、**評価の入力の正は常に manifest** とする。

**`standard` の位置づけ:** 単なるルールセットではなく、**そのサービスを本番に出すために満たすべき生存基準**として扱う（Readiness Standard の説明は [アーキテクチャ](architecture.md) および [Production Readiness の計測と評価](production-readiness-evaluation.md) と整合させる）。

**UX の方針:**

- **CI:** すでに manifest が揃っているなら `target` は省略し、`charge`（不足 evidence の自動補完）→ `fire` のみでもよい。初回はローカルで `target` を実行して manifest をコミットする運用を想定する。
- **フラグでの `target`:** 対話なしで `--manifest` / `--standard` / `--service` だけ更新するモードも用意できる（スクリプト・自動化向け）。
- `--standard` には **短い alias** を許してもよい（例: `exchange-critical` → `ec`）。組織の運用に合わせて定義する。

**実装との関係:** 現リポジトリでは **`hado target`** が manifest の `service` / `standard` を書き込める。判定は **`hado evaluate`** が一括で行う。将来は `evaluate` を **`fire` のエイリアス**にする、`evaluate` にフェーズ選択を付ける、など実装で決める。いずれにせよ **論理フェーズは target → charge → fire** とドキュメントで固定する。

詳細な責務分担・データの流れは [アーキテクチャ](architecture.md) を参照する。

## プロダクト原則

1. Automated Production Readiness Review

   PRR を自動化し、CI/CD、ローカル開発、Platform Portal から実行できるようにする。

2. Production Readiness as Code

   リリース準備状態の期待値を、Notion やスプレッドシートや暗黙知ではなく、レビュー可能なコードとして表現する。

3. Reliability as Code

   SLO、アラート、Runbook、所有者、ロールバック計画、監査要件、DR 期待値を、リリース判定に使えるコード管理された情報として扱う。

4. Module / Plugin-based and language-agnostic

   特定の言語、監視ベンダー、IDP、CI システムに依存しない。HADO 内部の拡張は Module、外部システムから HADO を呼び出す接続口は Plugin として整理する。

5. Self-service guardrails

   SRE が中央集権的に止めるための道具ではなく、各チームが自分たちでリリース可能状態を作れるためのガードレールにする。

6. Pass/fail よりも理由と次の行動

   ただ `FAILED` と出すのではなく、なぜ止まったのか、次に何を直すべきかを返す。

## 用語

初期案では `config` と `profile` という名前を使っていたが、HADO の概念としては少し直感的ではない。現時点では、次の用語に寄せる。

```text
HADO Manifest
  評価対象サービスが「自分は何者で、どの証拠をどこから読むか」を宣言するファイル。
  例: hado.yaml

Readiness Standard
  サービスが満たすべきリリース準備基準。
  例: critical-api, web-service, exchange-grade
```

つまり、サービス側に置く `hado.yaml` は単なる設定ファイルではなく「サービスの readiness manifest」である。一方、`critical-api` のような再利用可能な基準は「profile」ではなく「standard」と呼ぶ。

```text
旧: config  -> 新: manifest
旧: profile -> 新: standard
```

## Non-goals

- HADO は CI/CD オーケストレータではない。
- HADO は SonarQube クローンではない。
- HADO は汎用静的解析ツールではない。
- HADO は監視データベースではない。
- HADO はあらゆるリリースを中央承認制にするものではない。
- MVP 時点で全カテゴリを網羅しない。
- MVP 時点で AI 分析は必須にしない。
- MVP 時点で Web UI は必須にしない。

## ドキュメント言語

初期ドキュメントは日本語で管理する。公開範囲やコントリビューション方針が固まった段階で、英語ドキュメントへの移行を検討する。
