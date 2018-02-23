package server

import (
	"bytes"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/eternnoir/tgsticmaker/maker"
	"github.com/eternnoir/tgsticmaker/termmaker"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	Maker *maker.StickerMaker
}

func (ser *Server) Start(addr string) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Any("/**", func(c echo.Context) error {
		return c.File("StickerMaker/dist/index.html")
	}, middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "StickerMaker/dist",
	}))
	a := e.Group("/api")
	a.POST("/stickerset", ser.CreateSitckerSet)
	a.GET("/stickerset/:name", ser.CheckStickerExist)
	a.POST("/sticker", ser.AppendStickerToSet)

	return e.Start(addr)
}

func (ser *Server) CheckStickerExist(c echo.Context) error {
	stickerName := c.Param("name")
	if stickerName == "" {
		return c.JSON(http.StatusUnprocessableEntity, BasicResult{Error: "Sticker Name is empty"})
	}
	result, err := ser.Maker.CheckStickerExist(stickerName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BasicResult{Error: err})
	}
	return c.JSON(http.StatusOK, BasicResult{Result: result})
}

func (ser *Server) CreateSitckerSet(c echo.Context) error {
	userId := c.FormValue("userId")
	stickerName := c.FormValue("stickerName")
	stickerTitle := c.FormValue("stickerTitle")
	userId64, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, BasicResult{Error: "User ID must be int"})
	}
	log.Printf("Receive create sticker req. %s %s %s", userId, stickerName, stickerTitle)
	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	stickerList := make([]maker.Sticker, len(files))
	for i, file := range files {
		pngByte, err := func(file *multipart.FileHeader) ([]byte, error) {
			// Source
			src, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer src.Close()
			log.Println(file.Filename)
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, src); err != nil {
				return nil, c.JSON(http.StatusInternalServerError, err)
			}

			img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
			if err != nil {
				return nil, c.JSON(http.StatusInternalServerError, err)
			}
			log.Printf("Start to resize %s\n", file.Filename)
			pngByte, err := maker.ResizeImage(img)
			if err != nil {
				return nil, c.JSON(http.StatusInternalServerError, err)
			}
			return pngByte, nil
		}(file)
		if err != nil {
			return err
		}
		stickerList[i] = maker.Sticker{
			Img:   pngByte,
			Emoji: termmaker.RandomEmoji(),
		}
	}
	log.Println("All resize process done")
	req := &maker.CreateStickerSetReq{
		UserId:   userId64,
		Name:     stickerName,
		Title:    stickerTitle,
		Stickers: stickerList,
	}

	_, err = ser.Maker.InitStickerSet(req)
	go ser.Maker.AddStickerToSetFormOffset(req, 1)
	return c.JSON(http.StatusOK, BasicResult{
		Result: map[string]interface{}{
			"stickerName": ser.Maker.StickerFullName(stickerName),
		},
		Error: err,
	})
}

func (ser *Server) AppendStickerToSet(c echo.Context) error {
	return nil
}

type BasicResult struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error,omitempty"`
}
