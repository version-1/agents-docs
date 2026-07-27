# カスタムスキルの命名

display_name には ケバブケース（例: `cmd-run`）を使用してください。
prefix はスキルのカテゴリを表すことが多いですが、必須ではありません。

- code-*: コード調査、設計、レビューなどコード関連のスキル
- cmd-*: コマンド実行やシステム操作などのスキル
- format-*: 特定の出力形式、テンプレート、レポート構造へ情報を整理するスキル
- role-*: agent のロールを定義するスキル
- knowledge-*: 知識ベースやドキュメントを表すスキル

出力フォーマット型 skill の `format-*` prefix など、カテゴリ別の命名方針は [`docs/guide/skill-category.md`](guide/skill-category.md) を参照してください。

## Role skill と agent/playbook の構造

role 系は、委譲の入口、実行主体、作業手順を分けて定義します。

```text
role-<name> → <name> agent → role-<name>-playbook
```

- `role-<name>` は対応する agent を起動する入口です。frontmatter の `description` には、他の role と区別できる発火場面と対象外の場面を記載します。本文は agent 名を示すだけに留め、作業手順や実行時の判断基準は書きません。
- `codex/agents/<name>.toml` と `claude/agents/<name>.md` は agent のモデル、利用可能なツール、簡潔な説明を定義します。agent 本文には `role-<name>-playbook` を使って作業することだけを書きます。
- `role-<name>-playbook` は agent が実行する作業手順、判断基準、併用 skill、出力・検証の期待値を定義します。

新しい role を追加・変更するときは、この 3 層の整合を同時に確認し、routing 条件は入口の `description`、モデルと利用可能なツールは agent 定義、作業手順と実行時の判断基準は playbook に置いてください。同じ内容を複数の層へ重複して書かないでください。
