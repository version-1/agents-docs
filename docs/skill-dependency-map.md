# Skill Dependency Map

このドキュメントは、internal skill 同士の明示的な依存関係を整理するためのマップです。

ここでの依存は、`SKILL.md` の本文や description に「併用する」「優先する」「使う」「参照する」と明記されている関係だけを対象にします。似た目的を持つだけの関係や、エージェントロール上の暗黙の関係は含めません。

## 全体図

```mermaid
graph TD
    cmd_create_pr["cmd-create-pr"] --> cmd_commit["cmd-commit"]
    cmd_create_pr --> role_reviewer["role-reviewer"]
    cmd_create_pr --> code_review["code-review"]
    cmd_create_pr --> format_pr_description["format-pr-description"]

    cmd_start_branch["cmd-start-branch"] --> cmd_dispatch_agent["cmd-dispatch-agent"]
    cmd_start_branch --> cmd_rmbranch["cmd-rmbranch"]

    beautify_commit["beautify-commit"] --> cmd_commit
    beautify_commit --> grill_me["grill-me (external)"]

    role_reviewer --> code_review
    role_implementer["role-implementer"] --> coding["coding"]

    coding --> code_naming["code-naming"]
    coding --> code_review

    component_design["component-design"] --> coding
    component_design --> code_test["code-test"]
    component_design --> code_review
    component_design --> code_next_dev["code-next-developer-review"]

    code_naming --> code_typo["code-typo"]
    code_naming --> code_review
    code_naming --> coding
    code_naming --> code_next_dev

    code_quality["code-quality-review"] --> code_review
    code_quality --> code_next_dev

    code_review --> coding
    code_review --> code_quality

    code_next_dev --> code_review

    code_simplifier["code-simplifier"] --> code_quality
    code_simplifier --> code_review
    code_simplifier --> code_next_dev
```

## ワークフロー系

| Skill | 依存先 | 関係 |
|---|---|---|
| `cmd-create-pr` | `cmd-commit` | 未コミット変更や High 指摘対応後の追加差分を commit するときに使う。 |
| `cmd-create-pr` | `role-reviewer` | PR 作成・更新前のレビューゲートとして使う。 |
| `cmd-create-pr` | `code-review` | PR 前レビューで必要に応じて併用する。 |
| `cmd-create-pr` | `format-pr-description` | PR description を作るときに使う。 |
| `cmd-start-branch` | `cmd-dispatch-agent` | 不要ブランチ整理を別 agent に投げるときに使う。 |
| `cmd-start-branch` | `cmd-rmbranch` | 不要ブランチ整理 agent の依頼内容として使う。 |
| `beautify-commit` | `cmd-commit` | 整理後の commit 作成で使う。 |
| `beautify-commit` | `grill-me` | 分割方針が曖昧なときに整理方針を詰めるため使う。 |

## ロール系

| Skill | 依存先 | 関係 |
|---|---|---|
| `role-reviewer` | `code-review` | Reviewer として検証するとき、通常レビュー観点も併せて参照する。 |
| `role-implementer` | `coding` | Implementer として実装するときの基本方針として併用する。 |

## コード作業系

| Skill | 依存先 | 関係 |
|---|---|---|
| `coding` | `code-naming` | 命名判断で迷う場合に併用する。 |
| `coding` | `code-review` | レビューのみが目的の場合は `code-review` を優先する。 |
| `component-design` | `coding` | UI 設計後に実装へ進む場合に併用する。 |
| `component-design` | `code-test` | テスト設計やテスト追加で併用する。 |
| `component-design` | `code-review` | 実装済み差分の欠陥レビューでは優先する。 |
| `component-design` | `code-next-developer-review` | 次の開発者の読みやすさを確認する場合に併用する。 |
| `code-naming` | `code-typo` | typo や spelling の検査だけなら優先する。 |
| `code-naming` | `code-review` | バグ、仕様違反、セキュリティ、データ整合性のレビューなら優先する。 |
| `code-naming` | `coding` | 実装変更やリネーム作業まで行う場合に併用する。 |
| `code-naming` | `code-next-developer-review` | 次に開発する人の理解しやすさ全体を見る場合に併用する。 |
| `code-quality-review` | `code-review` | 即時欠陥の検出が主目的なら優先する。 |
| `code-quality-review` | `code-next-developer-review` | 次の開発者の迷いやすさが主目的なら優先する。 |
| `code-review` | `coding` | 具体的なコード変更の実装方法を相談するときに優先する。 |
| `code-review` | `code-quality-review` | 将来的な変更容易性やコード品質を重点的に見るときに優先する。 |
| `code-next-developer-review` | `code-review` | 欠陥や重大リスクを見つけた場合に通常レビュー観点として分ける。 |
| `code-simplifier` | `code-quality-review` | 品質レポートや quality gate の観点から改善候補を拾う。 |
| `code-simplifier` | `code-review` | 可読性、効率性、テスト容易性などのレビュー観点として使う。 |
| `code-simplifier` | `code-next-developer-review` | 保守性や次の開発者の理解しやすさの観点として使う。 |

## 更新ルール

skill を追加または更新したときに、他 skill を明示的に参照する文を増やした場合は、このドキュメントも更新します。

更新時は次を確認します。

- `SKILL.md` に `$skill-name`、`` `skill-name` ``、または plain text で依存先が書かれているか。
- 依存が「必ず使う」「必要に応じて併用」「代替として優先」のどれに近いか。
- Mermaid の全体図と表の両方に同じ関係が載っているか。
- 外部 skill への依存は `(external)` と明記しているか。
