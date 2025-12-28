---
name: language-ruby
description: Ruby 言語でのコーディングに関するベストプラクティスとスタイルガイドを提供します。
---

# Ruby コーディングガイドライン

- コードの可読性と一貫性を最優先に考慮する
- Ruby の慣習に従い、コミュニティで広く受け入れられているスタイルを採用する
- 意味のある変数名とメソッド名を使用し、コードの意図を明確にする
- 適切なインデント（2スペース）と空白を使用して、コードの構造を明確にする
- メソッドは単一の責務を持つように設計し、、長すぎないようにする
- コメントは必要最低限にとどめ、コード自体が意図を伝えるようにする
- エラーハンドリングは適切に行い、例外が発生した場合の挙動を明確にする
- テストコードは RSpec を使用し、コードの品質を保証する
- パフォーマンスを考慮し、不要な計算やメモリ使用を避ける
- セキュリティベストプラクティスに従い、脆弱性を防止する
- メタプログラミングは最低限にする

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
