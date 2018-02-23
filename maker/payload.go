package maker

type CreateStickerSetReq struct {
	UserId   int64     `json:"ser_id"`
	Name     string    `json:"name"`
	Title    string    `json:"title"`
	Stickers []Sticker `json:"stickers"`
}

type Sticker struct {
	Img   []byte `json:"img"`
	Emoji string `json:"emoji"`
}
