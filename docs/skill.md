
# Skill の類型

## 1. ワークフロー型スキル

## 2. 判断基準・評価型スキル

## 3. ガードレール・制約型スキル

## 4. 思想・スタイル型スキル

## 5. ロール型スキル

## 6. 出力フォーマット型スキル

## 7. 知識圧縮型スキル

## 8. スコアリング・査定型スキル

---

# 上位分類

## A. 行動定義型

## B. 判断定義型

## C. 制約定義型

---


## 1. ワークフロー型スキル

手順や実行フローを定義するスキル。

### 目的

モデルに対して、再現可能な手順を順番に実行させる。

### 特徴

- 手順が段階的
- 実行順序がある
- 中間確認がある
- 成果物が定義されている

### 例

- bug-investigation
- incident-response
- migration-checklist
- debugging-flow
- release-process
- refactoring-flow

### 典型構造

1. 状況収集
2. 前提確認
3. 実行
4. 検証
5. レポート出力

---

## 2. 判断基準・評価型スキル

品質評価や意思決定基準を定義するスキル。

### 目的

レビューや判断の基準を統一する。

### 特徴

- 原則
- ヒューリスティック
- トレードオフ判断
- 優先順位
- severity 判定

### 例

- next-developer-review
- maintainability-review
- architecture-review
- api-review
- production-readiness-review

### 典型構造

- Philosophy
- Review principles
- Review perspectives
- Severity guidelines
- Refactoring guidance

---

## 3. ガードレール・制約型スキル

禁止事項や安全制約を定義するスキル。

### 目的

危険・低品質・非推奨な行動を防ぐ。

### 特徴

- 明示的禁止
- 安全制約
- 非推奨事項
- リスク軽減

### 例

- security-guardrails
- migration-safety
- production-safe-sql
- ai-coding-rules

### 例

- 不必要な全ファイル書き換えを禁止
- 破壊的変更を黙って導入しない
- 過剰抽象化を避ける
- 隠れた副作用を避ける

---

## 4. 思想・スタイル型スキル

チームや設計思想を定義するスキル。

### 目的

出力を特定の設計思想や文化に揃える。

### 特徴

- アーキテクチャ方針
- コーディング哲学
- 設計価値観
- トレードオフ方針

### 例

- engineering-principles
- clean-code-philosophy
- architecture-principles
- pragmatic-engineering

### 例

- 賢さより明快さを優先する
- 魔法より explicit を優先する
- 将来安全に変更できることを優先する
- 副作用を見えるようにする

---

## 5. ロール型スキル

特定の役割や視点を与えるスキル。

### 目的

モデルの振る舞いや思考スタイルを変える。

### 特徴

- 視点ベース
- 専門的立場
- 優先順位が異なる
- 擬似人格的

### 例

- planner
- verifier
- skeptic-reviewer
- architect
- performance-specialist
- tech-lead

### 典型的振る舞い

- 前提を疑う
- 代替案を検討する
- 正しさを検証する
- スケーラビリティを重視する
- リスク低減を重視する

---

## 6. 出力フォーマット型スキル

レスポンス構造を標準化するスキル。

### 目的

出力形式を統一し、再利用しやすくする。

### 特徴

- テンプレート
- フォーマット規則
- 構造化出力
- 一貫したセクション

### 例

- adr-writer
- incident-report-writer
- pr-review-format
- changelog-generator

### 例

- Summary
- Findings
- Risks
- Recommendations
- Follow-ups

---

## 7. 知識圧縮型スキル

特定分野の知識を圧縮して持たせるスキル。

### 目的

専門知識を再利用可能な形で提供する。

### 特徴

- ドメイン知識
- ベストプラクティス
- アンチパターン
- 専門用語

### 例

- postgres-performance
- rails-conventions
- aws-lambda-patterns
- ddd-tactical-patterns
- graphql-schema-design

### 典型内容

- 推奨パターン
- アンチパターン
- 運用観点
- パフォーマンス観点

---

## 8. スコアリング・査定型スキル

評価基準に基づき採点するスキル。

### 目的

品質評価や優先順位付けを行う。

### 特徴

- スコアリング
- 重み付け
- severity 判定
- readiness 評価

### 例

- production-readiness-review
- architecture-scorecard
- seniority-review
- scalability-assessment

### 典型出力

- スコア
- 強み
- 弱み
- 致命的リスク
- 推奨改善案

---

# 上位分類

多くの skill は、大きく以下の3種類に分類できる。

---

## A. 行動定義型

「何をするか」を定義する。

### 例

- ワークフロー型
- ロール型

---

## B. 判断定義型

「どう考え、どう評価するか」を定義する。

### 例

- レビュー系
- 思想系
- 査定系

---

## C. 制約定義型

「何を避けるべきか」を定義する。

### 例

- ガードレール系
- セーフティルール
- 本番制約

---

# 複合型スキル

実務で強い skill は、複数カテゴリを組み合わせていることが多い。

例: next-developer-review

含まれるもの:

- 判断基準
- 設計思想
- ガードレール
- severity 判定
- 出力フォーマット

つまり、単一カテゴリではなく複合型である。


