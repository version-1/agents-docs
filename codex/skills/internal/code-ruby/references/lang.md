
# Ruby 言語 コーディングガイドライン

- 新たな実装にメタプログラミングを使用しない 
  - 例外的なケースを除き、メタプログラミングはコードの可読性を損なう可能性があるため、使用を避けるべきです。
  - コードの可読性を優先するために、明示的なコードを書くことが推奨されます。
- 早帰リターンを使用する
    - 条件が満たされない場合は、早期にリターンすることで、コードのネストを減らし、可読性を向上させることができます。
    - 例: 
        ```ruby
        def example_method(value)
        return unless value.valid?
        # 続きの処理
        end
        ```
# Ruby on Rails コーディングガイドライン

- コントローラーは薄く保つ
    - コントローラーはリクエストの処理とレスポンスの生成に集中し、ビジネスロジックはモデルやサービスオブジェクトに移すべきです。
    - これにより、コードの再利用性が向上し、テストが容易になります。
- モデルは単一責任を持つべき
    - モデルはデータの管理とビジネスロジックの一部を担当しますが、過度に複雑なモデルは避けるべきです。
    - 複雑なビジネスロジックはサービスオブジェクトに移すことで、モデルの責任を明確に保ち、コードの可読性と保守性を向上させることができます。
- ビューはロジックを最小限に保つ
    - ビューはデータの表示に集中し、複雑なロジックはヘルパーメソッドやプレゼンターオブジェクトに移すべきです。
    - これにより、ビューのコードがシンプルになり、保守性が向上します。
- コールバックをスキップするメソッドに注意
    - コールバックをスキップするメソッド（例: `update_column`）は、コールバックやバリデーションを無視するため、非推奨です。
    - 使用する際は、人間に確認を求めてください。
    - 注意深く検討し、データの整合性が損なわれないようにする必要があります。
- ActiveRecord クエリの最適化
    - N+1 クエリを避けるために、`includes` や `eager_load` を使用して関連データを事前にロードすることが推奨されます。
    - クエリのパフォーマンスを向上させるために、必要なデータのみを選択するように心がけるべきです。
    - 例: 
        ```ruby
        # N+1 クエリの例
        users = User.all
        users.each do |user|
        puts user.posts.count
        end

        # 最適化されたクエリの例
        users = User.includes(:posts).all
        users.each do |user|
        puts user.posts.size
        end
        ```


# 推奨ライブラリ

- Web フレームワーク: Ruby on Rails
- テストフレームワーク: RSpec
- テストデータ生成: FactoryBot
- 静的コード解析: RuboCop
- ドキュメンテーション: YARD
- タスクランナー: Rake
- デバッグ: Byebug
- HTTP クライアント: Faraday
- データベースマイグレーション: ActiveRecord Migrations
- 環境変数管理: dotenv
- ロギング: Logger (標準ライブラリ) または Lograge
- 認証: Devise
- キャッシュ: Redis (redis-rb)
- ジョブキュー: Sidekiq
- バックグラウンドジョブ: ActiveJob (Rails 標準)
- JSON 処理: JSON (標準ライブラリ)
- 日付・時間操作: ActiveSupport::TimeWithZone (Rails 標準)
- ファイルアップロード: CarrierWave または ActiveStorage (Rails 標準)
- メール送信: ActionMailer (Rails 標準) または Mailgun
- 認可: Pundit または CanCanCan
- API ドキュメンテーション: Swagger UI または rswag
- パフォーマンス監視: NewRelic または ScoutAPM
- エラートラッキング: Sentry または Rollbar
- バージョン管理: Git
