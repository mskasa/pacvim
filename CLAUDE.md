# CLAUDE.md — PacVim 開発ガイド

このファイルはAI駆動開発向けのプロジェクトガイドです。
コードを変更する前に必ずこのファイルを読んでください。

---

## プロジェクト概要

**PacVim** は、パックマンのゲームルールで遊びながら Vim の操作を学べる Go 製のゲームです。

### 経緯

もともと `github.com/nsf/termbox-go` を使ったターミナルゲームとして実装されていましたが、
以下の理由から Ebitengine を用いて1から書き直しています。

- termbox のセルバッファをゲーム状態として流用していたため、テストが書けず拡張も困難だった
- WebAssembly 配布（ブラウザでプレイ）を実現したかった
- Vim コマンドをさらに追加するにあたり、クリーンな設計基盤が必要だった

設計判断の詳細は `docs/decisions/` の ADR を参照してください。

### 技術スタック

| 項目 | 採用技術 |
|---|---|
| 言語 | Go 1.22+ |
| ゲームエンジン | `github.com/hajimehoshi/ebiten/v2` |
| 入力 | `ebiten/v2/inpututil` |
| テキスト描画 | `ebiten/v2/text/v2` |
| テスト | `go test` 標準ライブラリ |

> `github.com/nsf/termbox-go` はこのプロジェクトに存在しません。
> termbox に関するコードやコメントを書かないでください。

---

## ディレクトリ構造

```
pacvim/
├── main.go              # エントリポイント（ebiten.RunGame のみ）
├── game.go              # ebiten.Game インターフェースの実装
├── state/
│   ├── gamestate.go     # GameState struct（ゲームの全状態）
│   ├── player.go        # Player struct・Vim コマンドのロジック
│   ├── enemy.go         # Enemy インターフェース・Builder・Strategy
│   ├── stage.go         # Stage struct・ステージ定義
│   ├── map.go           # Grid・Cell・マップ読み込み
│   └── files/
│       └── stage/       # マップファイル（map01.txt〜）※go:embed のため state/ 配下に配置
├── input/
│   └── input.go         # キー入力 → Command 型への変換
├── renderer/
│   └── renderer.go      # Ebitengine を使った描画ロジック
├── files/
│   └── fonts/           # 埋め込みフォント（renderer パッケージから参照予定）
├── docs/
│   ├── decisions/       # ADR（Architecture Decision Records）
│   └── design/          # 設計ドキュメント
└── Makefile
```

---

## アーキテクチャの核心原則

### 原則1. state パッケージは ebiten を import しない

ゲームの状態は `state` パッケージの純粋な Go struct として持ちます。
`ebiten` への依存は `renderer` パッケージにのみ許可します。

```
state パッケージ           renderer パッケージ
┌──────────────────┐      ┌───────────────────────┐
│  GameState       │ ───> │  Renderer             │
│  Player          │      │  ebiten.Image を操作   │
│  Enemy           │      │  state を読む（書かない）│
│  Grid / Cell     │      └───────────────────────┘
│                  │
│  ebiten を        │
│  import しない    │
└──────────────────┘
```

`state` パッケージの関数・メソッドは `ebiten` の型を引数・戻り値に含めてはいけません。
この制約により `go test` で単体テストが書けます。

### 原則2. Update() は状態を変える。Draw() は状態を読むだけ。

```go
// game.go
func (g *Game) Update() error {
    cmd, ch := g.input.Read()
    return g.state.Update(cmd, ch)
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.renderer.Draw(screen, g.state) // state を変更しない
}
```

`Draw()` 内で `state` を変更するコードを書いてはいけません。

### 原則3. 入力は Command 型に変換してから state に渡す

`state` パッケージは `ebiten.Key` を直接知りません。
`input.Handler.Read()` が毎フレーム呼ばれ、キー状態を `Command` 型に変換して返します。

```go
// state/player.go は Command だけを受け取る
func (p *Player) Apply(cmd input.Command, ch rune, stage *Stage) { ... }
```

---

## 主要な型定義

### Grid と Cell（`state/map.go`）

```go
type CellKind int

const (
    CellSpace       CellKind = iota
    CellWall
    CellBoundary
    CellApple
    CellAppleEaten
    CellPoison
)

type Grid struct {
    Cells  [][]Cell
    Width  int
    Height int
}

// 範囲外は壁として返す（呼び出し側でガード不要）
func (g *Grid) At(x, y int) Cell { ... }

func (g *Grid) Set(x, y int, kind CellKind) { ... }

func (g *Grid) IsWalkable(x, y int) bool {
    k := g.At(x, y).Kind
    return k == CellSpace || k == CellApple || k == CellAppleEaten || k == CellPoison
}
```

### GameState（`state/gamestate.go`）

```go
type GameState struct {
    Stages    []Stage
    StageIdx  int
    Player    *Player
    Enemies   []Enemy
    Life      int
    EnemyTick int
    Phase     GamePhase
}

type GamePhase int

const (
    PhaseOpening GamePhase = iota
    PhasePlaying
    PhaseStageClear
    PhaseGameOver
    PhaseGameClear
)

func (gs *GameState) Update(cmd input.Command, ch rune) error { ... }
```

### Player（`state/player.go`）

```go
type Player struct {
    X, Y        int
    Score       int
    TargetScore int
    State       PlayerState
    inputNum    int      // カウント蓄積（3w の "3" など）
    inputG      bool     // "g" 入力待ち（gg コマンド用）
    lastFind    *findCmd // f/t の繰り返し用（; , コマンド）
}

type PlayerState int

const (
    PlayerAlive PlayerState = iota
    PlayerDead
    PlayerWon
)

func (p *Player) Apply(cmd input.Command, ch rune, stage *Stage) { ... }

// walk: 経路上の全セルで判定（リンゴ連続取得・敵接触で死亡）
func (p *Player) walkTo(targetX, targetY int, stage *Stage) { ... }

// jump: 目的地のみ判定（敵・壁を飛び越える）
func (p *Player) jumpTo(targetX, targetY int, stage *Stage) { ... }

func (p *Player) repeatCount() int {
    if p.inputNum == 0 { return 1 }
    return p.inputNum
}

func (p *Player) resetInput() {
    p.inputNum = 0
    p.inputG = false
}
```

### Command（`input/input.go`）

```go
type Command int

const (
    CmdNone Command = iota
    CmdNum            // 数字キー（0〜9、inputNum に蓄積）
    CmdMoveLeft       // h
    CmdMoveDown       // j
    CmdMoveUp         // k
    CmdMoveRight      // l
    CmdWordForward    // w
    CmdWordEnd        // e
    CmdWordBack       // b
    CmdLineBegin      // 0
    CmdLineEnd        // $
    CmdLineFirstWord  // ^
    CmdFileTop        // gg
    CmdFileBottom     // G
    CmdFindForward    // f{char} — ch に対象文字が入る
    CmdFindBack       // F{char}
    CmdTillForward    // t{char}
    CmdTillBack       // T{char}
    CmdRepeatFind     // ;
    CmdRepeatFindRev  // ,
    CmdBlockForward   // }
    CmdBlockBack      // {
    CmdHigh           // H
    CmdMiddle         // M
    CmdLow            // L
    CmdWordEndBack    // ge
    CmdQuit           // q
)

type Handler struct {
    awaitChar bool // f/t 入力後、次フレームで文字を受け取る
    pendingCmd Command
    prevG      bool // g 入力後、次フレームで g を待つ
}

// Read は毎フレーム呼ばれる。2ストロークの状態管理はここで完結させる。
func (h *Handler) Read() (Command, rune) { ... }
```

### Enemy（`state/enemy.go`）

```go
type EnemyKind int

const (
    EnemyHunter EnemyKind = iota
    EnemyGhost
)

type Enemy interface {
    Position() (x, y int)
    SetPosition(x, y int)
    Think(p *Player, grid *Grid) (x, y int)
    Move(x, y int, grid *Grid)
    HasCaptured(p *Player) bool
    Kind() EnemyKind
}

type Strategy interface {
    Eval(p *Player, x, y int) float64
}

// 実装: AssaultStrategy（直線距離最小化）、TrickyStrategy（20%ランダム）
```

---

## Vim コマンドの実装方針

### walk と jump の使い分け

| タイプ | コマンド | 挙動 |
|---|---|---|
| walk | h j k l w e b f t F T | 経路上の全セルで判定。リンゴを連続取得できる |
| jump | 0 $ ^ gg G H M L { } | 目的地のみ判定。壁・敵を飛び越える |

### カウント入力

`inputNum` に蓄積し、コマンド実行後に `resetInput()` で必ずリセットします。
`0` 単独は `CmdLineBegin` として扱い、`inputNum > 0` のときの `0` だけ数字として蓄積します。
この振り分けは `input.Handler.Read()` 側で行います。

### 2ストロークコマンド

`input.Handler` 内で状態を持ちます。`state` パッケージに入力状態を漏らしません。

- `gg`: 1フレーム目で `prevG = true`、2フレーム目で `g` が来たら `CmdFileTop` を返す
- `f{char}`: 1フレーム目で `awaitChar = true`、2フレーム目で文字を `ch` に乗せて返す
- `ge`: `g` + `e` のシーケンスで `CmdWordEndBack` を返す

---

## 描画の方針（`renderer/renderer.go`）

```go
const (
    TileSize     = 16
    OffsetX      = 40  // 行番号エリアの幅
    OffsetY      = 0
    ScreenWidth  = 800
    ScreenHeight = 600
)

func gridToScreen(gx, gy int) (float64, float64) {
    return float64(OffsetX + gx*TileSize), float64(OffsetY + gy*TileSize)
}
```

### Draw() 内の描画順

1. 背景
2. グリッド（壁・境界・リンゴ・毒・通路）
3. 敵
4. プレイヤー
5. UI（スコア・ライフ・レベル・操作ヒント）

---

## ステージの定義（`state/stage.go`）

```go
type Stage struct {
    Level        int
    MapPath      string
    Grid         *Grid
    HunterConfig EnemyConfig
    GhostConfig  *EnemyConfig  // nil のステージもある
    GameSpeed    time.Duration
}

func InitStages() []Stage { ... }
```

### マップファイルの文字規則（`state/files/stage/*.txt`）

```
+   境界線（マップの外枠）→ CellBoundary
!   障害物（縦・横・斜め接続部）→ CellWall
-   障害物（水平線）→ CellWall
|   障害物（垂直線）→ CellWall
o   リンゴ（収集対象）→ CellApple
X   毒 → CellPoison
P   プレイヤー初期位置（読み込み後は CellSpace）
H   ハンター初期位置（読み込み後は CellSpace）
G   ゴースト初期位置（読み込み後は CellSpace）
    スペース（通路）→ CellSpace
```

---

## 敵の移動タイミング

Ebitengine のループは TPS=60 固定のため、フレームカウンタで `GameSpeed` を実現します。

```go
// GameState.Update() 内
gs.EnemyTick++
if time.Duration(gs.EnemyTick)*(time.Second/60) >= gs.Stage().GameSpeed {
    gs.moveEnemies()
    gs.EnemyTick = 0
}
```

---

## テストの書き方

`state` パッケージは `ebiten` に依存しないため、`go test` で単体テストが書けます。
`Grid` をテスト用に組み立てて `Player.Apply()` を呼ぶパターンを基本とします。

```go
func TestPlayerWalkRight(t *testing.T) {
    grid := buildGrid([]string{
        "++++",
        "+ Po",
    })
    stage := &state.Stage{Grid: grid}
    p := &state.Player{X: 2, Y: 1, State: state.PlayerAlive}
    p.Apply(input.CmdMoveRight, 0, stage)
    if p.X != 3 || p.Score != 0 {
        t.Errorf("got X=%d Score=%d", p.X, p.Score)
    }
    p.Apply(input.CmdMoveRight, 0, stage)
    if p.X != 4 || p.Score != 1 {
        t.Errorf("expected apple eaten, got X=%d Score=%d", p.X, p.Score)
    }
}
```

---

## 新しいコマンドを追加する手順

1. `input/input.go` の `Command` 型に定数を追加する
2. `input.Handler.Read()` でキーを `Command` にマッピングする
3. `state/player.go` の `Apply()` に `case` を追加してロジックを実装する
4. `state/player_test.go` にテストを追加する
5. 設計判断が生じた場合は `kizami adr` で ADR を作成する

---

## ドキュメント管理（kizami）

このプロジェクトは **[kizami](https://github.com/mskasa/kizami)** を使って設計判断を記録します。

### ADR を書くタイミング

- ライブラリ・フレームワークの採用・変更
- パッケージ構成やアーキテクチャの変更
- 複数ファイルにまたがる設計判断
- 既存の設計を廃止・置き換える変更

1ファイル内に収まる理由はコードコメントで十分です。ADR は不要です。

### 基本コマンド

```bash
kizami adr "タイトル"          # ADR を作成
kizami adr --ai "タイトル"     # AI にドラフトを生成させる（git add 後に実行）
kizami design "タイトル"       # 設計ドキュメントを作成
kizami list                    # 一覧
kizami search "キーワード"     # 検索
kizami blame <file>            # ファイルを参照する ADR を逆引き
kizami audit                   # Related Files とコードの乖離を検出
```

### ステータス

| ステータス | 意味 |
|---|---|
| `Draft` | 作成直後・実装前 |
| `Active` | 実装済み・現在有効 |
| `Inactive` | 無効化（代替なし） |
| `Superseded by <slug>` | 別の ADR に置き換えられた |

### ADR 一覧

| ファイル | 内容 |
|---|---|
| `0001-rewrite-from-scratch-with-ebitengine.md` | termbox 版を捨て Ebitengine で1から書き直した理由 |
| `0002-separate-state-and-renderer.md` | state パッケージを ebiten から切り離した理由 |
| `0003-abstract-input-as-command-type.md` | キー入力を Command 型で抽象化した理由 |

---

## 開発コマンド

```bash
make run          # 実行
make test         # テスト
make build        # macOS バイナリ
make build-win    # Windows バイナリ
make wasm         # WebAssembly ビルド（docs/wasm/ に出力）
```

---

## 実装ロードマップ

### Phase 1：state パッケージ（termbox 不要・テスト可能）

- [x] `state/map.go`：`Grid` / `Cell` / マップ読み込み
- [x] `state/player.go`：`Player` と基本コマンド（h j k l w e b 0 $ ^ gg G）
- [x] `state/enemy.go`：`Enemy` / Builder / Strategy（Assault / Tricky）
- [x] `state/stage.go`：`Stage` とステージ定義
- [x] `state/gamestate.go`：`GameState.Update()`
- [x] `state/*_test.go`：上記のユニットテスト

### Phase 2：Ebitengine への接続

- [x] `go.mod`：`github.com/hajimehoshi/ebiten/v2` を追加
- [x] `input/input.go`：キー入力 → `Command` 変換
- [x] `renderer/renderer.go`：グリッド / プレイヤー / 敵 / UI の描画
- [x] `game.go`：`ebiten.Game` の実装
- [x] `main.go`：`ebiten.RunGame` で起動

### Phase 3：機能拡張

- [ ] `f F t T ; ,` コマンド
- [ ] `{ }` コマンド
- [ ] `H M L` コマンド
- [ ] `ge gE` コマンド
- [ ] スプライト画像・SE
- [ ] WebAssembly 配布
