

docs/ja/ 配下のドキュメントを読んでコードの設計・実装を行ってください。

特に指示がない限りは、日本語での受け答えをお願いします。

## ディレクトリ構成

```
.
├── Agents.md
├── docs
│   └── ja
│       └── agents
│           ├── roles              エージェントの役割(人格)別ドキュメント
│           │   ├── doc.md
│           │   ├── implementer.md
│           │   ├── planner.md
│           │   ├── role.md
│           │   └── verifier.md
│           └── skills
│               ├── architectures  アーキテクチャ別のスキル
│               ├── languages      プログラミング言語別のスキル
│               ├── practice       コードレビュー、テスト、デプロイなどのガイドライン
│               └── skills.md      スキル作成ガイドライン
└── README.md                      プロジェクト説明
```


## 役割別ドキュメント

- `roles/role.md` : エージェントの役割全般に関するガイドライン
- `roles/planner.md` : Planner エージェントの役割に関するガイドライン
- `roles/implementer.md` : Implementer エージェントの役割に関するガイドライン
- `roles/verifier.md` : Verifier エージェントの役割に関するガイドライン
- `roles/doc.md` : Doc エージェントの役割に関するガイドライン
