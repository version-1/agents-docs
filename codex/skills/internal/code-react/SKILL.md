---
name: code-react
description: React / JSX / TSX / .jsx / .tsx 実装を書く前に、コンポーネント設計、Atomic Design を補助観点にした UI 分割、hooks、props、state、UI 実装の責務分離を確認するときに使う。
---

# Code React

## 目的

React 実装で、コンポーネント境界、状態、表示ロジックを整理して既存 UI に馴染む判断をするための参照 skill。

## 使い方

- まず [../code-general/SKILL.md](../code-general/SKILL.md) を読む。
- React 実装では [references/react.md](references/react.md) を読む。
- TypeScript の型設計も必要な場合は [../code-ts/references/lang.md](../code-ts/references/lang.md) も読む。

## コンポーネント分割

Atomic Design は、UI を小さな部品から組み立てるための補助観点として使う。
ただし `atoms` / `molecules` / `organisms` の分類を目的化しない。

優先する判断軸:

- state と表示を分離できているか。
- ドメイン知識を持つ component と汎用 UI component が混ざっていないか。
- props が呼び出し側にとって自然で安定しているか。
- feature / page に閉じるべき component を共通化しすぎていないか。
- 共通化は、具体的な重複と再利用先が見えてから行う。
