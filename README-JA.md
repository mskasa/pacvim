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

- [PacVim の起動方法](#PacVim-の起動方法)
- [PacVim のルール](#PacVim-のルール)
  - [ゲーム画面](#ゲーム画面)
  - [オブジェクトについて](#オブジェクトについて)
  - [ゲームの状態について](#ゲームの状態について)
- [プレイヤーの操作方法](#操作方法)
  - [動作種別について](#動作種別について)
- [ライセンス](#ライセンス)
- [著者](#著者)

<!-- /TOC -->

## PacVim の起動方法

以下のバイナリファイルをダブルクリックして起動してください。

- Windows: [./bin/win/pacvim.exe](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/win)
- Mac: [./bin/mac/pacvim](https://github.com/masahiro-kasatani/pacvim/tree/master/bin/mac)

## PacVim のルール

PacVim はパックマンのルールを踏襲しています。

### ゲーム画面

![ゲーム画面](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/screen.png)

### オブジェクトについて

| オブジェクト名 |                                                                                                                                                         表示                                                                                                                                                         | 補足説明                     |
| :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :--------------------------- |
| りんご         |                                               ![りんご（未）](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_1.png) ![りんご（済）](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/apple_2.png)                                                | 食べると緑色になります       |
| 毒             |                                                                                                           ![毒](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/poison.png)                                                                                                           | -                            |
| 障害物         | ![障害物１](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_1.png) ![障害物２](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_2.png) ![障害物３](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/wall_3.png) | -                            |
| プレイヤー     |                                                                                                       ![プレイヤー](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/player.png)                                                                                                       | -                            |
| 敵（ハンター） |                                                                                                        ![ハンター](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/hunter.png)                                                                                                        | -                            |
| 敵（ゴースト） |                                                                                                        ![ゴースト](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/ghost.png)                                                                                                         | 障害物をすり抜けられる敵です |

### ゲームの状態について

| 状態           | 遷移条件                            |
| :------------- | :---------------------------------- |
| ステージクリア | すべてのりんごを食べる              |
| ステージ失敗   | 敵に捕まる or 毒を食べる            |
| ゲームクリア   | すべてのステージをクリアする        |
| ゲームオーバー | ライフが 0 の状態でステージ失敗する |

## プレイヤーの操作方法

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

### 動作種別について

- `walk`

  - `walk` は目的地まで、 1 マスずつ一瞬で移動するイメージです。そのため、敵や障害物、りんごとの当たり判定が適用されます。一気にりんごを食べたいときに使いましょう。

    - 例： `w` を入力した場合

      ![walkの例](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-w.gif)

- `jump`

  - `jump` は目的地まで、間を飛び越えて一瞬で到達するイメージです。そのため、敵や障害物、りんごとの当たり判定が適用されません。敵や障害物を避けて移動したいときに使いましょう。

    - 例： `$` を入力した場合

      ![jumpの例](https://raw.githubusercontent.com/masahiro-kasatani/pacvim/readme-images/files/readme-doller.gif)

## ライセンス

MIT

## 著者

[笠谷昌弘](https://masahiro-kasatani.github.io/portfolio/)
