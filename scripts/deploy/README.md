## deploy

`scripts/deploy` は plugin marketplace と、agent、設定、rule、補助コマンドなどの非-skill 資産を配布する Go コマンドです。skill は `~/.codex/skills` や `~/.claude/skills` へコピーせず、plugin marketplace として生成・配布します。

skill plugin の生成、配布、Codex / Claude への登録は [plugin 配布](../../docs/plugin-distribution.md) を参照してください。

## 利用方法

通常はリポジトリ root で次を実行します。

```bash
make deploy-dry-run
make deploy
```

`make deploy-dry-run` は非-skill 資産と plugin marketplace の配布予定を確認します。`make deploy` は同じ対象を配布します。skill plugin だけを操作する場合は次を使います。

```bash
make build-plugins
make deploy-plugins-dry-run
make deploy-plugins
```

plugin 配布は `~/.local/share/agents-docs/marketplace` を backup 後に置き換えます。host の marketplace 登録、plugin の追加・有効化、cache の操作は deploy の責務ではありません。

## deploy CLI

直接実行する場合は、`deploy.json` に定義した plugin marketplace と非-skill 資産を扱います。

```bash
go run ./scripts/deploy/cmd -config=deploy.json -dry-run
go run ./scripts/deploy/cmd -config=deploy.json
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `-config` | Yes | コピー元とコピー先を書いた JSON 設定ファイル |
| `-tag` | No | 指定タグを持つ item だけを配布する |
| `-dry-run` | No | 実際にはコピーせず、予定と件数を出力する |
| `-no-color` | No | ANSI カラー出力を無効にする |

`source` の相対パスはコマンドのカレントディレクトリ、`destination` の相対パスは設定ファイルのあるディレクトリから解決します。`~/` はホームディレクトリに展開します。

```json
{
  "items": [
    {
      "source": "codex/agents",
      "tags": ["agents"],
      "destination": "~/.codex/agents",
      "replace": true
    },
    {
      "source": "codex/AGENTS.md",
      "tags": ["instructions"],
      "destination": "~/.codex/AGENTS.md"
    }
  ]
}
```

コピー先が存在する場合は先に backup を作成します。`replace: true` では backup 後にコピー先を洗い替えます。`preserve` は `replace: true` のディレクトリで保持する直下名を指定します。`mergeJSON: true` は JSON ファイルのトップレベルキーを既存設定へマージします。`template`、`exclude`、`tags` の詳細は `deploy.json` とテストを参照してください。

## 開発時の確認

```bash
cd scripts/deploy
go test ./...
```

公開仕様を変更するときは、この README、[DESIGN.md](DESIGN.md)、および [plugin 配布](../../docs/plugin-distribution.md) の該当箇所を同時に更新します。
