---
name: cmd-create-pr
description: GitHub Pull Request を安全な手順で作成または更新するときに使う。ユーザーが「PR 作って」「PR 出して」「pull request 作成して」「この変更を PR にして」「pr を出す手順を進めて」などを依頼した場合は必ず使う。差分確認、検証、commit、最新 base branch への追従、role-reviewer による PR 前レビュー、High 指摘の自動対応、safe-git-push、PR description 作成、gh pr create / edit までの順序を整理し、未確認の変更や直接 git push を避ける。
---

# Create PR

## 目的

GitHub Pull Request を出す前後の作業を、安全で再現しやすい手順に揃える。

作業ツリーにユーザーや別エージェントの変更が混在している前提で、PR に含める差分、検証結果、PR 本文、push 方法を確認しながら進める。

## 基本方針

- PR は「読んだ差分」「必要な検証」「意図が分かる commit」「レビューしやすい description」が揃ってから作る。
- ユーザーの未確認変更を勝手に commit、push、PR に含めない。
- PR 作成前に `git fetch origin` で base branch を更新し、head branch を最新 base に追従させる。古い base のままレビュー、push、PR 作成を進めない。
- `git push` は直接実行しない。prompt なしで push する場合は、引数なしの `~/.codex/bin/safe-git-push` を使う。
- `safe-git-push` が拒否した場合や見つからない場合は、拒否理由を確認してユーザーに方針を確認する。
- 既存 PR がある場合は、新規作成ではなく更新を検討する。
- PR を新規作成または更新する前に `$role-reviewer` で差分をレビューし、High severity の指摘は PR 作成前に自動対応する。
- PR 本文を作るときは `$format-pr-description` を使う。
- commit が必要な場合は `$cmd-commit` を使う。

## 事前確認

最初に現在の状態を確認する。

```bash
git status --short --branch
git remote -v
git branch --show-current
```

次に、PR の向きと既存 PR の有無を確認する。

```bash
gh pr status
```

必要に応じて default branch を確認する。

```bash
git symbolic-ref refs/remotes/origin/HEAD
```

base branch を更新する。

```bash
git fetch origin
```

## 手順

1. 現在のブランチ、base branch、remote を確認する。
2. `git status --short --branch` で staged / unstaged / untracked を確認する。
3. `git diff` と `git diff --staged` を読み、PR に含めるべき差分だけか判断する。
4. 未コミット変更がある場合は、PR に含める範囲を決めて `$cmd-commit` で commit する。
5. `git fetch origin` で base branch を最新化する。
6. current branch が最新 base を含んでいるか確認する。
   - `git merge-base --is-ancestor origin/<base> HEAD` が成功するなら、head は最新 base に追従済み。
   - 成功しない場合は、まず `git log --oneline HEAD..origin/<base>` と `git log --oneline origin/<base>..HEAD` で base 側と head 側の差分を確認する。
   - current branch に独自 commit がない場合だけ、`git merge --ff-only origin/<base>` で fast-forward してよい。
   - current branch に独自 commit がある場合は、merge commit 作成、rebase、作業の積み直しのどれを使うかユーザーへ確認する。確認なしに rebase や merge commit を作らない。
7. 追従後に変更内容に合う検証を実行する。文書だけの変更なら、Markdown のリンク、見出し、差分確認で十分な場合がある。
8. `git log --oneline <base>..HEAD` と `git diff --stat <base>...HEAD` で PR 差分を確認する。
9. `$role-reviewer` で PR 差分をレビューする。レビュー対象は `<base>...HEAD` の差分、実行済み検証、未コミット変更の有無、最新 base への追従状況。
10. High severity の指摘がある場合は、PR 作成前に自動で対応する。対応後は必要な検証を再実行し、必要に応じて `$cmd-commit` で追加 commit を作る。
11. High 指摘が残っていないことを確認するまで `$role-reviewer` の確認を繰り返す。
12. `~/.codex/bin/safe-git-push` で current branch を push する。
13. `$format-pr-description` で PR description を作る。
14. 既存 PR がなければ `gh pr create`、既存 PR があれば `gh pr edit` で description を更新する。
15. `gh pr view` で URL、base、head、state を確認する。
16. PR URL、含めた commit、検証結果、レビュー結果、base 追従状況、未対応事項を報告する。

## PR 前レビュー

PR を出す直前のレビューは、公開前に明らかな品質リスクを潰すためのゲートとして扱う。

レビュー時は `$role-reviewer` を使い、必要に応じて `$code-review` も併用する。確認対象は次のとおり。

- `<base>...HEAD` の差分
- 実行済みの検証結果
- 作業ツリーに未コミット変更が残っていないか
- head branch が最新 base branch に追従済みか
- PR description に書くべき既知の制約や未対応事項

レビュー結果は severity を明示して扱う。

- High: PR 作成前に自動対応する。対応できない場合は PR 作成を止め、理由と選択肢をユーザーへ確認する。
- Medium / Low: 影響と緊急度を判断し、必要なら対応する。PR に進む場合は PR description または最終報告に残す。
- None: そのまま PR 作成へ進む。

High 指摘に対応した場合は、修正差分を読み、関連する検証を再実行し、追加 commit を作ってから再レビューする。レビューが同じ High 指摘を繰り返す場合や、安全な自動修正の範囲を超える場合は、無理に進めずユーザーへ判断を戻す。

## PR 作成コマンド

PR 本文はファイルに書いてから渡す。

```bash
gh pr create --base <base-branch> --head <current-branch> --title "<title>" --body-file <body-file>
```

既存 PR を更新する場合:

```bash
gh pr edit <pr-number> --body-file <body-file>
```

PR 作成後の確認:

```bash
gh pr view <pr-number> --json number,title,url,headRefName,baseRefName,state
```

## PR Description

PR description には最低限、次を含める。

- Summary
- 主な変更点
- やったこと
- 動作確認
- レビュー時に見てほしい点
- 未対応事項や相談ごと

本文の形式は `$format-pr-description` のテンプレートを優先する。

## 判断ルール

- 差分がない場合は PR を作らず、作るものがないことを報告する。
- untracked file は内容を確認してから commit 対象にする。
- 秘密情報、ローカル設定、生成キャッシュ、バックアップ、ログを PR に含めない。
- base branch が不明な場合は、origin の default branch を確認する。不明なまま進めない。
- branch が default branch の場合は、作業ブランチを作るか確認する。default branch から直接 PR を作ろうとしない。
- head branch が最新 base branch を含んでいない場合は、PR 作成前に追従する。fast-forward できない場合は、確認なしに merge commit や rebase を行わない。
- 検証が失敗した場合は、失敗内容を報告し、PR を出すか修正するか確認する。ただしユーザーが明示した場合は、失敗を PR 本文に書いたうえで作成してよい。
- `$role-reviewer` で High severity の指摘が出た場合は、PR を作る前に修正し、再検証と再レビューを行う。安全に自動対応できない High 指摘は、PR 作成を止めてユーザーに確認する。
- 既存 PR がある場合は、同じ branch で重複 PR を作らない。

## 報告形式

PR 作成または更新後は、次を簡潔に報告する。

```text
PR を作成しました: <url>

- base: <base-branch>
- head: <head-branch>
- commits: <主な commit>
- base 追従: <最新 base に追従済み / 確認できなかった理由>
- 検証: <実行した検証、または未実行の理由>
- レビュー: <role-reviewer の結果、High 対応の有無>
- 補足: <未対応事項や注意点があれば>
```

PR を作らなかった場合は、止めた理由と次に必要な操作を明示する。
