## deploy

任意のファイルやディレクトリを、設定ファイルに書いた配置先へコピーするコマンドです。

### 使い方

リポジトリルートから実行する場合:

```bash
make deploy-docs-dry-run
make deploy-docs
```

`make deploy-docs-dry-run` は `scripts/deploy/deploy.example.json` を使って、コピー予定の内容だけを表示します。
`make deploy-docs` は同じ設定ファイルを使って実際にコピーします。

`scripts/deploy` ディレクトリで直接実行する場合:

```bash
go run ./cmd -config=deploy.json -dry-run
go run ./cmd -config=deploy.json
```

`-dry-run` を付けると、実際にはコピーせずに作成予定のディレクトリとコピー予定のファイルを出力します。

ビルド済みバイナリを使う場合:

```bash
make build-deploy
./bin/deploy -config=./scripts/deploy/deploy.example.json -dry-run
./bin/deploy -config=./scripts/deploy/deploy.example.json
```

### オプション

| フラグ | 必須 | 説明 |
|---|---|---|
| `-config` | Yes | コピー元とコピー先を書いた JSON 設定ファイル |
| `-dry-run` | No | 実際にはコピーせず、予定される `MKDIR` / `COPY` を出力する |

### 設定ファイル

設定ファイルは JSON です。相対パスは設定ファイルのあるディレクトリから解決されます。`~/` から始まるパスはホームディレクトリへ展開されます。

```json
{
  "items": [
    {
      "source": "../../out/.codex/skills",
      "destination": "~/.codex/skills"
    },
    {
      "source": "../../out/.codex/agents",
      "destination": "~/.codex/agents"
    },
    {
      "source": "../../out/.codex/Agents.md",
      "destination": "~/.codex/Agents.md"
    },
    {
      "source": "../../codex/config.toml",
      "destination": "~/.codex/config.toml"
    }
  ]
}
```

`items` は上から順番に処理されます。ディレクトリもファイルも同じ `source` / `destination` 形式で指定できます。

### dry-run の出力例

```text
DRY-RUN item[0] dir  /repo/out/.codex/skills -> /Users/me/.codex/skills
MKDIR    /Users/me/.codex/skills
COPY     /repo/out/.codex/skills/example/SKILL.md -> /Users/me/.codex/skills/example/SKILL.md
DRY-RUN item[1] file /repo/codex/config.toml -> /Users/me/.codex/config.toml
COPY     /repo/codex/config.toml -> /Users/me/.codex/config.toml
```

### コピー仕様

- `source` がファイルの場合、`destination` を配置先ファイルパスとしてコピーします。
- `source` がディレクトリの場合、ディレクトリの中身を `destination` ディレクトリ配下へコピーします。
- 既存ファイルは上書きします。
- コピー先にある余分なファイルは削除しません。
- 通常ファイルとディレクトリ以外はスキップします。
