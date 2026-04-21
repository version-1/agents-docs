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

`-dry-run` を付けると、実際にはコピーせずにバックアップ、削除、作成予定のディレクトリとコピー予定のファイルを出力します。

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
| `-dry-run` | No | 実際にはコピーせず、予定される `REMOVE` / `MKDIR` / `COPY` / `SKIP` を出力する |

### 設定ファイル

設定ファイルは JSON です。`source` の相対パスは deploy コマンドを実行したカレントディレクトリから解決されます。`destination` の相対パスは設定ファイルのあるディレクトリから解決されます。`~/` から始まるパスはホームディレクトリへ展開されます。

```json
{
  "items": [
    {
      "source": "out/.codex/skills",
      "destination": "~/.codex/skills",
      "replace": true,
      "exclude": [
        "**/.DS_Store",
        "**/*.tmp"
      ]
    },
    {
      "source": "out/.codex/agents",
      "destination": "~/.codex/agents",
      "replace": true
    },
    {
      "source": "out/.codex/Agents.md",
      "destination": "~/.codex/Agents.md"
    },
    {
      "source": "codex/config.toml",
      "destination": "~/.codex/config.toml"
    }
  ]
}
```

`items` は上から順番に処理されます。ディレクトリもファイルも同じ `source` / `destination` 形式で指定できます。

コピー先に既存のファイルまたはディレクトリがある場合は、コピー前にバックアップします。バックアップは 1 回の実行につき 1 つのタイムスタンプ付きディレクトリにまとめて作成され、`destination` の絶対パス構造を再現します。バックアップ先は設定ファイルと同じディレクトリ配下の `.deploy-backups/<timestamp>/` です。

リポジトリルートから `./bin/deploy -config=./scripts/deploy/deploy.example.json` を実行する場合、`source` は `out/.codex/skills` のようにリポジトリルート基準で指定します。

`replace` は省略可能です。`true` の場合、コピー前に `destination` を削除してから配置します。`false` または未指定の場合は既存ファイルを上書きするだけで、コピー先にある余分なファイルは残します。

`exclude` は省略可能です。コピー元の中で除外したいファイルやディレクトリを glob で指定します。

- パターンは `source` からの相対パスに対して評価します。
- パス区切りは `/` で指定します。
- `*` は `/` を含まない任意の文字列に一致します。
- `?` は `/` を含まない任意の 1 文字に一致します。
- `**` は `/` を含む任意の文字列に一致します。
- `/` を含まないパターンはファイル名・ディレクトリ名だけにも一致します。
- ディレクトリに一致した場合、その配下はまとめてスキップします。

例:

```json
{
  "items": [
    {
      "source": "./src",
      "destination": "./dest",
      "replace": true,
      "exclude": [
        "*.tmp",
        "**/*.log",
        "cache/**"
      ]
    }
  ]
}
```

### dry-run の出力例

```text
DRY-RUN item[0] dir  /repo/out/.codex/skills -> /Users/me/.codex/skills
BACKUP   /Users/me/.codex/skills -> /repo/scripts/deploy/.deploy-backups/20260421-142600/Users/me/.codex/skills
REMOVE   /Users/me/.codex/skills
MKDIR    /Users/me/.codex/skills
COPY     /repo/out/.codex/skills/example/SKILL.md -> /Users/me/.codex/skills/example/SKILL.md
SKIP     /repo/out/.codex/skills/example/debug.tmp
DRY-RUN item[1] file /repo/codex/config.toml -> /Users/me/.codex/config.toml
COPY     /repo/codex/config.toml -> /Users/me/.codex/config.toml
```

### コピー仕様

- `source` がファイルの場合、`destination` を配置先ファイルパスとしてコピーします。
- `source` がディレクトリの場合、ディレクトリの中身を `destination` ディレクトリ配下へコピーします。
- コピー先に既存のファイルまたはディレクトリがある場合、コピー前にバックアップします。
- バックアップ先は `.deploy-backups/<timestamp>/` 配下で、`destination` の絶対パス構造を再現します。
- 実行時には `BACKUP <destination> -> <backup path>` を出力します。
- 既存ファイルは上書きします。
- `replace` が `false` または未指定の場合、コピー先にある余分なファイルは削除しません。
- `replace` が `true` の場合、コピー前に `destination` を削除します。
- 通常ファイルとディレクトリ以外はスキップします。
- `exclude` に一致したファイルやディレクトリはスキップします。
