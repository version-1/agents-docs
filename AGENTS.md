これは agents-docs レポジトリ用の AGENTS.md です。
Codex / Claude 向けの agent、skill、設定、プロンプトを管理します。skill は plugin marketplace として、その他の資産は `make deploy` でローカル環境へ配布します。

エージェント作業の標準フローは [`WORKFLOW.md`](WORKFLOW.md) を参照してください。

## ディレクトリ構成

```text
.
├── AGENTS.md                      このファイル（リポジトリの案内）
├── WORKFLOW.md                    エージェント作業の標準ワークフロー
├── README.md                      リポジトリ概要と基本コマンド
├── Makefile                       deploy ビルドと配布コマンド
├── deploy.json                    plugin marketplace と非-skill 資産の配布設定
├── skill-groups.json              plugin group の skill 構成
├── external-skills.json           外部 skill の取得元と hash lock
├── bin                            生成されたバイナリ置き場
│   └── deploy                     deploy 用バイナリ
├── codex                          Codex CLI 用の配布元
│   ├── AGENTS.md                  Codex に配布する AGENTS.md
│   ├── config.toml                Codex CLI 設定
│   ├── agents                     Codex agent 定義
│   │   └── <role>.toml            各 agent の説明、指示、nickname 候補
│   ├── prompts                    再利用プロンプト
│   ├── rules                      Codex ルール
├── skills                         Codex / Claude 共通の skill 正本
│   └── internal                   この repo で管理する skill
├── claude                         Claude 向けの配布元
│   ├── CLAUDE.md                  Claude に配布する指示
│   ├── agents                     Claude agent 定義
│   └── settings.json              Claude 設定
├── docs                           設計ドキュメント
│   ├── design.md
│   ├── design-skill.md
│   ├── skill-library.md           skill の用途一覧
│   └── guide                      ガイド類
│       └── skill-category.md      skill の類型ガイド
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
- `make build-plugins` で全 skill plugin marketplace を生成します。
- `make deploy-plugins-dry-run` で marketplace の配布予定を確認します。
- `make deploy-plugins` で marketplace を配布します。
- `make deploy[-dry-run]` は非-skill 資産と plugin marketplace を配布します。
- `make apply` は `make deploy` の別名です。

## ブランチ運用

- worktree で作業するブランチには、`worktrees/<番号>/` prefix を付けます（例: `worktrees/1/feat/refine-ui`）。
- worktree 以外で作業するブランチには、`chore/`、`fix/`、`feat/`、`docs/`、`refactor/` のいずれかの prefix を付けます。
- 作業ブランチをデフォルトブランチの最新状態へ追従させるときは、merge や `git rebase` ではなく `grape rebase` を使用します。

## 配布時の確認

通常は、次の順で実行します。

```bash
make build-plugins
make deploy-plugins-dry-run
make deploy-plugins
```

plugin の利用には配布後、Codex / Claude での marketplace 登録と plugin 有効化が必要です。手順は [`docs/plugin-distribution.md`](docs/plugin-distribution.md) を参照してください。

## 変更時の注意

- `skills/internal/<name>/SKILL.md` が共通 skill の本体です。
- 内部 skill の `SKILL.md` は、frontmatter を含む UTF-8 のファイルサイズで 5 KB（5,000 bytes）以下に収めます。
- `skills/internal/<name>/agents/openai.yaml` がある場合は、UI 表示用メタデータも内容に合わせます。
- skill の設計・命名時は [`docs/design-skill.md`](docs/design-skill.md) の方針も確認します。
- skill の追加・更新時は [`docs/skill-library.md`](docs/skill-library.md) の一覧も更新します。
- `codex/config.toml` の `[agents.<role>]` にある `nickname_candidates` は、覚えやすい英語の人名を 3 件ずつ設定します。
- `skill-groups.json` は plugin group の membership、`external-skills.json` は外部 skill の取得元と hash lock を管理します。
- `deploy.json` は plugin marketplace と非-skill 資産の配布先を管理します。
- `claude/agents/implemnter.md` は現状のファイル名です。リネームする場合は配布設定や参照も合わせて確認します。

## テスト実行時の補足

Go テストでキャッシュ権限エラーが出る場合は、以下の環境変数を指定して実行してください。

```bash
GOMODCACHE=/tmp/gomodcache GOCACHE=/tmp/gocache GOTOOLCHAIN=local go test ./...
```
