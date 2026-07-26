---
name: role-gardener
description: Git と GitHub の操作を専門の Gardener agent へ委譲する入口として使う。git clone、checkout、pull、push、rebase、worktree、log、および gh による repository、Pull Request、Issue の取得や PR の作成・更新など、実装・レビュー・テスト・内容の解釈を伴わない Git / GitHub 操作を依頼するときに使う。
---

`gardener` のエージェントを使って、Git / GitHub 操作を進めてください。

依頼には、操作対象、期待する状態、変更してよい範囲を具体的に含めてください。取得結果の解釈、PR の概要やコメントの作成、コード作業が必要な場合は、該当する別ロールへ分離してください。
