# Ebitengineのキーリピート実装方針

- Date: 2026-03-30
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

旧 termbox 版では `termbox.PollEvent()` がブロッキング呼び出しであり、OS レベルのキーリピートがそのままイベントとして流れてきていた。
Ebitengine へ移行後は TPS=60 の固定ループ内でポーリングする方式に変わり、`inpututil.IsKeyJustPressed` は押した瞬間の1フレームのみ true を返す。
そのため移動キー（h j k l w e b）を押し続けても最初の1回しか反応しないという問題が発生した。

## Decision

`input.Handler` にリピート状態（`repeatKey ebiten.Key`、`repeatFrames int`）を持たせ、ソフトウェアでキーリピートを実装する。
- 押した瞬間は `inpututil.IsKeyJustPressed` で即座に発火する。
- 押し続けて `repeatInitialDelay`（18フレーム ≈ 300ms）経過後、`repeatInterval`（4フレーム ≈ 60ms）ごとに繰り返し発火する。
- リピート対象は移動系コマンド（h j k l w e b）のみとし、g・q・f/t 系・数字キーは対象外とする。

## Consequences

**メリット**
- キーを押し続けると自然に連続移動でき、ゲームの操作感が向上する。
- リピートのタイミング（初期遅延・間隔）を定数で調整しやすい。
- `state` パッケージへの影響がなく、`input.Handler` 内で完結する。

**トレードオフ**
- リピートのタイミング定数はゲームの TPS（60）に依存する。TPS を変更する場合は定数の見直しが必要。
- OS のキーリピート設定と独立しているため、ユーザーの OS 設定とは異なるリピート速度になる。

## Alternatives Considered

- **OS のキーリピートをそのまま使う**: Ebitengine のポーリング方式では OS リピートイベントを直接取得できないため採用不可。
- **`ebiten.IsKeyPressed` のみを使う**: 押している間ずっと true になるため毎フレーム発火してしまい、カーソルが速すぎる。初期遅延を設けるリピート処理が必要。

## Related Files

- `CLAUDE.md`
- `input/input.go`
