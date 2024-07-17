package gui

import (
	"access-mask/object"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.design/x/clipboard"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"log"
	"math"
	"strconv"
)

const (
	cellHorizontalMargin = 30
	cellVerticalMargin   = 20

	descriptionCellVerticalSize = 250
	valueCellVerticalSize       = 20

	cellBorderWidth = 2

	descriptionShiftHorizontal = 14
	descriptionShiftVertical   = 8

	titleVerticalMargin = 10

	lineSpacing = 15
	valueScale  = 3
	nameScale   = 3
	cellCount   = 32
)

type gui struct {
	font text.Face

	object object.Type

	description [cellCount]image.Rectangle
	cell        [cellCount]image.Rectangle

	value uint32
}

func New() ebiten.Game {
	instance := &gui{
		font:   text.NewGoXFace(basicfont.Face7x13),
		object: object.File,
		value:  uint32(0xFFFFFFFF),
	}
	_ = instance.importValue(clipboard.Read(clipboard.FmtText))
	return instance
}

func (gui *gui) Update() error {
	width, _ := ebiten.WindowSize()
	cursorX, cursorY := ebiten.CursorPosition()
	availableWidth := width - cellHorizontalMargin*2
	cellWidth := availableWidth / cellCount
	cellStartX := width/2 - availableWidth/2
	for currentCell := 0; currentCell < cellCount; currentCell++ {
		gui.description[currentCell] = image.Rect(
			cellStartX+currentCell*cellWidth,
			cellVerticalMargin,
			cellStartX+(currentCell+1)*cellWidth,
			cellVerticalMargin+descriptionCellVerticalSize,
		)
		gui.cell[currentCell] = image.Rect(
			gui.description[currentCell].Min.X,
			gui.description[currentCell].Max.Y,
			gui.description[currentCell].Max.X,
			gui.description[currentCell].Max.Y+valueCellVerticalSize,
		)

		if cursorX >= gui.cell[currentCell].Min.X &&
			cursorY >= gui.cell[currentCell].Min.Y &&
			cursorX <= gui.cell[currentCell].Max.X &&
			cursorY <= gui.cell[currentCell].Max.Y {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				gui.value ^= uint32(1 << (cellCount - currentCell - 1))
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		gui.object = (gui.object + object.Types - 1) % object.Types
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		gui.object = (gui.object + 1) % object.Types
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			clipboard.Write(clipboard.FmtText, []byte(fmt.Sprintf("0x%x", gui.value)))
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			if err := gui.importValue(clipboard.Read(clipboard.FmtText)); err != nil {
				log.Println("could not import value from clipboard:", err)
			}
		}
	}

	if !ebiten.IsKeyPressed(ebiten.KeyAlt) {
		gui.value = gui.value & object.Mask[gui.object]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) {
		gui.value = 0
	}

	return nil
}

func (gui *gui) Draw(screen *ebiten.Image) {
	width, height := ebiten.WindowSize()

	screen.Clear()

	options := text.DrawOptions{}

	for currentCellIndex := 0; currentCellIndex < cellCount; currentCellIndex++ {
		cellList := []*image.Rectangle{&gui.description[currentCellIndex], &gui.cell[currentCellIndex]}
		for isValue, currentCell := range cellList {
			vector.StrokeRect(screen,
				float32(currentCell.Min.X), float32(currentCell.Min.Y),
				float32(currentCell.Size().X), float32(currentCell.Size().Y),
				cellBorderWidth,
				colornames.White,
				false)

			if isValue == 0 {
				options.GeoM.Reset()
				options.GeoM.Rotate(-90 * math.Pi / 180)
				options.GeoM.Scale(2, 2)
				options.GeoM.Translate(
					float64(currentCell.Min.X+currentCell.Size().X/2-descriptionShiftHorizontal),
					float64(currentCell.Max.Y-descriptionShiftVertical),
				)
				text.Draw(screen,
					object.Description[gui.object][currentCellIndex],
					gui.font,
					&options,
				)
			} else {
				if gui.value&uint32(1<<(cellCount-currentCellIndex-1)) != 0 {
					vector.DrawFilledRect(screen,
						float32(currentCell.Min.X)+cellBorderWidth, float32(currentCell.Min.Y)+cellBorderWidth,
						float32(currentCell.Size().X)-cellBorderWidth*2, float32(currentCell.Size().Y)-cellBorderWidth*2,
						colornames.White,
						false)
				}
			}
		}
	}

	value := fmt.Sprintf("0x%x\n\nCtrl+C to copy\nCtrl+V to paste", gui.value)
	options.GeoM.Reset()
	options.PrimaryAlign = text.AlignCenter
	options.SecondaryAlign = text.AlignStart
	options.LineSpacing = lineSpacing
	options.GeoM.Scale(valueScale, valueScale)
	options.GeoM.Translate(float64(width/2), float64(height/2))
	text.Draw(screen, value, gui.font, &options)

	title := fmt.Sprintf("<- %s ->", object.Name[gui.object])
	options.GeoM.Reset()
	options.PrimaryAlign = text.AlignCenter
	options.GeoM.Scale(nameScale, nameScale)
	_, textHeight := text.Measure(title, gui.font, 0)
	options.GeoM.Translate(float64(width/2), float64(height-int(textHeight)*nameScale-titleVerticalMargin))
	text.Draw(screen, title, gui.font, &options)
}

func (gui *gui) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (gui *gui) importValue(rawValue []byte) error {
	if value, err := strconv.ParseUint(string(rawValue), 0, 32); err == nil {
		gui.value = uint32(value)
		return nil
	} else {
		return err
	}
}
