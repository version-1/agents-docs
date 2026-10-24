---
name: role-feature-finder-playbook
description: Issue Finder agent がコードベース内の利用者課題、未完了仕様、不自然な制約を根拠として機能提案を発見するときに使う分野別の探索観点、提案判定、重要度基準、誤検知の除外条件を定義する。
---

# Feature Finder Playbook

`issue-finder-playbook` と併用する。

## 探索観点

- README、仕様、roadmap、TODO に明示された未完了の利用者価値
- UI、API、CLI の途中で途切れる workflow や手作業による迂回
- 類似機能間の不自然な非対称性、公開 contract の利用不能な余地
- tests、examples、support notes が示す反復的な利用者負担
- 既存 domain model で表現済みだが利用者へ提供されていない能力

## 提案判定

「誰が」「どの状況で」「何に困るか」をリポジトリ内の根拠から説明し、既存の責務や contract と自然に接続できるものだけを課題にする。外部動向は必要な場合の補強材料に限り、提案理由の代わりにしない。

## 重要度

- Critical: 該当なし。機能提案に Critical を付けない。
- High: 中核 workflow の完遂を妨げる未提供機能で、広い利用者影響が確認できる。
- Medium: 明確な利用者負担や反復作業を解消する。
- Low: 限定的だが根拠のある利便性・一貫性の改善になる。

## 除外

- 流行、競合の存在、実装可能性だけに基づくアイデア
- 対象利用者と解決する問題を説明できない提案
- 既存機能で無理なく達成できる要求
- 不具合修正や内部だけの保守性改善
