---
name: modification-design
description: 実装前に、既存コードや仕様を根拠として、変更後のシステム構造を Modification Design に整理するときに使う。ユーザーが「変更後の設計」「実装前の設計」「何をどう変えるか」「責務、interface、contract、依存関係、影響、リスクを整理して」などを求めた場合は必ず使う。Plan のような実装手順やタスクリストではなく、変更対象となるモジュール、コンポーネント、API、型、データの変更後の責務と境界をレビュー可能な自然言語で定義する。コード、擬似コード、diff は出力しない。
---

# Modification Design

実装手順ではなく、変更後の責務、公開 contract、依存方向、影響をレビューできる自然言語の設計として定義する。実装計画は `role-planner`、React の詳細な component tree は `component-design`、実装済み差分の欠陥検出は `code-review` を優先する。

指定されたコード、仕様、issue、設計文書、公開 interface を先に確認し、現状を推測で断定しない。設計を左右する不明点は確認する。全体を止めない不明点は仮定または未解決事項、その影響先と risk を明示する。

## 原則

- 既存システムの説明へ広げず変更点を扱い、`lego-programming` に沿って一つの変更理由を持つ部品に分ける。
- 各対象の入力、出力、所有状態、副作用、依存先を示し、公開 interface と data contract を実装詳細から分離する。
- 追加・削除・変更・変更なしを区別する。互換性や migration があるなら、移行中の旧新 contract を示す。
- 実装手段を不必要に固定しない。

## 出力形式

次の章をこの順で出力し、該当なしも「変更なし」または「該当なし」と書く。

1. **変更概要**: 何を、なぜ、どの構造・振る舞いへ変えるか。
2. **変更対象**: ファイル一覧でなく、既存名を使った責務単位のレゴ。新規は新規と示す。
3. **Modification Design**: 対象ごとの現在、変更後、変更があれば状態・副作用、依存先・呼出元・委譲・依存方向。責務移動は移動元の削除と移動先の追加を対応させる。
4. **Interface Changes**: API、CLI、公開関数・class・props、export ごとの追加・削除・変更・変更なし、互換性。例は field 名、型、必須性、意味を自然言語または表で示す。
5. **Data / Contract Changes**: payload、DTO/Entity、DB・永続化、event/message、設定・形式の変更有無。変更 field は名称、型、必須性、意味、default・制約・互換性、producer/consumer を必要な範囲で示す。
6. **Dependency Changes**: 前後の依存方向、新設・削除・委譲、外部 I/O 境界。循環依存と下位実装への依存を避ける。
7. **Impact**: API 利用者、既存データ、tests、各クライアント、外部サービス・運用、必要なら observability/security/performance。影響なしには短い根拠を添える。
8. **Risks**: 今回の対象に結びつく互換性、migration/rollback、contract 不整合、性能・可用性、権限境界、二重管理、未解決事項の影響。

コード、擬似コード、diff、タスク順、工数、担当、DoD は書かない。未確認の名称を事実にせず、contract の変更有無は必ず明記する。完了時には、責務・依存、interface/data contract と互換性、影響・risk が具体的な変更対象に結びつき、実装開始前に判断可能かを確認する。
