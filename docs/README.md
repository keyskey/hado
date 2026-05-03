# Project HADO ドキュメント

Project HADO は、サービスを本番環境へリリースする前に「本番で生き残れる状態か」を判定するための、オープンソースの Production Readiness プラットフォームです。

初期ドキュメントは日本語で管理します。公開範囲やコントリビューション方針が固まった段階で、英語ドキュメントへの移行を検討します。

## ドキュメント一覧

1. [プロジェクト概要](overview.md)

   Project HADO が目指すゴール、コンセプト、プロダクト原則、用語、non-goals など、最もハイレベルな前提をまとめる。

2. [アーキテクチャ](architecture.md)

   HADO 本体の論理アーキテクチャ、技術選定、開発者体験、HADO Manifest、Readiness Standard、Module / Plugin アーキテクチャ、およびこのリポジトリの `internal` パッケージ構成をまとめる。

3. [Production Readiness の計測と評価](production-readiness-evaluation.md)

   Production Readiness をどのカテゴリ・gate・evidence・decision として評価するかをまとめる。

4. [開発計画とロードマップ](roadmap.md)

   MVP、開発ステップ、First Implementation Bias、フェーズごとの成果物をまとめる。

5. [実装状況](implementation-status.md)

   コードとの対応を手で保守する。実装変更時は Cursor Rule `hado-implementation-docs` と Skill `hado-doc-sync` に従い更新する。

6. [Go C1 カバレッジ計測ツール](gobce.md)

   最初から別リポジトリとして開発する `gobce` の目的、スコープ、HADO 連携方針をまとめる。

7. [Infrastructure Readiness とマニフェスト設計](infrastructure-readiness-and-manifest-design.md)

   本体を薄く保ちつつ、Readiness Standard を組織・サービスタイプ・実行基盤・Tier ごとに細かく定義できるようにし、HADO Manifest のスキーマを安定させるための設計指針（腐敗防止層、Standard の分割・合成、Manifest の形）をまとめる。

8. [未解決課題](open-design-decisions.md)

   Open Design Decisions、未解決の命名・設計・実装方針をまとめる。

9. [ローカル開発コマンド](local-development.md)

   `Makefile` で提供している `lint` / `format` / `test` 系コマンドと事前準備をまとめる。
