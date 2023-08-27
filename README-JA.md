# PacVim \~Vim 学習ゲーム\~ ʕ◔ϖ◔ʔ

![pacvim](https://github.com/masahiro-kasatani/pacvim/blob/readme-images/files/readme.png?raw=true)

<p align="right">
The Go gopher was designed by <a href="https://go.dev/blog/gopher" target="_blank">Renée French</a>.
</p>

[![test](https://github.com/masahiro-kasatani/pacvim/actions/workflows/test.yaml/badge.svg)](https://github.com/masahiro-kasatani/pacvim/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/masahiro-kasatani/pacvim)](https://goreportcard.com/report/github.com/masahiro-kasatani/pacvim)
[![codecov](https://codecov.io/gh/masahiro-kasatani/pacvim/branch/master/graph/badge.svg?token=KZ2LVX4GCT)](https://codecov.io/gh/masahiro-kasatani/pacvim)

| [English](https://github.com/masahiro-kasatani/pacvim/blob/master/README.md) | 日本語 |

<!-- TOC -->

## 目次

- [PacVim を遊びたい方へ](#pacvim-を遊びたい方へ)
  - [PacVim の起動方法](#pacvim-の起動方法)
  - [PacVim のルール](#pacvim-のルール)
    - [ゲーム画面](#ゲーム画面)
    - [オブジェクトについて](#オブジェクトについて)
    - [ゲームの状態について](#ゲームの状態について)
  - [プレイヤーの操作方法](#操作方法)
    - [動作種別について](#動作種別について)
- [PacVim を開発したい方へ](#pacvim-を開発したい方へ)
  - [開発用コマンド](#開発用コマンド)
  - [実行用オプション](#実行用オプション)
  - [PacVim の改良方法](#pacvim-の改良方法)
    - [ステージマップの追加方法](#ステージマップの追加方法)
    - [敵の種類の追加方法](#敵の種類の追加方法)
    - [敵の思考ロジックの追加方法](#敵の思考ロジックの追加方法)
- [ライセンス](#ライセンス)
- [著者](#著者)

<!-- /TOC -->

## PacVim を遊びたい方へ

### PacVim の起動方法

以下のバイナリファイルをダブルクリックして起動してください。

- Windows: [./bin/win/pacvim.exe](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/win)
- Mac: [./bin/mac/pacvim](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/mac)

### PacVim のルール

PacVim はパックマンのルールを踏襲しています。

#### ゲーム画面

![ゲーム画面](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/screen.png)

#### オブジェクトについて

| オブジェクト名 |                                                                                                                                                         表示                                                                                                                                                         | 補足説明                     |
| :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------- |
| りんご         |                                               ![りんご（未）](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_1.png) ![りんご（済）](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_2.png)                                                | 食べると緑色になります       |
| 毒             |                                                                                                           ![毒](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/poison.png)                                                                                                           | -                            |
| 障害物         | ![障害物１](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_1.png) ![障害物２](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_2.png) ![障害物３](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_3.png) | -                            |
| プレイヤー     |                                                                                                       ![プレイヤー](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/player.png)                                                                                                       | -                            |
| 敵（ハンター） |                                                                                                        ![ハンター](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/hunter.png)                                                                                                        | -                            |
| 敵（ゴースト） |                                                                                                        ![ゴースト](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/ghost.png)                                                                                                         | 障害物をすり抜けられる敵です |

#### ゲームの状態について

| 状態           | 遷移条件                            |
| :------------- | :---------------------------------- |
| ステージクリア | すべてのりんごを食べる              |
| ステージ失敗   | 敵に捕まる or 毒を食べる            |
| ゲームクリア   | すべてのステージをクリアする        |
| ゲームオーバー | ライフが 0 の状態でステージ失敗する |

### プレイヤーの操作方法

| キー | 動作種別 | 動作                                   |
| :--: | :------- | :------------------------------------- |
| `h`  | `walk`   | 左へ 1 マス移動する                    |
| `j`  | `walk`   | 下へ 1 マス移動する                    |
| `k`  | `walk`   | 上へ 1 マス移動する                    |
| `l`  | `walk`   | 右へ 1 マス移動する                    |
| `w`  | `walk`   | 次の単語の先頭に移動する               |
| `e`  | `walk`   | 次の単語の末尾に移動する               |
| `b`  | `walk`   | 前の単語の先頭に移動する               |
| `0`  | `jump`   | 現在の行の先頭に移動する               |
| `$`  | `jump`   | 現在の行の末尾に移動する               |
| `^`  | `jump`   | 現在の行の最初の単語の先頭に移動する   |
| `gg` | `jump`   | 最初の行の最初の単語の先頭に移動する   |
| `G`  | `jump`   | 最後の行の最初の単語の先頭に移動する   |
| `NG` | `jump`   | N 行目の行の最初の単語の先頭に移動する |
| `q`  | -        | ゲームをやめる                         |

#### 動作種別について

- `walk`

  - `walk` は目的地まで、 1 マスずつ一瞬で移動するイメージです。そのため、敵や障害物、りんごとの当たり判定が適用されます。一気にりんごを食べたいときに使いましょう。

    - 例： `w` を入力した場合

      ![walkの例](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-w.gif)

- `jump`

  - `jump` は目的地まで、間を飛び越えて一瞬で到達するイメージです。そのため、敵や障害物、りんごとの当たり判定が適用されません。敵や障害物を避けて移動したいときに使いましょう。

    - 例： `$` を入力した場合

      ![jumpの例](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-doller.gif)

## PacVim を開発したい方へ

### 開発用コマンド

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
> make をインストールしていない場合、Makefile を参考にコマンドを実行してください。<br>
>
> - 例：MacOS でビルドをする場合<br>
>   - `go build -o bin/mac/pacvim .`

### 実行用オプション

```sh
./pacvim -h
Usage of ./pacvim:
  -level int
    	Level at the start of the game. (1-2) (default 1)
  -life int
    	Remaining lives. (0-5) (default 3)
```

- 例：残機 5 でレベル 3 からスタートしたい場合
  - `go run . -level 3 -life 5`

### PacVim の改良方法

#### ステージマップの追加方法

[参考コミット](https://github.com/masahiro-kasatani/pacvim/commit/ab3afdd377e3ac83e0b05b279096f3bcbdd5a26f)

#### 敵の種類の追加方法

[参考コミット](https://github.com/masahiro-kasatani/pacvim/commit/6c5f88a32b7ffe73bd640717f0470407578c65d0)

#### 敵の思考ロジックの追加方法

[参考コミット](https://github.com/masahiro-kasatani/pacvim/commit/b0f405ff0be4dc3143579536f89aa30c83c608b6)

## ライセンス

MIT

## 著者

[笠谷昌弘](https://masahiro-kasatani.github.io/portfolio/)
