# Skill の分類ガイド

このドキュメントは、skill を体系的に分類し、`docs/skill-library.md` を一貫して更新するためのガイドです。

skill は単一の性質だけでできているとは限りません。そのため、このガイドでは「上位分類」と「詳細類型」を分け、一覧へ載せるときは主用途に最も近い分類を選びます。

## 分類の全体像

| 上位分類 | 判断基準 | 含まれる主な類型 |
|---|---|---|
| A. 行動定義型 | モデルに「何をするか」を定義する。 | ワークフロー型、ロール型、出力フォーマット型 |
| B. 判断定義型 | モデルに「どう考え、どう評価するか」を定義する。 | 判断基準・評価型、思想・スタイル型、スコアリング・査定型 |
| C. 制約定義型 | モデルに「何を避けるべきか」を定義する。 | ガードレール・制約型 |
| D. 知識提供型 | モデルに「何を知っておくべきか」を提供する。 | 知識圧縮型 |

## 上位分類

### A. 行動定義型

モデルの行動、手順、役割、出力形式を定義する skill です。

「この状況では何をどの順番で行うか」「どの役割として振る舞うか」「どの形式で出力するか」を揃えたいときに使います。

### B. 判断定義型

モデルの評価基準、設計思想、優先順位、採点基準を定義する skill です。

「何を良いとみなすか」「どのリスクを重く見るか」「どう比較して判断するか」を揃えたいときに使います。

### C. 制約定義型

モデルが避けるべき行動、危険な操作、品質を下げる判断を定義する skill です。

「してはいけないこと」「確認なしに進めてはいけないこと」「安全のために守るべき境界」を明確にしたいときに使います。

### D. 知識提供型

特定分野の前提知識、ベストプラクティス、アンチパターン、専門用語を圧縮して提供する skill です。

「この領域では何を知っている前提で考えるべきか」を揃えたいときに使います。

## 詳細類型

### 1. ワークフロー型スキル

手順や実行フローを定義する skill です。

#### 目的

モデルに対して、再現可能な手順を順番に実行させる。

#### 特徴

- 手順が段階的
- 実行順序がある
- 中間確認がある
- 成果物が定義されている

#### skill の例

- bug-investigation
- incident-response
- migration-checklist
- debugging-flow
- release-process
- refactoring-flow

#### 典型構造

1. 状況収集
2. 前提確認
3. 実行
4. 検証
5. レポート出力

### 2. 判断基準・評価型スキル

品質評価や意思決定基準を定義する skill です。

#### 目的

レビューや判断の基準を統一する。

#### 特徴

- 原則
- ヒューリスティック
- トレードオフ判断
- 優先順位
- severity 判定

#### skill の例

- next-developer-review
- maintainability-review
- architecture-review
- api-review
- production-readiness-review

#### 典型構造

- Philosophy
- Review principles
- Review perspectives
- Severity guidelines
- Refactoring guidance

### 3. ガードレール・制約型スキル

禁止事項や安全制約を定義する skill です。

#### 目的

危険、低品質、非推奨な行動を防ぐ。

#### 特徴

- 明示的禁止
- 安全制約
- 非推奨事項
- リスク軽減

#### skill の例

- security-guardrails
- migration-safety
- production-safe-sql
- ai-coding-rules

#### 制約の例

- 不必要な全ファイル書き換えを禁止する
- 破壊的変更を黙って導入しない
- 過剰抽象化を避ける
- 隠れた副作用を避ける

### 4. 思想・スタイル型スキル

チームや設計思想を定義する skill です。

#### 目的

出力を特定の設計思想や文化に揃える。

#### 特徴

- アーキテクチャ方針
- コーディング哲学
- 設計価値観
- トレードオフ方針

#### skill の例

- engineering-principles
- clean-code-philosophy
- architecture-principles
- pragmatic-engineering

#### 方針の例

- 賢さより明快さを優先する
- 暗黙的な処理より明示的な処理を優先する
- 将来安全に変更できることを優先する
- 副作用を見えるようにする

### 5. ロール型スキル

特定の役割や視点を与える skill です。

#### 目的

モデルの振る舞いや思考スタイルを変える。

#### 特徴

- 視点ベース
- 専門的立場
- 優先順位が異なる
- 擬似人格的

#### skill の例

- planner
- verifier
- skeptic-reviewer
- architect
- performance-specialist
- tech-lead

#### 典型的な振る舞い

- 前提を疑う
- 代替案を検討する
- 正しさを検証する
- スケーラビリティを重視する
- リスク低減を重視する

### 6. 出力フォーマット型スキル

レスポンス構造を標準化する skill です。

#### 目的

出力形式を統一し、再利用しやすくする。

#### 特徴

- テンプレート
- フォーマット規則
- 構造化出力
- 一貫したセクション

#### skill の例

- adr-writer
- incident-report-writer
- pr-review-format
- changelog-generator

#### 出力セクションの例

- Summary
- Findings
- Risks
- Recommendations
- Follow-ups

### 7. 知識圧縮型スキル

特定分野の知識を圧縮して持たせる skill です。

#### 目的

専門知識を再利用可能な形で提供する。

#### 特徴

- ドメイン知識
- ベストプラクティス
- アンチパターン
- 専門用語

#### skill の例

- postgres-performance
- rails-conventions
- aws-lambda-patterns
- ddd-tactical-patterns
- graphql-schema-design

#### 典型内容

- 推奨パターン
- アンチパターン
- 運用観点
- パフォーマンス観点

### 8. スコアリング・査定型スキル

評価基準に基づき採点する skill です。

#### 目的

品質評価や優先順位付けを行う。

#### 特徴

- スコアリング
- 重み付け
- severity 判定
- readiness 評価

#### skill の例

- production-readiness-review
- architecture-scorecard
- seniority-review
- scalability-assessment

#### 典型出力

- スコア
- 強み
- 弱み
- 致命的リスク
- 推奨改善案

## 複合型スキルの扱い

実務で使う skill は、複数の類型を組み合わせていることが多いです。

たとえば `next-developer-review` は、判断基準、設計思想、ガードレール、severity 判定、出力フォーマットを含む複合型です。ただし、主な用途は「レビュー判断を行うこと」なので、`docs/skill-library.md` では判断定義型の判断基準・評価型に置きます。

## skill-library.md への分類ルール

`docs/skill-library.md` を更新するときは、次のルールで分類します。

1. 1つの skill は、主分類に一度だけ掲載する。
2. 主分類は、ユーザーがその skill を呼び出す主目的で決める。
3. 複合要素は、重複掲載ではなく概要や使う場面に含める。
4. 行動を始めるための skill は、原則として行動定義型に置く。
5. 評価やレビューの基準を揃える skill は、原則として判断定義型に置く。
6. 禁止事項や安全境界が主目的の skill は、原則として制約定義型に置く。
7. 特定分野の知識提供が主目的の skill は、原則として知識提供型に置く。
