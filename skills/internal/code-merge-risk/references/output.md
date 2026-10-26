# 出力契約 v1.0

Markdown レポートの後に、同一の結果を JSON コードブロックで出力する。スコア・状態・根拠を一致させる。JSON は以下のキーをすべて含め、列挙値を翻訳しない。根拠は対象 SHA のファイル・行へのリンク、検証ログ、GitHub の取得結果など、追跡できる位置を示す。

## Markdown の順序

1. 総合点（保留時は「未確定」）、判定、PR URL、head/base SHA、確認時刻。
2. 判断の要約。
3. 6項目の点数と調査根拠の表。未評価・評価途中の項目は「未確定」。
4. リスク詳細：ID、項目、障害シナリオ、発生条件、根拠、防御策、深刻さ、発生の見込み、点数、対処案。
5. 実行した検証、未実行の検証、情報不足、必須条件の状態と根拠。

## JSON の型

- `schema_version`: 固定文字列 `"1.0"`。
- `pr_url`: PR URL の文字列。
- `head_sha`, `base_sha`: 評価対象の完全 SHA。取得不能の場合のみ null とし hold にする。
- `evaluated_at`: 最終確認時の UTC 時刻（ISO 8601）。
- `assessment_status`: `complete` または `hold`。
- `overall_score`: 0〜100 の整数。hold の場合は必ず null。
- `decision`: `auto_merge_candidate` / `human_review` / `hold`。
- `summary`: 判断要約の文字列。
- `categories`: 以下の6 IDを各1件含む配列。
  - `id`: correctness / security / data / performance / compatibility / operations。
  - `score`: 完了項目は整数、未完了項目は null。
  - `rationale`: 調査範囲と項目点の理由。
  - `evidence`: 根拠を示す文字列配列。
  - `risk_ids`: 関連するリスク ID の配列。リスクなしは空配列。
- `risks`: リスクオブジェクトの配列。特定なしは空配列。
  - `id`, `category`, `scenario`, `trigger`: 文字列。category は上記6 IDのいずれか。
  - `evidence`, `mitigations`: 根拠と確認済み防御策の文字列配列。
  - `severity`: minor / moderate / major / critical。
  - `likelihood`: low / plausible / high。
  - `score`: 採点表の整数。
  - `recommendation`: 対処案の文字列。
- `checks`: 検証オブジェクトの配列。
  - `name`, `target_sha`, `evidence`: 文字列。target_sha は不明なら null。
  - `status`: pass / fail / not_run。
  - `reason`: 結果または未実行理由と評価への影響の文字列。
- `unknowns`: 情報不足オブジェクトの配列。なければ空配列。
  - `detail`, `impact`, `resolution`: 不足内容、判断への影響、解消方法の文字列。
  - `blocks_assessment`: 評価を保留にするかの真偽値。
- `merge_gates`: 次の5 IDを各1件含む配列。
  - `id`: required_ci / no_conflicts / not_draft / required_approvals / review_threads。
  - `status`: pass / fail / unknown / not_applicable。
  - `evidence`: 状態と要求条件を裏付ける文字列。

## 整合性確認

- complete では SHA と6項目点が揃い、総合点は項目点の最大値。リスクなしで調査完了した項目は0点。
- 各完了項目点は category がその項目と一致するリスク点の最大値。所属リスクがなければ0点。risk_ids は存在する ID だけを参照し、他項目のリスクへの関連参照は項目点に含めない。
- hold では assessment_status と decision をともに hold、総合点を null にする。確認済みリスクは残す。
- blocks_assessment が true の情報不足を残したまま complete にしない。
- auto_merge_candidate は complete、総合点50以下、全必須条件が pass/not_applicable の場合だけ。
- fail/unknown の必須条件を点数へ加算せず、候補判定から除外する。
