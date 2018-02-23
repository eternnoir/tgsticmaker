package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/eternnoir/tgsticmaker/maker"
	"github.com/eternnoir/tgsticmaker/server"
	"github.com/eternnoir/tgsticmaker/termmaker"
	"github.com/fatih/color"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	cli "gopkg.in/urfave/cli.v1"
)

var (
	aDebug   = false
	aPORT    = ""
	aTgToken = ""
	aBotName = ""
)

var flags = []cli.Flag{
	cli.BoolFlag{
		Name:        "debug, d",
		Usage:       "Debug mode.",
		EnvVar:      "DEBUG",
		Destination: &aDebug,
	},
	cli.StringFlag{
		Name:        "port, p",
		Usage:       "Server listen port. e.g. 8458",
		Value:       "8458",
		EnvVar:      "PORT",
		Destination: &aPORT,
	},
	cli.StringFlag{
		Name:        "token, t",
		Usage:       "Telegram bot token",
		EnvVar:      "TOKEN",
		Destination: &aTgToken,
	},
	cli.StringFlag{
		Name:        "botname, b",
		Usage:       "Telegram bot name",
		EnvVar:      "BOT",
		Destination: &aBotName,
	},
}

func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "TgSticmaker"
	app.Usage = "Telegram sticker maker"

	app.Flags = flags
	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Start Customer Web Backend.",
			Action: start,
			Flags:  app.Flags,
		},
		{
			Name:   "make",
			Usage:  "Make sticker set",
			Action: make,
			Flags:  app.Flags,
		},
	}

	app.Action = start
	return app
}

func start(c *cli.Context) error {
	bot, err := tgbotapi.NewBotAPI(aTgToken)
	if err != nil {
		log.Panic(err)
	}
	stickerMaker := maker.New(bot, aBotName)
	ser := &server.Server{
		Maker: stickerMaker,
	}
	return ser.Start("0.0.0.0:" + aPORT)
}

func make(c *cli.Context) error {
	getToken()
	getBotName()
	bot, err := tgbotapi.NewBotAPI(aTgToken)
	if err != nil {
		log.Panic(err)
	}
	stickerMaker := maker.New(bot, aBotName)
	termmak := &termmaker.Maker{
		Maker: stickerMaker,
	}
	termmak.Start()
	return nil
}

func main() {
	app := InitApp()
	if err := app.Run(os.Args); err != nil {
		log.Panic(err)
	}
}

func getToken() string {
	if aTgToken == "" {
		color.Green("Please input tg token")
		aTgToken = readLine()
	}
	color.Green("TG Token is: %s", aTgToken)
	return aTgToken
}

func getBotName() string {
	if aBotName == "" {
		color.Green("Please input tg bot name")
		aBotName = readLine()
	}
	color.Green("TG Bot name is: %s", aBotName)
	return aBotName
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}
