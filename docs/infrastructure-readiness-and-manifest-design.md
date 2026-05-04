# Infrastructure Readiness とマニフェスト設計

Status: 設計ドラフト（実装の前提を固定するための指針）
Date: 2026-04-29（CLI フェーズとの対応追記: 2026-05-03）

## 目的

Infrastructure Readiness は、実行基盤（Kubernetes / ECS / Cloud Run / App Engine など）、IaC の作法、組織ポリシー、サービス形態（Web API / バッチ等）、Tier によって **求める gate と evidence の取り方が分岐しやすい** 領域である。一方で HADO では次を同時に満たしたい。

1. **本体のロジックは薄く汎用的**に保ち、ベンダーや言語、IaC 形式の詳細を本体に持ち込まない。
2. **HADO Manifest のスキーマは、拡張モジュールやツールの増加に伴って構造が増殖しない**（安定した形のままシンプルに保つ）。
3. **Readiness Standard は、組織・サービスタイプ・実行基盤・Tier などに応じて細かく多数のパターンを定義できる**。
4. ツールや API の **入出力フォーマットの差・陳腐化しやすい知識**は、**Adapter / Analyzer / Module** などの境界に閉じ込める（腐敗防止層）。

この文書は、そのトレードオフを解くための **設計指針とデータの流れ** を固定する。細部のフィールド名や YAML の完全な schema は、実装フェーズで [未解決課題](open-design-decisions.md) と併せて詰める。

---

## 設計指針の要約

| 関心事 | 置き場 | 増え方 |
| --- | --- | --- |
| 「何を満たせばよいか」（gate、閾値、必須度） | Readiness Standard（YAML、複数ファイル可） | 組織・基盤・タイプ・Tier ごとにファイルが増える |
| 「このサービスは誰で、証拠はどこから取るか」 | HADO Manifest（**安定した少数のブロック**） | モジュール数が増えても **トップレベルの形は変えない** |
| プロデューサ固有のファイル形式 → 正規化 metric | Adapter（in-process が既定） | adapter 種別の列挙は増えるが、manifest は `adapter` + `path`（または同等の参照）の繰り返しに留める |
| リポジトリ内の成果物の解析（IaC / マニフェストの深い意味） | Analyzer（プロセス内ライブラリまたは別プロセス module） | 実装とリリース単位は本体から分離 |
| 外部 API / SaaS / クラウド制御面 | Module（別プロセス想定） | SDK や認証は module 側 |

**本体が知る語彙**は次に限定する。

- Readiness Standard が参照する **gate ID** と、その gate に紐づく **threshold / required** などの宣言。
- 評価時に集約された **正規化済みの評価コンテキスト**（例: gate ID → 数値または真偽、欠損は「未供給」）。
- 上記を照合する **汎用ゲート評価**（「閾値以上か」「必須文字列が非空か」など、ドメイン固有の意味は持たない）。

ドメイン固有の意味（「この YAML が PDB を満たすか」等）は **Standard の gate の説明文と module / analyzer の実装**に置き、本体は「`infra.pdb_min_available` の actual が standard の min 以上か」といった **機械的な比較**に徹する。

---

## なぜマニフェストを増やさないのか

インフラの差分が manifest に入り込むと、典型的に次が起きる。

- `kubernetes`, `ecs`, `cloud_run` ごとに **トップレベルのキーが増える**。
- 新しい module ごとに **専用ブロックが足される**。
- バージョンアップのたびに **サービスリポジトリ側の hado.yaml を書き換える**頻度が上がる。

これは「マニフェストがサービスの入口である」という利点を損ないやすい。そこで **カスタマイズの重心を Standard 側に寄せ**、Manifest には **「どの standard を使うか」「証拠の参照（パス・adapter・module 起動に必要な最小の束）」**だけを載せる。

インフラの詳細（何をチェックするかの列挙）は **Standard の束（複数 YAML の合成結果）**として表現し、Manifest は **その束を指す 1 本の論理参照**（ファイルパス、または将来の catalog ID）で足りるようにする。

---

## Readiness Standard の多パターン化

組織・サービスタイプ・実行基盤・Tier ごとに細かいパターンを持てるようにするため、Standard は次のような **分割・合成**を許容する設計とする（実装は「合成の解決」と「最終 gate 一覧の生成」に落ちる）。

### 推奨する分割の考え方

```text
組織ベースライン（全サービス共通の最低ライン）
  + サービスアーキタイプ（web-api / batch / data-pipeline など）
  + 実行基盤パック（k8s / ecs / cloud-run など）
  + Tier 上乗せ（critical では追加 gate、low では緩和）
  → 評価時に適用する「解決済み standard」
```

合成の具体手法（`import` / `extends` / ディレクトリ規約で自動レイヤリングなど）は [未解決課題](open-design-decisions.md) で詰める。重要なのは **「複雑さは standard リポジトリ側の YAML の組み合わせに閉じ、Manifest の形は変えない」**ことである。

### Gate ID の命名

インフラ領域の gate は、**安定した ID 空間**（例: `infra.*`, `platform.*`）をドキュメント化し、組織独自 gateは **`com.example.*` のような名前空間付き ID**で表現できるようにする（本体は ID を不透明なキーとして扱い、意味は Standard と module の契約に委ねる）。

---

## 安定した HADO Manifest の形（方針）

Manifest は **サービスのアイデンティティ**と **証拠への参照の集合**に徹し、次のようなブロックの組み合わせに **長期的に固定**する（フィールド名は実装と同期させるが、**トップレベルの概念ブロックは増やさない**ことを原則とする）。

```text
version
service        # 名前、owner、tier、language、論理カタログ参照 など（既存方針を踏襲）
standard       # 適用する readiness standard の参照（1 本または合成の入口）
evidence         # 「どの adapter / どのパス」「どの module に何を渡すか」の宣言
modules          # （任意）module runner 導入後も、manifest では「起動する module の一覧と version 固定」程度に留める
```

### `evidence` の安定化の鍵: 参照と束ね方

プラットフォーム固有の構造は `evidence` の **直下にキーを増やして表現しない**。代わりに、次のいずれか（または併用）で表現する。

1. **Adapter 参照のリスト**（coverage で既に近い形）  
   `adapter` + `path`（＋必要なら `id` や `labels` など最小メタデータ）の繰り返し。新しいツールが増えても **配列要素の種類が増えるだけ**で、スキーマの木構造は変わらない。

2. **Evidence bundle（成果物の束）**  
   CI が生成する `evidence-bundle.json` のような **1 ファイル**に、複数の成果物パスやハッシュを載せ、Manifest からはその **束へのパス 1 本**だけを指す。詳細なキーは bundle 側の schema で進化させ、HADO Manifest は不変に近づける。

3. **Module への入力**  
   Module に渡すのは **RunRequest 内の evidence 参照**であり、Manifest 側は「どの module をどの version で起動するか」と「リポジトリルートから見た入力の束」程度に留める。クラウド固有のパラメータは **module 設定ファイル**（別ファイル、または bundle 内）に逃がし、Manifest はそのファイルへのパスを 1 本持つだけにできる。

---

## 腐敗防止層の役割分担

```text
[外部世界: API / SaaS / クラウド / 生ファイル]
        |
        v
   Module または Analyzer（別プロセス可）
   ・認証・SDK・API のバージョン差を吸収
   ・プロデューサ形式を内部表現に落とす
        |
        v
   Adapter（必要ならここでもう一段）
   ・単一ファイル形式（terraform show JSON, helm template, k8s yaml など）を
     「gate ID に紐づけ可能な metric / boolean」へ
        |
        v
[正規化 EvaluationContext]
        |
        v
   HADO 本体（gate evaluator）
   ・threshold / required / severity の照合のみ
```

- **Module**: 外部システムや重い依存を隔離。言語非依存の拡張点。
- **Analyzer**: リポジトリ内の IaC / マニフェストを読み、**組織独自のルール**を含めたいときの主戦場。別プロセスの module として配布してもよい。
- **Adapter**: 「ファイル形式 → 既知の metric ID」への写像。本体に近いが、**新しい adapter の追加で manifest のトップレベルは変えない**。

---

## Infrastructure Readiness における評価の流れ（論理）

1. **Manifest** から `standard` 参照と `evidence` / `modules` の参照を読む。
2. **Standard 解決**で、組織・タイプ・基盤・Tier を反映した **gate 一覧**を得る。
3. **Orchestration**（本体の薄い層）が、必要な adapter / module を起動し順序づけ、**EvaluationContext**（gate ID → 値）を埋める。
4. **Gate evaluator** が Standard と EvaluationContext だけを見て pass/fail を決める。

CLI の **target / charge / fire** に写すと、1〜2 は主に **target**（manifest に残す照準・standard 解決・**evidence スキャフォールド**）、3 は **charge**（manifest を読み書きしながら evidence 参照を埋める）、4 は **fire** に相当する。**正本は manifest** とし、`.hado/context.json` のような中間ファイルを必須にしない方針は [概要](overview.md) を参照する。詳細は [アーキテクチャ](architecture.md) の「CLI の 3 段階」を参照する。

ここで 3 の orchestration は「どの gate にどの module が責任を持つか」の **宣言**（Standard 側の annotation、または別の small registry）で駆動できるようにし、本体に **k8s / ECS の分岐**を書かない。

---

## 明示的なトレードオフ

- **Manifest を極限まで薄くすると**、初見の人は「このサービスで何がチェックされるか」を Manifest だけ読んでは分かりにくい。対策として、評価結果や dry-run で **解決済み standard の gate 一覧を出力**できるようにする価値が高い。
- **Standard の合成**が強力になるほど、**同名 gate の上書き規則**や **順序**の仕様が必要になる。これは standard フォーマット側の設計課題として切り出す。

---

## 実装・未決事項との関係

- 現行の Go 実装は coverage / operations にフォーカスした最小 evaluator である。本文書の **EvaluationContext と薄い evaluator** は、インフラ gate を追加するときの **拡張方向**として整合させる。
- Standard の合成構文、module catalog、RunRequest に載せる evidence の正確な形は [未解決課題](open-design-decisions.md) に残し、本書は **責務の割り当てとマニフェスト安定性の原則**を優先して記述する。

---

## 関連ドキュメント

- [アーキテクチャ](architecture.md)（論理コンポーネント、module 契約）
- [Production Readiness の計測と評価](production-readiness-evaluation.md)
- [未解決課題](open-design-decisions.md)
