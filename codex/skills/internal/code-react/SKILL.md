---
name: code-react
description: React / JSX / TSX / .jsx / .tsx 実装を書く前に、コンポーネント設計、Atomic Design を補助観点にした UI 分割、hooks、props、state、UI 実装の責務分離を確認するときに使う。
---

# Code React

## 目的

React 実装で、コンポーネント境界、状態、表示ロジックを整理して既存 UI に馴染む判断をするための参照 skill。

## 使い方

- まず [../code-general/SKILL.md](../code-general/SKILL.md) を読む。
- React 実装では references/react.md を必ず読む。
  - このファイルは、この SKILL.md と同じディレクトリの references/ 配下にある。
  - 見つからない場合は ~/.claude/skills/code-react/references/react.md を試す。
  - それでも見つからない場合はユーザーに報告し、references なしで進めない。
- TypeScript の型設計も必要な場合は code-ts の references/lang.md も必ず読む。
  - このファイルは ~/.claude/skills/code-ts/references/lang.md にある。
  - 見つからない場合はユーザーに報告し、references なしで進めない。

## component-design との使い分け

この skill は、React 実装中に component 境界、hooks / props / state、Atomic Design の補助観点、責務分離を確認するために使う。
画面構成、大きな JSX 分割、component tree、実装順序を先に整理する必要がある場合は、`component-design` を先に使う。

## コンポーネント分割

Atomic Design は、UI を小さな部品から組み立てるための補助観点として使う。
ただし `atoms` / `molecules` / `organisms` の分類を目的化しない。
`atoms` / `molecules` は design system 層として扱い、ドメイン知識を持たせない。
この層では design token、variant、layout、interaction state など、見た目と基本操作だけを表現する。
業務ルール、API 由来のデータ構造、ユーザーや注文などのドメイン語彙が入る component は feature / domain 側へ置く。

優先する判断軸:

- state と表示を分離できているか。
- ドメイン知識を持つ component と汎用 UI component が混ざっていないか。
- `atoms` / `molecules` が design system の見た目、状態、操作だけを表現しているか。
- props が呼び出し側にとって自然で安定しているか。
- feature / page に閉じるべき component を共通化しすぎていないか。
- 共通化は、具体的な重複と再利用先が見えてから行う。
