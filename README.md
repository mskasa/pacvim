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

- [For those who want to play with PacVim](#for-those-who-want-to-play-with-pacvim)
  - [How to start PacVim](#how-to-start-pacvim)
  - [PacVim Rules](#pacvim-rules)
    - [Game screen](#game-screen)
    - [About objects](#about-objects)
    - [About the state of the game](#about-the-state-of-the-game)
  - [Player Controls](#player-controls)
    - [About action type](#about-action-type)
- [For those who want to develop PacVim](#for-those-who-want-to-develop-pacvim)
  - [Commands for development](#commands-for-development)
  - [Execution options](#execution-options)
  - [How to customize PacVim](#how-to-customize-pacvim)
    - [How to add a stage map](#how-to-add-a-stage-map)
    - [How to add enemy types](#how-to-add-enemy-types)
    - [How to add the enemy's thought logic](#how-to-add-the-enemys-thought-logic)
- [License](#license)
- [Author](#author)

<!-- /TOC -->

## For those who want to play with PacVim

### How to start PacVim

PacVim is started by double-clicking on the binary file below.

- Windows: [./bin/win/pacvim.exe](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/win)
- Mac: [./bin/mac/pacvim](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/mac)

### PacVim Rules

PacVim follows the rules of Pac-Man.

#### Game screen

![game screen](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/screen.png)

#### About objects

| Object name   |                                                                                                                                                         Display                                                                                                                                                         | Supplementary explanation               |
| :------------ | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :-------------------------------------- |
| apple         |                                                       ![apple1](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_1.png) ![apple2](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_2.png)                                                       | This turns green when eaten.            |
| poison        |                                                                                                          ![poison](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/poison.png)                                                                                                           | -                                       |
| obstacles     | ![obstacle1](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_1.png) ![obstacle2](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_2.png) ![obstacle3](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_3.png) | -                                       |
| player        |                                                                                                          ![player](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/player.png)                                                                                                           | -                                       |
| Enemy(hunter) |                                                                                                          ![hunter](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/hunter.png)                                                                                                           | -                                       |
| Enemy(ghost)  |                                                                                                           ![ghost](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/ghost.png)                                                                                                            | Enemies that can slip through obstacles |

#### About the state of the game

| State         | To transition to the left state |
| :------------ | :------------------------------ |
| Stage clear   | Eat all apples                  |
| Stage failure | Caught by enemy or eat poison   |
| Game clear    | Clear all stages                |
| Game over     | Stage failure with 0 life.      |

### Player Controls

|    Key    | Action type | Action                                                             |
| :-------: | :---------- | :----------------------------------------------------------------- |
| `h`, `Nh` | `walk`      | move left (If `Nh`, repeat N times)                                |
| `j`, `Nj` | `walk`      | move down (If `Nj`, repeat N times)                                |
| `k`, `Nk` | `walk`      | move up (If `Nk`, repeat N times)                                  |
| `l`, `Nl` | `walk`      | move right (If `Nl`, repeat N times)                               |
| `w`, `Nw` | `walk`      | move forward to next word beginning (If `Nw`, repeat N times)      |
| `e`, `Ne` | `walk`      | move forward to next word ending (If `Ne`, repeat N times)         |
| `b`, `Nb` | `walk`      | move backward to previous word beginning (If `Nb`, repeat N times) |
|    `0`    | `jump`      | move to the beginning of the current line                          |
|    `$`    | `jump`      | move to the end of the current line                                |
|    `^`    | `jump`      | move to the beginning of the first word on the current line        |
|   `gg`    | `jump`      | move to the beginning of the first word on the first line          |
|    `G`    | `jump`      | move to the beginning of the first word on the last line           |
|   `NG`    | `jump`      | move to the beginning of the first word on the nth line            |
|    `q`    | -           | quit the game                                                      |

#### About action type

- `walk`

  - `walk` is the image of moving one square at a time to the destination in an instant. Therefore, hit detection with enemies, obstacles, and apples is applied. Use it when you want to eat apples all at once.

    - e.g. If you type `w`.

      ![walk example](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-w.gif)

- `jump`

  - `jump` is the image of jumping between to the destination and reaching it in an instant. Therefore, hit detection with enemies, obstacles, and apples is not applied. Use it when you want to move to avoid enemies or obstacles.

    - e.g. If you type `$`.

      ![jump example](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-doller.gif)

## For those who want to develop PacVim

### Commands for development

```sh
make help
Usage:
    make <command>

Commands:
    fmt
        go fmt
    lint
        golangci-lint run
    deps
        go mod tidy
    test
        go test
    cover
        create cover.html
    build
        Make a macOS executable binary
    build-win
        Make a Windows executable binary
    clean
        Remove binary files
```

> **Note**
>
> If you have not installed make, refer to the Makefile and execute the command.<br>
>
> - e.g. When building on MacOS<br>
>   - `go build -o bin/mac/pacvim .`

### Execution options

```sh
./pacvim -h
Usage of ./pacvim:
  -level int
    	Level at the start of the game. (1-2) (default 1)
  -life int
    	Remaining lives. (0-5) (default 3)
```

- e.g. If you want to start from level 3 with 5 lives.
  - `go run . -level 3 -life 5`

### How to customize PacVim

#### How to add a stage map

[Reference commit](https://github.com/masahiro-kasatani/pacvim/commit/ab3afdd377e3ac83e0b05b279096f3bcbdd5a26f)

#### How to add enemy types

[Reference commit](https://github.com/masahiro-kasatani/pacvim/commit/6c5f88a32b7ffe73bd640717f0470407578c65d0)

#### How to add the enemy's thought logic

[Reference commit](https://github.com/masahiro-kasatani/pacvim/commit/b0f405ff0be4dc3143579536f89aa30c83c608b6)

## License

MIT

## Author

[Masahiro Kasatani](https://masahiro-kasatani.github.io/portfolio/)
