<p align="center">
  <b>PacVim is a Vim learning game ʕ◔ϖ◔ʔ</b>
</p>

![pacvim](https://github.com/masahiro-kasatani/pacvim/blob/readme-images/files/readme.png?raw=true)

<p align="right">
The Go gopher was designed by Renée French.
</p>

---

## Building and running
### How to build(Mac OS)
```
# mac
$ go build -o bin/mac/pacvim ./main.go
# win
$ GOOS=windows GOARCH=amd64 go build -o bin/win/pacvim.exe ./main.go
```

### How to run
Double click to run binary
* ./bin/win/pacvim.exe
* ./bin/mac/pacvim

or execute 'go run'
```
$ go run main.go
$ go run main.go -initiaLevel 5
```

## How to play
![objects](https://user-images.githubusercontent.com/61332083/103471133-790e3a80-4dbf-11eb-96fd-a6525766b5f5.png)
### Objects
| char | what it does |
|:-:|:-:|
| o  | food  |
| X  | poison  |
| G  | ghost  |

### Game status
| state | conditions |
|:-:|:-:|
| Game Clear  | Clear all 10 stages  |
| Game Over  | Fail a total of 4 times  |
| Stage Win  | Eat all food  |
| Stage Lose  | Eat poison<br>Caught by a ghost  |

### Operation
| key | what it does |
|:-:|:-:|
| q  | quit the game  |
| h  | move left  |
| j  | move down  |
| k  | move up  |
| l  | move right  |
| w  | move forward to next word beginning  |
| e  | move forward to next word ending  |
| b  | move backward to next word beginning  |
| $  | move to the end of the line  |
| 0  | move to the beginning of the line  |
| gg/1G  | move to the beginning of the first line  |
| NG  | move to the beginning of the line given by N  |
| G  | move to the beginning of the last line  |
| ^  | move to the first word at the current line  |

#### "w,e,b" ： Move at once.
![w](https://user-images.githubusercontent.com/61332083/103471176-3c8f0e80-4dc0-11eb-97b4-b20b905b55c4.gif)

#### "$,0,gg,NG,G,^" ： Warp
![dollar](https://user-images.githubusercontent.com/61332083/103471203-ba531a00-4dc0-11eb-9533-2d6be6e8b962.gif)

### How to play with your own map
Replace map file (./files/stage/)

## License
MIT

## Author
Masahiro Kasatani
