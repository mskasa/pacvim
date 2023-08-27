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
  - [Game screen](#game-screen)
  - [About objects](#about-objects)
  - [About the state of the game](#about-the-state-of-the-game)
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

### Game screen

![game screen](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/screen.png)

### About objects

| Object name   |                                                                                                                                                         Display                                                                                                                                                         | Supplementary explanation               |
| :------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-------------------------------------- |
| apple         |                                                       ![apple1](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_1.png) ![apple2](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_2.png)                                                       | This turns green when eaten.            |
| poison        |                                                                                                          ![poison](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/poison.png)                                                                                                           | -                                       |
| obstacles     | ![obstacle1](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_1.png) ![obstacle2](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_2.png) ![obstacle3](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_3.png) | -                                       |
| player        |                                                                                                          ![player](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/player.png)                                                                                                           | -                                       |
| Enemy(hunter) |                                                                                                          ![hunter](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/hunter.png)                                                                                                           | -                                       |
| Enemy(ghost)  |                                                                                                           ![ghost](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/ghost.png)                                                                                                            | Enemies that can slip through obstacles |

### About the state of the game

| State         | To transition to the left state |
| :------------ | :------------------------------ |
| Stage clear   | Eat all apples                  |
| Stage failure | Caught by enemy or eat poison   |
| Game clear    | Clear all stages                |
| Game over     | Stage failure with 0 life.      |

## Player Controls

| Key  | Action type | Action                                                      |
| :--: | :---------- | :---------------------------------------------------------- |
| `h`  | `walk`      | move left                                                   |
| `j`  | `walk`      | move down                                                   |
| `k`  | `walk`      | move up                                                     |
| `l`  | `walk`      | move right                                                  |
| `w`  | `walk`      | move forward to next word beginning                         |
| `e`  | `walk`      | move forward to next word ending                            |
| `b`  | `walk`      | move backward to previous word beginning                    |
| `0`  | `jump`      | move to the beginning of the current line                   |
| `$`  | `jump`      | move to the end of the current line                         |
| `^`  | `jump`      | move to the beginning of the first word on the current line |
| `gg` | `jump`      | move to the beginning of the first word on the first line   |
| `G`  | `jump`      | move to the beginning of the first word on the last line    |
| `NG` | `jump`      | move to the beginning of the first word on the nth line     |
| `q`  | -           | quit the game                                               |

### About action type

- `walk`

  - `walk` is the image of moving one square at a time to the destination in an instant. Therefore, hit detection with enemies, obstacles, and apples is applied. Use it when you want to eat apples all at once.

    - e.g. If you type `w`.

      ![walk example](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-w.gif)

- `jump`

  - `jump` is the image of jumping between to the destination and reaching it in an instant. Therefore, hit detection with enemies, obstacles, and apples is not applied. Use it when you want to move to avoid enemies or obstacles.

    - e.g. If you type `$`.

      ![jump example](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-doller.gif)

## License

MIT

## Author

[Masahiro Kasatani](https://masahiro-kasatani.github.io/portfolio/)
