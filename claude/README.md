# claude ディレクトリ

Claude Code 向けの設定、エージェント定義、コマンド、スキルを配置するディレクトリです。

## ディレクトリ構成

```text
claude
├── CLAUDE.md       Claude Code に渡す共通の開発指針
├── settings.json   Claude Code の設定
├── agents          役割別エージェント定義
├── commands        スラッシュコマンド定義
└── skills          タスク別・実践別のスキル定義
```

## 各ディレクトリ

### agents

Claude Code で利用する役割別エージェントの定義を配置します。

```text
agents
├── documenter.md   ドキュメント作成を担当するエージェント
├── implemnter.md   実装を担当するエージェント
├── planner.md      計画・設計を担当するエージェント
├── reviewer.md     レビューを担当するエージェント
└── scout.md        調査を担当するエージェント
```

### commands

Claude Code で呼び出すコマンド定義を配置します。

```text
commands
├── hello.md        動作確認用のコマンド
└── rmbranch.md     Git ブランチ削除用のコマンド
```

### skills

Claude Code のスキル定義を配置します。各スキルはディレクトリ単位で管理し、`SKILL.md` に利用条件や手順を記述します。

```text
skills
├── code-architect  アーキテクチャ設計・改善
├── code-refactor   リファクタリング
├── code-review     コードレビュー
├── code-tracer     コードの呼び出し経路調査
├── code-typo       typo 指摘
├── coding          コーディング全般
├── documenting     ドキュメント作成
└── testing         テスト・検証
```
