package directus

import (
	"encoding/json"

	"github.com/thisisdevelopment/go-dockly/xerrors/iferr"
)

type AuthResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

type DirectusSchema struct {
	Data interface{} `json:"data"`
}

func (d *DirectusSchema) Convert() string {
	data, err := json.Marshal(d.Data)
	iferr.Panic(err)
	return string(data)
}
