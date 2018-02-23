package termmaker

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"image"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"image/png"

	_ "image/jpeg"

	"github.com/eternnoir/tgsticmaker/maker"
	"github.com/fatih/color"
	"github.com/nfnt/resize"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var (
	ModeFullAuto = "Full Auto mode"
	ModeManual   = "Manual mode"
)

var EmojiList []string

func init() {
	EmojiList = make([]string, 0)
	emoji := [][]int{
		{128513, 128591},
	}
	for _, value := range emoji {
		for x := value[0]; x < value[1]; x++ {
			str := html.UnescapeString("&#" + strconv.Itoa(x) + ";")
			EmojiList = append(EmojiList, str)
		}
	}

}

type Maker struct {
	Maker *maker.StickerMaker
}

func (maker *Maker) Start() {
	mode := ""
	prompt := &survey.Select{
		Message: "What you want to do?:",
		Options: []string{ModeFullAuto, ModeManual},
	}
	survey.AskOne(prompt, &mode, nil)
	switch mode {
	case ModeFullAuto:
		maker.FullAuto()
	case ModeManual:
	default:
		panic("No mode")
	}
}

func (m *Maker) FullAuto() {
	stickerName := Ask("Sticker set name")
	userId := Ask("User Id")
	title := Ask("Sticker Title")
	stickerPngFolder := Ask("Image folder path")
	color.Green("Sticker set name: %s\nUserId: %s\nImageFolder: %s\n", stickerName, userId, stickerPngFolder)
	files, err := ioutil.ReadDir(stickerPngFolder)
	if err != nil {
		color.Red(err.Error())
		panic(err)
	}

	stickerList := make([]maker.Sticker, len(files))
	for i, f := range files {
		s, err := CreateStickerByFile(stickerPngFolder, f, true)
		if err != nil {
			panic(err)
		}
		stickerList[i] = s
	}
	userId64, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		panic(err)
	}
	req := &maker.CreateStickerSetReq{
		UserId:   userId64,
		Name:     stickerName,
		Title:    title,
		Stickers: stickerList,
	}
	color.Green("Start to create sticker set.")
	_, err = m.Maker.CreateSitckerSet(req, true)
	if err != nil {
		panic(err)
	}
}

func CreateStickerByFile(path string, file os.FileInfo, randomEmoji bool) (maker.Sticker, error) {
	fmt.Println(file.Name())
	fImg, err := os.Open(filepath.Join(path, file.Name()))
	defer fImg.Close()
	if err != nil {
		return maker.Sticker{}, err
	}
	img, format, err := image.Decode(fImg)
	if err != nil {
		return maker.Sticker{}, err
	}
	color.Cyan("Load %s format %s\n", file.Name(), format)
	m := resize.Resize(256, 256, img, resize.NearestNeighbor)
	m = resize.Resize(512, 512, m, resize.NearestNeighbor)
	buf := new(bytes.Buffer)
	err = (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(buf, m)
	if err != nil {
		return maker.Sticker{}, err
	}
	fmt.Println(buf.Len())
	ret := maker.Sticker{
		Img:   buf.Bytes(),
		Emoji: getRandomEmoji(),
	}
	color.Cyan("Convert %s %s success.", file.Name(), ret.Emoji)
	return ret, nil
}

func Ask(msg string) string {
	resp := ""
	prompt := &survey.Input{
		Message: msg,
	}
	survey.AskOne(prompt, &resp, MustHaveValue)
	return resp
}

func YorN(msg string) bool {
	ret := false
	prompt := &survey.Confirm{
		Message: msg,
	}
	survey.AskOne(prompt, &ret, nil)
	return ret
}

func GetStickerSetName() string {
	name := ""
	prompt := &survey.Input{
		Message: "Sticker name",
	}
	survey.AskOne(prompt, &name, MustHaveValue)
	return name
}

func MustHaveValue(val interface{}) error {
	// since we are validating an Input, the assertion will always succeed
	if str, ok := val.(string); !ok || len(str) < 1 {
		return errors.New("This response cannot be empty.")
	}
	return nil
}

func getRandomEmoji() string {
	return EmojiList[rand.Intn(len(EmojiList))]
}

func RandomEmoji() string {
	return getRandomEmoji()
}
