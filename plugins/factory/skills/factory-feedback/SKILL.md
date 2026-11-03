---
name: factory-feedback
description: factory の実行で見つかった問題や改善案を、共通 factory スキルの改善 feedback として version-1/agents-docs の GitHub Issue にする。プロジェクト固有の運用ルールの保存には factory-recap を使う。
---

# Factory Feedback

factory の手順・役割契約・判断基準の改善に必要な根拠を整理し、`version-1/agents-docs` に Issue を作成する。利用先プロジェクトの実装 Issue とは分ける。

## 根拠を整理する

- ユーザーの指摘、会話、参照できる管理 Issue・PR・実行結果から、発生条件、期待した振る舞い、実際の振る舞い、影響を整理する。観測事実と推測を分け、未確認事項を断定しない。
- 関連する factory の本文・references を読み、改善対象の手順や契約を特定する。実行側の逸脱か、スキルの不足・曖昧さかを区別する。対策案は受け入れ条件とともに示し、根拠が弱ければ検討案として残す。
- プロジェクトだけに適用すべき知見は `factory-recap` の対象として分ける。共通改善とローカル改善が混在する場合も、Issue には共通改善に必要な内容だけを含める。

## Issue を作成する

1. `gh issue list --repo version-1/agents-docs --state all --search <検索語>` などで関連 Issue を確認する。重複なら既存 URL と差分を報告し、同じ Issue を再作成しない。既存 Issue の編集・コメントは別途依頼がある場合に行う。
2. タイトルに factory の問題と改善対象を表す。本文には背景・再現条件、期待と実際、根拠、改善案、受け入れ条件を必要な範囲で記載する。参照不能な情報はその旨を示す。
3. 送信内容から認証情報、個人情報、非公開のコード・ログ・URL を除き、共有可能な最小の再現例にする。非公開情報の共有が必要なら、その内容と送信先を提示して確認する。
4. `gh repo view version-1/agents-docs` で宛先を確認する。現在の remote に依存せず、`gh issue create --repo version-1/agents-docs --title <タイトル> --body-file <本文ファイル>` を使う。呼び出しは Issue 作成の依頼として扱うが、実行環境の確認・外部送信規則を守る。実行対象の指定を意味する `factory` ラベルや `Factory goal:` は自動付与しない。
5. 成功時は作成された URL と要点を報告する。失敗時は草案と理由を残す。応答が不明な場合は作成済みか確認してから再試行し、重複投稿しない。

このスキルでは factory 本体の修正、実装開始、PR 作成を行わない。
