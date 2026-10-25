---
name: cmd-commit
description: Git リポジトリで変更内容を確認し、適切な粒度とメッセージで安全に commit を作成するときに使う。ユーザーが「commitして」「コミット作って」「差分をコミット」「変更を分けてコミット」「この変更をコミット」などを依頼した場合に使う。git status / diff / log を確認し、未確認の変更や unrelated changes を勝手に含めず、必要に応じて commit 粒度、stage 対象、commit message を整理する。
---

# Commit

## 目的

ユーザーや他エージェントの変更が混在する作業ツリーで、意図に合う範囲だけを stage し、読める履歴を作る。

## ワークフロー

1. `git status --short` で変更ファイル、staged / unstaged / untracked の状態を確認する。
2. `git diff` と `git diff --staged` を読み、必要に応じて `git log --oneline -n 10` で既存履歴の形式を確認する。
3. 変更の目的、影響範囲、無関係な変更の有無と commit 粒度を判断する。複数 commit が必要なら、理由と分割案を短く提示して確認する。
4. commit に含めるファイルだけを stage する。確認せずに `git add .` しない。
5. `git diff --staged` で最終的に入る内容を確認する。
6. 既存履歴の言語と形式に合わせて commit message を作る。
7. `git commit` を実行する。
8. `git status --short` で残っている変更を確認し、commit hash、message、残変更、検証状況を報告する。

## 絶対ルール

- ユーザーが明示していない変更を revert、削除、上書き、commit しない。未確認または無関係な変更が混在している場合は対象ファイルだけを stage する。
- untracked file は内容と必要性を確認してから stage する。
- 秘密情報、ローカル設定、生成キャッシュ、バックアップ、ログを含めない。
- pre-commit hook や commit 失敗が起きたら、出力を読み、変更されたファイルや失敗理由を確認してから次を判断する。
- 署名、amend、force push、履歴改変はユーザーが明示した場合だけ行う。

## Commit 粒度

- 1つの目的にまとまる変更（従属する実装、テスト、文書を含む）は 1 commit にし、無関係な修正、別機能、単独のフォーマット変更、生成物更新は分ける。
- 分割したほうがよいが判断が曖昧な場合は、commit 候補と含めるファイルを示して確認する。
- ユーザーが「全部コミット」と言っても、明らかな不要物や秘密情報は含めず確認する。

## Commit Message

- ユーザー指定、または既存履歴の言語、prefix、Conventional Commits などの形式を優先する。
- 判断材料がなければ日本語で書く。
- 1行目は変更の目的を簡潔に書く。
- commit message には必ず本文を付け、理由、主な変更、検証内容を書く。
- 「修正」「更新」「調整」だけの曖昧な件名を避ける。
- prefix、タイトル、本文フォーマットのデフォルトが必要な場合は `references/commit-message.md` を読む。

## 確認と報告

commit 後は、次を簡潔に報告する。

- 作成した commit hash と commit message
- commit に含めた主な変更
- 実行した検証、または未実行の理由
- 残っている未コミット変更の有無

commit しなかった場合は、止めた理由、確認が必要な点、次に実行すべき操作を明示する。
