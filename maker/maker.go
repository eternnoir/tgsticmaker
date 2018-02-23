package maker

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"

	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nfnt/resize"
)

var PNGSIZELIMIT = 510 * 1000

type StickerMaker struct {
	api     *tgbotapi.BotAPI
	botName string
}

func New(api *tgbotapi.BotAPI, botname string) *StickerMaker {
	return &StickerMaker{
		api:     api,
		botName: botname,
	}
}

func (s *StickerMaker) CheckStickerExist(name string) (bool, error) {
	sn := formatStickerName(name, s.botName)
	resp, err := s.api.GetStickerSet(sn)
	if err != nil {
		return false, err
	}
	if resp.Name == "" {
		return false, nil
	}
	return true, nil
}

func (s *StickerMaker) CreateSitckerSet(req *CreateStickerSetReq, skipCreateFail bool) (bool, error) {
	if len(req.Stickers) < 1 {
		return false, errors.New("Stickers is empty")
	}
	initResp, err := s.InitStickerSet(req)
	if err != nil {
		log.Println(err.Error())
		if !skipCreateFail {
			return false, err
		}
	}
	if !initResp && !skipCreateFail {
		return false, errors.New("Cannot init sticker set")
	}
	err = s.AddStickerToSetFormOffset(req, 1)
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *StickerMaker) InitStickerSet(req *CreateStickerSetReq) (bool, error) {
	sn := formatStickerName(req.Name, s.botName)
	firstSticker := req.Stickers[0]
	createNewConfig := tgbotapi.CreateNewStickerSetConfig{
		BaseStickerConfig: tgbotapi.BaseStickerConfig{
			UserId:     req.UserId,
			Name:       sn,
			PngSticker: firstSticker.Img,
			Emojis:     firstSticker.Emoji,
		},
		Title: req.Title,
	}
	return s.api.CreateNewStickerSet(createNewConfig)
}

func (s *StickerMaker) AddStickerToSetFormOffset(req *CreateStickerSetReq, offset int64) error {
	sn := formatStickerName(req.Name, s.botName)
	stickers := req.Stickers[offset:]
	for _, st := range stickers {
		fmt.Println("Add sticker to set.")
		cfg := tgbotapi.AddStickerToSetConfig{
			BaseStickerConfig: tgbotapi.BaseStickerConfig{
				UserId:     req.UserId,
				Name:       sn,
				PngSticker: st.Img,
				Emojis:     st.Emoji,
			},
		}
		_, err := s.api.AddStickerToSet(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StickerMaker) StickerFullName(name string) string {
	return formatStickerName(name, s.botName)
}

func formatStickerName(name, botname string) string {
	return fmt.Sprintf("%s_by_%s", name, botname)
}

func ResizeImage(img image.Image) ([]byte, error) {
	size := uint(512)
	offset := uint(10)
	m := img
	var ret []byte
	for {
		m = resize.Resize(size, size, img, resize.NearestNeighbor)
		m = resize.Resize(512, 512, m, resize.NearestNeighbor)
		buf := new(bytes.Buffer)
		err := (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(buf, m)
		if err != nil {
			return nil, err
		}
		if len(buf.Bytes()) > PNGSIZELIMIT {
			size = size - offset
		} else {
			ret = buf.Bytes()
			break
		}

		if size < 50 {
			return nil, errors.New("image too large")
		}
	}
	return ret, nil
}
