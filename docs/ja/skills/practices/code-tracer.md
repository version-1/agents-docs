---
name: code-tracer
description: 特定のモジュールの呼び出し経路を追跡・解析するためのスキル
---

# SKILL: Go呼び出し経路可視化（Codex CLI）

## 目的
特定の Go の **構造体 / 関数 / メソッド**について、呼び出し経路（Call Path / Call Graph）を **根拠付き**で可視化する。

- 「どこから呼ばれているか（Callers）」と「何を呼んでいるか（Callees）」を整理
- interface / 動的ディスパッチが絡む箇所は **確実（確定）** と **推定（候補）** を分離
- 出力は Markdown（Mermaid）を第一候補、必要に応じて dot / png を生成

---

## 入力（ユーザーが指定すること）
最低限（必須）
- `symbol`: 対象シンボル（例：`pkg.Func` / `(*pkg.Type).Method` / `pkg.Type`）
- `direction`: `callers` or `callees` or `both`

任意（デフォルトあり）
- `depth`: 2（1〜5 推奨）
- `scope`: パッケージ/ディレクトリ範囲（例：`./internal/...`、`./...`）
- `format`: `mermaid`（既定）/ `dot` / `png`
- `exclude`: `vendor`, `*_test.go`, `mock`, `generated` など

---

## 出力（成果物）
- `/tmp/callpath_YYYYMMDDHHMMSS.md`
  - YYYYMMDDHHMMSS は実行時刻
  - 対象シンボル、解析条件、結果要約
  - Mermaid 図（フローチャート or グラフ）
  - **確実 / 推定** の凡例と注釈
  - “根拠”セクション（gopls / callgraph 由来のメモ）

任意
- `callpath.dot`（Graphviz 用）
- `callpath.png`（dot から生成）

---

## 依存ツール
### 必須（推奨）
- Go toolchain（`go`）
- `gopls`（参照/呼び出し階層の根拠取得に利用）

### 任意（精度・視覚化強化）
- `go-callvis`（静的 callgraph を広く取得する用途）
- Graphviz（`dot` コマンド：png 生成に利用）

> 注意: ツールの有無に応じて「できる範囲」で段階的に出力する。
> 例: `gopls` だけでも “参照箇所一覧 + 近傍の呼び出し関係” は作れる。

---

## 解析戦略（重要：推測禁止）
本スキルは「grep 推測」ではなく、次の根拠に基づく。

1. **gopls の参照（references）**  
   - シンボルが参照されている箇所（= 呼び出し元の強い候補）を取得
2. **gopls の呼び出し階層（call hierarchy）**  
   - callers/callees を深さ付きで取得（可能なら）
3. **go-callvis 等の callgraph**（任意）  
   - pointer analysis を含む広い call graph を補助として使う  
4. **interface / reflection / DI**  
   - 静的に確定できない場合は「推定」として候補集合を示す  
   - “確実な線” と “推定の線” を図・文章で分ける

---

## 実行手順（Codex にやらせること）
### Step 0: 条件整理
- 対象シンボルの表記揺れをなくす（できるだけフル修飾）
- `direction` と `depth` を決める（既定: both + depth=2）
- `scope` を決める（既定: `./...` だが大規模なら狭める）

### Step 1: gopls で根拠収集
- references を集める
- 可能なら call hierarchy も集める
- テスト/生成コードなどを除外して整形

### Step 2: callgraph で補助（任意）
- パッケージスコープを絞って call graph を作成
- dot を得て、必要なら png 化

### Step 3: `callpath.md` に統合
- 要約（結論）
- Mermaid 図
- “確実/推定” の区分
- 根拠（取得コマンドや出力の要点）

---

## Mermaid 図の表現ルール
- **確実（確定）エッジ**: 実線
- **推定（候補）エッジ**: 破線（Mermaid の `-.->` を使用）
- ノード名は短くしつつ、必要なら `pkg.Type.Method` を併記
- 重要経路は太字注釈（文章側で）

例（イメージ）
```mermaid
flowchart TD
  A[handler.Handle] --> B[service.Do]
  B --> C[repo.Find]
  B -.-> D[(iface).Call]:::maybe
  classDef maybe stroke-dasharray: 5 5;
```



