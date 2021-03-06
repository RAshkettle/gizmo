package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/go-gl/mathgl/mgl32"
)

type ui struct {
	canvas             *pixelgl.Canvas
	miniMapSprite      *pixel.Sprite
	miniMapFrameCanvas *pixelgl.Canvas
	miniMapCanvas      *pixelgl.Canvas
	miniMapScale       float64
	middleTextStr      string
	deathScreenTimer   float64
	uPos               mgl32.Vec2
	lifeCanvas         *pixelgl.Canvas
	playerLife         float64
}

func (u *ui) create() {
	u.canvas = pixelgl.NewCanvas(pixel.R(0, 0, float64(wViewMax), float64(wViewMax)))
	u.canvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})

	u.miniMapCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 100, 100))
	u.miniMapCanvas.SetUniform("uPos", &u.uPos)
	u.miniMapCanvas.SetUniform("uTime", &global.uTime)
	u.miniMapCanvas.SetFragmentShader(fragmentShaderMinimap)

	img, _, _, _ := loadTexture(fmt.Sprintf("%v%v", wAssetObjectsPath, "minimap.png"))
	pic := pixel.PictureDataFromImage(img)
	u.miniMapSprite = pixel.NewSprite(pic, pic.Bounds())
	u.miniMapFrameCanvas = pixelgl.NewCanvas(pic.Bounds())

	u.updatePlayerLife()
}

// Mini map
func (u *ui) updateMiniMap() {
	u.miniMapScale = 0.25 / (float64(global.gWorld.width) / 1024)
	pos := global.gPlayer.getPosition()
	//canvas := pixelgl.NewCanvas(pixel.R(0, 0, 1, 1))
	//canvas.Clear(pixel.RGBA{1.0, 0, 0, 0.5})

	offsetX := float64(global.gWorld.width/2) * u.miniMapScale
	offsetY := float64(global.gWorld.height/2) * u.miniMapScale
	//offset_x2 := offset_x //- float64(global.gWorld.width/2)*u.miniMapScale
	//offset_y2 := offset_y //- float64(global.gWorld.height/2)*u.miniMapScale

	u.uPos = mgl32.Vec2{float32(offsetX / 2), float32(offsetY / 2)}

	bounds := u.miniMapFrameCanvas.Bounds()
	//u.miniMapSprite.Draw(u.miniMapFrameCanvas, pixel.IM.Moved(pixel.V(offset_x2, offset_y2)).ScaledXY(pixel.V(0.8, 0.8), pixel.V(0.8, 0.8)))
	u.miniMapSprite.Draw(u.miniMapFrameCanvas, pixel.IM.Moved(pixel.V(bounds.Max.X/2, bounds.Max.Y/2)).ScaledXY(pixel.V(0.6, 0.6), pixel.V(0.6, 0.6)))

	//global.gWorld.bgSprite.Draw(u.miniMapCanvas, pixel.IM.ScaledXY(pixel.V(u.miniMapScale, u.miniMapScale), pixel.V(u.miniMapScale/2, u.miniMapScale/2)).Moved(pixel.V(u.miniMapScale*pos.X+offset_x-float64(global.gWorld.width/2)*u.miniMapScale, u.miniMapScale*pos.Y+offset_y-float64(global.gWorld.height/2)*u.miniMapScale)))
	bounds = global.gWorld.bgSprite.Frame()
	global.gWorld.bgSprite.Draw(u.miniMapCanvas,
		pixel.IM.ScaledXY(
			pixel.V(u.miniMapScale, u.miniMapScale),
			pixel.V(u.miniMapScale/2, u.miniMapScale/2),
		).Moved(
			pixel.V((bounds.Max.X*u.miniMapScale)/2+pos.X*u.miniMapScale, (bounds.Max.Y*u.miniMapScale)/2+pos.Y*u.miniMapScale)))
	//canvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(u.miniMapScale*pos.X+offset_x-float64(global.gWorld.width/2)*u.miniMapScale, u.miniMapScale*pos.Y+offset_y-float64(global.gWorld.height/2)*u.miniMapScale)))
	//canvas.Draw(u.canvas, pixel.IM.Moved(pixel.V(offset_x, offset_y)))
}

func (u *ui) setMiddleText(text string) {
	//u.middleTextStr = text
	//u.middleText.Clear()
	//u.middleText.WriteString(fmt.Sprintf("%v", text))
}

func (u *ui) updateFPS(fps int) {
	//u.fps.Clear()
	//u.fps.WriteString(fmt.Sprintf("FPS: %v", strconv.Itoa(fps)))
}

func (u *ui) updatePlayerLife() {
	if u.playerLife != global.gPlayer.life {
		u.playerLife = global.gPlayer.life
		u.lifeCanvas = global.gFont.write(fmt.Sprintf("Life: %v", u.playerLife))
	}
}

func (u *ui) draw(dt float64) {

	// Draw death screen
	color := pixel.RGBA{}
	if global.gPlayer.life == 0 {
		u.deathScreenTimer += dt
		u.setMiddleText(wDeathScreenText)

		red := u.deathScreenTimer / 5
		if red > 0.5 {
			red = 0.5
		}
		color = pixel.RGBA{R: red, G: 0, B: 0, A: u.deathScreenTimer / 5}
		u.canvas.Clear(color)
		u.canvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+u.canvas.Bounds().Max.X/2, global.gCamera.pos.Y+u.canvas.Bounds().Max.Y/2)))
	} else {
		u.deathScreenTimer = 0
	}
	u.updatePlayerLife()

	u.miniMapCanvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})
	u.miniMapFrameCanvas.Clear(pixel.RGBA{R: 0, G: 0, B: 0, A: 0})

	//u.updateMiniMap()

	//offsetX := (float64(global.gVariableConfig.WindowWidth) / global.gCamera.zoom) - u.lifeCanvas.Bounds().Max.X
	offsetX := (u.lifeCanvas.Bounds().Max.X / global.gCamera.zoom)
	offsetY := (float64(global.gVariableConfig.WindowHeight) / global.gCamera.zoom) - u.lifeCanvas.Bounds().Max.Y
	lifeScale := float64(global.gVariableConfig.WindowWidth) / float64(global.gVariableConfig.WindowHeight) / 5
	u.lifeCanvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, lifeScale).Moved(pixel.V(global.gCamera.pos.X+offsetX, global.gCamera.pos.Y+offsetY)))

	//bounds := u.miniMapFrameCanvas.Bounds()
	//u.miniMapFrameCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+bounds.Max.X/2, global.gCamera.pos.Y+bounds.Max.Y/2)))
	//miniMapScale := float64(global.gVariableConfig.WindowWidth) / float64(global.gVariableConfig.WindowHeight) / 5
	//u.miniMapFrameCanvas.Draw(global.gWin, pixel.IM.Scaled(pixel.ZV, miniMapScale).Moved(pixel.V(global.gCamera.pos.X+offsetX, global.gCamera.pos.Y+offsetY)))
	//bounds = u.miniMapCanvas.Bounds()
	//u.miniMapCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X, global.gCamera.pos.Y)))
	//u.miniMapCanvas.Draw(global.gWin, pixel.IM.Moved(pixel.V(global.gCamera.pos.X+wViewMax/2.0, global.gCamera.pos.Y+wViewMax/2.0)))

}
