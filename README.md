# PacVim is a Vim learning game ʕ◔ϖ◔ʔ

![pacvim](https://github.com/masahiro-kasatani/pacvim/blob/readme-images/files/readme.png?raw=true)

<p align="right">
The Go gopher was designed by <a href="https://go.dev/blog/gopher" target="_blank">Renée French</a>.
</p>

[![test](https://github.com/masahiro-kasatani/pacvim/actions/workflows/test.yaml/badge.svg)](https://github.com/masahiro-kasatani/pacvim/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/masahiro-kasatani/pacvim)](https://goreportcard.com/report/github.com/masahiro-kasatani/pacvim)
[![codecov](https://codecov.io/gh/masahiro-kasatani/pacvim/branch/master/graph/badge.svg?token=KZ2LVX4GCT)](https://codecov.io/gh/masahiro-kasatani/pacvim)

| English | [日本語](https://github.com/masahiro-kasatani/pacvim/blob/master/README-JA.md) |

<!-- TOC -->

## Table of Contents

- [How to start PacVim](#how-to-start-pacvim)
- [PacVim Rules](#pacvim-rules)
- [Player Controls](#player-controls)
  - [About action type](#about-action-type)
- [License](#license)
- [Author](#author)

<!-- /TOC -->

## How to start PacVim

PacVim is started by double-clicking on the binary file below.

- Windows: [./bin/win/pacvim.exe](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/win)
- Mac: [./bin/mac/pacvim](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/mac)

## PacVim Rules

PacVim follows the rules of Pac-Man.

TODO ゲーム画面の画像貼り、各オブジェクトの説明を入れる

| State         | To transition to the left state |
| :------------ | :------------------------------ |
| Stage clear   | Eat all apples                  |
| Stage failure | Caught by enemy or eat poison   |
| Game clear    | Clear all stages                |
| Game over     | Stage failure with 0 life.      |

## Player Controls

|  key  | action                                                      | action type |
| :---: | :---------------------------------------------------------- | :---------- |
|   h   | move left                                                   | walk        |
|   j   | move down                                                   | walk        |
|   k   | move up                                                     | walk        |
|   l   | move right                                                  | walk        |
|   w   | move forward to next word beginning                         | walk        |
|   e   | move forward to next word ending                            | walk        |
|   b   | move backward to previous word beginning                    | walk        |
|   0   | move to the beginning of the current line                   | jump        |
|   $   | move to the end of the current line                         | jump        |
|   ^   | move to the beginning of the first word on the current line | jump        |
| gg/1G | move to the beginning of the first word on the first line   | jump        |
|   G   | move to the beginning of the first word on the last line    | jump        |
|  NG   | move to the beginning of the first word on the nth line     | jump        |
|   q   | quit the game                                               | -           |

### About action type

- walk

- jump

## License

MIT

## Author

[Masahiro Kasatani](https://masahiro-kasatani.github.io/portfolio/)
