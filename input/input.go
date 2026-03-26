// Package input はキー入力を Command 型に変換する責務を持つ。
// Command 型は state パッケージに渡され、ebiten.Key には依存しない。
// input.Handler による ebiten キーマッピングは Phase 2 で実装する。
package input

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
	CmdFindForward          // f{char} — ch に対象文字が入る
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
