---
name: code-refactor
description: コードレビュー結果、code quality report、git diff、PR 差分、指定ファイルをもとに、挙動を変えずにコードを簡略化・リファクタリングするためのスキル。可読性、保守性、テスト容易性、複雑度、効率性、技術的負債の指摘を実装に反映し、変更範囲と検証結果を簡潔に報告するときに使う。ユーザーが quality report や code-quality-review の結果をもとに改善してほしいと言った場合も使う。
---

# Code Refactor

## 目的

レビュー指摘や既存差分を起点に、挙動を維持したままコードを簡略化する。
読みやすさ、変更しやすさ、テストしやすさを上げることを優先し、必要以上に大きな設計変更は避ける。
code quality report がある場合は、品質上の findings と quality gate をリファクタリング対象の優先順位として扱う。
レビュー指摘がない場合は、可読性、保守性、テスト容易性、複雑度、効率性の観点で指摘を抽出し、反映する。

## 入力

次のいずれかを受け取る。

- レビュー対象
- コードレビュー結果
- code quality report、quality gate、scorecard
- `git diff`、`gh pr diff`、コミットハッシュなどの差分
- 対象ファイル、ディレクトリ、関数名
- 期待する制約、検証コマンド、避けたい変更
- ループ回数。指定がなければ最大 3 回まで、レビュー、リファクタ、検証を繰り返す。

入力が曖昧で変更範囲を安全に決められない場合だけ、最小限の質問をする。

## 進め方

1. 対象差分と関連コードを読み、現状の挙動、公開契約、テストを把握する。
2. 既存のレビュー指摘や quality report があれば、可読性、保守性、テスト容易性、複雑度、効率性、検証不足に分類する。
3. 必要に応じて `code-review`, `code-next-developer-review`, `code-quality-review` の観点を使い、差分中心に問題を洗い出す。
   - サブエージェントが利用でき、変更範囲が大きい場合は、観点ごとにサブエージェントを起動してレビューを集める。
   - 小さな変更、またはサブエージェントを使えない環境では、自分で同じ観点を確認する。
4. 変更ごとに「リファクタ分類」から少なくとも 1 つのカテゴリを選び、必要な reference を読む。
5. 挙動維持を前提に、選んだ分類の目的に合う小さく戻しやすい単位からリファクタリングする。
6. 既存テストまたは適切な検証コマンドを実行する。実行できない場合は理由と残リスクを明記する。
7. レビュー指摘や quality gate が残る場合は、指定されたループ回数を上限に 3 から 6 を繰り返す。指定がなければ最大 3 回まで繰り返す。
8. 変更したファイル、分類ごとの要約、検証結果、残るリスクを報告する。

## リファクタ分類

リファクタでは、次の分類のいずれかを必ず使う。1 つの変更に複数カテゴリが当てはまる場合は、主目的のカテゴリを 1 つ選び、補助的なカテゴリを必要に応じて併記する。

| カテゴリ | 内容 | 詳細 |
|---|---|---|
| `Simplify` | 複雑さを減らす | [references/simplify.md](references/simplify.md) |
| `Split` | 責務を分割する | [references/split.md](references/split.md) |
| `Extract` | 共通化・切り出し | [references/extract.md](references/extract.md) |
| `Inline` | 過剰な抽象化を戻す | [references/inline.md](references/inline.md) |
| `Rename` | 名前を改善する | [references/rename.md](references/rename.md) |
| `Move` | 適切な場所へ移動する | [references/move.md](references/move.md) |
| `Encapsulate` | 隠蔽する | [references/encapsulate.md](references/encapsulate.md) |
| `Generalize` | 共通化・抽象化する | [references/generalize.md](references/generalize.md) |
| `Specialize` | 抽象化をやめる | [references/specialize.md](references/specialize.md) |
| `Replace` | より良い実装へ置き換える | [references/replace.md](references/replace.md) |
| `Remove` | 不要なものを削除する | [references/remove.md](references/remove.md) |
| `Reorganize` | 構造を整理する | [references/reorganize.md](references/reorganize.md) |

### Reference の読み方

- 変更前に、採用するカテゴリの reference を読む。
- 迷った場合は、より局所的で挙動変更リスクが低いカテゴリを優先する。
- 複数カテゴリを使う場合は、実際に編集するカテゴリの reference だけ読む。
- reference の例をそのまま当てはめず、既存コードの設計、命名、テスト構造に合わせる。

## 出力

最終報告は次の順で簡潔に書く。実施したリファクタは分類ごとにまとめ、未使用カテゴリは書かない。

1. 変更した内容
2. リファクタ分類ごとにやったこと
3. 変更ファイルと主要な行
4. 検証結果
5. 改善された点、Benefit
6. quality report / quality gate に対する対応状況
7. 残るリスク、または未実行の検証

### 出力例

```text
変更した内容:
...

リファクタ分類:
- Simplify: 早期 return に変更し、ネストを 1 段減らした。
- Extract: 重複していた日付変換を小さな関数に切り出した。
- Rename: 責務が伝わるように一時変数名を変更した。

検証:
...
```
