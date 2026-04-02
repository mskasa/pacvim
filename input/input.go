// Package input はキー入力を Command 型に変換する責務を持つ。
// Command 型は state パッケージに渡される。ebiten.Key への依存はこのパッケージに閉じる。
package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Command はプレイヤーへの操作コマンドを表す。
type Command int

const (
	CmdNone          Command = iota
	CmdNum                  // 数字キー（0〜9、inputNum に蓄積）
	CmdMoveLeft             // h
	CmdMoveDown             // j
	CmdMoveUp               // k
	CmdMoveRight            // l
	CmdWordForward          // w
	CmdWordEnd              // e
	CmdWordBack             // b
	CmdLineBegin            // 0（単独）
	CmdLineEnd              // $
	CmdLineFirstWord        // ^
	CmdFileTop              // gg
	CmdFileBottom           // G
	CmdFindForward          // f{char}
	CmdFindBack             // F{char}
	CmdTillForward          // t{char}
	CmdTillBack             // T{char}
	CmdRepeatFind           // ;
	CmdRepeatFindRev        // ,
	CmdBlockForward         // }
	CmdBlockBack            // {
	CmdHigh                 // H
	CmdMiddle               // M
	CmdLow                  // L
	CmdWordEndBack          // ge
	CmdQuit                 // q
	CmdConfirm              // Enter
)

// キーリピートのタイミング定数（TPS=60 基準）
const (
	repeatInitialDelay = 18 // 最初のリピートまでのフレーム数（約300ms）
	repeatInterval     = 4  // リピート間隔のフレーム数（約60ms）
)

// Handler はキー入力を毎フレーム読み取り Command に変換する。
// 2ストロークコマンド（gg, ge, f/t）の状態はここで管理する。
// 移動キーのリピート状態もここで管理する。
type Handler struct {
	prevG        bool       // g を受け取った後、次のキーを待つ
	awaitChar    bool       // f/F/t/T の後、対象文字を待つ
	pendingCmd   Command    // awaitChar 中の保留コマンド
	repeatKey    ebiten.Key // 現在リピート中のキー
	repeatFrames int        // repeatKey を押し続けたフレーム数
}

// PendingDisplay は現在の入力待ち状態を表示用文字列で返す。
// prevG モード中は "g"、awaitChar モード中は "f"/"F"/"t"/"T" を返す。
func (h *Handler) PendingDisplay() string {
	if h.prevG {
		return "g"
	}
	if h.awaitChar {
		switch h.pendingCmd {
		case CmdFindForward:
			return "f"
		case CmdFindBack:
			return "F"
		case CmdTillForward:
			return "t"
		case CmdTillBack:
			return "T"
		}
	}
	return ""
}

// Read は1フレーム分のキー入力を読み取り、(Command, rune) を返す。
// CmdNum のとき rune は押された数字文字（0〜9）。0 単独か数字蓄積中かの判断は state 側で行う。
// CmdFindForward/Back/TillForward/Back のとき rune は対象文字。
// それ以外の場合 rune は 0。
func (h *Handler) Read() (Command, rune) {
	// awaitChar モード: f/F/t/T の対象文字を待つ
	if h.awaitChar {
		chars := ebiten.AppendInputChars(nil)
		if len(chars) > 0 {
			cmd := h.pendingCmd
			h.awaitChar = false
			h.pendingCmd = CmdNone
			return cmd, chars[0]
		}
		return CmdNone, 0
	}

	shift := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)

	// prevG モード: gg または ge の2文字目を待つ
	if h.prevG {
		if !shift && inpututil.IsKeyJustPressed(ebiten.KeyG) {
			h.prevG = false
			return CmdFileTop, 0
		}
		if !shift && inpututil.IsKeyJustPressed(ebiten.KeyE) {
			h.prevG = false
			return CmdWordEndBack, 0
		}
		// G キーが Shift より 1 フレーム先に検出された場合、次フレームで
		// Shift+G として判定できるよう IsKeyPressed で確認してフリーズを防ぐ。
		if shift && ebiten.IsKeyPressed(ebiten.KeyG) {
			h.prevG = false
			return CmdFileBottom, 0
		}
		return CmdNone, 0
	}

	// 数字キー（0〜9）は常に CmdNum として返す。
	// 0 単独か数字蓄積中かの判断は state 側（Player.Apply）で行う。
	digits := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5,
		ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for i, k := range digits {
		if !shift && inpututil.IsKeyJustPressed(k) {
			return CmdNum, rune('1' + i)
		}
	}
	if !shift && inpututil.IsKeyJustPressed(ebiten.Key0) {
		return CmdNum, '0'
	}

	if shift {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyG):
			return CmdFileBottom, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyH):
			return CmdHigh, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyM):
			return CmdMiddle, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyL):
			return CmdLow, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyF):
			h.awaitChar = true
			h.pendingCmd = CmdFindBack
			return CmdNone, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyT):
			h.awaitChar = true
			h.pendingCmd = CmdTillBack
			return CmdNone, 0
		case inpututil.IsKeyJustPressed(ebiten.Key4): // $
			return CmdLineEnd, 0
		case inpututil.IsKeyJustPressed(ebiten.Key6): // ^
			return CmdLineFirstWord, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyLeftBracket): // {
			return CmdBlockBack, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyRightBracket): // }
			return CmdBlockForward, 0
		}
	} else {
		// 移動キーはリピート処理を経由して発火する
		repeatableKeys := []struct {
			key ebiten.Key
			cmd Command
		}{
			{ebiten.KeyH, CmdMoveLeft},
			{ebiten.KeyJ, CmdMoveDown},
			{ebiten.KeyK, CmdMoveUp},
			{ebiten.KeyL, CmdMoveRight},
			{ebiten.KeyW, CmdWordForward},
			{ebiten.KeyE, CmdWordEnd},
			{ebiten.KeyB, CmdWordBack},
		}
		for _, rk := range repeatableKeys {
			if ebiten.IsKeyPressed(rk.key) {
				if inpututil.IsKeyJustPressed(rk.key) {
					h.repeatKey = rk.key
					h.repeatFrames = 0
					return rk.cmd, 0
				}
				if h.repeatKey == rk.key {
					h.repeatFrames++
					elapsed := h.repeatFrames - repeatInitialDelay
					if elapsed >= 0 && elapsed%repeatInterval == 0 {
						return rk.cmd, 0
					}
				}
				return CmdNone, 0
			}
		}
		// リピート中のキーが離されたらリセット
		h.repeatKey = 0
		h.repeatFrames = 0

		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyG):
			h.prevG = true
			return CmdNone, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyF):
			h.awaitChar = true
			h.pendingCmd = CmdFindForward
			return CmdNone, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyT):
			h.awaitChar = true
			h.pendingCmd = CmdTillForward
			return CmdNone, 0
		case inpututil.IsKeyJustPressed(ebiten.KeySemicolon):
			return CmdRepeatFind, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyComma):
			return CmdRepeatFindRev, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyQ):
			return CmdQuit, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
			return CmdConfirm, 0
		}
	}

	return CmdNone, 0
}
