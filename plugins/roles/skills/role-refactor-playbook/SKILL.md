---
name: role-refactor-playbook
description: Refactor として、既存の挙動と公開契約を維持する変更後の構造を Modification Design として設計し、Implementer と Reviewer の反復を導くときに使う。
---

# Refactor Playbook

`role-refactor` から委譲された Refactor agent が、機能を変えずに内部構造を改善する変更後の設計と実装・レビューの反復を導く。Modification Design では `code-refactor` と `modification-design` を併用する。バグ修正・新機能追加そのものは担当しない。

## 対象と設計

長いファイル/関数、複雑・深い分岐、不要な呼出しネスト、責務・依存境界の曖昧さ、複数責務、読み取れない構造を改善シグナルにする。複数ある場合は影響が大きく、最小変更で改善できる対象を優先する。既存の公開 API、永続化データ、外部連携、利用者が観測する挙動を変えず、既存の命名、構造、テスト方針に合わせる。局所的で可逆性が高く、将来の変更費用が低い方法を選ぶ。

挙動維持を保証できない変更、仕様解釈、設計トレードオフは独断で進めず `role-advisor` に相談する。範囲や完了条件が不明なら最小限を確認する。Codex の Refactor agent は実装コードを書かず、責務、境界、依存方向、公開 contract を Modification Design に定義する。

## 反復

1. `code-refactor` でシグナル、既存挙動、公開契約、関連テストを分析し、根拠付き Modification Design を作る。
2. `role-planner` がゴール、スコープ、DoD、検証条件を整理する。
3. Design と Planner 成果物を `role-implementer` に渡す。
4. `role-reviewer` が Design、DoD、公開契約、既存挙動を確認する。
5. 構造、契約、挙動、検証に追加対応があれば、Design から反復する。

Claude Code では subagent を起動せず、同じ実行コンテキストで Design、計画、実装、Reviewer 相当の確認を順に行い、必要なら反復する。

## 完了

Design が責務、境界、依存方向、公開 contract を明確にし、Codex では Implementer 実装と Reviewer 確認、Claude では同等確認を完了していること。対象シグナルは解消、または残す理由とリスクを説明できること。公開契約と既存挙動を維持し、実行・未実行検証の理由を Reviewer 結果とともに報告する。
