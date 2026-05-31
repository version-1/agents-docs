---
name: ci-fix
description: GitHub Actions や CI の失敗を調査し、原因を切り分けて修正するときに使う。ユーザーが「CI 落ちてる」「GitHub Actions 直して」「workflow が失敗」「checks が red」「test failure を直して」「lint/build が CI でだけ落ちる」などを依頼した場合は必ず使う。gh CLI、workflow ログ、ローカル再現、最小修正、再検証、PR 反映までを安全な順序で進める。
---

# CI Fix

## 目的

GitHub Actions / CI の失敗を、ログの症状だけで場当たり的に直さず、失敗箇所、再現条件、原因、修正、再検証を根拠付きでつなげる。

CI は複数の job、matrix、環境差、キャッシュ、依存サービス、権限、secret、外部 API にまたがるため、最初に観測範囲を狭める。修正は可能な限りローカルで再現できる形に落とし込み、CI 専用の一時しのぎを避ける。

## 基本方針

- まず失敗した check、workflow、job、step、ログの該当箇所を特定する。
- ログ全体を貼り付けず、エラー、直前の実行コマンド、環境差が分かる範囲を読む。
- 同じエラーが複数 job に出ている場合は、最初の根本原因を優先する。
- ローカル再現できる失敗は、CI を再実行する前にローカルで直す。
- CI でだけ失敗する場合は、OS、shell、working directory、env、permissions、cache、services、secrets、toolchain version を疑う。
- 実装修正が必要な場合は `$coding` を併用する。
- テストや検証コマンドの追加、flake 対策、回帰テストが必要な場合は `$code-test` を併用する。
- PR 作成や更新まで進める場合は `$cmd-create-pr` に引き継ぐ。

## 事前確認

最初に現在位置と PR / check の状況を確認する。

```bash
git status --short --branch
git branch --show-current
gh pr status
```

PR 番号が分かる場合は、check 一覧を確認する。

```bash
gh pr checks <pr-number>
```

PR 番号が不明な場合は、現在ブランチの PR を確認する。

```bash
gh pr view --json number,title,url,headRefName,baseRefName,statusCheckRollup
```

## 調査手順

1. 失敗している workflow / job / step を特定する。
2. 失敗ログから、実行コマンド、エラーメッセージ、直前の setup step、matrix 条件を抜き出す。
3. 失敗が次のどれに近いか分類する。
   - Code failure: テスト、lint、typecheck、build、format が実装差分で壊れている。
   - CI definition failure: workflow YAML、working directory、permissions、path filter、artifact、cache key、service container が壊れている。
   - Environment drift: Node / Go / Ruby / Python / OS / shell / dependency version の差で壊れている。
   - Secret or permission failure: token、secret、OIDC、GitHub permissions、fork PR 制約で壊れている。
   - Flake or external dependency: timing、並列実行、ネットワーク、外部 API、rate limit、順序依存で壊れている。
4. 関連する workflow ファイル、実行コマンド、直近差分を読む。
5. ローカル再現コマンドを決める。CI の step と同じ working directory、env、toolchain をできるだけ揃える。
6. 最小修正を実装する。
7. ローカル検証を実行する。
8. 必要に応じて CI を再実行し、同じ失敗が解消したことを確認する。
9. 修正内容、根本原因、検証結果、残るリスクを報告する。

## GitHub Actions ログ確認

`gh pr checks` で失敗 check の名前を確認する。

```bash
gh pr checks <pr-number>
```

詳細が必要な場合は run を特定する。

```bash
gh run list --branch <branch-name> --limit 10
gh run view <run-id> --json name,status,conclusion,event,headBranch,headSha,jobs
```

ログを読む場合は、失敗 job に絞る。大量ログを読む前に job 名と step 名を特定する。

```bash
gh run view <run-id> --job <job-id> --log
```

ログから見るポイント:

- `Run ...` の実行コマンド
- `working-directory`
- `env`
- toolchain setup の version
- 最初に失敗した assertion、compile error、lint error、permission error
- 後続の cascade failure ではなく root failure

## 修正の判断基準

### Code failure

- CI step と同じコマンドをローカルで実行する。
- テスト失敗は、失敗したテスト名、入力、期待値、実際値を確認する。
- lint / typecheck は、指摘の意味を理解してから修正する。
- 実装を変える場合は `$coding` を使い、既存仕様と公開契約を壊さない。
- バグ修正なら、再発しやすい領域に回帰テストを追加するか検討する。

### CI definition failure

- workflow YAML の構文、indent、job dependency、permissions、shell、working directory を確認する。
- path filter や matrix 条件が意図した job を動かしているか確認する。
- cache は高速化のための補助として扱い、cache hit 前提でしか通らない構成にしない。
- artifact や generated file は、作成 step と参照 step の job 境界を確認する。

### Environment drift

- CI とローカルの toolchain version を比較する。
- version を pin するか、既存の `.tool-versions`、`go.mod`、`packageManager`、lockfile、Dockerfile に合わせる。
- 依存更新で直す場合は、lockfile 変更の影響を確認する。
- CI だけの workaround ではなく、プロジェクトの標準環境を明確にする。

### Secret or permission failure

- secret 値そのものを表示しない。
- secret が必要な workflow か、fork PR で secret が渡らないケースか確認する。
- `permissions` が不足している場合は、必要最小限の権限だけ追加する。
- token や secret の新規発行、設定変更が必要な場合は、ユーザーに判断を戻す。

### Flake or external dependency

- 失敗が再実行で消えるかだけで終わらせず、flake の原因を探す。
- sleep、時刻、乱数、並列実行、共有リソース、外部 API、rate limit を確認する。
- テスト側で同期、固定 seed、依存の分離、タイムアウトの明確化を検討する。
- 再実行が妥当な場合も、flake の観測結果を報告に残す。

## CI 再実行

CI 再実行は、ログ確認とローカル修正後に行う。

```bash
gh run rerun <run-id> --failed
```

再実行する前に確認すること:

- 未コミット変更が CI に反映されない状態ではないか。
- push が必要な場合、直接 `git push` せず、リポジトリの安全手順に従っているか。
- 外部サービスや課金リソースを消費する workflow ではないか。

## 完了条件

- 失敗した workflow / job / step と root cause を説明できる。
- 修正差分が root cause に対応している。
- 関連するローカル検証が通っている。
- CI の再実行または代替根拠で、同じ失敗が解消したことを確認している。
- CI でだけ確認できる項目が未確認の場合は、その理由と次の確認手順を残している。

## 報告形式

最終報告は次を簡潔に含める。

```text
CI 修正結果:
- 失敗箇所: <workflow / job / step>
- 原因: <root cause>
- 修正: <変更内容>
- 検証: <ローカル検証と CI 再実行結果>
- 残るリスク: <flake、外部依存、未確認事項があれば>
```

PR に反映する場合は、commit、push、PR description 更新を `$cmd-create-pr` の手順へつなげる。
