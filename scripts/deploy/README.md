## deploy

任意のファイルやディレクトリを、設定ファイルに書いた配置先へコピーするコマンドです。

### 使い方

リポジトリルートから実行する場合:

```bash
make deploy-dry-run
make deploy
```

`make deploy-dry-run` は repo root の `deploy.json` と `external-skills.json` を使って、コピー予定の内容だけを表示します。
`make deploy` は同じ設定ファイルを使って実際にコピーします。

Go コマンドで直接実行する場合:

```bash
go run ./scripts/deploy/cmd -config=deploy.json -external-skills=external-skills.json -dry-run
go run ./scripts/deploy/cmd -config=deploy.json -external-skills=external-skills.json
```

`-dry-run` を付けると、実際にはコピーせずにバックアップ、削除、作成予定のディレクトリとコピー予定のファイルを出力します。

ビルド済みバイナリを使う場合:

```bash
make build-deploy
./bin/deploy -config=./deploy.json -external-skills=./external-skills.json -dry-run
./bin/deploy -config=./deploy.json -external-skills=./external-skills.json
```

### オプション

| フラグ | 必須 | 説明 |
|---|---|---|
| `-config` | Yes | コピー元とコピー先を書いた JSON 設定ファイル |
| `-external-skills` | No | 外部 skill を取得して配布する JSON 設定ファイル |
| `-dry-run` | No | 実際にはコピーせず、item ごとの予定と件数サマリを出力する |
| `-no-color` | No | ANSI カラー出力を無効にする |

### 設定ファイル

設定ファイルは JSON です。`source` の相対パスは deploy コマンドを実行したカレントディレクトリから解決されます。`destination` の相対パスは設定ファイルのあるディレクトリから解決されます。`~/` から始まるパスはホームディレクトリへ展開されます。

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
    },
    {
      "source": "codex/agents",
      "destination": "~/.codex/agents",
      "replace": true
    },
    {
      "source": "codex/Agents.md",
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

リポジトリルートから `./bin/deploy -config=./deploy.json` を実行する場合、`source` は `codex/skills` のようにリポジトリルート基準で指定します。

`replace` は省略可能です。`true` の場合、コピー前に `destination` を削除してから配置します。`false` または未指定の場合は既存ファイルを上書きするだけで、コピー先にある余分なファイルは残します。

`flatten` は省略可能です。`true` の場合、`source` 配下から `SKILL.md` を持つディレクトリを探し、そのディレクトリを `destination/<ディレクトリ名>` に配置します。例えば `source/internal/role-planner/SKILL.md` は `destination/role-planner/SKILL.md` として配置されます。`flatten` はディレクトリの `source` でのみ使用できます。同じディレクトリ名の skill が複数見つかった場合は、上書きせずエラーにします。

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
[DRY-RUN] item[0] dir
  source:      /repo/codex/skills
  destination: /Users/me/.codex/skills
  backup: /repo/scripts/deploy/.deploy-backups/20260421-142600/Users/me/.codex/skills
  replace: remove existing destination
  summary: 18 copied, 12 dirs, 1 skipped

[DRY-RUN] item[1] file
  source:      /repo/codex/config.toml
  destination: /Users/me/.codex/config.toml
  summary: 1 copied, 0 dirs, 0 skipped
```

通常の出力は見出し、バックアップ、洗い替え、サマリを色分けします。CI やログ保存などで ANSI エスケープを避けたい場合は `-no-color` を指定します。

### 外部 skill 設定

`-external-skills` を指定すると、外部の skill をネットワーク経由で取得して `destination` に配布します。未指定の場合、従来どおり `-config` の内容だけを処理します。

```json
[
  {
    "name": "grill-me",
    "url": "https://github.com/mattpocock/skills/tree/main/skills/productivity/grill-me",
    "type": "git",
    "destination": [
      "~/.codex/skills/external/grill-me",
      "~/.claude/skills/grill-me"
    ]
  }
]
```

`type` は現在 `git` のみ対応します。`url` は `https://github.com/<owner>/<repo>/tree/<ref>/<path>` 形式だけを受け付けます。dry-run でも取得と `SKILL.md` の存在確認を行うため、URL が不正、取得できない、または取得先に `SKILL.md` がない場合はエラーになります。

外部 skill 同士の `name` 重複、内部 skill と外部 skill の同名衝突、外部 skill の `destination` 重複は上書きせずエラーにします。

### コピー仕様

- `source` がファイルの場合、`destination` を配置先ファイルパスとしてコピーします。
- `source` がディレクトリの場合、ディレクトリの中身を `destination` ディレクトリ配下へコピーします。
- `flatten` が `true` の場合、`source` 配下の `SKILL.md` を持つディレクトリだけを `destination` 直下へフラットにコピーします。
- コピー先に既存のファイルまたはディレクトリがある場合、コピー前にバックアップします。
- バックアップ先は `.deploy-backups/<timestamp>/` 配下で、`destination` の絶対パス構造を再現します。
- 実行時には `backup: <backup path>` を出力します。
- 既存ファイルは上書きします。
- `replace` が `false` または未指定の場合、コピー先にある余分なファイルは削除しません。
- `replace` が `true` の場合、コピー前に `destination` を削除します。
- 通常ファイルとディレクトリ以外はスキップします。
- `exclude` に一致したファイルやディレクトリはスキップします。
