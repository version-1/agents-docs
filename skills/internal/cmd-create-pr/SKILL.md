---
name: cmd-create-pr
description: GitHub Pull Request を安全な手順で作成または更新するときに使う。ユーザーが「PR 作って」「PR 出して」「pull request 作成して」「この変更を PR にして」「pr を出す手順を進めて」などを依頼した場合は必ず使う。差分確認、検証、commit、最新 base branch への追従、role-reviewer による PR 前レビュー、High 指摘の自動対応、grape push、PR description 作成、gh pr create / edit までの順序を整理し、未確認の変更や直接 git push を避ける。
---

# Create PR

ユーザーや他エージェントの変更が混在し得るため、読んだ差分、必要な検証、意図が分かる commit、レビューしやすい本文を揃える。ユーザーの未確認変更を勝手に commit、push、PR に含めない。commit は `$cmd-commit`、本文は `$format-pr-description`、既存 PR の編集は global 規則に従い `safe-gh-edit` を使う。

## 必須の安全条件

- 最初に branch、remote、status、staged/unstaged/untracked の差分、既存 PR、base を確認する。untracked は内容確認後だけ含め、秘密情報、ローカル設定、生成キャッシュ、バックアップ、ログは含めない。
- base は `git fetch origin` 後に確認し、head が最新 base を含むまで古い base のままレビュー、push、PR 作成を進めない。base 不明なら origin の default branch を確認し、不明なまま進めない。
- default branch からは直接 PR を作らず、作業 branch 作成または変更退避を確認する。
- 追従していない場合、独自 commit がなければ fast-forward のみ可。独自 commit がある場合は、merge commit、`grape rebase`、積み直しの方針をユーザーに確認する。rebase 開始は `grape rebase` のみで、競合後だけ `git rebase --continue|--abort|--skip` を使える。`grape` が拒否・不在なら Git へ自動フォールバックせず、理由を確認して方針を聞く。
- push は `grape push` のみ。同じ branch の既存 PR があれば重複作成しない。
- 検証失敗時は内容を報告して作成か修正か確認する（ユーザーが明示的に許可した場合のみ本文へ記載して作成可）。

## 手順

1. `git status --short --branch`、`git remote -v`、`git branch --show-current`、`gh pr status`、必要なら `git symbolic-ref refs/remotes/origin/HEAD` を確認する。
2. `git diff` と `git diff --staged` を読み、含める範囲を決め、未コミット分は `$cmd-commit` で commit する。
3. `git fetch origin` 後、`git merge-base --is-ancestor origin/<base> HEAD` で追従を確認する。失敗時は双方の log を読み、上記安全条件に従う。
4. 適切な検証を実行し、`git log --oneline <base>..HEAD` と `git diff --stat <base>...HEAD` で PR 差分を確認する。
5. `$role-reviewer` で `<base>...HEAD`、検証、未コミット変更、base 追従、既知制約をレビューする。High は安全に自動修正し、差分確認・再検証・追加 commit・再レビューを行う。High が残る、反復する、安全な自動修正を越える場合は PR を止め、理由と選択肢を確認する。Medium/Low は必要なら対応し、進める場合は本文または最終報告へ残す。
6. High がないことを確認後に `grape push`、本文ファイルを作成し、既存 PR は安全な edit、新規は `gh pr create --base <base> --head <branch> --title "<title>" --body-file <file>` を使う。`gh pr view` で URL、base、head、state を確認する。

本文は Summary、主な変更、実施内容、動作確認、レビューしてほしい点、未対応事項/相談を含める。

## 報告

PR URL、base/head、主な commit、base 追従、検証、レビューと High 対応、未対応事項を簡潔に報告する。作らなければ停止理由と次に必要な操作を明示する。
