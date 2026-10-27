# Skill Library

このドキュメントは、`skills/internal` と手管理 plugin で管理している skill の用途を素早く把握するための一覧です。

現在の active internal skill の `SKILL.md` 本文は、必須契約と発火条件を維持したまま、frontmatter を含む UTF-8 で 5,000B 以下に整理済みです。

group と個別 plugin の skill 実体は `skills/internal/<skill-name>/SKILL.md` にあります。手管理 plugin の skill 実体は `plugins/<plugin>/skills/<skill-name>/SKILL.md` にあります。配布時は plugin 内の `skills/<skill-name>` へ配置されるため、skill 本文の相対 resource 参照も plugin 内で解決されます。外部 skill は `external-skills.json` の URL と hash lock に従い、marketplace build 時に取得・検証されます。

## Plugin groups と個別選択

group の正本は `skill-groups.json` です。build は group plugin と各 skill の個別 `skill-<name>` plugin を生成し、`plugins/` の手管理 plugin も検証して取り込みます。`includes` を持つ group は参照先 group の skill を合成して同梱します。相対 resource を参照する個別 plugin には必須の closure が同梱されます。具体的な生成・登録・移行手順は [plugin 配布](plugin-distribution.md) を参照してください。

`tasq-cli` と `tq-orchestrator` は group に含めず、各 `skill-tasq-cli` と `skill-tq-orchestrator` plugin として利用します。

| グループ | 用途 | 代表 skill / 数 |
|---|---|---|
| `core` | 基本的な agent 作業をまとめて利用 | `git`、`delivery`、`documentation` を合成 / 6 |
| `git` | コミットと grape による Git 操作 | `cmd-commit`、`grape-usage` / 2 |
| `delivery` | PR 作成と説明整備 | `cmd-create-pr`、`format-pr-description` / 2 |
| `agent-orchestration` | agent への分割・委譲 | `cmd-batch`、`cmd-dispatch-agent` / 2 |
| `code-foundations` | 共通実装原則、命名、コード追跡 | `code-general`、`code-naming`、`format-tracer` / 4 |
| `lego-programming` | 責務分離と部品境界の設計・実装 | `lego-programming` / 1 |
| `design` | UI構成と変更設計 | `component-design`、`modification-design` / 2 |
| `frontend` | CSS、React、TypeScript の実装 | `code-react`、`code-ts` / 4 |
| `backend` | Go、Ruby/Rails の実装 | `code-go`、`code-ruby` / 2 |
| `testing` | テスト戦略、仕様保証・変更耐性の評価、既存ケースの整理 | `code-test`、`format-test-cases` / 2 |
| `code-review` | 仕様、品質、保守性、テスト差分のレビュー | `code-review`、`code-refactor` / 5 |
| `documentation` | 文書と図解の作成 | `documenting`、`format-rich-html-diagram` / 2 |
| `authoring` | skill 作成、改善、作業知見の再利用 | `skill-creator`、`agent-retro` / 3 |
| `office` | PDF、プレゼン、スプレッドシート | `pdf`、`pptx`、`xlsx` / 3 |
| `thinking` | 計画・判断の問い直し | `grilling` / 1 |
| `technical-writing` | 日本語技術文書の執筆・推敲 | `japanese-tech-writing` / 1 |

## 手管理 plugin

| Plugin | 用途 | skill / 数 |
|---|---|---|
| `finder` | 分野別の課題発見 | `issue-finder-playbook`、5 分野の finder と各 playbook / 11 |
| `roles` | finder 以外の専門 agent への委譲 | `role-advisor`、`role-documenter`、`role-gardener`、`role-implementer`、`role-planner`、`role-refactor`、`role-reviewer`、`role-scouter`、`role-worker` と各 playbook / 18 |

分類は [docs/guide/skill-category.md](guide/skill-category.md) の考え方に合わせています。複数の性質を持つ skill は、主な用途で分類しています。

skill 間の明示的な併用・優先関係は [docs/skill-dependency-map.md](skill-dependency-map.md) にまとめています。

## Internal Skills

### A. 行動定義型

「何をするか」を定義する skill です。ワークフロー型、ロール型、出力フォーマット型を中心に分類しています。

#### ワークフロー型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `agent-retro` | 直近のエージェント作業を振り返り、再利用可能なルールへ反映する。 | 遠回りした知見を skill、AGENTS.md、CLAUDE.md などに残したいとき。 |
| `cmd-batch` | 広範囲の調査、編集、検証を複数エージェントや並列作業へ分割する。 | 大規模変更、横断調査、複数担当への分担、統合手順の整理が必要なとき。 |
| `cmd-commit` | 変更内容を確認し、適切な粒度とメッセージで安全に Git commit を作成する。 | commit 作成を依頼されたとき。未確認の変更や unrelated changes を含めず、stage 対象と message を整理する。 |
| `cmd-create-pr` | GitHub Pull Request を安全な手順で作成または更新する。 | PR 作成、PR 提出、pull request 作成を依頼されたとき。差分確認、検証、commit、role-reviewer による PR 前レビュー、High 指摘の自動対応、grape push、PR description 作成、gh pr create / edit の順序を整理する。 |
| `cmd-dispatch-agent` | 指定された agent を起動し、結果を待たずにタスクを投げる。 | ユーザーが agent 起動、worker への委譲、投げっぱなし実行を明示したとき。明示がなくても、単純で不明瞭な点がない自己完結タスクを任せたいとき。 |
| `code-refactor` | 挙動を変えずにコードを簡略化、リファクタリングする。 | レビュー指摘、quality report、diff、指定ファイルをもとに可読性、保守性、テスト容易性を改善するとき。 |
| `code-tracer` | 特定シンボルや関数の callers / callees / both を根拠付きで追跡する。 | 呼び出し経路、影響範囲、依存関係を Markdown や Mermaid で可視化したいとき。 |
| `code-general` | 言語固有スキルが適用できない、または言語をまたぐ実装作業の共通原則。 | 既存構造、命名、責務境界に合わせて小さく安全な差分を作るとき。言語固有の判断では該当する言語別スキルを併用する。 |
| `component-design` | React / Next.js / TypeScript UI の実装前に、画面構成とコンポーネント境界を設計する。 | 新規画面、大きな JSX 分割、既存 UI へのまとまった機能追加を行う前。 |
| `documenting` | 読み手と目的に合わせて文書を作成、更新する。 | 構成を決め、文書全体の冗長さ、重複、整合性、分割の必要性を見直すとき。 |
| `lego-programming` | 実装対象や差分を高凝集・疎結合な部品として境界づけ、組み合わせ、見直す。 | 実装、リファクタリング、コードレビュー、UI 分割、モジュール設計で責務分離や再利用単位を確認したいとき。 |
| `modification-design` | 実装前に変更後の責務、公開 interface、contract、依存構造、影響、リスクを自然言語で設計する。 | Plan や実装手順ではなく、変更後のシステム構造をレビュー可能な形で固めたいとき。 |

#### ロール型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `issue-finder-playbook` | Issue Finder agent の調査、安全制約、品質基準、Markdown 出力 contract を定義する。 | 各 finder playbook と組み合わせて、コードを変更せず課題を発見・出力するとき。 |
| `role-bug-finder` | `issue-finder` agent へ不具合発見を委譲する。 | コードベースの再現可能な不具合、仕様違反、境界条件の欠陥を探すとき。 |
| `role-bug-finder-playbook` | 不具合発見の探索観点と判定基準を定義する。 | Issue Finder agent が不具合を調査するとき。 |
| `role-maintenance-finder` | `issue-finder` agent へ保守性課題の発見を委譲する。 | 責務混在、複雑性、重複、変更・テスト困難性を探すとき。 |
| `role-maintenance-finder-playbook` | 保守性課題の探索観点と判定基準を定義する。 | Issue Finder agent が保守性を調査するとき。 |
| `role-feature-finder` | `issue-finder` agent へ機能提案の発見を委譲する。 | リポジトリ内の根拠から利用者の未解決課題と機能案を探すとき。 |
| `role-feature-finder-playbook` | 根拠のある機能提案の探索観点と判定基準を定義する。 | Issue Finder agent が機能提案を調査するとき。 |
| `role-vulnerability-finder` | `issue-finder` agent へ脆弱性発見を委譲する。 | コード、設定、依存関係、権限境界の悪用可能な問題を安全に探すとき。 |
| `role-vulnerability-finder-playbook` | 脆弱性発見の安全境界、探索観点、判定基準を定義する。 | Issue Finder agent が脆弱性を調査するとき。 |
| `role-documentation-finder` | `issue-finder` agent へ文書課題の発見を委譲する。 | 文書の不整合、配置、導線、progressive disclosure の問題を探すとき。 |
| `role-documentation-finder-playbook` | 文書課題の探索観点と判定基準を定義する。 | Issue Finder agent が文書構造と整合性を調査するとき。 |
| `role-advisor` | `advisor` agent へ技術助言を委譲する。 | 技術判断や設計相談が必要なとき。 |
| `role-advisor-playbook` | Advisor agent の助言方針を定義する。 | `advisor` agent が技術助言を行うとき。 |
| `role-documenter` | `documenter` agent へ技術文書の作成・更新を委譲する。 | 調査から構成、作成、文書全体のレビューまで一貫した文書化が必要なとき。 |
| `role-documenter-playbook` | Documenter agent の文書化手順と判断基準を定義する。 | `documenter` agent が技術文書を作成・更新するとき。 |
| `role-gardener` | `gardener` agent へ Git / GitHub 操作を委譲する。 | 内容の解釈を伴わない repository、branch、worktree、PR、Issue の操作が必要なとき。 |
| `role-gardener-playbook` | Gardener agent の Git / GitHub 操作方針と責務境界を定義する。 | `gardener` agent が Git / GitHub の機械操作を行うとき。 |
| `role-implementer` | `implementer` agent へ実装を委譲する。 | 計画に沿ったコード実装が必要なとき。 |
| `role-implementer-playbook` | Implementer agent の実装方針を定義する。 | `implementer` agent が実装を行うとき。 |
| `role-planner` | `planner` agent へ計画作成を委譲する。 | 実装前の計画整理が必要なとき。 |
| `role-planner-playbook` | Planner agent の計画方針を定義する。 | `planner` agent が計画を作成するとき。 |
| `role-refactor` | `refactor` agent へ構造改善の設計を委譲する。 | 挙動を維持したリファクタリングが必要なとき。 |
| `role-refactor-playbook` | Modification Design、Implementer、Reviewer の反復手順を定義する。 | `refactor` agent が構造改善を設計・統括するとき。 |
| `role-reviewer` | `reviewer` agent へレビューを委譲する。 | 実装済み変更のレビューが必要なとき。 |
| `role-reviewer-playbook` | Reviewer agent のレビュー方針を定義する。 | `reviewer` agent がレビューを行うとき。 |
| `role-scouter` | `scouter` agent へ調査を委譲する。 | コードや設定の調査が必要なとき。 |
| `role-scouter-playbook` | Scouter agent の調査方針を定義する。 | `scouter` agent が調査を行うとき。 |
| `role-worker` | `worker` agent へ機械的な作業を委譲する。 | コマンド実行、定型操作、検証が必要なとき。 |
| `role-worker-playbook` | Worker agent の機械的な作業方針を定義する。 | `worker` agent がコマンド作業や検証を行うとき。 |

#### 出力フォーマット型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `format-pr-description` | GitHub Pull Request の description を定義済みフォーマットで作成、更新する。 | PR 番号、diff、コミット履歴から Summary、変更点、動作確認を整理するとき。 |
| `format-tracer` | コードポインタに添える呼び出し経路を、リンク付きの番号付きリストへ整形する。 | レビュー、不具合調査、実装説明で、対象処理へ至る経路をテキストで伝えたいとき。 |
| `format-rich-html-diagram` | アーキテクチャ、コンポーネント構造、データフロー、状態管理、シーケンス、処理フローを単一 `index.html` でリッチに可視化する。 | 開発理解、設計整理、コードリーディング、オンボーディング用に、ブラウザで開ける図解資料を作りたいとき。 |
| `format-test-cases` | 既存テストを、条件、入力、操作、期待結果、根拠を含む Markdown のテストケース表に整理する。 | テストコードの実際の検証内容を把握・共有したいとき。 |

### B. 判断定義型

「どう考え、どう評価するか」を定義する skill です。判断基準・評価型、思想・スタイル型、スコアリング・査定型を中心に分類しています。

#### 判断基準・評価型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-naming` | 関数名、クラス名、型名、変数名などコード要素の命名候補や改善案を出す。 | 責務、抽象度、既存語彙、検索性、ユーザーの命名の好みに沿って名前を比較・改善したいとき。 |
| `code-next-developer-review` | 次に開発する人が困らないかという観点で差分をレビューする。 | 実装者の前提を鵜呑みにせず、敵対的検証を行い、命名、責務境界、型、テストの読みやすさ、前提知識の残し方を確認するとき。 |
| `code-review` | 実装済み差分、PR、コミット、指定ファイルをリスク中心にレビューし、テスト差分があればケース単位の反例・不足・削除された保証も確認する。 | 実装者の前提を鵜呑みにせず、敵対的検証を行い、仕様違反、回帰、公開契約破壊、セキュリティ、データ整合性を確認するとき。 |
| `code-test` | テスト戦略を選び、テストが仕様を保証し、正常な変更に耐えるかを評価する。 | テスト設計、回帰テスト追加、既存テストのレビューで、期待値、診断性、再現性、独立性、性能を判断するとき。 |

#### 思想・スタイル型スキル

`code-naming` は判断基準・評価型にも近いですが、ユーザーの命名の好みを反映する思想・スタイル型の性質を持ちます。主分類は判断基準・評価型に置いています。

#### スコアリング・査定型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-quality-review` | code quality report として、定量シグナル、品質スコア、scorecard、quality gate を整理する。 | 実装者の前提を鵜呑みにせず、敵対的検証を行い、complexity、coverage、重複率、lint などを補助証拠にしつつ、merge 前対応と follow-up の判断材料にしたいとき。 |
| `code-merge-risk` | PR のマージ後の悪影響を0〜100点で採点し、根拠付き Markdown と JSON を出力する。 | 採点基準と評価手順を揃え、情報不足や必須条件を分離して将来の自動マージ候補を判定したいとき。マージは実行しない。 |

### C. 制約定義型

「何を避けるべきか」を定義する skill です。主にガードレール・制約型を分類します。

#### ガードレール・制約型スキル

現時点で、この類型を主用途にする internal skill はありません。

### D. 知識提供型

「何を知っておくべきか」を提供する skill です。主に知識圧縮型を分類します。

#### 知識圧縮型スキル

| Skill | 概要 | 使う場面 |
|---|---|---|
| `code-css` | CSS / `.css` / Sass / SCSS / `.scss` 実装向けの layout、responsive design、selector、cascade、保守しやすいスタイル設計の参照。 | CSS 実装を書く前に、スタイル設計の判断基準を確認するとき。 |
| `code-go` | Go / `.go` 実装と Go テスト向けの goroutine、interface、error handling、設計、テスト方針の参照。 | Go 実装や Go テストを書く前に、実装ガイドを確認するとき。 |
| `code-react` | React / JSX / TSX / `.jsx` / `.tsx` 実装向けの hooks、props、state、Atomic Design を補助観点にしたコンポーネント設計の参照。 | React 実装を書く前に、UI 実装ガイドを確認するとき。 |
| `code-ruby` | Ruby / Rails / `.rb` 実装と RSpec / Minitest テスト向けの Active Record、MVC、設計、テスト方針の参照。 | Ruby / Rails 実装やテストを書く前に、実装ガイドを確認するとき。 |
| `code-ts` | TypeScript / `.ts` / `.tsx` 実装向けの型設計、interface / type、strict typing、責務分離の参照。 | TypeScript 実装を書く前に、実装ガイドを確認するとき。 |

### E. Deprecated Internal Skills

`skills/internal/deprecated/` 配下の skill は履歴保持用です。plugin には含めません。

| Skill | 理由 |
|---|---|
| `cmd-agent-status` | active skill としては配布せず、履歴保持のため deprecated へ移動した。 |
| `cmd-rmbranch` | active skill としては配布せず、履歴保持のため deprecated へ移動した。 |
| `cmd-start-branch` | active skill としては配布せず、履歴保持のため deprecated へ移動した。 |
| `beautify-commit` | active skill としては配布せず、履歴保持のため deprecated へ移動した。 |
| `ci-fix` | active skill としては配布せず、履歴保持のため deprecated へ移動した。 |
| `code-test-review` | テスト差分レビューを `code-review` に統合したため、独立配布しない。 |
| `code-typo` | active skill としては配布しない。 |
| `role-doc` | active skill としては配布しない。 |

## External Skills

外部 skill は `external-skills.json` で取得元と hash lock を管理します。marketplace build が全外部 skill を検証・取り込み、配布先は plugin membership から決まります。

| Skill | 取得元 | 主な分類 | 概要 |
|---|---|---|---|
| `skill-creator` | `anthropics/skills` | 行動定義型 / ワークフロー型 | 新しい skill の作成、既存 skill の更新、メタデータや検証手順の整備に使う。 |
| `frontend-design` | `anthropics/skills` | 行動定義型 / ワークフロー型 | Web UI、ページ、コンポーネント、HTML/CSS/React などを高品質な frontend design として実装するときに使う。 |
| `grilling` | `mattpocock/skills` | 判断定義型 / 判断基準・評価型 | 計画や設計を厳しく質問し、曖昧さや判断漏れを潰すために使う。 |
| `empirical-prompt-tuning` | `mizchi/skills` | 判断定義型 / スコアリング・査定型 | skill やプロンプトを実験的に改善し、評価と反復で性能を詰めるために使う。 |
| `grape-usage` | `version-1/grape` | 行動定義型 / リファレンス型 | 安全な push、ポリシーで制限された rebase、worktree 操作など、grape コマンドの正確な利用方法を確認するときに使う。 |
| `japanese-tech-writing` | `k16shikano/fd287c3133457c4fd8f5601d34aa817d` | 判断定義型 / 判断基準・評価型 | 日本語の技術文書・書籍原稿を執筆・推敲するときに、段落構成、論証、読者の負荷、語り、表現の規範を適用する。 |
| `tasq-cli` | `version-1/tasq` | リファレンス型 | Tasq の typed CLI を使った課題、コメント、artifact の操作を確認するときに使う。個別 plugin `skill-tasq-cli` として利用する。 |
| `tq-orchestrator` | `version-1/tasq` | ワークフロー型 | Tasq の ready issue を監視し、subagent への委譲をオーケストレーションするときに使う。個別 plugin `skill-tq-orchestrator` として利用する。 |

## 運用メモ

- 新しい internal skill を追加したら、この一覧にも追記する。
- 追加時は [docs/guide/skill-category.md](guide/skill-category.md) を参照し、主用途に最も近い分類へ配置する。
- 複数カテゴリにまたがる skill は、重複掲載せず主分類に置き、必要に応じて概要に複合要素を明記する。
- `code-general` や `code-review` のようにガードレールや思想を含む skill でも、一覧では主な行動や判断の用途を優先して配置する。
- `skills/internal/<name>/agents/openai.yaml` がある場合は、UI 表示用の説明も `SKILL.md` と矛盾しないように更新する。
- 他 skill への明示的な依存や利用順序を追加したら、[docs/skill-dependency-map.md](skill-dependency-map.md) も更新する。
- 配布前は `make build-plugins` と `make deploy-plugins-dry-run` で生成対象と配布予定を確認する。
