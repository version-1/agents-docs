
# React コーディングガイドライン

- コンポーネントは小さく、単一の責任を持つように設計する。
- データと振る舞いを分離する。
- JSX 内でのロジックは最小限に抑え、必要な場合は関数やフックを使用して分割する。
- スタイルは CSS-in-JS ライブラリや CSS モジュールを使用して、コンポーネントごとに管理する。
- 状態管理には、必要に応じて React の Context API などのライブラリを使用する。
- イベントハンドラーは、明確な命名規則を使用して、何をするのかがわかるようにする。
- コンポーネントのプロパティには、TypeScript の型定義を使用して、型安全性を確保する。
- コンポーネントの再利用性を高めるために、必要に応じてジェネリクスを使用する。
- コンポーネントのライフサイクルを理解し、適切なフック（useEffect、useMemo、useCallback など）を使用してパフォーマンスを最適化する。
- コンポーネントのテストには、React Testing Library や Vitest を使用して、ユーザーの視点からのテストを行う。
- コンポーネントの命名は、PascalCase を使用し、意味のある名前を付ける。
- コンポーネントの構造は、論理的なグループに分け、必要に応じてフォルダを使用して整理する。
- コンポーネントの状態は、必要に応じてローカル状態とグローバル状態を適切に使い分ける。

## 非推奨API

- クラスコンポーネントは、React Hooks の導入以降は非推奨とされているため、関数コンポーネントを使用することが推奨される。
- `componentWillMount`、`componentWillReceiveProps`、`componentWillUpdate`などのライフサイクルメソッドは、React 16.3 以降は非推奨とされているため、使用しないことが推奨される。

# React アーキテクチャ

コンポーネントを実装する際は、ドメイン知識を持つ feature / domain 側の component と、
見た目と基本操作だけを表現する design system 側の component を分けることが推奨される。

## Core Principles

- まずは画面要件を満たす実装を優先する
- 同じ UI パターンが 2 回以上出たら共通化を検討する
- 3 回以上使われる、または今後増えそうなら共通コンポーネント化する
- 見た目だけでなく、振る舞い・状態・アクセシビリティも含めて共通化する
- ドメイン固有の意味を持つものは `features/` や `components/domain/` に置く
- 汎用 UI は `components/ui/` に置く
- props が複雑になりすぎる共通化は避ける

## ディレクトリ構成の例:

ディクレクトリ構成の例を以下に示す。
src/components/ui/ 以下に design system 側のコンポーネントを配置し、src/features/ 以下にドメイン知識を持つ feature / domain 側のコンポーネントを配置する。
これは一例であり、プロジェクトの規模や要件に応じて柔軟に変更することが推奨される。
ただし、ドメイン知識を持つ feature / domain 側のコンポーネントと、見た目と基本操作だけを表現する design system 側のコンポーネントを分けることは強く推奨される。

```
src/
├── components/              // 共通コンポーネント
│   └── ui/                  // design system 側のコンポーネント
│       ├── button/index.tsx
│       ├── dialog/index.tsx
│       ├── card/index.tsx
│       └── badge/index.tsx
│
├── features/                // ドメイン知識を持つ feature / domain 側のコンポーネント
│   ├── projects/
│   │   ├── components/
│   │   │   ├── card/index.tsx
│   │   │   ├── list/index.tsx
│   │   │   └── badge/index.tsx
│   │   ├── hooks/
│   │   └── api/
│   │
│   └── tasks/
│       ├── components/
│       └── hooks/
│
└── app/
    └── projects/
        └── page.tsx
```

### コンポーネントのディレクトリ設計について

Table コンポーネントを例に、コンポーネントのディレクトリ設計の一例を示す。
このように配置してディレクトリからコンポーネントの依存関係・責務を明確にすることが推奨される。

```tree
table/                         # 「Table」という UI の責務単位でコードを集約
├── index.tsx                 # 外部公開される Table コンポーネント
├── index.module.css          # Table 全体のスタイルを管理
├── header/                   # Header(子コンポーネント) に関連する実装をコロケーション
│   ├── index.tsx             # Header の描画責務
│   └── index.module.css      # Header のスタイルを局所化
└── content/                  # Content に関連する実装をコロケーション
    └── index.tsx             # Content の描画責務
```

- 一コンポーネント、一ファイルの原則を守ることが推奨される。
- 親コンポーネント以外から参照されない閉じたコンポーネントは、親コンポーネントと同じディレクトリに配置することが推奨される。
- table/index.tsx, table/index.module.css Table のように、コンポーネントの名前と同じファイル名を使用することが推奨される。


