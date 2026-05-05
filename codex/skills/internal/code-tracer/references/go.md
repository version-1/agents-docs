# Go Call Tracing

Go コードで呼び出し経路を追跡するときだけ、このリファレンスを読む。

## 対象の指定

- 関数: `pkg.Func`
- メソッド: `(*pkg.Type).Method` または `(pkg.Type).Method`
- 型: `pkg.Type`
- パッケージ範囲: `./...`、`./internal/...`、特定ディレクトリ

表記揺れがある場合は、定義箇所を確認してから完全修飾名へ寄せる。

## 依存ツール

優先して使う。

- Go toolchain（`go`）
- `gopls` の references / call hierarchy

必要に応じて使う。

- `go-callvis`
- Graphviz の `dot`
- `go list`
- `go test` によるパッケージ解決確認

ツールがない場合は、使える範囲の根拠を明記し、grep 由来の結果を確定経路として扱わない。

## 根拠収集

1. 定義箇所を特定する。
2. `gopls` references で参照箇所を集める。
3. 可能なら `gopls` call hierarchy で callers / callees を取得する。
4. 広い callgraph が必要なら `go-callvis` を補助として使う。
5. `go list` や既存テストで対象パッケージが解決できることを確認する。

## Go 固有の注意点

- interface 呼び出しは実装候補が複数ありうるため、確定経路と推定候補を分ける。
- DI、factory、wire、mock、test double 経由の接続は設定や初期化コードも確認する。
- reflection、plugin、init、goroutine、channel、callback は静的 callgraph だけで追えない場合がある。
- `_test.go`、mock、generated、vendor を含めるかは解析目的に応じて明示する。
- pointer receiver / value receiver の違いでメソッド集合が変わるため、対象メソッドの receiver を確認する。

## 推奨手順

1. `symbol`、`direction`、`depth`、`scope` を決める。未指定なら `both`、`depth=2`。
2. 定義と参照を `gopls` で確認する。
3. callers / callees を `gopls` call hierarchy で追う。
4. interface や DI が絡む場合は、実装候補と binding 箇所を別枠で整理する。
5. 必要なら `go-callvis` で補助的な callgraph を作る。
6. Mermaid 図では、確実な経路を実線、推定候補を破線で表す。

## 出力時に残す根拠

- 対象シンボルの定義ファイル。
- 使用した `gopls` / `go list` / `go-callvis` などの要点。
- 除外したファイル種別。
- 確定できなかった interface / reflection / DI の候補。
