
これは agents-docs レポジトリ用の Agents.md です。
汎用的な docs/ja/Agents.md が実際に使われるドキュメントになるので注意してください。

## ディレクトリ構成

```
.
├── Agents.md                      このファイル（リポジトリの案内）
├── bin                             生成ツールのバイナリ
├── codex                           Codex CLI 用の設定
├── docs                            ドキュメントのソース（実際に使うのは docs/ja 配下）
│   └── ja                          
│       ├── Agents.md               実運用向けの Agents.md
│       ├── agents                  エージェントの役割・共通スキル
│       │   ├── roles                 役割（人格）別ドキュメント
│       │   │   ├── doc.md              Doc ロール
│       │   │   ├── implementer.md      Implementer ロール
│       │   │   ├── planner.md          Planner ロール
│       │   │   ├── role.md             役割共通のガイドライン
│       │   │   └── verifier.md         Verifier ロール
│       │   └── skills                エージェント共通のスキルガイド
│       │       └── skills.md          スキル作成ガイドライン
│       └── skills                  スキル体系（各カテゴリ配下に定義）
│           ├── architectures        アーキテクチャ別スキル
│           │   └── <name>.md          各アーキテクチャのスキル
│           ├── languages            言語別スキル
│           │   └── <language>.md      各言語のスキル
│           ├── practices            実践ガイド（レビュー/テスト等）
│           │   └── <name>.md          各プラクティスのスキル
│           ├── roles                役割別スキル
│           │   └── <role>.md          各ロールのスキル
│           └── skills.md            スキル作成ガイドライン
├── out                             docs から生成された出力（直接編集しない）
└── scripts                         ドキュメント生成用スクリプト
```

## コマンド

- `make build-skillmaker` コマンドで skills 生成用のバイナリを生成します。
- `make deploy` で Codex / Claude 用のドキュメント・skill・agent を配布します。

## テスト実行時の補足

テスト実行でキャッシュ権限エラーが出る場合は、以下の環境変数を指定して実行してください。

```
GOMODCACHE=/tmp/gomodcache GOCACHE=/tmp/gocache GOTOOLCHAIN=local go test ./...
```
