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
)

// Handler はキー入力を毎フレーム読み取り Command に変換する。
// 2ストロークコマンド（gg, ge, f/t）の状態はここで管理する。
type Handler struct {
	prevG      bool    // g を受け取った後、次のキーを待つ
	awaitChar  bool    // f/F/t/T の後、対象文字を待つ
	pendingCmd Command // awaitChar 中の保留コマンド
	inNumMode  bool    // 数字蓄積中かどうか（0 の振り分けに使用）
}

// Read は1フレーム分のキー入力を読み取り、(Command, rune) を返す。
// CmdNum のとき rune は押された数字文字。
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
			h.inNumMode = false
			return cmd, chars[0]
		}
		return CmdNone, 0
	}

	shift := ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)

	// prevG モード: gg または ge の2文字目を待つ
	if h.prevG {
		if !shift && inpututil.IsKeyJustPressed(ebiten.KeyG) {
			h.prevG = false
			h.inNumMode = false
			return CmdFileTop, 0
		}
		if !shift && inpututil.IsKeyJustPressed(ebiten.KeyE) {
			h.prevG = false
			h.inNumMode = false
			return CmdWordEndBack, 0
		}
		return CmdNone, 0
	}

	// 数字キー（1〜9 は常に CmdNum）
	digits := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4, ebiten.Key5,
		ebiten.Key6, ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	for i, k := range digits {
		if !shift && inpututil.IsKeyJustPressed(k) {
			h.inNumMode = true
			return CmdNum, rune('1' + i)
		}
	}
	// 0: 数字蓄積中なら CmdNum、単独なら CmdLineBegin
	if !shift && inpututil.IsKeyJustPressed(ebiten.Key0) {
		if h.inNumMode {
			return CmdNum, '0'
		}
		return CmdLineBegin, 0
	}

	// 以降のコマンドは数字蓄積をリセット
	h.inNumMode = false

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
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyH):
			return CmdMoveLeft, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyJ):
			return CmdMoveDown, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyK):
			return CmdMoveUp, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyL):
			return CmdMoveRight, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyW):
			return CmdWordForward, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyE):
			return CmdWordEnd, 0
		case inpututil.IsKeyJustPressed(ebiten.KeyB):
			return CmdWordBack, 0
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
		}
	}

	return CmdNone, 0
}
