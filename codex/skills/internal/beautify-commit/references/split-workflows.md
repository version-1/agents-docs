# Split Workflows

commit 分割時の代表的な操作パターンをまとめる。
コマンドは状況に合わせて調整し、実行前に `git status --short` と対象範囲を確認する。

## 基準 ref との差分を整理して commit し直す

通常の入口はこのパターンを使う。
基準 ref は、ベースブランチまたは残したい履歴の終点 commit hash のどちらでもよい。
ユーザーが基準 ref を指定していない場合は、デフォルトブランチでよいか yes/no で確認してから進める。

1. 基準 ref を決める。
2. `git fetch` が必要か判断する。最新のリモート状態が必要なら、ユーザー確認や環境の権限に従って実行する。
3. 基準 ref がブランチなら、必要に応じて `git merge-base <base-ref> HEAD` で分岐点を確認する。
4. `git log --oneline --reverse <base-ref>..HEAD` で差分に含まれる commit を確認する。
5. `git diff --stat <base-ref>..HEAD` と `git diff <base-ref>..HEAD` で差分全体を読む。
6. 必要に応じて `grill-me` skill を使い、変更単位、依存順、統合すべき commit、分けるべき commit を詰める。
7. バックアップ branch を作る。
8. 基準 ref との差分に含まれる commit をいったん未コミット差分へ戻す。
9. 変更単位ごとに stage し、`cmd-commit` skill を使って commit する。
10. `git log --oneline --decorate <base-ref>..HEAD` と `git status --short` を確認する。

差分を未コミット差分へ戻す代表例:

```bash
git branch backup/<branch>-before-beautify-<date>
git reset <base-ref>
```

注意:

- `<base-ref>` は残したい履歴の終点として扱う。現在 branch の整理基準として正しいか確認する。
- `git reset <base-ref>` は対象範囲の commit を取り消し、変更内容を作業ツリーに戻す。
- push 済み履歴を書き換える場合は、実行前にユーザー確認を取る。
- 整理後の commit 作成は `cmd-commit` skill に従い、staged diff を毎回確認する。

## 未コミットの巨大差分を分割する

作業ツリーに大きな未コミット差分がある場合は、履歴を書き換えずに分割できる。

1. `git status --short` と `git diff --stat` で変更全体を確認する。
2. `git diff` を読み、変更単位を列挙する。
3. ファイル単位で分けられる場合は、対象ファイルだけ `git add <path>` する。
4. 同じファイル内に複数の意図が混在する場合は、`git add -p <path>` で hunk 単位に stage する。
5. `git diff --staged` で commit に入る内容を確認して commit する。
6. `git status --short` で残りを確認し、次の変更単位を stage する。

注意:

- hunk が大きすぎる場合は、`s` で split できるか試す。
- patch 分割が難しい場合は、対象ファイルを一時的に編集して 1 commit 分の差分だけ残し、commit 後に残りを戻す。
- 生成物や lockfile は、対応する実装変更と同じ commit に入れるか、生成物更新だけで意味がある場合に分ける。

## 直近 1 commit を分割する

直近 commit だけを分割する場合は、`git reset` で commit を未コミット差分に戻してから再 commit する。

```bash
git branch backup/<branch>-before-split-<date>
git reset HEAD^
```

その後は「未コミットの巨大差分を分割する」と同じ手順で stage / commit する。

注意:

- `git reset HEAD^` は commit だけを取り消し、変更は作業ツリーに残す。
- push 済みの直近 commit を分割すると履歴が変わる。force push が必要になる可能性をユーザーへ伝える。

## 連続した複数 commit をまとめて分割し直す

連続した commit 群を整理する場合は、対象範囲の base を決め、まとめて未コミット差分へ戻してから再構成する。

```bash
git branch backup/<branch>-before-split-<date>
git reset <base>
```

ここで `<base>` は残したい最後の commit を指す。
例えば直近 4 commit を作り直すなら `HEAD~4` を使う。

その後、変更単位ごとに stage / commit する。

注意:

- 元の commit 順が意味を持つ場合は、`git log --reverse <base>..HEAD` で順序を確認してから分割案を作る。
- 複数人と共有済みの branch では、reset による履歴書き換えを実行前に必ず確認する。
- reset 後に迷った場合は、バックアップ branch から元の状態を参照できる。

## 古い commit を分割する

直近ではない commit を分割する場合は、interactive rebase で対象 commit を `edit` にする。

```bash
git branch backup/<branch>-before-split-<date>
git rebase -i <target>^
```

rebase todo で分割したい commit を `edit` に変更する。
停止したら次を実行する。

```bash
git reset HEAD^
```

変更単位ごとに stage / commit し、分割が終わったら次を実行する。

```bash
git rebase --continue
```

注意:

- rebase 中に conflict が起きたら、差分を読み、解消後に `git add <resolved-files>`、`git rebase --continue` を実行する。
- conflict 解消で無関係な変更を混ぜない。
- rebase を中止する必要がある場合は `git rebase --abort` を使えるが、実行前に現在の状態を確認する。

## commit を並べ替えたり統合したりしながら分割する

複数 commit の一部を統合、一部を分割、一部を並べ替える場合は、interactive rebase を使う。

使う操作:

- `pick`: commit をそのまま残す。
- `reword`: message だけ変更する。
- `edit`: commit の中身を変更、分割する。
- `squash` / `fixup`: 前の commit に統合する。
- 行の順序変更: commit 順を変更する。

注意:

- 並べ替えは依存関係を壊しやすい。build や test が通る順序か確認する。
- 同じファイルの近い行を複数 commit が触っている場合、conflict が増える可能性がある。

## 安全確認

分割作業の前:

- `git status --short` が想定どおりか確認する。
- 対象 branch と対象範囲を言語化する。
- バックアップ branch を作る。
- push 済み履歴を書き換える場合はユーザー確認を取る。

分割作業の途中:

- 各 commit 前に `git diff --staged` を読む。
- 各 commit 後に `git status --short` で残りを見る。
- 迷ったら `git log --oneline --decorate -n 20` とバックアップ branch を確認する。

分割作業の後:

- `git log --oneline --decorate -n <必要数>` で commit 一覧を確認する。
- `git status --short` が想定どおりか確認する。
- 可能なら関連テスト、lint、build を実行する。
- push 済み履歴を書き換えた場合は、通常の push では失敗する可能性と force-with-lease が必要になり得ることを伝える。
