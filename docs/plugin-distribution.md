# Plugin 配布

このリポジトリは、全 skill を Codex と Claude の plugin marketplace として生成する。利用環境ごとの選択は marketplace の生成時ではなく、各 host の native plugin install で行う。これにより、一つの生成物から Codex と Claude が独立して必要な group / skill plugin を選べる。

## 調査記録と移行設計

従来は `codex/skills` を `~/.codex/skills`、同じ内容を `~/.claude/skills` へ直接コピーしていた。この方式では、環境ごとに選ぶ単位がなく、host が管理する plugin の状態とも分離していた。外部 skill の配布先も lock に混在していた。

移行後は、共通原本から全 plugin を一度だけ生成し、host ごとに native install を実行する。marketplace の生成・ローカル配布は host の登録・install・有効化を変更しない。host cache や設定ファイルを直接編集しない。

| 項目 | 旧方式 | 移行後 |
|---|---|---|
| group / 個別 skill 原本 | `codex/skills` | `skills/internal` |
| 手管理 plugin 原本 | なし | `plugins/<name>` |
| 外部 skill | URL、hash、配布先を設定 | URL と `treeHash` lock のみ |
| 選択単位 | 一括コピー | group plugin または skill plugin を host ごとに install |
| host への反映 | skill directory を直接置換 | marketplace を登録して native install |
| 既存 skill directory | deploy が置換 | 自動削除・自動置換しない |

## 入力と生成物

- `skills/internal/<name>/SKILL.md` は group と個別 skill plugin に生成する internal skill の正本である。
- `plugins/<name>/` は group に属さない手管理 plugin の正本である。`.codex-plugin/plugin.json`、`.claude-plugin/plugin.json`、`skills/<skill-name>/SKILL.md` を持たせる。build は manifest、skill frontmatter、symlink を検証し、source 内容から生成した version を出力 manifest に設定する。source manifest の `version` は hash から除外され、build 時に置き換わる。
- `external-skills.json` は外部 skill の URL と `treeHash` lock を管理し、destination は持たない。
- `skill-groups.json` は group と skill の membership を定義する。`includes` には他の group 名を指定でき、参照先の skill を合成した group plugin を作る。`includes` は host に別 plugin の install を指示するものではなく、合成 plugin 内へ skill を同梱する指定である。

既存の `groups` に追加する `core` 定義の抜粋です。参照先 group の定義は省略しています。

```json
{
  "groups": {
    "core": {
      "includes": [
        "git",
        "delivery",
        "documentation"
      ],
      "description": "基本作業をまとめて進める"
    }
  }
}
```

`make build-plugins` は全 group plugin と、`skills/internal/` と外部 skill の個別 plugin `skill-<name>` を `out/marketplace` に生成し、repo の `plugins/` にある手管理 plugin をそのまま取り込む。group plugin の skill は `plugins/<group>/skills/<name>` に、個別 plugin は `plugins/skill-<name>/skills/<name>` に置かれる。手管理 plugin は group と `skill-*` plugin の名前空間と衝突できない。各 plugin は Codex 用 `.codex-plugin/plugin.json` と Claude 用 `.claude-plugin/plugin.json` を持つ。

plugin の相対 resource 参照は必要な closure を同梱して解決する。特定 skill の個別 plugin だけを install しても、必須の相対参照で失敗しないことが契約である。group 間の依存は install 時に別 group / skill plugin を追加する。

生成先 root にのみ Codex 用 `.agents/plugins/marketplace.json` と Claude 用 `.claude-plugin/marketplace.json` が置かれる。Codex の local source は marketplace root 基準の `./plugins/<plugin>` を指す。GitHub 上のこのリポジトリ URL を marketplace として直接追加することはできず、必ず生成済み `out/marketplace` またはその配布コピーを追加する。

manifest version は plugin 内容の SHA 由来で `0.0.0-<12hex>` とする。内容が変われば version も変わるため、host 側の更新を検出できる。実際の更新手順とローカル marketplace の反映範囲は、利用する host の help で確認する。

## 処理の流れ

```text
skills/internal + external-skills.json + skill-groups.json + plugins/
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

外部 skill は marketplace build 時に全 14 件を取得し、hash を検証して取り込む。install 時の選択は host に移管するため、profile や selection による取得対象の絞り込みは行わない。生成される group と個別 plugin の一覧は `make build-plugins` の出力で確認する。

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
codex plugin add core@agents-docs
```

確認には `codex plugin list --marketplace agents-docs` または `codex plugin list --available --json` を使う。ローカル path の marketplace は `upgrade` の対象ではない。再生成後の反映方法は利用中の CLI の help で確認し、Codex を再起動して反映を確認する。

## Claude で使う

Claude Code では marketplace と install scope を明示する。

```bash
claude plugin marketplace add --scope user <repo>/out/marketplace
claude plugin install -s user core@agents-docs
```

scope は `user`、`project`、`local` から利用場所に合わせて選ぶ。更新は `claude plugin marketplace update agents-docs`、`claude plugin update -s <scope> <plugin>` で行い、restart が必要である。有効化は `claude plugin enable -s <scope> <plugin>` を使う。無効化は利用中の host help を確認して操作する。

## 選択例と依存の注意

`roles` は `plugins/roles` で管理する、finder 以外の `role-*` と対応する playbook をまとめる plugin である。`finder` は `plugins/finder` で管理する、課題発見用の `role-*-finder` と対応する playbook をまとめる plugin である。`core` は `git`、`delivery`、`documentation` をまとめて利用したい場合に向く。role を利用する場合は `roles@agents-docs` または `finder@agents-docs` を install する。Claude でも同じ plugin 名を install できる。`backend` や `code-foundations` など、`core` に含まれない専門 group は必要に応じて追加する。

group をまたぐ行動上の依存は host が自動 install しない。`includes` で skill を同梱することも、参照先 group plugin の host install を自動化するものではない。`delivery` を単独で install して `cmd-create-pr` を使う場合は、Git 操作の `git@agents-docs` とレビューの `roles@agents-docs`、`code-review@agents-docs` も追加する。`core` を install する場合も role を使うには `roles@agents-docs` を追加する。フロントエンドまたはバックエンドの実装では、共通原則の `code-foundations@agents-docs` も追加する。

従来の `planning-advice`、`implementation`、`review-roles`、`operations` plugin を利用している場合は、それらを `roles@agents-docs` に置き換える。旧 `documentation` plugin で Documenter agent への委譲を利用している場合も `roles@agents-docs` を追加する。`core@agents-docs` は `roles` を同梱しないため、role を使う場合は追加する。finder の role は `finder@agents-docs` に残る。

group plugin と、その group に含まれる `skill-<name>` を重複 install しない。同じ skill が二重に読み込まれるおそれがあるためである。`core` とその構成 group も重複 install しない。異なる group にまたがる依存がある場合は、依存先の group plugin または依存先 skill plugin も install する。個別 plugin の必須相対 resource は closure に含まれるが、別 plugin に属する行動上の依存までは自動 install しない。

### thinking から pocock-skills への移行

`thinking` group は削除し、既存の `grilling` を `pocock-skills` へ移した。`pocock-skills` は `grilling`、`grill-with-docs`、`to-spec`、`to-tickets` の4件を同梱する。marketplace を再生成・更新後、各 host で次の plugin を追加し、旧 `thinking@agents-docs` を削除または無効化して重複読み込みを避ける。host の設定や cache は deploy で自動変更されない。

```bash
codex plugin add pocock-skills@agents-docs
claude plugin install -s user pocock-skills@agents-docs
```

`skill-grilling@agents-docs` などの個別 plugin を利用している場合は、そのまま利用できる。group に切り替える場合は同じ skill の個別 plugin も削除または無効化する。

上流の `grill-with-docs` は `domain-modeling` を参照するが、この group には同梱しない。利用時は上流から別途用意する。`to-spec` と `to-tickets` はトラッカーとラベルの設定を前提とし、未設定時には上流の `setup-matt-pocock-skills` を案内する。これらの補助 skill は今回の登録対象外であり、host が自動 install するものではない。

補助 skill の取得元は、登録済み4件と同じ固定コミットの [domain-modeling](https://github.com/mattpocock/skills/tree/697d4ce9742da558fd1ba6697c8e9775e2e302dd/skills/engineering/domain-modeling) と [setup-matt-pocock-skills](https://github.com/mattpocock/skills/tree/697d4ce9742da558fd1ba6697c8e9775e2e302dd/skills/engineering/setup-matt-pocock-skills) で確認できる。[上流の Quickstart](https://github.com/mattpocock/skills/blob/697d4ce9742da558fd1ba6697c8e9775e2e302dd/README.md#quickstart-30-second-setup) の installer では、この2件と利用 host を選択し、導入後に `/setup-matt-pocock-skills` を実行する。installer は上流の最新版を取得するため、固定版が必要ならリンク先の skill ディレクトリ全体を取得して host の skill 配置先へ導入する。`pocock-skills` の4件は installer で重複選択しない。

## 旧 skill 配布からの移行と rollback

この変更は `~/.codex/skills` と `~/.claude/skills` を自動で削除・置換しない。Codex の `.system` や、利用者自身・他の配布元が管理する skill を保護するためである。

plugin の登録と必要な install を確認した後に限り、旧方式でこのリポジトリが管理していた skill を手動で退避する。退避前に対象を確認し、`.system` と出所が不明な skill は対象外にする。問題があれば退避したディレクトリを元の場所へ戻し、追加した plugin を削除または無効化して rollback する。

## 参考

- [OpenAI: Building plugins](https://developers.openai.com/plugins/build/plugins)
- [Claude Code: Plugin reference](https://code.claude.com/docs/en/plugins-reference)
- [Claude Code: Plugin marketplaces](https://code.claude.com/docs/en/plugin-marketplaces)
