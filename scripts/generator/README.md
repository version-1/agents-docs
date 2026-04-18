## generator

`docs/ja` 配下のドキュメントを Codex / Claude 向けの出力に変換するコマンドです。

### 使い方

```bash
# Codex 向け（ネスト構造を維持）
./bin/generator -input=./docs/ja -output=./out/.codex/ -mode=codex

# Claude 向け（フラット構造）
./bin/generator -input=./docs/ja -output=./out/.claude/ -mode=claude
```

`make gen-docs` で両方を一括生成できます。

### オプション

| フラグ | 必須 | 説明 |
|---|---|---|
| `-input` | Yes | 変換対象のルートディレクトリ（例: `./docs/ja`） |
| `-output` | Yes | 出力先ディレクトリ（例: `./out/.codex/`） |
| `-mode` | No | 出力モード: `codex`（デフォルト）または `claude` |

### 処理内容

1. 出力先ディレクトリの中身をクリアする
2. `Agents.md` と `agents/` をそのまま出力先にコピーする
3. `skills/` 配下の各 `.md` を `SKILL.md` として出力する

スキルの出力方式は `-mode` で切り替わります。

#### `-mode=codex`（デフォルト）

ソースのディレクトリ階層を維持して出力します。Codex のようにネストされたスキル構造をサポートするツール向けです。

```
docs/ja/skills/languages/go/go.md
  → out/.codex/skills/languages/go/go/SKILL.md
```

#### `-mode=claude`

カテゴリ階層を除去し、各スキルファイルの YAML フロントマターにある `name` フィールドをディレクトリ名として使用します。Claude のようにスキルをフラットに配置する必要があるツール向けです。

```
docs/ja/skills/languages/go/go.md   （name: language-go）
  → out/.claude/skills/language-go/SKILL.md
```

フロントマターに `name` フィールドがない場合はエラーになります。

### 除外ルール

以下のファイルはスキル生成の対象外です。

- `skills/skills.md` — スキル作成ガイドライン
- `Agents.md` — エージェント定義（コピー処理で別途扱う）

### アーキテクチャ

```
cmd/main.go           エントリポイント。CLI フラグを解析し、mode に応じた
                      SkillGenerator を選択して Generator に注入する。

internal/
  app/                オーケストレーション層
    generator.go        Generator 本体。処理の流れを制御する。
                        ① 入出力ディレクトリの検証
                        ② 出力先クリア
                        ③ Agents.md / agents/ のコピー
                        ④ SkillGenerator による skills 生成

  domain/             ドメイン層（ビジネスルール）
    skill.go            SkillGenerator インターフェースと 2 つの実装:
                        - pathRespectSkillGenerator: 階層を維持（Codex 用）
                        - flatSkillGenerator: フラットに展開（Claude 用）
    paths.go            出力パスの変換ルール
    frontmatter.go      YAML フロントマターの解析

  fsadapter/          ファイルシステム操作層
    filesystem.go       FileSystem インターフェース定義
    fsutil.go           汎用ファイル操作（DirExists, CopyDir, CopyFile 等）
    walk.go             WalkResources: skills ディレクトリの走査と書き出し
    filter.go           走査時のフィルタ（IsMarkdown, IsExcluded）

  infra/              インフラ層（外部依存の具象実装）
    osfs.go             OSFS: FileSystem インターフェースの OS 実装
```

#### 依存の方向

```
cmd → app → domain → fsadapter ← infra
```

- `app` は `domain.SkillGenerator` と `fsadapter.FileSystem` に依存する
- `domain` は `fsadapter` のインターフェースと走査関数を使ってスキルを生成する
- `infra` は `fsadapter.FileSystem` を実装する（`cmd` から注入）
- 各層は一方向にのみ依存し、循環しない
