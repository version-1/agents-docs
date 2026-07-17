# Replace

## 目的

同じ挙動を、より適切な構造や標準 API に置き換える。

## 使う場面

- if / switch が増え続け、分岐ごとの責務が独立している。
- 手書き処理を標準ライブラリや既存 helper で安全に置き換えられる。
- Map、polymorphism、table driven logic で意図が明確になる。

## 代表的な技法

- switch を Map / lookup table に置き換える。
- if 分岐を polymorphism や strategy に置き換える。
- 手書き parsing / formatting を標準 API に置き換える。
- mutable な処理をより安全な immutable 更新に置き換える。

## 注意点

- 置き換え前後で edge case、順序、エラー、nil / null の扱いが同じか確認する。
- polymorphism は分岐が増え続ける根拠がある場合に使う。
- 標準 API への置き換えでは locale、timezone、丸め、encoding の差に注意する。

## 報告観点

- 何を何に置き換えたか。
- 置き換えで明確になった責務やデータ構造。
- edge case の確認結果。
