# codex ディレクトリ

Codex CLI 向けの設定、エージェント定義、プロンプト、コマンド、スキルを配置するディレクトリです。

## ディレクトリ構成

```text
codex
├── config.toml     Codex CLI の設定
├── agents          役割別エージェント定義
├── commands        コマンドとして利用するスキル定義
├── prompts         プロンプト定義
└── skills          タスク別・実践別のスキル定義
```

## 各ディレクトリ

### agents

Codex CLI で利用する役割別エージェントの定義を配置します。

```text
agents
├── documenter.toml  ドキュメント作成を担当するエージェント
├── implementer.toml 実装を担当するエージェント
├── planner.toml     計画・設計を担当するエージェント
├── reviewer.toml    レビューを担当するエージェント
└── scount.toml      調査を担当するエージェント
```

### commands

Codex CLI のコマンドとして利用するスキル定義を配置します。

```text
commands
└── pr-description
    └── SKILL.md     PR 説明文生成用のコマンド
```

### prompts

Codex CLI で再利用するプロンプト定義を配置します。

```text
prompts
├── hello.md         動作確認用のプロンプト
└── rmbranch.md      Git ブランチ削除用のプロンプト
```

### skills

Codex CLI のスキル定義を配置します。各スキルはディレクトリ単位で管理し、`SKILL.md` に利用条件や手順を記述します。

```text
skills
├── code-architect  アーキテクチャ設計・改善
├── code-refactor   リファクタリング
├── code-review     コードレビュー
├── code-tracer-go  Go コードの呼び出し経路調査
├── code-typo       typo 指摘
├── coding          コーディング全般
├── documenting     ドキュメント作成
└── testing         テスト・検証
```
