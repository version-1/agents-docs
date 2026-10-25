---
name: issue-finder-playbook
description: Issue Finder agent がコードベースを変更せず、指定された finder 分野の課題を根拠付きで調査し、優先順位付けした Markdown を指定ディレクトリへ出力するときの共通手順、安全制約、品質基準、出力 contract を定義する。
---

# Issue Finder Playbook

対応する `role-*-finder-playbook` と併用する。分野別 playbook は課題判定、この playbook は調査と出力を扱う。

## 安全境界と入力

ソース、設定、既存文書を変更しない。書込みは指定出力先のみ、未指定なら `/tmp/issue-finder/`。canonical path を解決し、symlink を含まず解決後も `$TMPDIR` または `/tmp` 配下だけを使う。安全に書けなければ別場所へ無断出力せず報告する。非破壊検証だけを行い、生成物・更新を残すコマンドを避ける。秘密情報を収集・記録せず、外部サービスへ書き込まない。最新情報が必要な場合だけ読み取り専用の web/browser を使い、一次・公式情報を優先して URL と確認日を記録する。

finder 種別は分野別 playbook、対象は指定がなければリポジトリ全体、最大課題数は指定がなければ10件。対象・除外・未調査領域を成果物に書き、「全体」は実際の確認範囲の完全性保証ではない。

## 調査と判定

指示、仕様、構造、検証法を確認し、Git では開始時 status/diff を記録する。分野別観点から候補を集め、該当箇所、contract、呼び出し経路、テスト、必要な履歴を追う。再現、静的因果、仕様不整合など第三者が追える根拠を確認し、既存 Issue/TODO/課題文書との重複も確認する（GitHub Issue は接続・読取可能時のみ）。重要度、確信度、影響で優先付け、終了時に status/diff を比較する。新規変更が残れば完了扱いにしない。件数のために弱い候補を課題にしない。

課題は再現可能な事象、仕様矛盾、または追跡できる因果があり、影響と解決方針を説明できるものだけ。不足候補は `investigation-needed.md`、既存と同一なら `duplicates.md` に分ける。重要度は `Critical/High/Medium/Low`、確信度は `High/Medium/Low` と理由を必須にし、分野別基準に従う。

## 出力 contract

各課題を `<finder種別>-<対象名>-<YYYYMMDD-HHMMSS>-<slug>.md` にする。各要素は英小文字、数字、ハイフンへ正規化し、`.`、`..`、空白、absolute path、separator を含めず、出力先直下であることと slug の一意性を確認する。

各ファイルは、タイトル、finder 種別、概要、対象/範囲、追跡可能な根拠、発生条件、影響、重要度/理由、確信度/理由、重複確認、解決方針、未解決事項の順にする。上限外は `additional-candidates.md`、課題なしは `no-findings.md`（範囲、検証、未調査領域）を必要時だけ作る。

変更方針が一意で仕様判断不要かつ責務・公開 interface・contract・依存を根拠で確定できるときだけ `modification-design` の規定8章を完全に含め、境界判断では `lego-programming` を併用する。そうでなければ推測せず、選択肢・不足情報を未解決事項と `investigation-needed.md` に記録する。完了時は出力先、件数、補助ファイル、最重要、未調査領域、コードベース未変更を報告する。
