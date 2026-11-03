# 役割契約

これは factory のタスク別役割であり、同名の登録済み agent の存在を前提にしない。実行先のモデル規則と利用可能な agent に合わせて、役割と下記の契約を指示として渡す。モデルが指定されていて利用できない場合は、別バージョンを使ったと偽らず制約を報告する。

初期提案の `refinener` と `documentor` は互換のため役割名として維持する。役割ごとの受け取る情報、判断、報告・エスカレーションは以下の agent ガイドに従う。

## Issue ごとの担当と引き継ぎ

- 1 issue に1人の issue-runner を割り当て、完了まで責任と状態を維持させる。保留中は再開できる状態を保存する。verifier、refinener、QA はその issue の工程ごとに専任担当を起動し、複数 issue の成果物や判断を同じ担当の文脈に混在させない。
- 同一 issue は implementer、verifier、refinener、verifier による再検証、QA の順に進める。verifier が依頼する独立した reviewer は並行できる。issue 間は最大5件まで並行できるが、依存または変更範囲の衝突がある場合は直列化する。
- 全役割を常駐させない。工程を終えた担当は、次の担当または再開担当が作業を再現できる状態を残して終了し、agent 枠を解放する。再起動時も担当を入れ替えたことだけで改善の回数や保留の履歴をリセットしない。
- 引き継ぎには issue、受け入れ条件、worktree とブランチ、base/head SHA、成果物、実行した検証と結果、未解決指摘、再試行回数、次の工程または保留理由を含める。

## Agent ガイド

- [Pdm](agents-pdm.md): ゴール、追加指示、優先度、割り当て、停滞を統括する。
- [issue-runner](agents-issue-runner.md): issue ごとの成果物、担当、変更伝達、handoff を管理する。
- [implementer](agents-implementer.md): 実装・調査・検証と指摘修正を行う。
- [verifier](agents-verifier.md): 独立した検証と reviewer の結果集約を行う。
- [refinener](agents-refinener.md): 挙動を維持した品質改善を1巡行う。
- [QA](agents-qa.md): 品質、Ready、承認、自動統合を分離して判定する。
- [code-reviewer](agents-code-reviewer.md): コードの仕様適合、回帰、境界条件、公開契約、テストを確認する。
- [ui-reviewer](agents-ui-reviewer.md): 実表示、画面幅、状態、アクセシビリティを確認する。
- [ux-reviewer](agents-ux-reviewer.md): 操作導線、エラー・空・待機状態、回復を確認する。
- [security-reviewer](agents-security-reviewer.md): 信頼境界、漏えい、入力、権限、依存を確認する。
- [refactor](agents-refactor.md): 挙動を維持した構造改善を行う。
- [documentor](agents-documentor.md): 変更に必要な文書の整合と導線を整える。
