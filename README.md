
## Agents-docs

Codex / Claude などの AI エージェントに配布するルール、プロンプト、skill、agent 定義を管理するレポジトリです。

このリポジトリでは、主に次の項目を管理しています。

- Codex / Claude 向けの共通指示
- agent の役割定義とメタデータ
- skill の本体と UI 表示用メタデータ
- 再利用するプロンプトとルール
- ローカル環境へ配布するための設定と deploy ツール

## skill 一覧

- [docs/skill-library.md](docs/skill-library.md)
- [docs/skill-dependency-map.md](docs/skill-dependency-map.md)

## コマンド

- `make build-deploy` で deploy 用のバイナリを生成します。
- `make deploy-dry-run` で `deploy.json` と `external-skills.json` に基づくコピー予定を確認します。
- `make deploy` で `deploy.json` と `external-skills.json` に基づいて設定とスキルを配置します。

`make deploy` は Codex 用の補助コマンドも `~/.codex/bin` に配置します。
`safe-git-push` は agent が prompt なしで使うための安全な `git push` wrapper です。
`safe-gh-edit` は自分が作成した PR / Issue だけを prompt なしで編集するための `gh pr edit` / `gh issue edit` wrapper です。
`safe-local-curl` は localhost / loopback / private address 宛ての確認だけを prompt なしで行うための安全な `curl` wrapper です。

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

- deploy コマンドの詳細: [scripts/deploy/README.md](scripts/deploy/README.md)
- レポジトリ設計: [docs/design.md](docs/design.md)

## ライセンス

このリポジトリの内容は Apache License 2.0 のもとで利用できます。詳細は [LICENSE](LICENSE) を参照してください。

改変した内容を公開・配布する場合は、可能であればブログ、README、ドキュメントなどで、このリポジトリに多少なりとも触れてもらえるとうれしいです。このお願いはライセンス条件ではなく、任意のクレジットとしての希望です。
