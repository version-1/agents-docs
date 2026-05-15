# Skill Library

このドキュメントは、`codex/skills` で管理している skill の用途を素早く把握するための一覧です。

各 skill の実体は `codex/skills/internal/<skill-name>/SKILL.md` にあります。外部 skill は `external-skills.json` の設定に従って deploy 時に取得されます。

## Internal Skills

### Command

| Skill | 概要 | 使う場面 |
|---|---|---|
| `cmd-agent-status` | 現在のエージェント作業状況を短く報告する。 | ユーザーが「今どうなっている」「ステータス」「進捗」などを聞いたとき。 |
| `cmd-batch` | 広範囲の調査、編集、検証を複数エージェントや並列作業へ分割する。 | 大規模変更、横断調査、複数担当への分担、統合手順の整理が必要なとき。 |
| `cmd-dispatch-agent` | 指定された agent を起動し、結果を待たずにタスクを投げる。 | ユーザーが agent 起動、worker への委譲、投げっぱなし実行を明示したとき。 |
| `cmd-rmbranch` | `main` と `develop` を残し、不要なローカルブランチを安全に削除する。 | ローカルブランチ整理を依頼されたとき。未マージブランチは確認してから扱う。 |

### Code Work

| Skill | 概要 | 使う場面 |
|---|---|---|
| `coding` | 既存コードベースで機能追加、修正、バグ修正、テスト追加を実装するための基本方針。 | 実装作業全般。既存構造、命名、責務境界に合わせて差分を作るとき。 |
| `component-design` | React / Next.js / TypeScript UI の実装前に、画面構成とコンポーネント境界を設計する。 | 新規画面、大きな JSX 分割、既存 UI へのまとまった機能追加を行う前。 |
| `code-simplifier` | 挙動を変えずにコードを簡略化、リファクタリングする。 | レビュー指摘、diff、指定ファイルをもとに可読性や保守性を改善するとき。 |
| `code-test` | テスト設計、回帰テスト追加、テストコードレビュー、テスト戦略を整理する。 | 正常系、異常系、境界値、flake、モック方針を検討するとき。 |

### Review And Analysis

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-review` | 実装済み差分、PR、コミット、指定ファイルをリスク中心にレビューする。 | 仕様違反、回帰、公開契約破壊、セキュリティ、データ整合性を確認するとき。 |
| `code-next-developer-review` | 次に開発する人が困らないかという観点で差分をレビューする。 | 命名、責務境界、型、テストの読みやすさ、前提知識の残し方を確認するとき。 |
| `code-tracer` | 特定シンボルや関数の callers / callees / both を根拠付きで追跡する。 | 呼び出し経路、影響範囲、依存関係を Markdown や Mermaid で可視化したいとき。 |
| `code-typo` | コード、diff、ファイルパス、PR 文面などの typo を指摘する。 | 設計レビューではなく、スペルミスや表記ゆれだけを確認したいとき。 |

### Documentation

| Skill | 概要 | 使う場面 |
|---|---|---|
| `documenting` | README、ADR、Runbook、API Docs、開発者向け文書を作成、更新する。 | 変更理由、影響範囲、互換性、運用注意、戻し方を未来の開発者へ残すとき。 |
| `pr-description` | GitHub Pull Request の description を作成、更新する。 | PR 番号、diff、コミット履歴から Summary、変更点、動作確認を整理するとき。 |

### Role Skills

| Skill | 概要 | 使う場面 |
|---|---|---|
| `role-advisor` | 技術判断、設計相談、ベストプラクティスをテックリード視点で助言する。 | 実装を伴わず、アーキテクチャ、DB、DevOps、テスト、保守性を相談するとき。 |
| `role-doc` | Doc として、変更理由や設計判断など未来の理解コストを下げる情報を文書化する。 | README や手順書の形式より、何を言語化すべきかの判断が重要なとき。 |
| `role-implementer` | Implementer として、計画済みのゴールと完了条件に沿って実装する。 | Planner が整理した仕様に従い、小さく安全な差分を作るとき。 |
| `role-planner` | Planner として、ゴール、スコープ、タスク、完了条件、リスクを整理する。 | 実装前に方針を固め、Implementer と Reviewer が迷わない計画を作るとき。 |
| `role-reviewer` | Reviewer として、実装済み差分が仕様と品質要件を満たすか検証する。 | スコープ、公開契約、境界条件、責務違反、検証不足を確認するとき。 |
| `role-scouter` | Scouter として、コード、ドキュメント、設定、履歴を調査し判断材料を集める。 | 実装やレビュー判断の前に、不確実性を下げるための根拠収集をするとき。 |

### Agent Operations

| Skill | 概要 | 使う場面 |
|---|---|---|
| `agent-retro` | 直近のエージェント作業を振り返り、再利用可能なルールへ反映する。 | 遠回りした知見を skill、AGENTS.md、CLAUDE.md などに残したいとき。 |

## External Skills

外部 skill は `external-skills.json` で取得元と配布先を管理します。

| Skill | 取得元 | 概要 |
|---|---|---|
| `skill-creator` | `anthropics/skills` | 新しい skill の作成、既存 skill の更新、メタデータや検証手順の整備に使う。 |
| `grill-me` | `mattpocock/skills` | 計画や設計を厳しく質問し、曖昧さや判断漏れを潰すために使う。 |
| `empirical-prompt-tuning` | `mizchi/skills` | skill やプロンプトを実験的に改善し、評価と反復で性能を詰めるために使う。 |

## 運用メモ

- 新しい internal skill を追加したら、この一覧にも追記する。
- `codex/skills/internal/<name>/agents/openai.yaml` がある場合は、UI 表示用の説明も `SKILL.md` と矛盾しないように更新する。
- deploy 前は `make deploy-dry-run` で配布対象に含まれることを確認する。
