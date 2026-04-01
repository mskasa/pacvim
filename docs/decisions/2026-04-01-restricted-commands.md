# ステージごとに使用可能コマンドを制限する

- Date: 2026-04-01
- Type: ADR
- Status: Draft
- Author: masahiro.kasatani

## Context

`f/t` コマンドは walk だが、マップ設計だけでは `w` や `h/l` との差別化が困難であった。
また、3面でも「壁がなければ `h/l` だけでクリアできる」という問題があり、
マップ設計で学習を強制しようとすると不自然な制約が生まれやすかった。

「使えないキーをシステムで封印する」ことで、マップ設計に頼らず学習文脈を強制できるようにする。

## Decision

- `Stage` struct に `RestrictedCmds []input.Command` フィールドを追加する
- `GameState.Update()` 内で、受け取った `Command` が `RestrictedCmds` に含まれる場合は無視する
- 制限されているキーを押したときは UI にフィードバックを表示する（「このステージでは使えません」など）
- 各ステージの制限は以下のとおり

| 面 | 封印するコマンド |
|---|---|
| 1面 | なし |
| 2面 | h l（CmdMoveLeft・CmdMoveRight） |
| 3面 | h l w e b（CmdMoveLeft・CmdMoveRight・CmdWordForward・CmdWordEnd・CmdWordBack） |
| 4面 | w e b（CmdWordForward・CmdWordEnd・CmdWordBack） |
| 5面 | なし（全コマンド解禁） |

## Consequences

- マップ設計だけに頼らず、学習文脈をシステムで担保できる
- プレイヤーへの事前告知 UI が必要になる（封印されているキーを画面に表示する）
- 実際の vim にはない制約のため「vim の正確な再現」からは離れる設計である

## Alternatives Considered

- マップ設計のみで学習を強制する（不自然な壁配置が必要になり、ゲームとしての自然さが損なわれる）
- ヒント表示のみにとどめる（任意であるため学習の強制力がない）

## Related Files

- state/stage.go
- state/gamestate.go
- renderer/renderer.go
- CLAUDE.md
