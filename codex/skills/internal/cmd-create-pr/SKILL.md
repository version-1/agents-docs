---
name: cmd-create-pr
description: GitHub Pull Request を安全な手順で作成または更新するときに使う。ユーザーが「PR 作って」「PR 出して」「pull request 作成して」「この変更を PR にして」「pr を出す手順を進めて」などを依頼した場合は必ず使う。差分確認、検証、commit、safe-git-push、PR description 作成、gh pr create / edit までの順序を整理し、未確認の変更や直接 git push を避ける。
---

# Create PR

## 目的

GitHub Pull Request を出す前後の作業を、安全で再現しやすい手順に揃える。

作業ツリーにユーザーや別エージェントの変更が混在している前提で、PR に含める差分、検証結果、PR 本文、push 方法を確認しながら進める。

## 基本方針

- PR は「読んだ差分」「必要な検証」「意図が分かる commit」「レビューしやすい description」が揃ってから作る。
- ユーザーの未確認変更を勝手に commit、push、PR に含めない。
- `git push` は直接実行しない。prompt なしで push する場合は、引数なしの `~/.codex/bin/safe-git-push` を使う。
- `safe-git-push` が拒否した場合や見つからない場合は、拒否理由を確認してユーザーに方針を確認する。
- 既存 PR がある場合は、新規作成ではなく更新を検討する。
- PR 本文を作るときは `$pr-description` を使う。
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

## 手順

1. 現在のブランチ、base branch、remote を確認する。
2. `git status --short --branch` で staged / unstaged / untracked を確認する。
3. `git diff` と `git diff --staged` を読み、PR に含めるべき差分だけか判断する。
4. 未コミット変更がある場合は、PR に含める範囲を決めて `$cmd-commit` で commit する。
5. 変更内容に合う検証を実行する。文書だけの変更なら、Markdown のリンク、見出し、差分確認で十分な場合がある。
6. `git log --oneline <base>..HEAD` と `git diff --stat <base>...HEAD` で PR 差分を確認する。
7. `~/.codex/bin/safe-git-push` で current branch を push する。
8. `$pr-description` で PR description を作る。
9. 既存 PR がなければ `gh pr create`、既存 PR があれば `gh pr edit` で description を更新する。
10. `gh pr view` で URL、base、head、state を確認する。
11. PR URL、含めた commit、検証結果、未対応事項を報告する。

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

本文の形式は `$pr-description` のテンプレートを優先する。

## 判断ルール

- 差分がない場合は PR を作らず、作るものがないことを報告する。
- untracked file は内容を確認してから commit 対象にする。
- 秘密情報、ローカル設定、生成キャッシュ、バックアップ、ログを PR に含めない。
- base branch が不明な場合は、origin の default branch を確認する。不明なまま進めない。
- branch が default branch の場合は、作業ブランチを作るか確認する。default branch から直接 PR を作ろうとしない。
- 検証が失敗した場合は、失敗内容を報告し、PR を出すか修正するか確認する。ただしユーザーが明示した場合は、失敗を PR 本文に書いたうえで作成してよい。
- 既存 PR がある場合は、同じ branch で重複 PR を作らない。

## 報告形式

PR 作成または更新後は、次を簡潔に報告する。

```text
PR を作成しました: <url>

- base: <base-branch>
- head: <head-branch>
- commits: <主な commit>
- 検証: <実行した検証、または未実行の理由>
- 補足: <未対応事項や注意点があれば>
```

PR を作らなかった場合は、止めた理由と次に必要な操作を明示する。
