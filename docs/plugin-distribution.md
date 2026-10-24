# Plugin 配布

このリポジトリは、全 skill を Codex と Claude の plugin marketplace として生成する。利用環境ごとの選択は marketplace の生成時ではなく、各 host の native plugin install で行う。これにより、一つの生成物から Codex と Claude が独立して必要な group / skill plugin を選べる。

## 調査記録と移行設計

従来は `codex/skills` を `~/.codex/skills`、同じ内容を `~/.claude/skills` へ直接コピーしていた。この方式では、環境ごとに選ぶ単位がなく、host が管理する plugin の状態とも分離していた。外部 skill の配布先も lock に混在していた。

移行後は、共通原本から全 plugin を一度だけ生成し、host ごとに native install を実行する。marketplace の生成・ローカル配布は host の登録・install・有効化を変更しない。host cache や設定ファイルを直接編集しない。

| 項目 | 旧方式 | 移行後 |
|---|---|---|
| skill 原本 | `codex/skills` | `skills/internal` |
| 外部 skill | URL、hash、配布先を設定 | URL と `treeHash` lock のみ |
| 選択単位 | 一括コピー | group plugin または skill plugin を host ごとに install |
| host への反映 | skill directory を直接置換 | marketplace を登録して native install |
| 既存 skill directory | deploy が置換 | 自動削除・自動置換しない |

## 入力と生成物

- `skills/internal/<name>/SKILL.md` は internal skill の正本である。
- `external-skills.json` は外部 skill の URL と `treeHash` lock を管理し、destination は持たない。
- `skill-groups.json` は group と skill の membership を定義する。

`make build-plugins` は全 group plugin と、全 skill の個別 plugin `skill-<name>` を `out/marketplace` に生成する。group plugin の skill は `plugins/<group>/skills/<name>` に、個別 plugin は `plugins/skill-<name>/skills/<name>` に置かれる。各 plugin は Codex 用 `.codex-plugin/plugin.json` と Claude 用 `.claude-plugin/plugin.json` を持つ。

plugin の相対 resource 参照は必要な closure を同梱して解決する。特定 skill の個別 plugin だけを install しても、必須の相対参照で失敗しないことが契約である。group 間の依存は install 時に別 group / skill plugin を追加する。

生成先 root にのみ Codex 用 `.agents/plugins/marketplace.json` と Claude 用 `.claude-plugin/marketplace.json` が置かれる。Codex の local source は marketplace root 基準の `./plugins/<plugin>` を指す。GitHub 上のこのリポジトリ URL を marketplace として直接追加することはできず、必ず生成済み `out/marketplace` またはその配布コピーを追加する。

manifest version は plugin 内容の SHA 由来で `0.0.0-<12hex>` とする。内容が変われば version も変わるため、host 側の更新を検出できる。実際の更新手順とローカル marketplace の反映範囲は、利用する host の help で確認する。

## 処理の流れ

```text
skills/internal + external-skills.json + skill-groups.json
                              |
                              v
         build-plugins (全外部 skill を fetch / hash 検証)
                              |
                              v
                      out/marketplace
                              |
              +---------------+----------------+
              v                                v
  codex plugin marketplace add ...  claude plugin marketplace add ...
              |                                |
              v                                v
       必要な plugin を add              必要な plugin を install
```

外部 skill は marketplace build 時に全 11 件を取得し、hash を検証して取り込む。install 時の選択は host に移管するため、profile や selection による取得対象の絞り込みは行わない。現在の生成数は 20 group plugin と 68 個別 plugin（internal 57、external 11）の計 88 plugin である。

## 生成と任意のローカル配布

まず marketplace を生成する。

```bash
make build-plugins
```

固定のローカル path を使いたい場合は、生成物を `~/.local/share/agents-docs/marketplace` へ backup 後に置き換える。

```bash
make deploy-plugins-dry-run
make deploy-plugins
```

通常の `make deploy[-dry-run]` は plugin build と非-skill 資産の配布を行う。`TAG=skills` は plugin marketplace のみを対象にする。これらの操作は host の marketplace 登録や plugin install を実行しない。

### 失敗時の扱いと手動 rollback

build は staging で生成し、失敗時は既存の `out/marketplace` を保持する。配布は generic deploy の backup、replace、copy の順で行う。copy が失敗しても自動 restore はしないため、配布済み marketplace を戻す必要がある場合は backup から手動で復元する。

backup は deploy config と同じ場所の `.deploy-backups/<time>/<absolute destination>` にあり、destination の絶対パス構造を再現する。たとえば `~/.local/share/agents-docs/marketplace` の backup は、repo root の `.deploy-backups/<time>/Users/<user>/.local/share/agents-docs/marketplace` を探す。復元前に対象、backup 時刻、現在の配布先を確認し、backup の内容を配布先へ戻す。既存の配布先を置換する操作になるため、この文書は自動実行コマンドを提供しない。

## Codex で使う

生成先または配布済み marketplace を登録し、必要な plugin を追加する。Codex manifest は `policy.installation: "AVAILABLE"`、`authentication: "ON_INSTALL"`、`category: "Productivity"` を使う。

```bash
codex plugin marketplace add <repo>/out/marketplace
codex plugin add delivery@agents-docs
codex plugin add backend@agents-docs
```

確認には `codex plugin list --marketplace agents-docs` または `codex plugin list --available --json` を使う。ローカル path の marketplace は `upgrade` の対象ではない。再生成後の反映方法は利用中の CLI の help で確認し、Codex を再起動して反映を確認する。

## Claude で使う

Claude Code では marketplace と install scope を明示する。

```bash
claude plugin marketplace add --scope user <repo>/out/marketplace
claude plugin install -s user delivery@agents-docs
claude plugin install -s user backend@agents-docs
```

scope は `user`、`project`、`local` から利用場所に合わせて選ぶ。更新は `claude plugin marketplace update agents-docs`、`claude plugin update -s <scope> <plugin>` で行い、restart が必要である。有効化は `claude plugin enable -s <scope> <plugin>` を使う。無効化は利用中の host help を確認して操作する。

## 選択例と依存の注意

`delivery` は PR 作成、説明、CI 修正に、`backend` は Go または Ruby/Rails の実装に向く。たとえば Codex で Go の実装から PR 作成まで行うなら、`backend@agents-docs`、`code-foundations@agents-docs`、`delivery@agents-docs` を追加する。Claude でも同じ plugin 名を install できる。

group をまたぐ行動上の依存は host が自動 install しない。たとえば `delivery` の `cmd-create-pr` を使って PR 作成まで行うには、Git 操作の `git@agents-docs`、レビューの `review-roles@agents-docs` と `code-review@agents-docs` も追加する。フロントエンドまたはバックエンドの実装では、共通原則の `code-foundations@agents-docs` も追加する。必要な group / skill plugin は明示する。

group plugin と、その group に含まれる `skill-<name>` を重複 install しない。同じ skill が二重に読み込まれるおそれがあるためである。異なる group にまたがる依存がある場合は、依存先の group plugin または依存先 skill plugin も install する。個別 plugin の必須相対 resource は closure に含まれるが、別 plugin に属する行動上の依存までは自動 install しない。

## 旧 skill 配布からの移行と rollback

この変更は `~/.codex/skills` と `~/.claude/skills` を自動で削除・置換しない。Codex の `.system` や、利用者自身・他の配布元が管理する skill を保護するためである。

plugin の登録と必要な install を確認した後に限り、旧方式でこのリポジトリが管理していた skill を手動で退避する。退避前に対象を確認し、`.system` と出所が不明な skill は対象外にする。問題があれば退避したディレクトリを元の場所へ戻し、追加した plugin を削除または無効化して rollback する。

## 参考

- [OpenAI: Building plugins](https://developers.openai.com/plugins/build/plugins)
- [Claude Code: Plugin reference](https://code.claude.com/docs/en/plugins-reference)
- [Claude Code: Plugin marketplaces](https://code.claude.com/docs/en/plugin-marketplaces)
