---
name: lego-programming
description: コーディングをする時のマインドセットを提供する。
---

コーディングを行いながら設計をする際に
- レゴ(レゴの定義を参照する)を作ること
- レゴを組み合わせてシステムを作ること
を意識しながらコーディングを行ってください。

# レゴの定義

- モジュール
- コンポーネント
- ライブラリなど

必要な機能を満たした高凝集・疎結合なモジュールを表します。

## 例

## バックエンド

- Controller
- UseCase / Service
- Repository
- Domain Model / Entity
- Value Object
- Factory
- Domain Service
- Policy
- Specification
- Validator
- Mapper
- DTO
- Serializer
- Adapter
- Client（API Client、SDK）
- Gateway
- Middleware
- Event Handler
- Message Consumer / Producer
- Utility（責務が明確なもの）

## フロントエンド

- Atoms
- Molecules
- Organisms
- Template
- UI Component
- Custom Hook
- Store
- Context
- API Client
- Repository
- Presenter
- ViewModel
- Form
- Layout
- Theme

## 共通

- モジュール
- コンポーネント
- パッケージ
- ライブラリ
- SDK
- プラグイン
- CLI
- Workflow
- ジョブ

# レゴの境界

- クリーンアーキテクチャ
- レイヤードアーキテクチャ
- DDD
- Atomic Design
- デザインパターン
などのアーキテクチャのベストプラクティスをもとに境界を作ります。

# 作業手順

1. レゴの境界を作成して、定義する
   - レイヤリング、コンポーネント分割、 依存関係の方向性、責務の分離などを意識する
   - インタフェースを定義すると同義
2. レゴを作る
   - レゴの境界に沿って、必要な機能を満たす高凝集・疎結合なモジュールを作る
3. レゴを組み合わせてシステムを作る
   - 作成したレゴを組み合わせて、必要な機能を満たすシステムを作る
4. レゴの境界を見直す
   - 作成したシステムを見直して、レゴの境界が適切かどうかを確認する
   - 必要に応じて、レゴの境界を修正する

1. が明快に決まる場合には、 2. 3. は並列に作業してください。


# 例

## 例 1

### NG: ページごとにコンポーネントを切る設計

```
- pages/
    - HomePage/
        - HomePage.tsx
        - HomePageHeader.tsx
        - HomePageFooter.tsx
    - AboutPage/
        - AboutPage.tsx
        - AboutPageHeader.tsx
        - AboutPageFooter.tsx
```

### OK: Atomic デザインを意識して、レゴを作り配る設計

```
- components/
    - atoms/
        - Button.tsx
        - Input.tsx
    - molecules/
        - Form.tsx
        - Card.tsx
    - organisms/
        - Header.tsx
        - Footer.tsx
```

## 例 2

### NG: 1つのレゴに複数の責務を書く

```text
UserController
├─ HTTP処理
├─ バリデーション
├─ ビジネスロジック
├─ DBアクセス
└─ メール送信
```

❌ 1つのモジュールが複数の責務を持っている。

---

### OK: レゴごとに責務を分離する

```text
UserController
        │
        ▼
CreateUserUseCase
    ├── UserRepository
    └── Mailer
```

```text
presentation/
└── UserController

application/
└── CreateUserUseCase

domain/
└── User

infrastructure/
├── UserRepository
└── Mailer
```

✅ 各レゴは1つの責務だけを持ち、インタフェースを介して組み合わせる。
- Controller: HTTPの責務
- UseCase: ビジネスロジック
- Repository: データアクセス
- Mailer: 外部サービス連携

# 完了条件

以下をすべて満たしたら完了とする。

## 1. レゴの境界が定義されている

- 実装前または実装初期に、主要なモジュール・コンポーネント・レイヤーの責務が説明されている
- 各レゴの入力・出力・公開インタフェースが明確である
- 依存関係の方向が一貫している
- UI、ビジネスロジック、データ取得、外部サービス連携などの責務が混ざっていない

## 2. 高凝集・疎結合になっている

- 1つのレゴが1つの明確な責務を持っている
- 他のレゴの内部実装に依存していない
- 再利用可能な処理がページ固有・画面固有の場所に閉じ込められていない
- 変更理由が異なる処理が同じファイル・同じ関数・同じコンポーネントに混在していない

## 3. レゴを組み合わせてシステムが成立している

- 小さなレゴを組み合わせて、要求された機能が動作している
- 上位のレゴは下位のレゴを組み合わせるだけで、細かい実装詳細を持ちすぎていない
- 新しい機能追加時に、既存レゴの再利用・拡張で対応できる構造になっている

## 4. 境界の見直しが行われている

- 実装後に、レゴの責務・名前・配置・依存関係を見直している
- 不自然に肥大化したレゴがない
- 重複した責務を持つレゴがない
- 必要に応じて分割・統合・命名変更を行っている

## 5. 説明可能である

- なぜその境界でレゴを分けたのか説明できる
- どのレゴを再利用できるか説明できる
- 今後の変更に対して、どこを変更すればよいか説明できる
