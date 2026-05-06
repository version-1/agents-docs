# 設計

## ディレクトリ構成

このレポジトリでは、実際に AI に配布するドキュメントを `codex/` と `claude/` 配下に配置します。

`docs/` 配下は設計や運用方針など、レポジトリ自体の補助ドキュメントを置く場所です。

## deploy 設定

`deploy.json` の `items` に、コピー元とコピー先を設定します。

```json
{
  "items": [
    {
      "source": "codex/skills",
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

## replace による洗い替え

`replace` はコピー先を洗い替えするかどうかを item ごとに指定するフラグです。

```json
{
  "items": [
    {
      "source": "codex/agents",
      "destination": "~/.codex/agents",
      "replace": true
    },
    {
      "source": "codex/config.toml",
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
  source:      /repo/codex/agents
  destination: /Users/me/.codex/agents
  backup: /repo/.deploy-backups/20260421-142600/Users/me/.codex/agents
  replace: remove existing destination
  summary: 8 copied, 4 dirs, 0 skipped
```

# カスタムスキルの命名

display_name には ケバブケース（例: `cmd-run`）を使用してください。
prefix はスキルのカテゴリを表すことが多いですが、必須ではありません。

- code-*: コード調査、設計、レビューなどコード関連のスキル
- cmd-*: コマンド実行やシステム操作などのスキル
- role-*: agent のロールを定義するスキル
