# deploy Design

このディレクトリは二つの配布責務を持つ。`deploy` は plugin marketplace と、agent、設定、rule、補助コマンドなどの非-skill 資産を `deploy.json` に従って配布する。`plugins` は group・個別 skill と手管理 plugin から marketplace を生成する。

利用者向けの手順は [README.md](README.md)、skill plugin の選択・host 登録・移行は [docs/plugin-distribution.md](../../docs/plugin-distribution.md) を正本とする。

## 責務境界

| 層 | 責務 | 持ち込まないもの |
|---|---|---|
| `cmd` | CLI flags を読み、各 runner を起動する | host 固有の設定変更 |
| `internal/deploy` | marketplace と非-skill 資産の backup、replace、copy、report | skill の解決・plugin manifest |
| `internal/plugins` | skill group 化、個別 plugin、manifest と marketplace の生成 | Codex / Claude への登録・有効化 |
| external skill loader | 全外部 skill の取得と hash lock 検証 | 恒久的な host cache |
| file operations | コピー、backup、replace、dry-run report | 設定解釈 |

host の marketplace 登録、plugin install / enable、host cache の更新、host 設定ファイルの直接編集は builder / deploy の責務外である。marketplace を生成・配布しても有効化状態は変えない。

## Plugin build flow

```text
skills/internal + skill-groups.json + external-skills.json + plugins/
                       |
                       v
 group / individual skill sources + verified external skills + managed plugins
                       |
                       v
 out/marketplace/plugins/<plugin>
                       |
                       v
 marketplace manifests and summary.json
                       |
                       v
 backup + replace ~/.local/share/agents-docs/marketplace
```

group の各 skill と外部 skill は `skill-<name>` 個別 plugin にも含める。個別 plugin は必須相対 resource の closure も含める。`plugins/<name>` の手管理 plugin は source manifest、skill frontmatter、symlink を検証して取り込み、個別 `skill-<name>` plugin は生成しない。Codex manifest の local source は marketplace root から見た `./plugins/<plugin>` である。

外部 skill は `external-skills.json` の URL と hash lock を使い、全件を一時領域へ取得して検証する。外部 skill の destination は lock に保存しない。

## Non-skill deployment flow

`deploy.Runner` は次の順に処理する。

1. `deploy.json` を読み込み、source と destination を解決する。
2. destination が存在すれば backup を作る。
3. `replace` が指定されていれば destination を洗い替える。
4. copy と report を実行する。dry-run は書き込まず予定だけを報告する。

`source` の相対パスは実行時の working directory、`destination` の相対パスは設定ファイルのあるディレクトリを基準にする。この公開契約を変更する場合は integration test と README を更新する。

## 検証

変更後は少なくとも次を確認する。

```bash
cd scripts/deploy
go test ./...
cd ../..
make build-plugins
make deploy-plugins-dry-run
```

公開契約を変更する場合は README、DESIGN、[plugin 配布](../../docs/plugin-distribution.md)、該当するテストを同時に更新する。
