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
| walk | h j k l w e b f t F T | 経路上の **全ステップ** で `checkCell` を呼ぶ。リンゴを連続取得、毒・敵に接触したら即死 |
| jump | 0 $ ^ gg G H M L { } | **目的地のみ** `checkCell` を呼ぶ。経路上の壁・敵・毒は無視する |

実装上のルール:
- walk は `walkTo` を繰り返し呼ぶ。`walkTo` は内部で `checkCell(stage, enemies)` を呼ぶ。
- jump は `jumpTo` を直接呼ぶ。`jumpTo` は `IsWalkable` で目的地を確認してから `checkCell(stage, enemies)` を呼ぶ。
- `Apply` の引数に `enemies []Enemy` を渡すことで、walk 中の各ステップで敵接触を判定できる。
- `checkCell` は敵 → セル種別の順に判定する。敵に接触した時点で即リターンする。

```go
// walk: 1ステップ移動して判定
func (p *Player) walkTo(targetX, targetY int, stage *Stage, enemies []Enemy) bool

// jump: 目的地へ直接移動して判定（経路上のセル・敵は無視）
func (p *Player) jumpTo(targetX, targetY int, stage *Stage, enemies []Enemy)

// checkCell: 現在地の敵接触・セル種別を判定
func (p *Player) checkCell(stage *Stage, enemies []Enemy)
```

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

## 入力実装の注意事項（`input/input.go`）

### `IsKeyJustPressed` と `IsKeyPressed` の使い分け

Ebitengine の入力 API には2種類あります。

| API | 挙動 | 用途 |
|---|---|---|
| `inpututil.IsKeyJustPressed(key)` | 押した瞬間の1フレームのみ true | 1回だけ発火させたいコマンド |
| `ebiten.IsKeyPressed(key)` | 押している間ずっと true | リピート処理の継続判定 |

移動キーに `IsKeyJustPressed` のみを使うと、押し続けても最初の1回しか反応しません。
**移動キーは必ずリピート処理を経由して発火させてください。**

### キーリピートの実装

`Handler` にリピート状態（`repeatKey`・`repeatFrames`）を持たせ、以下の挙動を実現しています。

- 押した瞬間 → 即座に発火（`IsKeyJustPressed` で検出）
- 押し続けて `repeatInitialDelay` フレーム（約300ms）経過後 → `repeatInterval` フレーム（約60ms）ごとに繰り返し発火
- キーを離したら → `repeatKey` と `repeatFrames` をリセット

### リピート対象のキー

| 区分 | キー | 理由 |
|---|---|---|
| **対象（連続発火）** | `h` `j` `k` `l` `w` `e` `b` | 押し続けでカーソルを連続移動させたい |
| **対象外（1回のみ）** | `g`（gg の1打目）、`q`、`f/F/t/T`（文字待ち）、数字キー | 誤操作を防ぐ・2ストロークの状態管理を壊さない |

### 数値入力（カウントプレフィックス）の責務分離

`input` と `state` の責務は明確に分離します。

**`input` 側の責務：**
数字キー（`0`〜`9`）が押されたら常に `CmdNum` を返す。`0` 単独か数字蓄積中かの判断は行わない。

```go
// 0 も含めて常に CmdNum として返す
if !shift && inpututil.IsKeyJustPressed(ebiten.Key0) {
    return CmdNum, '0'
}
```

**`state` 側の責務：**
`CmdNum` を受け取ったとき、`0` かつ `inputNum == 0` であれば `CmdLineBegin` として処理する。
それ以外は `inputNum` に桁を加算して蓄積する。

```go
case input.CmdNum:
    if ch == '0' && p.inputNum == 0 {
        p.jumpTo(p.lineBeginX(stage), p.Y, stage, enemies)
        break
    }
    p.inputNum = p.inputNum*10 + int(ch-'0')
    return
```

各コマンド実行時は `repeatCount()` で繰り返し回数を取得し、実行後は必ず `resetInput()` を呼ぶ。

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
    Theme        string   // 学習テーマ（例: "基本移動"）
    Commands     []string // このステージで練習するコマンドの説明（PhaseReady 画面に表示）
}

func InitStages() []Stage { ... }
```

### ステージの学習テーマ

| Level | Theme | 練習コマンド | 封印コマンド |
|---|---|---|---|
| 1 | Basic Movement | `h` `l` `j` `k` | なし |
| 2 | Word Motion | `w` `e` `b` | `h` `l` |
| 3 | Line & File Motion | `0` `$` `^` `gg` `G` | なし |
| 4 | Find Character | `f{c}` `t{c}` `;` `,` | `w` `e` `b` |
| 5 | All Commands | 全コマンド | なし |

`Theme` と `Commands` は `PhaseReady`（準備画面）に表示される。
新しいステージを追加するときはこれらのフィールドも設定すること。

### ステージ設計の原則

#### 原則1. そのステージで学ぶコマンドを多用しないとクリアできないこと

テーマコマンド（`Theme` / `Commands`）を使わなければ全リンゴを取得できない配置にすること。
テーマ外のコマンドだけでクリアできる抜け道を作らないこと。

確認手順として、**テーマコマンドを封印した状態でプレイし、クリアできないことを確認する**。

**例外：1面・2面はこの原則を適用しない。**
移動コマンド自体がテーマであるため、テーマコマンドを封印すると移動そのものができなくなる。
1・2面は「テーマコマンドを使うと効率よくクリアできる」設計を目指せばよい。

**コマンドの封印はマップ設計より優先度が高い強制手段である。**
マップ設計だけでは差別化が困難なコマンド（例：`f/t` と `h/l`）は、封印によって学習文脈をシステムで担保する。
封印するコマンドは `Stage.RestrictedCmds []input.Command` フィールドで定義する（「ステージの学習テーマ」テーブルの「封印コマンド」列を参照）。
封印されているコマンドを押したときは UI にフィードバックを表示し（「このステージでは使えません」など）、プレイヤーに事前告知すること。
新しいステージを追加するときは `RestrictedCmds` を明示的に設定すること（封印なしの場合も意図的に空にしていることを確認すること）。

#### 原則2. ゲームとして適切な難易度になっていること

以下の要素で難易度を調整できる。

- `GameSpeed`（敵の移動間隔）
- 敵の数と種類
- リンゴの配置密度
- マップの広さと壁の量

難易度の目安：

- **1〜2面**：Vim 初心者でも詰まらず進める。敵は少なく動作は遅め。
- **3〜4面**：テーマコマンドを意識しないと取りにくいリンゴが増える。
- **5面**：全コマンドを駆使する必要がある総合面。

#### 原則3. 敵の特性を活かしていること

- **Hunter**：壁を通り抜けられない。壁で区切られたエリアに配置すると追い詰め効果が高い。
- **Ghost**：壁を通り抜けられる。壁がないステージでは Hunter と同じ動作になるため意味がない。

ゴーストを使う場合は必ずマップに壁を配置すること。
敵の配置はプレイヤーが「逃げる」「回避する」動作が発生するよう意識すること。

### ステージを変更するときのチェックリスト

- [ ] テーマコマンドを封印した状態でプレイして、クリアできないことを確認した（1・2面は除く）
- [ ] `GameSpeed`・敵の数・リンゴの配置密度・マップの広さが、そのステージの難易度目安に合っている
- [ ] ゴーストを使う場合、マップに壁が配置されている

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
| `2026-03-30-ebitengine.md` | Ebitengine のキーリピート実装方針（ソフトウェアリピートの採用理由） |
| `2026-03-30-num-input.md` | 数値入力（カウントプレフィックス）の責務分離（0キーの解釈を state 側で行う） |
| `2026-03-30-stage-learning-metadata.md` | ステージへの学習メタデータ埋め込み方針（Theme・Commands を Stage struct に持たせる） |
| `2026-03-30-stage-learning-metadata.md` | ステージへの学習メタデータ埋め込み方針（Theme・Commands を Stage struct に持たせる） |

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
- [ ] ステージごとの使用可能コマンド制限
- [ ] スプライト画像・SE
- [ ] WebAssembly 配布
