# Commit Message Format

commit message の体裁に迷う場合は、このデフォルトを使う。
既存履歴やユーザー指定に別の形式がある場合は、そちらを優先する。

## Language

既存履歴やユーザー指定に合わせる。
判断できない場合は日本語テンプレートを使う。

- 既存履歴が英語中心なら英語テンプレートを使う。
- 既存履歴が日本語中心なら日本語テンプレートを使う。
- ユーザーが言語を指定した場合は、その言語のテンプレートを使う。
- prefix は日本語テンプレートでも英語テンプレートでも英語のまま使う。

## Prefix

既存履歴やユーザー指定から形式を判断できない場合は、次の prefix を使う。
迷う場合は、変更の主目的で選ぶ。

| Prefix | 使う場面 |
|---|---|
| `feat` | ユーザー向け機能、agent / skill / command の新規追加 |
| `fix` | バグ修正、誤動作、壊れている挙動の修正 |
| `docs` | README、設計文書、手順書、説明文だけの変更 |
| `test` | テストの追加、修正、テストデータ更新 |
| `refactor` | 挙動を変えない構造整理、責務分離、命名整理 |
| `chore` | 配布設定、依存、生成物、メタデータ、開発補助の変更 |
| `ci` | CI / GitHub Actions / 自動化パイプラインの変更 |
| `build` | ビルド設定、パッケージング、ツールチェーンの変更 |
| `perf` | 性能改善 |
| `style` | フォーマット、空白、lint 指摘など挙動に影響しない表記変更 |

複数に見える場合は、最終的にユーザーへ提供する価値に最も近いものを選ぶ。
例えば skill 本体と一覧ドキュメントを同時に追加する場合、主目的が skill 追加なら `feat` を使う。

## Japanese Template

### Title

形式を指定されていない場合は、次の形にする。
タイトルだけの commit は作らず、必ず本文も付ける。

```text
<prefix>: <変更の目的を短く説明する>
```

例:

```text
feat: add cmd-commit skill
docs: update skill library categories
fix: prevent unrelated files from being committed
```

タイトルを書くときは、次を守る。

- prefix は英語、説明は日本語で書く。
- 末尾に句点を付けない。
- できるだけ 72 文字以内に収める。
- ファイル名ではなく、変更の目的を書く。
- `update`、`fix`、`change` だけのタイトルにしない。

例:

```text
feat: cmd-commit skill を追加
docs: skill 一覧を更新
fix: 無関係な変更が commit に混ざらないようにする
```

### Body

本文は必ず付ける。
空行を挟んで次の形式にする。

```text
<prefix>: <変更の目的を短く説明する>

Why:
- <変更が必要になった理由を日本語で書く>

What:
- <主な変更点を日本語で書く>
- <影響範囲を日本語で書く>

Verification:
- <実行した検証を日本語で書く>
```

検証していない場合は、`Verification` に未実行と理由を書く。

```text
Verification:
- 未実行。ドキュメントのみの変更。
```

日本語で本文を書く場合も見出しは `Why`、`What`、`Verification` のままにする。
この形式は履歴を機械的に読みやすくするためのデフォルトであり、既存履歴に別の本文形式がある場合はそちらを優先する。

## English Template

### Title

Use this format when the repository history is mainly English or the user asks for English.

```text
<prefix>: <short imperative summary>
```

Examples:

```text
feat: add cmd-commit skill
docs: update skill library categories
fix: prevent unrelated files from being committed
```

Title rules:

- Keep the prefix in English.
- Write the summary in English.
- Do not end with a period.
- Keep it within 72 characters when possible.
- Describe the purpose, not just the file name.
- Avoid vague titles such as `update`, `fix`, or `change`.

### Body

Always include a body.
Use this format after a blank line.

```text
<prefix>: <short imperative summary>

Why:
- <why the change is needed>

What:
- <main change>
- <impact or scope>

Verification:
- <verification performed>
```

If verification was not run, say so with the reason.

```text
Verification:
- Not run (documentation-only change).
```
