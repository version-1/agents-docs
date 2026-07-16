# Remove

## 目的

不要なコード、設定、props、CSS、依存を削除し、保守対象を減らす。

## 使う場面

- dead code、未使用 export、未使用 props がある。
- 古い分岐や feature flag が不要になっている。
- CSS、設定、テスト helper が使われていない。

## 代表的な技法

- dead code の削除。
- 未使用 import / export / props / state の削除。
- 不要な CSS class や style の削除。
- 使われていない dependency や config の削除。

## 注意点

- 動的参照、reflection、template、設定ファイルから使われていないか確認する。
- public API の削除は互換性に影響する。
- テストだけが使っている helper は、テスト意図も見直す。

## 報告観点

- 削除した対象と未使用と判断した根拠。
- 動的参照や公開契約の確認範囲。
- 削除後に実行した検証。
