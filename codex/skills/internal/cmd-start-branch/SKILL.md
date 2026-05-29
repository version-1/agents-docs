---
name: cmd-start-branch
description: Git リポジトリで新しい作業を始める前に、最新のデフォルトブランチへ追従し、タスク内容に合う作業ブランチを作成して、不要なローカルブランチ整理を非同期に依頼するときに使う。ユーザーが「ブランチ切って」「作業を始める準備して」「新しいタスク用ブランチを作って」「最新 main から始めたい」などを依頼した場合はこの skill を使い、ブランチ名を報告してからタスク詳細の入力を待つ。
---

# Start Branch

## 目的

新しい作業に入る前の Git 準備を安全に行う。
最新のデフォルトブランチから意味のある作業ブランチを作り、並行して不要なローカルブランチ整理を任せる。
準備が終わったら、ブランチ名をユーザーへ伝え、実装や調査には進まずタスク詳細を待つ。

## 絶対ルール

- 作業ツリーに未コミット変更がある場合は、内容を確認し、ブランチ切り替えで上書きや混入が起きないか判断する。
- ユーザーの変更を revert、削除、上書きしない。
- デフォルトブランチを確認してから更新する。`main` と決め打ちしない。
- 最新化には通常 `git fetch` と、必要なら `git pull --ff-only` を使う。履歴改変、merge commit 作成、force 操作はしない。
- ブランチ名はタスク内容から短く具体的にする。曖昧なら安全な仮名を使い、ユーザーが後で変えられるようにする。
- 不要ブランチ整理は `$cmd-dispatch-agent` で別 agent に投げ、依頼内容に `$cmd-rmbranch` を含める。結果は待たない。
- ブランチ作成後は、ブランチ名を出力し、タスク詳細が入力されるのを待つ。まだ実装、編集、調査を始めない。

## 手順

1. `git status --short --branch` で現在のブランチと作業ツリー状態を確認する。
2. デフォルトブランチを確認する。
   - まず `git symbolic-ref refs/remotes/origin/HEAD` を見る。
   - 取れない場合は `git remote show origin` を使う。
   - それでも不明なら、`main` と `develop` の存在を確認し、判断できなければユーザーへ確認する。
3. 最新のデフォルトブランチへ追従する。
   - `git fetch origin` を実行する。
   - デフォルトブランチへ移動できる状態なら checkout する。
   - ローカルのデフォルトブランチが behind の場合は `git pull --ff-only` で進める。
   - ローカル差分や競合リスクがある場合は止めて、ブランチ作成前に状況を報告する。
4. タスク内容からブランチ名を作る。
   - 「ブランチ種別の選び方」の表から、指示内容に最も合う種別を選ぶ。
   - 形式は `<type>/<short-task>` にする。
   - 英小文字、数字、ハイフンを使う。
   - 例: `feat/add-start-branch-skill`, `fix/login-timeout`, `refactor/split-auth-service`
5. `git checkout -b <branch-name>` で作業ブランチを作成する。
   - 同名ブランチがある場合は、既存ブランチの位置と状態を確認し、上書きせず別名を提案する。
6. `$cmd-dispatch-agent` を使って、別 agent に `$cmd-rmbranch` を依頼する。
   - 依頼には「ローカルブランチのみ」「main/develop/current は削除しない」「リモートは触らない」「未マージは強制削除せず候補報告」を含める。
   - agent の結果は待たず、起動情報だけ報告する。
7. 最後に次の形式で報告し、タスク詳細の入力を待つ。

```text
作業ブランチを作成しました: `<branch-name>`
不要ブランチ整理は別 agent に依頼済みです。結果は待たずに進めています。

このブランチで進めるタスクの詳細を送ってください。
```

## ブランチ名の判断

- ユーザーがタスク名や Issue 番号を示している場合は、それを優先する。
- Issue 番号がある場合は `feat/123-short-task` のように先頭へ入れる。
- まだ内容が薄い場合は `chore/start-task-branch` のような仮名にし、後で必要なら rename を提案する。
- 既存ブランチと衝突する場合は末尾に短い識別子を足す。

## ブランチ種別の選び方

指示内容の主目的から 1 つ選ぶ。
複数に見える場合は、ユーザーにとって主な成果物になる変更を優先する。
迷う場合は `chore` ではなく、より具体的な種別を選ぶ。

| Type | 使う場面 | 例 |
|---|---|---|
| `feat` | 新機能、既存機能へのユーザー向け能力追加 | `feat/report-filter`, `feat/42-export-csv` |
| `fix` | バグ修正、不具合回避、期待と違う挙動の修正 | `fix/login-timeout`, `fix/price-rounding` |
| `refactor` | 挙動を変えない構造整理、責務分割、命名整理 | `refactor/split-auth-service`, `refactor/repository-contract` |
| `test` | テスト追加、テスト修正、fixture や検証手順の整備 | `test/add-order-cases`, `test/stabilize-api-spec` |
| `perf` | 性能改善、メモリ削減、クエリ最適化 | `perf/cache-rate-lookup`, `perf/reduce-render-cost` |
| `style` | フォーマット、lint、見た目だけのコード整形。UI 見た目変更は `feat` または `fix` を優先する | `style/format-go-files`, `style/lint-imports` |
| `ci` | CI、ワークフロー、ビルドパイプライン、リリース自動化 | `ci/add-test-workflow`, `ci/cache-go-build` |
| `docs` | README、設計メモ、Runbook など文書だけの変更 | `docs/update-runbook`, `docs/add-skill-guide` |
| `chore` | 依存更新、設定整理、生成物更新、運用上の雑務。上記に当てはまるならそちらを優先する | `chore/update-deps`, `chore/start-task-branch` |

## 報告に含めること

- 作成したブランチ名
- 追従したデフォルトブランチ名
- 不要ブランチ整理 agent の起動情報
- タスク詳細を待っていること

## 避けること

- ブランチ作成後に、ユーザーがまだ詳細を送っていないタスクを推測して進めること。
- `git reset --hard`、force push、rebase、未確認の削除を行うこと。
- ローカル変更が混在しているのに確認せずデフォルトブランチへ移動すること。
- ブランチ整理 agent の完了を待って、開始手順を止めること。
