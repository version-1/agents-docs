
## Agents-docs

Codex / Claude などの AI エージェントに配布するルール、プロンプト、skill、agent 定義を管理するレポジトリです。

このリポジトリでは、主に次の項目を管理しています。

- Codex / Claude 向けの共通指示
- agent の役割定義とメタデータ
- skill の本体と UI 表示用メタデータ
- 再利用するプロンプトとルール
- skill plugin marketplace を生成・配布するツール

## skill 一覧

- [docs/skill-library.md](docs/skill-library.md)
- [docs/skill-dependency-map.md](docs/skill-dependency-map.md)

## Skill plugin の利用

skill の正本は `skills/internal/` です。全 group と全 skill の個別 plugin を marketplace に生成し、Codex / Claude では必要な plugin だけを native install します。生成・配布と host への登録・有効化は別の段階です。詳細は [plugin 配布](docs/plugin-distribution.md) を参照してください。

```bash
make build-plugins
make deploy-plugins-dry-run
make deploy-plugins
```

## コマンド

- `make build-deploy` で非-skill 資産を配布する deploy 用バイナリを生成します。
- `make deploy-dry-run` で非-skill 資産と plugin marketplace の配布予定を確認します。
- `make deploy` で非-skill 資産と plugin marketplace を配布します。
- `make deploy-plugins[-dry-run]` は plugin marketplace だけを扱います。`TAG=skills` は marketplace 配布を対象にします。

`make deploy` は Codex 用の補助コマンドも `~/.codex/bin` に配置します。
push と rebase には、ポリシーに基づいて Git 操作を制限する `grape` を使います。
`safe-gh-edit` は自分が作成した PR / Issue だけを prompt なしで編集するための `gh pr edit` / `gh issue edit` wrapper です。
`safe-local-curl` は localhost / loopback / private address 宛ての確認だけを prompt なしで行うための安全な `curl` wrapper です。

## デプロイ手順

1. `make build-plugins` を実行します。
2. 必要なら `make deploy-plugins-dry-run`、`make deploy-plugins` を実行します。
3. [plugin 配布](docs/plugin-distribution.md) に従い、Codex / Claude で marketplace を登録して必要な plugin を有効化します。

## 設定ファイル

- `skill-groups.json`: plugin group と含める skill を定義します。
- `external-skills.json`: 外部 skill の取得元と hash lock を定義します。
- `deploy.json`: plugin marketplace と非-skill の agent / 設定などの配布先を定義します。

## 関連ドキュメント

- deploy コマンドの詳細: [scripts/deploy/README.md](scripts/deploy/README.md)
- レポジトリ設計: [docs/design.md](docs/design.md)

## ライセンス

このリポジトリの内容は Apache License 2.0 のもとで利用できます。詳細は [LICENSE](LICENSE) を参照してください。

改変した内容を公開・配布する場合は、可能であればブログ、README、ドキュメントなどで、このリポジトリに多少なりとも触れてもらえるとうれしいです。このお願いはライセンス条件ではなく、任意のクレジットとしての希望です。
