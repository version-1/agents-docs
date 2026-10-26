# Agents-docs

Codex / Claude などの AI エージェントに配布するルール、プロンプト、skill、agent 定義を管理するレポジトリです。

このリポジトリでは、主に次の項目を管理しています。

- Codex / Claude 向けの共通指示
- agent の役割定義とメタデータ
- skill の本体と UI 表示用メタデータ
- 再利用するプロンプトとルール
- skill plugin marketplace を生成・配布するツール

## Skill plugin のインストール

Codex CLI または Claude Code がインストールされた環境で、リポジトリのルートから実行します。まず、全 group、`skills/internal/` と外部 skill の個別 plugin、手管理 plugin を含む marketplace を生成します。

### 1. marketplace を生成

```bash
make build-plugins
```

基本作業用の `core` plugin には、`git`、`delivery`、`documentation` が含まれます。finder 以外の role を使う場合は、`roles` plugin を別途 install します。

### 2. 利用する環境にインストール

**Codex**

```bash
codex plugin marketplace add "$PWD/out/marketplace"
codex plugin add core@agents-docs
```

**Claude Code**

```bash
claude plugin marketplace add --scope user "$PWD/out/marketplace"
claude plugin install -s user core@agents-docs
```

marketplace を登録済みなら、必要な plugin を追加する最後のコマンドだけを実行します。利用する plugin 名は `core` から置き換えられます。`core` を install した場合は、その構成 group や同じ skill の個別 plugin を重複して install しないでください。追加後は Codex または Claude Code を再起動します。

### 固定パスへ配布する場合

生成した marketplace を固定パスで利用する場合は、次で配布します。

```bash
make deploy-plugins-dry-run
make deploy-plugins
```

登録時は、上記の `"$PWD/out/marketplace"` を `"$HOME/.local/share/agents-docs/marketplace"` に置き換えます。配布コマンドは marketplace の登録や plugin のインストールを行いません。

plugin の選択・依存関係・更新・移行については [plugin 配布](docs/plugin-distribution.md) を参照してください。

## 配布コマンド

- `make build-deploy` で非-skill 資産を配布する deploy 用バイナリを生成します。
- `make deploy-dry-run` で非-skill 資産と plugin marketplace の配布予定を確認します。
- `make deploy` で非-skill 資産と plugin marketplace を配布します。
- `make deploy-plugins[-dry-run]` は plugin marketplace だけを扱います。`TAG=skills` は marketplace 配布を対象にします。

`make deploy` は Codex 用の補助コマンドも `~/.codex/bin` に配置します。
push と rebase には、ポリシーに基づいて Git 操作を制限する `grape` を使います。
`safe-gh-edit` は自分が作成した PR / Issue だけを prompt なしで編集するための `gh pr edit` / `gh issue edit` wrapper です。
`safe-local-curl` は localhost / loopback / private address 宛ての確認だけを prompt なしで行うための安全な `curl` wrapper です。

## Skill の正本と設定ファイル

- `skills/internal/`: 共通 skill の正本です。
- `plugins/`: group に含めない手管理 plugin の正本です。`plugins/roles/skills/` には finder 以外の role skill を置きます。
- `skill-groups.json`: plugin group と含める skill を定義します。
- `external-skills.json`: 外部 skill の取得元と hash lock を定義します。
- `deploy.json`: plugin marketplace と非-skill の agent / 設定などの配布先を定義します。

## 関連ドキュメント

- skill の用途一覧: [docs/skill-library.md](docs/skill-library.md)
- skill の依存関係: [docs/skill-dependency-map.md](docs/skill-dependency-map.md)
- plugin の配布・更新・移行: [docs/plugin-distribution.md](docs/plugin-distribution.md)
- deploy コマンドの詳細: [scripts/deploy/README.md](scripts/deploy/README.md)
- レポジトリ設計: [docs/design.md](docs/design.md)

## ライセンス

このリポジトリの内容は Apache License 2.0 のもとで利用できます。詳細は [LICENSE](LICENSE) を参照してください。

改変した内容を公開・配布する場合は、可能であればブログ、README、ドキュメントなどで、このリポジトリに多少なりとも触れてもらえるとうれしいです。このお願いはライセンス条件ではなく、任意のクレジットとしての希望です。
