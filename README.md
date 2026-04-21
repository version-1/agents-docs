
## Agents-docs

AI にコードを書かせる際のルールやコンテキストをまとめるレポジトリです。

## ディレクトリについて

docs 配下にはスキルなどの元となるドキュメントが置かれていますが、
これは skill などを作る時の参考用で実際に AI に与えるドキュメントは codex/, agents/ 配下に配置しています。

## コマンド

- `make build-skillmaker` コマンドで skills 生成用のバイナリを生成します。
- `make build-deploy` で deploy 用のバイナリを生成します。
- `make gen-docs` で環境で使用するドキュメントを生成します。
- `make deploy-codex-docs` で生成済みのドキュメントを `~/.codex` に反映します。
- `make deploy-docs-dry-run` で `scripts/deploy/deploy.example.json` に基づくコピー予定を確認します。
- `make deploy-docs` で `scripts/deploy/deploy.example.json` に基づいて生成済みドキュメントを配置します。
  - generator の詳細は `scripts/generator/README.md` を参照してください。
  - deploy の詳細は `scripts/deploy/README.md` を参照してください。

## デプロイ手順

1. `make gen-docs` を実行して `docs/ja` から生成物を作成します。
2. `make deploy-docs-dry-run` を実行して、コピー予定のファイルとディレクトリを確認します。
3. `make deploy-docs` を実行して、Codex / Claude 用の設定とドキュメントを反映します。

生成と反映は一続きの手順なので、通常はこの2コマンドを続けて実行します。

```bash
make gen-docs
make deploy-docs-dry-run
make deploy-docs
```

## deploy 設定

`scripts/deploy/deploy.example.json` の `items` に、コピー元とコピー先を設定します。

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
    }
  ]
}
```

- `source` がファイルの場合、`destination` を配置先ファイルパスとしてコピーします。
- `source` がディレクトリの場合、ディレクトリの中身を `destination` ディレクトリ配下へコピーします。
- `source` の相対パスは deploy コマンドを実行したカレントディレクトリから解決します。
- `destination` の相対パスは設定ファイルのあるディレクトリから解決します。
- コピー先に既存のファイルまたはディレクトリがある場合、コピー前にバックアップします。
- バックアップ先は設定ファイルと同じディレクトリ配下の `.deploy-backups/<timestamp>/` です。
- バックアップは 1 回の実行につき 1 つのタイムスタンプ付きディレクトリにまとめ、`destination` の絶対パス構造を再現します。
- 実行時には `backup: <backup path>` を出力します。
- `replace` は省略可能で、`true` の場合はコピー前に `destination` を削除します。
- `replace` が `false` または未指定の場合、コピー先にある余分なファイルは残します。
- `exclude` は省略可能で、`source` からの相対パスに対する glob として評価します。
- `exclude` では `*` / `?` / `**` を指定できます。
- dry-run では実際にはコピーせず、item ごとの予定と件数サマリを出力します。
- 通常の出力は ANSI カラーで色分けされます。色を無効化する場合は `-no-color` を指定します。

### replace による洗い替え

`replace` はコピー先を洗い替えするかどうかを item ごとに指定するフラグです。

```json
{
  "items": [
    {
      "source": "../../out/.codex/agents",
      "destination": "~/.codex/agents",
      "replace": true
    },
    {
      "source": "../../codex/config.toml",
      "destination": "~/.codex/config.toml",
      "replace": false
    }
  ]
}
```

- `replace: true` の場合、コピー前に `destination` を `RemoveAll` 相当で削除します。
- `source` がディレクトリの場合、`destination` ディレクトリを削除してから、`source` の中身を配置します。
- `source` がファイルの場合、`destination` ファイルを削除してからコピーします。
- `replace: false` または未指定の場合、コピー先の余分なファイルやディレクトリは削除されません。
- `exclude` は `replace` 後のコピー対象に対して適用されます。`replace: true` の場合、除外されたファイルはコピーされませんが、既存の `destination` は先に削除されます。

dry-run では削除は実行されず、削除予定として `REMOVE` が出力されます。

```text
[DRY-RUN] item[0] dir
  source:      /repo/out/.codex/agents
  destination: /Users/me/.codex/agents
  backup: /repo/scripts/deploy/.deploy-backups/20260421-142600/Users/me/.codex/agents
  replace: remove existing destination
  summary: 8 copied, 4 dirs, 0 skipped
```
