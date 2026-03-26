package renderer

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed sprites/sprites.png
var spritesData []byte

// スプライトシート上の各スプライトの切り出し矩形。
// キャラクター（Gopher/Hunter/Ghost）は 32×32、オブジェクトは 16×16。
var (
	// Gopher（プレイヤー）フレーム 0〜3
	rectGopher0 = image.Rect(0, 0, 32, 32)
	rectGopher1 = image.Rect(32, 0, 64, 32)
	rectGopher2 = image.Rect(64, 0, 96, 32)
	rectGopher3 = image.Rect(96, 0, 128, 32)

	// Hunter（ハンター）フレーム 0〜3
	rectHunter0 = image.Rect(0, 32, 32, 64)
	rectHunter1 = image.Rect(32, 32, 64, 64)
	rectHunter2 = image.Rect(64, 32, 96, 64)
	rectHunter3 = image.Rect(96, 32, 128, 64)

	// Ghost（ゴースト）フレーム 0〜3
	rectGhost0 = image.Rect(0, 64, 32, 96)
	rectGhost1 = image.Rect(32, 64, 64, 96)
	rectGhost2 = image.Rect(64, 64, 96, 96)
	rectGhost3 = image.Rect(96, 64, 128, 96)

	// オブジェクト（16×16）
	rectApple      = image.Rect(0, 96, 16, 112)
	rectAppleEaten = image.Rect(16, 96, 32, 112)
	rectPoison     = image.Rect(32, 96, 48, 112)
	rectWall       = image.Rect(48, 96, 64, 112)
)

// spriteSheet はスプライトシート全体と各スプライトの切り出し済み画像を保持する。
type spriteSheet struct {
	img          *ebiten.Image
	gopherFrames [4]*ebiten.Image
	hunterFrames [4]*ebiten.Image
	ghostFrames  [4]*ebiten.Image
	apple        *ebiten.Image
	appleEaten   *ebiten.Image
	poison       *ebiten.Image
	wall         *ebiten.Image
}

// newSpriteSheet はスプライトシートを埋め込みデータからロードして返す。
func newSpriteSheet() (*spriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(spritesData))
	if err != nil {
		return nil, err
	}
	eimg := ebiten.NewImageFromImage(img)

	ss := &spriteSheet{img: eimg}
	ss.gopherFrames = [4]*ebiten.Image{
		ss.subImage(rectGopher0),
		ss.subImage(rectGopher1),
		ss.subImage(rectGopher2),
		ss.subImage(rectGopher3),
	}
	ss.hunterFrames = [4]*ebiten.Image{
		ss.subImage(rectHunter0),
		ss.subImage(rectHunter1),
		ss.subImage(rectHunter2),
		ss.subImage(rectHunter3),
	}
	ss.ghostFrames = [4]*ebiten.Image{
		ss.subImage(rectGhost0),
		ss.subImage(rectGhost1),
		ss.subImage(rectGhost2),
		ss.subImage(rectGhost3),
	}
	ss.apple = ss.subImage(rectApple)
	ss.appleEaten = ss.subImage(rectAppleEaten)
	ss.poison = ss.subImage(rectPoison)
	ss.wall = ss.subImage(rectWall)

	return ss, nil
}

// subImage は指定矩形のスプライトを切り出して返す。
func (ss *spriteSheet) subImage(r image.Rectangle) *ebiten.Image {
	return ss.img.SubImage(r).(*ebiten.Image)
}
