# codex ディレクトリ

Codex CLI 向けの設定、エージェント定義、プロンプト、コマンド、スキルを配置するディレクトリです。

## ディレクトリ構成

```text
codex
├── config.toml     Codex CLI の設定
├── agents          役割別エージェント定義
├── bin             Codex CLI から利用する補助コマンド
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
└── format-pr-description
    └── SKILL.md     PR 説明文生成用のコマンド
```

### bin

Codex CLI から prompt なしで実行させる補助コマンドを配置します。

- `safe-git-push`: 引数なし専用の安全な push wrapper。`main` / `master` / detached HEAD / 不正な branch / 複数 push URL を拒否し、`git push origin HEAD:<current-branch>` だけを実行します。
- `safe-gh-edit`: 自分が作成した PR / Issue だけを編集できる `gh pr edit` / `gh issue edit` wrapper。対象の author と認証中の GitHub user が一致しない場合は拒否します。
- `safe-local-curl`: localhost / loopback / private address 宛ての URL だけを実行できる `curl` wrapper。よく使う読み取り系オプションだけを許可し、外部 URL や複数 URL を拒否します。

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
├── code-tracer     コードの呼び出し経路調査
├── code-typo       typo 指摘
├── code-general    言語をまたぐ実装の共通原則
├── documenting     ドキュメント作成
└── testing         テスト・検証
```
