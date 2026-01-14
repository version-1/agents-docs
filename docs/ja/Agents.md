
agents/ 配下のドキュメントを読んでコードの設計・実装を行ってください。

特に指示がない限りは、日本語での受け答えをお願いします。

## ディレクトリ構成

```
.
├── agents
│   ├── roles              エージェントの役割(人格)別ドキュメント
│   │   ├── doc.md
│   │   ├── implementer.md
│   │   ├── planner.md
│   │   ├── role.md
│   │   └── verifier.md
├── skills
│   ├── architectures アーキテクチャ別のスキル
│   ├── languages     プログラミング言語別のスキル
│   ├── practices     コードレビュー、テスト、デプロイなどのガイドライン
│   ├── roles         エージェントの役割(人格)別ドキュメント
│   └── skills.md     スキル作成ガイドライン
└── Agents.md

```


## 役割別ドキュメント

相対パスは codex なら ~/.codex` ディレクトリを基準としています。

- `agents/roles/role.md` : エージェントの役割全般に関するガイドライン
- `agents/roles/planner.md` : Planner エージェントの役割に関するガイドライン
- `agents/roles/implementer.md` : Implementer エージェントの役割に関するガイドライン
- `agents/roles/verifier.md` : Verifier エージェントの役割に関するガイドライン
- `agents/roles/doc.md` : Doc エージェントの役割に関するガイドライン

## 言語別ドキュメント

- `agents/skills/languages/<language>.md` : 各プログラミング言語に関するスキルガイドライン
