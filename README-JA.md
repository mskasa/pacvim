# PacVim ~Vim 学習ゲーム~ ʕ◔ϖ◔ʔ

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

TODO ゲーム画面の画像貼り、各オブジェクトの説明を入れる

| 状態           | 遷移条件                            |
| :------------- | :---------------------------------- |
| ステージクリア | すべてのりんごを食べる              |
| ステージ失敗   | 敵に捕まる or 毒を食べる            |
| ゲームクリア   | すべてのステージをクリアする        |
| ゲームオーバー | ライフが 0 の状態でステージ失敗する |

## プレイヤーの操作方法

| キー  | 動作                                   | 動作種別 |
| :---: | :------------------------------------- | :------- |
|   h   | 左へ 1 マス移動する                    | walk     |
|   j   | 下へ 1 マス移動する                    | walk     |
|   k   | 上へ 1 マス移動する                    | walk     |
|   l   | 右へ 1 マス移動する                    | walk     |
|   w   | 次の単語の先頭に移動する               | walk     |
|   e   | 次の単語の末尾に移動する               | walk     |
|   b   | 前の単語の先頭に移動する               | walk     |
|   0   | 現在の行の先頭に移動する               | jump     |
|   $   | 現在の行の末尾に移動する               | jump     |
|   ^   | 現在の行の最初の単語の先頭に移動する   | jump     |
| gg/1G | 最初の行の最初の単語の先頭に移動する   | jump     |
|   G   | 最後の行の最初の単語の先頭に移動する   | jump     |
|  NG   | N 行目の行の最初の単語の先頭に移動する | jump     |
|   q   | ゲームをやめる                         | -        |

### 動作種別について

- walk

- jump

## ライセンス

MIT

## 著者

[笠谷昌弘](https://masahiro-kasatani.github.io/portfolio/)
