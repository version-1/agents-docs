---
name: beautify-commit
description: Git リポジトリで、ベースブランチまたは残したい履歴の終点 commit との差分に含まれる commit や巨大な未整理差分を、意味のある変更単位ごとに整理して commit し直すときに使う。ユーザーが「commit を分割して」「巨大なコミットを分けて」「連続コミットを整理して」「履歴をきれいにして」「この差分を変更単位でコミットして」「ベースブランチとの差分を整理して」などを依頼した場合に使う。基準 ref の指定を確認し、指定がなければデフォルトブランチでよいか yes/no で確認する。git log / diff / status を確認し、必要に応じて grill-me skill で分割方針を詰め、整理後の commit 作成には cmd-commit skill を使う。
---

# Split Commit

## 目的

ベースブランチ、または残したい履歴の終点 commit との差分に含まれる commit や大きな未整理差分を、レビューしやすく戻しやすい変更単位へ整理して commit し直す。
履歴を書き換える可能性があるため、基準 ref、対象範囲、push 状態、バックアップ、残す変更を明確にしてから操作する。

## 最初に確認すること

1. ユーザーが基準 ref を指定しているか確認する。基準 ref はベースブランチでも、残したい履歴の終点 commit hash でもよい。
2. 指定がない場合は、リポジトリのデフォルトブランチを調べ、「基準 ref は `<default-branch>` でよいですか？ yes/no」で確認する。
3. `git status --short` で作業ツリーが clean か、未コミット変更があるか確認する。
4. `git branch --show-current`、`git log --oneline --decorate -n 20`、必要に応じて upstream を確認する。
5. 基準 ref がブランチなら必要に応じて `git merge-base <base-ref> HEAD` を確認し、差分範囲を決める。
6. `git log --oneline --reverse <base-ref>..HEAD` と `git diff --stat <base-ref>..HEAD` で、差分に含まれる commit と変更ファイルを確認する。
7. 対象が「未コミットの巨大差分」「基準 ref 以降の連続 commit」「古い commit を含む履歴」のどれか分類する。
8. 分割案を作るため、`git diff <base-ref>..HEAD`、`git show --stat`、`git show --name-status` を読む。

詳細な操作パターンが必要な場合は `references/split-workflows.md` を読む。

## 絶対ルール

- ユーザーが明示していない変更を捨てない。
- 基準 ref が未指定の場合、デフォルトブランチでよいか yes/no で確認してから進める。
- 履歴を書き換える操作をするときは、対象 branch と対象 commit 範囲を明確にする。
- push 済み、共有済み、PR 作成済みの履歴を書き換える場合は、実行前にユーザー確認を取る。
- `git reset --hard`、`git rebase -i`、`git commit --amend`、`git push --force` は、必要性と影響を説明してから使う。
- 分割作業の前に、可能なら `backup/<branch>-before-split-<date>` のようなバックアップ branch を作る。
- 読んでいない差分を stage しない。`git add .` ではなく、ファイル単位または patch 単位で stage する。
- 分割後は、それぞれの commit が単独で目的を説明できるか確認する。

## 分割粒度

- 1 commit は 1 つの意図に絞る。
- コード、テスト、ドキュメントが同じ意図に従属するなら同じ commit に含める。
- 機能追加、バグ修正、リファクタ、フォーマット、生成物更新、設定変更は混ぜない。
- 後続 commit が前の commit に依存する場合は、依存順に並べる。
- レビュー単位として説明できない commit は、さらに分けるか隣接 commit と統合する。

## 基本ワークフロー

1. 対象範囲を決める。
2. 基準 ref との差分全体を読み、変更単位の候補を作る。
3. 分割案が曖昧、粒度が大きい、責務境界に迷う場合は、必要に応じて `grill-me` skill を使って整理方針を詰める。
4. 分割案を作り、必要ならユーザーに確認する。
5. バックアップ branch を作る。
6. 基準 ref との差分に含まれる commit をいったん未コミット差分へ戻し、整理して commit し直せる状態にする。
7. 変更単位ごとに stage 対象を作る。
8. 整理後の commit 作成では `cmd-commit` skill を使い、staged diff、commit message、残変更を確認しながら commit する。
9. すべて分割したら `git log --oneline --decorate <base-ref>..HEAD` と `git status --short` を確認する。
10. 実行した検証、作成した commit、残変更、履歴書き換えの有無を報告する。

## Commit Message

- 既存履歴の prefix、言語、Conventional Commits などの形式に合わせる。
- 分割後の各 commit message は、その commit 単独の目的を書く。
- 元の巨大 commit の message をそのまま流用しない。
- commit 作成時は `cmd-commit` skill を使う。
- 体裁に迷う場合は、`cmd-commit` skill の `references/commit-message.md` が利用可能なら参照する。

## 報告

完了時は次を簡潔に報告する。

- 分割前の対象 commit / 範囲
- 使用した基準 ref
- 作成したバックアップ branch
- 作成した commit 一覧
- 残っている未コミット変更の有無
- 実行した検証、または未実行の理由
- force push など、ユーザー側で必要な次操作があるか
