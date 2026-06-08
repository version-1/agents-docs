# Skill Library

このドキュメントは、`codex/skills` で管理している skill の用途を素早く把握するための一覧です。

各 skill の実体は `codex/skills/internal/<skill-name>/SKILL.md` にあります。外部 skill は `external-skills.json` の設定に従って deploy 時に取得されます。

分類は [docs/guide/skill-category.md](guide/skill-category.md) の考え方に合わせています。複数の性質を持つ skill は、主な用途で分類しています。

skill 間の明示的な併用・優先関係は [docs/skill-dependency-map.md](skill-dependency-map.md) にまとめています。

## Internal Skills

### A. 行動定義型

「何をするか」を定義する skill です。ワークフロー型、ロール型、出力フォーマット型を中心に分類しています。

#### ワークフロー型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `agent-retro` | 直近のエージェント作業を振り返り、再利用可能なルールへ反映する。 | 遠回りした知見を skill、AGENTS.md、CLAUDE.md などに残したいとき。 |
| `cmd-agent-status` | 現在のエージェント作業状況を短く報告する。 | ユーザーが「今どうなっている」「ステータス」「進捗」などを聞いたとき。 |
| `cmd-batch` | 広範囲の調査、編集、検証を複数エージェントや並列作業へ分割する。 | 大規模変更、横断調査、複数担当への分担、統合手順の整理が必要なとき。 |
| `cmd-commit` | 変更内容を確認し、適切な粒度とメッセージで安全に Git commit を作成する。 | commit 作成を依頼されたとき。未確認の変更や unrelated changes を含めず、stage 対象と message を整理する。 |
| `cmd-create-pr` | GitHub Pull Request を安全な手順で作成または更新する。 | PR 作成、PR 提出、pull request 作成を依頼されたとき。差分確認、検証、commit、role-reviewer による PR 前レビュー、High 指摘の自動対応、safe-git-push、PR description 作成、gh pr create / edit の順序を整理する。 |
| `cmd-dispatch-agent` | 指定された agent を起動し、結果を待たずにタスクを投げる。 | ユーザーが agent 起動、worker への委譲、投げっぱなし実行を明示したとき。明示がなくても、単純で不明瞭な点がない自己完結タスクを任せたいとき。 |
| `cmd-rmbranch` | `main` と `develop` を残し、不要なローカルブランチを安全に削除する。 | ローカルブランチ整理を依頼されたとき。未マージブランチは確認してから扱う。 |
| `cmd-start-branch` | 最新のデフォルトブランチから作業ブランチを作り、不要ブランチ整理を非同期に依頼する。 | 新しい作業を始める前に「ブランチ切って」「作業開始用ブランチを作って」などを依頼されたとき。ブランチ名を報告してタスク詳細を待つ。 |
| `beautify-commit` | ベースブランチまたは基準 commit との差分を、意味のある変更単位の commit へ安全に整理する。 | commit 分割、履歴整理、大きすぎる差分の再 commit、ベースブランチや commit hash を基準にした差分整理、interactive rebase や reset を伴う整理を依頼されたとき。 |
| `ci-fix` | GitHub Actions / CI の失敗を調査し、原因切り分けから修正、再検証まで進める。 | CI、GitHub Actions、checks、workflow、test / lint / build failure の修正を依頼されたとき。 |
| `code-simplifier` | 挙動を変えずにコードを簡略化、リファクタリングする。 | レビュー指摘、quality report、diff、指定ファイルをもとに可読性、保守性、テスト容易性を改善するとき。 |
| `code-test` | テスト設計、回帰テスト追加、テストコードレビュー、テスト戦略を整理する。 | 正常系、異常系、境界値、flake、モック方針を検討するとき。 |
| `code-tracer` | 特定シンボルや関数の callers / callees / both を根拠付きで追跡する。 | 呼び出し経路、影響範囲、依存関係を Markdown や Mermaid で可視化したいとき。 |
| `code-general` | 言語固有スキルが適用できない、または言語をまたぐ実装作業の共通原則。 | 既存構造、命名、責務境界に合わせて小さく安全な差分を作るとき。言語固有の判断では該当する言語別スキルを併用する。 |
| `component-design` | React / Next.js / TypeScript UI の実装前に、画面構成とコンポーネント境界を設計する。 | 新規画面、大きな JSX 分割、既存 UI へのまとまった機能追加を行う前。 |
| `documenting` | README、ADR、Runbook、API Docs、開発者向け文書を作成、更新する。 | 変更理由、影響範囲、互換性、運用注意、戻し方を未来の開発者へ残すとき。 |

#### ロール型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `role-advisor` | 技術判断、設計相談、ベストプラクティスをテックリード視点で助言する。 | 実装を伴わず、アーキテクチャ、DB、DevOps、テスト、保守性を相談するとき。 |
| `role-doc` | Doc として、変更理由や設計判断など未来の理解コストを下げる情報を文書化する。 | README や手順書の形式より、何を言語化すべきかの判断が重要なとき。 |
| `role-implementer` | Implementer として、計画済みのゴールと完了条件に沿って実装する。 | Planner が整理した仕様に従い、小さく安全な差分を作るとき。 |
| `role-planner` | Planner として、ゴール、スコープ、タスク、完了条件、リスクを整理する。 | 実装前に方針を固め、Implementer と Reviewer が迷わない計画を作るとき。 |
| `role-reviewer` | Reviewer として、実装済み差分が仕様と品質要件を満たすか検証する。 | スコープ、公開契約、境界条件、責務違反、検証不足を確認するとき。 |
| `role-scouter` | Scouter として、コード、ドキュメント、設定、履歴を調査し判断材料を集める。 | 実装やレビュー判断の前に、不確実性を下げるための根拠収集をするとき。 |

#### 出力フォーマット型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `format-pr-description` | GitHub Pull Request の description を定義済みフォーマットで作成、更新する。 | PR 番号、diff、コミット履歴から Summary、変更点、動作確認を整理するとき。 |
| `format-rich-html-diagram` | アーキテクチャ、コンポーネント構造、データフロー、状態管理、シーケンス、処理フローを単一 `index.html` でリッチに可視化する。 | 開発理解、設計整理、コードリーディング、オンボーディング用に、ブラウザで開ける図解資料を作りたいとき。 |

### B. 判断定義型

「どう考え、どう評価するか」を定義する skill です。判断基準・評価型、思想・スタイル型、スコアリング・査定型を中心に分類しています。

#### 判断基準・評価型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-naming` | 関数名、クラス名、型名、変数名などコード要素の命名候補や改善案を出す。 | 責務、抽象度、既存語彙、検索性、ユーザーの命名の好みに沿って名前を比較・改善したいとき。 |
| `code-next-developer-review` | 次に開発する人が困らないかという観点で差分をレビューする。 | 命名、責務境界、型、テストの読みやすさ、前提知識の残し方を確認するとき。 |
| `code-review` | 実装済み差分、PR、コミット、指定ファイルをリスク中心にレビューする。 | 仕様違反、回帰、公開契約破壊、セキュリティ、データ整合性を確認するとき。 |
| `code-typo` | コード、diff、ファイルパス、PR 文面などの typo を指摘する。 | 設計レビューではなく、スペルミスや表記ゆれだけを確認したいとき。 |

#### 思想・スタイル型スキル

`code-naming` は判断基準・評価型にも近いですが、ユーザーの命名の好みを反映する思想・スタイル型の性質を持ちます。主分類は判断基準・評価型に置いています。

#### スコアリング・査定型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-quality-review` | code quality report として、定量シグナル、品質スコア、scorecard、quality gate を整理する。 | complexity、coverage、重複率、lint などを補助証拠にしつつ、merge 前対応と follow-up の判断材料にしたいとき。 |

### C. 制約定義型

「何を避けるべきか」を定義する skill です。主にガードレール・制約型を分類します。

#### ガードレール・制約型スキル

現時点で、この類型を主用途にする internal skill はありません。

### D. 知識提供型

「何を知っておくべきか」を提供する skill です。主に知識圧縮型を分類します。

#### 知識圧縮型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-css` | CSS 実装向けのレイアウト、レスポンシブ、保守しやすいスタイル設計の参照。 | CSS 実装を書く前に、スタイル設計の判断基準を確認するとき。 |
| `code-go` | Go 実装と Go テスト向けの言語仕様、設計、テスト方針の参照。 | Go 実装や Go テストを書く前に、実装ガイドを確認するとき。 |
| `code-react` | React 実装向けのコンポーネント設計、状態管理、UI 実装の参照。 | React 実装を書く前に、UI 実装ガイドを確認するとき。 |
| `code-ruby` | Ruby / Rails 実装とテスト向けの設計、テスト方針の参照。 | Ruby / Rails 実装やテストを書く前に、実装ガイドを確認するとき。 |
| `code-ts` | TypeScript 実装向けの型設計、責務分離、既存コードへの合わせ方の参照。 | TypeScript 実装を書く前に、実装ガイドを確認するとき。 |

## External Skills

外部 skill は `external-skills.json` で取得元と配布先を管理します。

| Skill | 取得元 | 主な分類 | 概要 |
|---|---|---|---|
| `skill-creator` | `anthropics/skills` | 行動定義型 / ワークフロー型 | 新しい skill の作成、既存 skill の更新、メタデータや検証手順の整備に使う。 |
| `grill-me` | `mattpocock/skills` | 判断定義型 / 判断基準・評価型 | 計画や設計を厳しく質問し、曖昧さや判断漏れを潰すために使う。 |
| `empirical-prompt-tuning` | `mizchi/skills` | 判断定義型 / スコアリング・査定型 | skill やプロンプトを実験的に改善し、評価と反復で性能を詰めるために使う。 |

## 運用メモ

- 新しい internal skill を追加したら、この一覧にも追記する。
- 追加時は [docs/guide/skill-category.md](guide/skill-category.md) を参照し、主用途に最も近い分類へ配置する。
- 複数カテゴリにまたがる skill は、重複掲載せず主分類に置き、必要に応じて概要に複合要素を明記する。
- `code-general` や `code-review` のようにガードレールや思想を含む skill でも、一覧では主な行動や判断の用途を優先して配置する。
- `codex/skills/internal/<name>/agents/openai.yaml` がある場合は、UI 表示用の説明も `SKILL.md` と矛盾しないように更新する。
- 他 skill への明示的な依存や利用順序を追加したら、[docs/skill-dependency-map.md](skill-dependency-map.md) も更新する。
- deploy 前は `make deploy-dry-run` で配布対象に含まれることを確認する。
