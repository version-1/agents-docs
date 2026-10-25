---
name: format-rich-html-diagram
description: コード、アーキテクチャ、React コンポーネント構造、データフロー、状態管理、シーケンス、処理フロー、依存関係、ディレクトリ階層などを、SVG、2D、3D、アニメーション、CDN ライブラリを使ったリッチな単一 index.html として可視化・整理するときに使う。開発者がブラウザで index.html だけを開けば全体像、詳細、関係性、判断材料を理解できる成果物を作る必要がある場合は必ず使う。Markdown だけでは構造や流れを把握しづらい設計説明、コードリーディング結果、実装計画、レビュー補助資料、オンボーディング資料に使う。
---

# Format Rich HTML Diagram

コードや設計を、ブラウザで開く単一 `index.html` として説明する。対象は理解・整理・レビュー用の資料であり、コード変更ではない。

## 出力契約

- 成果物は原則 `index.html` 一つ。CSS と小さな JavaScript はインラインにし、依存追加、ビルド、設定変更はしない。
- 外部ライブラリは下記の許可 CDN だけ、フォントは Noto Sans JP と JetBrains Mono だけを必要最小限に使う。未列挙ライブラリ、フォント、外部画像、analytics、tracking は使わない。追加が必要なら作成前に URL、用途、理由、影響を示して確認を得る。
- 図だけに重要情報を閉じ込めず、要約テキストまたは表を隣接させる。装飾のための 3D・アニメーションは避ける。

| 用途 | 選択肢 |
|---|---|
| 軽量な図、注釈、凡例 | inline SVG |
| flow / sequence / state | Mermaid `https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs` |
| tree、sankey、独自 SVG | D3 `https://cdn.jsdelivr.net/npm/d3@7/+esm` |
| 依存グラフ | Cytoscape `https://cdn.jsdelivr.net/npm/cytoscape@3/dist/cytoscape.esm.min.mjs` |
| 大量要素の 2D | PixiJS `https://cdn.jsdelivr.net/npm/pixijs@8/dist/pixi.min.js` |
| 意味のある 3D | Three `https://cdn.jsdelivr.net/npm/three@0.173.0/build/three.module.js` |
| 状態・順序の説明 | Anime `https://cdn.jsdelivr.net/npm/animejs/dist/bundles/anime.umd.min.js` |
| コード表示 | Prism `https://cdn.jsdelivr.net/npm/prismjs@1/components/prism-core.min.js` |
| Markdown 断片 | marked `https://cdn.jsdelivr.net/npm/marked@12/+esm` |
| SVG zoom | svg-pan-zoom `https://cdn.jsdelivr.net/npm/svg-pan-zoom@3/dist/svg-pan-zoom.min.js` |

Google Fonts は Noto Sans JP `https://fonts.googleapis.com/css2?family=Noto+Sans+JP:wght@400;500;700&display=swap` と JetBrains Mono `https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap` だけを使う。両方を使う場合も同じ stylesheet にまとめる。

単純な図は inline SVG、構文で十分な flow は Mermaid、選択・フィルタする大規模グラフは Cytoscape、tree/sankey は D3 を優先する。Canvas、3D、アニメーションは密度・奥行き・時間変化が理解を実際に助けるときだけ使う。CDN を使うなら表示制限を Appendix に残す。

## 内容と UI

読者の問いと対象範囲を確認し、関連するファイル・型・処理・テストを根拠に、全体像、構造、flow、state、contract、risk を整理する。必要な章だけを選び、Overview、Architecture、Structure、Flow、State、Contracts、Risks/Open Questions、Appendix を基本にする。Appendix には根拠、使用 CDN、制限を置く。

章へ移動できるナビゲーションを置き、図と説明を近くに配置する。色だけで意味を伝えずラベル・凡例も付ける。大きな図は responsive にし、必要に応じて zoom または横スクロールを使う。アニメーションは再生・停止・リセットと `prefers-reduced-motion` を、3D は操作説明・凡例・2D 要約を、Canvas/WebGL は初期化失敗時の説明を備える。

完全な HTML として `<!doctype html>`、viewport を含める。SVG は `viewBox`、`role="img"`、`aria-labelledby`、`title`/`desc` を使う。Mermaid は `securityLevel: "strict"` を基本とし、未信頼テキストは escape する。図の元データは Mermaid source、JavaScript object、または JSON script として HTML に残す。直接開くと壊れるリンクへ依存しない。

## 完成時

ナビゲーション、図、表、本文、リンクをブラウザ表示前提で確認する。最終報告には、`index.html` のパス、主な可視化、使用 CDN、確認した表示、残る制限・未確認点を短く示す。
