これは agents-docs レポジトリ用の AGENTS.md です。
Codex / Claude 向けの agent、skill、設定、プロンプトを管理し、`make deploy` でローカル環境へ配布します。

## ディレクトリ構成

```text
.
├── AGENTS.md                      このファイル（リポジトリの案内）
├── README.md                      リポジトリ概要と基本コマンド
├── Makefile                       deploy ビルドと配布コマンド
├── deploy.json                    内部ファイルの配布設定
├── external-skills.json           外部 skill の取得・配布設定
├── bin                            生成されたバイナリ置き場
│   └── deploy                     deploy 用バイナリ
├── codex                          Codex CLI 用の配布元
│   ├── AGENTS.md                  Codex に配布する AGENTS.md
│   ├── config.toml                Codex CLI 設定
│   ├── agents                     Codex agent 定義
│   │   └── <role>.toml            各 agent の説明、指示、nickname 候補
│   ├── prompts                    再利用プロンプト
│   ├── rules                      Codex ルール
│   └── skills                     Codex skill 定義
│       ├── internal               この repo で管理する skill
│       │   ├── cmd-*              コマンド系 skill
│       │   ├── code-*             コード調査・設計・レビュー系 skill
│       │   ├── role-*             agent ロール用 skill
│       │   └── <name>             その他の internal skill
│       └── external               外部 skill の配置先
├── claude                         Claude 向けの配布元
│   ├── CLAUDE.md                  Claude に配布する指示
│   ├── agents                     Claude agent 定義
│   └── settings.json              Claude 設定
├── docs                           設計ドキュメント
│   ├── design.md
│   ├── design-skill.md
│   ├── skill-library.md           skill の用途一覧
│   └── memo                       メモ類
│       └── skill.md
└── scripts
    └── deploy                     deploy コマンドの Go 実装
        ├── cmd
        ├── internal
        ├── go.mod
        └── README.md
```

`.deploy-backups/` と `scripts/deploy/.deploy-backups/` は deploy 実行時のバックアップ出力です。通常は直接編集しません。

## コマンド

- `make build-deploy` で deploy 用のバイナリを生成します。
- `make deploy-dry-run` で `deploy.json` と `external-skills.json` に基づくコピー予定を確認します。
- `make deploy` で Codex / Claude 用の設定、agent、skill、prompt、rule を配布します。
- `make apply` は `make deploy` の別名です。

## 配布時の確認

通常は、次の順で実行します。

```bash
make deploy-dry-run
make deploy
```

`deploy-dry-run` でコピー予定、バックアップ先、削除予定、スキップ件数を確認してから `deploy` してください。

## 変更時の注意

- `codex/skills/internal/<name>/SKILL.md` が Codex skill の本体です。
- `codex/skills/internal/<name>/agents/openai.yaml` がある場合は、UI 表示用メタデータも内容に合わせます。
- skill の追加・更新時は [`docs/skill-library.md`](docs/skill-library.md) の一覧も更新します。
- `codex/agents/*.toml` の `nickname_candidates` は覚えやすい英語の人名を 3 件ずつ設定します。
- `deploy.json` は内部ファイルの配布先、`external-skills.json` は外部 skill の取得先を管理します。
- `claude/agents/implemnter.md` と `codex/agents/scount.toml` は現状のファイル名です。リネームする場合は配布設定や参照も合わせて確認します。

## テスト実行時の補足

Go テストでキャッシュ権限エラーが出る場合は、以下の環境変数を指定して実行してください。

```bash
GOMODCACHE=/tmp/gomodcache GOCACHE=/tmp/gocache GOTOOLCHAIN=local go test ./...
```
