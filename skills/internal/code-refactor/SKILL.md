---
name: code-refactor
description: コードレビュー結果、code quality report、git diff、PR 差分、指定ファイルをもとに、挙動を変えずにコードを簡略化・リファクタリングするためのスキル。可読性、保守性、テスト容易性、複雑度、効率性、技術的負債の指摘を実装に反映し、変更範囲と検証結果を簡潔に報告するときに使う。ユーザーが quality report や code-quality-review の結果をもとに改善してほしいと言った場合も使う。
---

# Code Refactor

レビュー指摘や既存差分を起点に、公開契約と挙動を維持して、読み・変更・テストしやすくする。quality report があれば findings と quality gate を優先順位にし、なければ対象から改善候補を抽出する。必要以上の設計変更は避ける。

## 入力と手順

対象はレビュー結果、quality report、差分、ファイル・関数、制約・検証方法。安全な範囲を決められないときだけ質問する。

1. 現状の挙動、公開契約、関連テストを把握する。
2. 指摘を可読性、保守性、テスト容易性、複雑度、効率性、検証不足に整理し、`code-next-developer-review` と `code-quality-review` ごとに subagent を起動して、観点を分けたレビューを集める。
3. 変更ごとに分類を選び、対応する reference を読んで、小さく戻しやすい単位で実施する。
4. 既存テストまたは適切な検証を実行し、できなければ理由と残リスクを報告する。

## 分類

`Simplify`、`Split`、`Extract`、`Inline`、`Rename`、`Move`、`Encapsulate`、`Generalize`、`Specialize`、`Replace`、`Remove`、`Reorganize` のいずれかを各変更に選ぶ。主目的を一つにし、実際に使う分類の `references/<分類小文字>.md` だけを読む。迷えば局所的で挙動変更リスクの低い分類を優先し、reference の例は既存設計・命名・テスト構造に合わせる。

## 報告

変更内容、分類ごとの実施内容、変更ファイルと主要行、検証結果、改善点、quality report/gate への対応、残リスクまたは未実行検証を簡潔に示す。未使用の分類は書かない。
