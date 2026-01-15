## generator

docs/ja 配下のドキュメントを Codex 向けの出力に変換するコマンドです。

### 使い方

```
./bin/generator -input=./docs/ja -output=./out/.codex/
```

### オプション

- `-input` : 変換対象のルートディレクトリ（必須）
- `-output` : 出力先ディレクトリ（例: `./out/.codex/`）

### 出力内容

- `Agents.md` と `agents/` はそのまま出力先にコピーします
- `skills/` 配下の各 `.md` は `SKILL.md` として出力します
