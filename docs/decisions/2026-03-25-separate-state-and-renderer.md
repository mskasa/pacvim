# 0002: state パッケージを描画ライブラリから分離する

- Date: 2026-03-25
- Status: Active
- Author: mskasa

## Context

Ebitengine で1から書き直すにあたり、旧実装と同じ轍を踏まないよう設計方針を定める必要があった。

旧実装の根本的な問題は「ゲームの状態が描画ライブラリ（termbox）のバッファと一体化していた」ことである。
Ebitengine に移行しても同じ設計を踏襲すれば、今度は `ebiten.Image` のピクセルバッファを
ゲーム状態として使うことになりかねず、同じ問題が再発する。

また Ebitengine は `Update()` と `Draw()` を毎秒60回呼び出す設計になっており、
この2つの責務を明確に分離することが Ebitengine の正しい使い方でもある。

## Decision

ゲームの状態を **`state` パッケージの純粋な Go struct** として持ち、
描画への依存を **`renderer` パッケージ** に完全に閉じ込める。

ルール:

- `state` パッケージは `github.com/hajimehoshi/ebiten/v2` を import しない
- `renderer` パッケージは `state` を読むが、書き込まない
- `Update()` は `state.GameState.Update()` を呼ぶ（状態変更）
- `Draw()` は `renderer.Renderer.Draw()` を呼ぶ（描画のみ、状態変更なし）

この制約により `state` パッケージは `go test` で単体テストが書ける。

## Consequences

**得るもの**

- `state` パッケージが純粋な Go のため、`go test` で単体テストが書ける
- 描画ライブラリを将来変えても `state` は影響を受けない
- `Update()` と `Draw()` の責務が明確になり、バグの原因箇所が特定しやすい

**失うもの・トレードオフ**

- パッケージが3層（state / input / renderer）に分かれるため、小規模なうちは冗長に見える
- パッケージ間の依存方向を意識して設計する必要がある

## Alternatives Considered

**`game.go` にすべてをまとめる**
シンプルだが `Update()` と `Draw()` の境界が崩れやすく、旧実装と同じ問題が再発する。
テストも書きにくくなる。

**ECS（Entity-Component-System）アーキテクチャを採用する**
PacVim の規模では過剰。state / input / renderer の3層分離で十分。

## Related Files

state/gamestate.go
state/player.go
state/enemy.go
state/map.go
renderer/renderer.go
game.go
