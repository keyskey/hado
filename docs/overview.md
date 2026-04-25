# Project HADO 概要

Status: 初期コンセプトノート
Date: 2026-04-25

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
