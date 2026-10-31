# DB定義・バージョン管理

大目標: 定期的に、設計・レビュー時の認識と実コードの振る舞いのズレを見つけ、認識を修正する。その変化を再確認できるよう、問題の版と当時の回答を保持する。

## Version 1

SQLの正本は [001.sql](../scripts/migrations/001.sql)。`PRAGMA user_version` を適用済み番号とし、ランタイムの `VERSION` と対応する。以下の列はSQLの制約・型と併せて参照する。

| テーブル | 列 | 意味 |
|---|---|---|
| repositories | id, identity | git common directoryの実パス。同一repoのworktreeで共有 |
| questions | id, repository_id | 問題ID、所有repoへのFK |
| questions | content | 問題JSON。選択肢・単一正解・解説・scope・根拠ファイルのSHA-256を含む |
| questions | fingerprint | 根拠以外の内容のSHA-256。同一repoで一意 |
| questions | replaces | 置き換え元questions.idへのFK。新版の初期進捗は独立 |
| questions | status, created_at | active／deleted／superseded、UTC作成時刻 |
| progress | question_id | questions.idへのFK、1対1 |
| progress | presentations | 保存ありで提示した累計。未回答も含む |
| progress | streak, mastered | 連続正解0〜3、習得フラグ |
| progress | due_at | UTCの次回復習日時。NULLは新規・リセット直後または習得済み |
| progress | generation | resetで加算する世代。リセット前の提示からの更新を拒否 |
| attempts | id, question_id | ランダム提示ID、questions.idへのFK |
| attempts | presented_at, ordinal, generation | UTC出題日時、その問題で何回目か、提示時の世代 |
| attempts | snapshot, choices | 提示時の全問題JSON、A〜Dに対応する選択肢ID配列 |
| attempts | answered_at, selected, correct, feedback | UTC回答日時、選択したID、正誤、採点時の解説・進捗JSON。未回答はすべてNULL |

各操作はパラメータ化SQLを利用する。提示・回答・新版登録は `BEGIN IMMEDIATE` の単一トランザクションで更新する。FKを有効化し、問題の物理削除はCLIに提供しない。練習はattempts・progressを更新しない。

## 復習規則

初回正解は回答日時+24時間、2連続正解は+7日、3連続正解で習得。不正解は連続正解0・回答日時+24時間。習得は通常選定から除外しresetで戻せる。これは忘却曲線を参考にした固定間隔であり、個人ごとの最適間隔を推定するものではない。

## 変更手順

1. 適用済みSQLは編集せず `002.sql` 以降の連番を追加し、この定義とランタイムの `VERSION` を更新。
2. 旧版→最新版の移行、失敗時ロールバック、未来バージョン拒否、既存履歴保持をテスト。
3. 初期の空DBを除き、移行前にSQLite backup APIで `quiz.sqlite3.backup-<UUID>` を同じディレクトリに保存。WALの確定データも含める。バックアップ失敗時は移行しない。
4. すべての未適用SQLを単一トランザクションで実行。SQL欠落・エラーは停止し旧版を保持。既知の番号より新しいDB、バージョン0で既存テーブルがあるDBには書き込まない。

バックアップは自動削除しない。復旧は実行中のクイズを停止してからバックアップの別コピーを `--db` で検証し、元DBの置換はユーザーが明示的に判断する。
