
## Agents-docs

AI にコードを書かせる際のルールやコンテキストをまとめるレポジトリです。

## コマンド

- `make build-deploy` で deploy 用のバイナリを生成します。
- `make deploy-dry-run` で `deploy.json` と `external-skills.json` に基づくコピー予定を確認します。
- `make deploy` で `deploy.json` と `external-skills.json` に基づいて設定とスキルを配置します。

## デプロイ手順

1. `make deploy-dry-run` を実行して、コピー予定のファイルとディレクトリを確認します。
2. `make deploy` を実行して、Codex / Claude 用の設定とスキルを反映します。

通常はこの2コマンドを続けて実行します。

```bash
make deploy-dry-run
make deploy
```

## 設定ファイル

- `deploy.json`: 内部の skill / agent / 設定ファイルの配布先を定義します。
- `external-skills.json`: GitHub など外部 URL から取得する skill を定義します。

## 関連ドキュメント

- deploy コマンドの詳細: `scripts/deploy/README.md`
- レポジトリ設計: `docs/design.md`
