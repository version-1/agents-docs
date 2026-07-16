# Reorganize

## 目的

ディレクトリ、依存関係、ファイル構成を整理し、コードベース全体の見通しをよくする。

## 使う場面

- 関連ファイルが散らばっていて変更箇所を見つけにくい。
- 依存方向が読み取りにくい。
- ディレクトリ構成が現在の責務やドメインと合っていない。

## 代表的な技法

- feature / domain / layer 単位での配置整理。
- barrel export や index file の整理。
- 循環依存の解消。
- テスト、fixture、mock の配置整理。
- ドキュメントや README の参照先整理。

## 注意点

- 大規模な移動はレビュー負荷が高い。挙動変更と構造変更を混ぜすぎない。
- path alias、build config、test runner、storybook などの参照を確認する。
- Reorganize だけで責務が改善しない場合は、Split や Move と組み合わせる。

## 報告観点

- 整理した構造と狙い。
- 解消した依存関係や発見しやすさの問題。
- build / test / import 解決の確認結果。
