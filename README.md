
## Agents-docs

AI にコードを書かせる際のルールやコンテキストをまとめるレポジトリです。

## コマンド

- `make build-skillmaker` コマンドで skills 生成用のバイナリを生成します。
- `make gen-docs` で環境で使用するドキュメントを生成します。
- `make deploy-codex-docs` で生成済みのドキュメントを `~/.codex` に反映します。

## デプロイ手順

1. `make gen-docs` を実行して `docs/ja` から生成物を作成します。
2. `make deploy-codex-docs` を実行して `~/.codex` に反映します。
