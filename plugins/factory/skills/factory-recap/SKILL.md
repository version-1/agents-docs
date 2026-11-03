---
name: factory-recap
description: factory の利用経験や振り返りから、利用先プロジェクト内の factory 運用ガイドライン skill を作成・更新する。共通 factory スキルへの改善 Issue 作成には factory-feedback を使う。
---

# Factory Recap

factory を使うプロジェクトの判断基準・検証方法・運用上の知見を、次の実行で再利用できる project スコープの skill にする。会話の要約だけを残すのではなく、どの条件で何を判断・実行するかを明記する。

## 対象と保存先

- 対象は factory の実行先プロジェクト。現在の作業場所、AGENTS.md / CLAUDE.md、既存のローカル skill、管理 Issue・PR とユーザーの指摘を確認する。対象プロジェクトを一意に決められなければ確認する。
- 既存の factory 運用ガイドライン skill があれば更新する。新規はプロジェクトルートの `.agents/skills/factory-project-guidelines/SKILL.md` を既定とし、Claude 専用など既存の配置規約があればそれに従う。ユーザーのホームやインストール済み plugin、agents-docs の共通 factory は変更しない。
- 新しい配置が factory から発見できない場合は、プロジェクトの AGENTS.md / CLAUDE.md の適切な入口にパスと適用条件を短く追記する。本文を複製しない。

## 知見をガイドラインにする

1. 実際に確認できた成功・失敗・ユーザーの修正から、再利用できるルールを抽出する。適用条件、行動、完了の確認方法、必要な根拠・参照先を整理する。推測や未合意の方針は確定ルールにせず、確認事項として報告する。
2. 対象プロジェクト特有の仕様判断、Issue 分割、検証コマンド、レビュー観点、承認が必要な境界を必要な範囲で記載する。factory 本体や既存ルールの転載、単発の作業記録、秘密情報は残さない。共通 factory の改善は `factory-feedback` の候補として分け、Issue を自動投稿しない。
3. 既存の有効なルールを保ち、重複や矛盾は根拠を照合して統合する。ユーザーの明示した方針を優先し、解消できない矛盾だけを確認する。ローカルルールで factory の品質条件や承認、ツール権限を免除しない。
4. 新規 skill の frontmatter は `name: factory-project-guidelines` とし、description に対象プロジェクトと「factory の開始・再開・Issue 計画・検証時に使う」ことを記載する。既存 skill の名前と適用範囲は尊重する。長い詳細は references に分け、入口から参照する。
5. 完成した skill を全文で読み、frontmatter、参照先、既存規則との整合、factory からの発見経路を確認する。利用可能なら skill validator も実行する。変更ファイルと追加・更新した判断基準、保留事項を報告する。

ガイドライン作成・更新までを行い、factory の再開、実装、commit・配布は依頼がある場合に行う。
