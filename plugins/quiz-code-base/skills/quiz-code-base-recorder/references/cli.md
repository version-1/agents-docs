# CLI・問題入力契約

大目標は、定期的に設計・レビュー時の認識と実コードの振る舞いのズレを発見し、認識を修正すること。AIは意味の解析・作問・解説、CLIは永続化・選定・採点・日程計算を担当する。

## 呼び出し

`python3 <recorder>/scripts/quiz.py --repo /absolute/repository --db /absolute/quiz.sqlite3 COMMAND`

`--db` は省略可能。既定は `~/.local/share/quiz-code-base/quiz.sqlite3`。グローバルオプションはコマンドより前に置く。標準出力はJSON、エラーは標準エラーのJSONと終了コード1。SQLの直接実行は不要。

| コマンド | 入力 | 結果 |
|---|---|---|
| `init` | なし | DB作成・移行とスキーマ番号 |
| `add` | `--input question.json` | ID、作成有無。同じ内容は重複登録しない |
| `select` | `--scope src/foo --scope db --count 5 [--new-only] [--exclude 1 2]` | 選定問題、`needs_review` のID、不足数 |
| `present` | `--id 1 [--practice]` | 提示票・問題文・並べ替え済みA〜D・回数表示 |
| `answer` | `--presentation presentation.json --choice A` | 採点・解説・次回予定 |
| `list` | `[--all]` | 問題の要約と状態（正解を含まない） |
| `inspect` | `--id 1` | 正解と解説を含む全内容 |
| `history` | `[--id 1] [--limit 20]` | 新しい順の出題・回答履歴と当時の問題 |
| `refresh` | `--id 1` | 意味に影響しない変更の根拠ハッシュ再採取 |
| `revise` | `--id 1 --input revised.json` | 旧版を保管して新版を登録 |
| `delete` | `--id 1` | 論理削除 |
| `reset` | `--id 1` | 習得・削除を解除し連続正解0、即時復習対象 |

`--scope` はリポジトリ相対のファイルまたはディレクトリ。`.` は全体。複数指定はOR。問題側のscopeは出題テーマが属する具体的なファイルを指定する。選定はscopeの完全一致またはディレクトリ配下で照合するため、機能名は事前にパスへ解決する。絶対パス・`..` を含むパスは拒否する。

選定は、期限到来の既出問題→未出題問題の順で、予定日時・IDによって決定する。`--new-only` は記録上の出題回数0だけを返す。根拠変更のある問題は選定しない。`needs_review` は最大要求数まで返すので、再確認後に再選定する。1回の学習で同じIDを出さないよう `--exclude` を渡す。

## 問題JSON

```json
{
  "scope": ["src/models/order.py"],
  "prompt": "この更新経路では、注文削除時に明細の削除を何が保証するか？",
  "choices": [
    {"id": "db", "text": "DBの外部キー制約", "explanation": "外部キーのON DELETE CASCADEが適用される。"},
    {"id": "model", "text": "モデルの削除フック", "explanation": "この経路ではモデルフックを通らない。"},
    {"id": "api", "text": "APIのバリデーション", "explanation": "削除連動を保証する検証はAPIにない。"},
    {"id": "none", "text": "連動せず明細が残る", "explanation": "DBの削除連動があるため残らない。"}
  ],
  "correct_id": "db",
  "explanation": "削除の保証はアプリのフックではなくDB制約にある。",
  "review_point": "モデルを迂回する削除経路もDB制約で保証されるか確認する。",
  "evidence": [
    {"path": "db/schema.sql", "locator": "order_items.order_id の外部キー定義"},
    {"path": "src/models/order.py", "locator": "delete_order"}
  ],
  "documentation_suggestion": ""
}
```

例の内容を対象コード未確認で流用しない。全選択肢の説明は必須。根拠ファイルは実在し、解決後もリポジトリ内にある必要がある。SHA-256はCLIが自動付加する。文書との不一致があるときは `documentation_suggestion` に箇所と整合の提案を書き、比較した文書も根拠に含める。全根拠のハッシュを再利用時に確認する。新規のファイル追加だけでは検知できない依存関係もあるので、生成時に制約や呼び出し経路の集約ファイルも根拠に含める。

## 提示・回答

`present` の出力JSONを一時ファイルに保存し、その `prompt`、`choices`、`display` だけをユーザーへ提示する。`presentation` は内部の提示票でありユーザーに回答用として編集させない。回答後に同じファイルを `answer` へ渡す。ファイルを新たに組み立て直さない。回答処理後は一時ファイルを保持する必要はない。

通常の提示票はDB上の出題と照合し、保存した順序で採点する。二重回答・より新しい出題がある古い回答・リセット前の回答・根拠変更後の回答は成績を更新せず拒否する。練習用の提示票はDBに保存せず、採点のみ行う。練習は初出判定にも影響しない。同じ問題の練習回答を繰り返しても履歴は増えない。

同じ内容の `add` は非supersededの既存IDを返すが、削除・習得・根拠変更を解除しない。必要なら明示的に再確認・resetする。`revise` は他の非superseded問題と重複したらエラーとし旧版を変更しない。履歴版と同じ内容への回帰、根拠ファイルだけの移動も新版登録できる。全内容・根拠が同一のreviseは拒否する。refreshは意味が同じ場合だけ利用し、旧出題スナップショットを書き換えない。

これは単一ローカルユーザーの学習ツールであり、試験の不正防止基盤ではない。回答前の正解情報をチャットに漏らさない責務は対話側にある。
