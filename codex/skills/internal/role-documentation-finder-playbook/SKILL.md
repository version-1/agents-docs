---
name: role-documentation-finder-playbook
description: Issue Finder agent がドキュメントの欠落、重複、矛盾、情報配置、導線、読み手や粒度の混在、progressive disclosure の問題を発見するときに使う分野別の探索観点、課題判定、重要度基準を定義する。
---

# Documentation Finder Playbook

`issue-finder-playbook` と併用する。

## 探索観点

- code、設定、CLI、API と文書化された手順・contract の不一致
- 同じ事実の重複管理、版の違い、相互に矛盾する指示
- README、guide、reference、runbook 間の導線不足や孤立文書
- 初学者向け概要への詳細混入、reference に必要な完全性の欠落
- 読み手、目的、前提知識、情報粒度が一文書内で混在する構成
- 最初から過剰な情報を読ませる、必要な詳細へ到達できないなど progressive disclosure の破綻
- 古い名称、壊れた link、未記載の前提、実行不能な example

## 課題判定

対象読者が誤判断する、必要情報へ到達できない、誤った操作をする、または複数文書の同期コストが発生する具体的状況を説明できるものだけを課題にする。軽微な誤字は意味や検索性へ影響する場合だけ含める。

## 重要度

- Critical: 文書に従う通常操作が不可逆なデータ損失や重大な security incident を招く。
- High: 中核手順や公開 contract の誤りにより、広い利用者が失敗する。
- Medium: 利用者が行き詰まる、誤解する、または保守者が重複更新を強いられる。
- Low: 局所的な導線、検索性、理解コストに明確な問題がある。

## 除外

- 読解や操作へ影響しない表現上の好み
- 読み手と目的に適した意図的な情報重複
- code や設定を確認せずに断定した陳腐化
