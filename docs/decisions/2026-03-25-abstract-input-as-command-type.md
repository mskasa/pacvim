# 0003: キー入力を Command 型で抽象化する

- Date: 2026-03-25
- Status: Active
- Author: mskasa

## Context

Ebitengine でキー入力を読むには `ebiten.Key` 型を使う。

```go
if inpututil.IsKeyJustPressed(ebiten.KeyH) { ... }
```

この判定を `state/player.go` の中で直接行うと `state` パッケージが `ebiten` に依存することになり、
ADR-0002 の「state パッケージは ebiten を import しない」という原則に違反する。

また PacVim の Vim コマンドには複数フレームにまたがる入力シーケンスがある。

- `gg`: `g` → `g` の2ストローク
- `f{char}`: `f` → 任意の文字の2ストローク
- `ge`: `g` → `e` の2ストローク
- カウント入力: `3` → `w` のように数字を蓄積してからコマンドを実行

この状態管理をゲームロジック側（`state` パッケージ）に持ち込むと、
Vim コマンドの処理とキー入力の解析が混在して複雑になる。

## Decision

キー入力は `input` パッケージで `Command` 型に変換してから `state` パッケージに渡す。

```go
// input パッケージが定義する型
type Command int

const (
    CmdNone Command = iota
    CmdMoveLeft
    CmdFindForward  // f{char} — 文字は rune で別途渡す
    // ...
)
```

`input.Handler.Read()` が毎フレーム `(Command, rune)` を返す。
`gg` / `f{char}` / `ge` などの複数フレーム入力の状態管理は `input.Handler` 内で完結させる。

`state.Player.Apply(cmd Command, ch rune, stage *Stage)` はキーを知らず、
`Command` と対象文字（`f/t` 系のみ使用）だけを受け取る。

## Consequences

**得るもの**

- `state` パッケージが `ebiten.Key` を知らないため、テストで `Command` を直接渡せる
- 複数フレーム入力の状態管理が `input` パッケージに閉じ込められ、ゲームロジックがシンプルになる
- キーバインドを変更する場合、`input` パッケージだけを変えればよい

**失うもの・トレードオフ**

- Vim コマンドの数だけ `Command` 定数が増える
- `input` パッケージと `state` パッケージの間に変換レイヤーが生まれる

## Alternatives Considered

**`state/player.go` で直接 `ebiten.Key` を読む**
シンプルだが `state` が `ebiten` に依存し ADR-0002 に違反する。
テスト時にキーイベントをシミュレートする別の仕組みが必要になる。

**`ebiten.Key` をインターフェースでラップして注入する**
テスト可能にはなるが、インターフェースの定義・実装・モックの管理が煩雑になる。
`Command` 型への変換の方がシンプルで意図が明確。

## Related Files

input/input.go
state/player.go
state/gamestate.go
game.go
