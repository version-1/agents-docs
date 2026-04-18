
ツールに応じて下記のドキュメントを必要に応じて読んでコードの設計・実装を行ってください。

```
# codex の場合
~/.codex/agents/
# claude の場合
~/.claude/agents/
```

特に指示がない限りは、日本語での受け答えをお願いします。

## ディレクトリ構成

```
~/.codex
├── agents
│   ├── roles              エージェントの役割(人格)別ドキュメント
│   │   ├── doc.md
│   │   ├── implementer.md
│   │   ├── planner.md
│   │   ├── role.md
│   │   ├── scout.md
│   │   └── verifier.md
│   ├── codestyle          コードスタイルの指針
│   │   ├── codestyle.md
│   │   ├── languages
│   │   │   └── go.md
│   │   └── layers
│   │       ├── controller.md
│   │       ├── entity.md
│   │       ├── model.md
│   │       ├── repository.md
│   │       └── serialize.md
│   └── testing            テスト方針
│       └── testing.md
├── skills
│   ├── architectures アーキテクチャ別のスキル
│   │   └── <name>.md
│   ├── languages     プログラミング言語別のスキル
│   │   └── <language>.md
│   ├── practices     コードレビュー、テスト、デプロイなどのガイドライン
│   │   └── <name>.md
│   ├── roles         エージェントの役割(人格)別ドキュメント
│   │   ├── <role>.md
│   │   └── scout.md
│   └── skills.md     スキル作成ガイドライン
└── Agents.md

```


## 役割別ドキュメント

相対パスは codex なら `~/.codex` ディレクトリを基準としています。

- `agents/roles/role.md` : エージェントの役割全般に関するガイドライン
- `agents/roles/planner.md` : Planner エージェントの役割に関するガイドライン
- `agents/roles/implementer.md` : Implementer エージェントの役割に関するガイドライン
- `agents/roles/scout.md` : Scout エージェントの役割に関するガイドライン
- `agents/roles/verifier.md` : Verifier エージェントの役割に関するガイドライン
- `agents/roles/doc.md` : Doc エージェントの役割に関するガイドライン

## 言語別ドキュメント

- `skills/languages/<language>.md` : 各プログラミング言語に関するスキルガイドライン
