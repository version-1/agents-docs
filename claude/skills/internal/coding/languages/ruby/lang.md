
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
