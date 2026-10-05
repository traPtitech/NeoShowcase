package web

import (
	"net/http"
)

func NewH2CClient() *http.Client {
	protocols := new(http.Protocols)
	protocols.SetHTTP2(true)
	protocols.SetUnencryptedHTTP2(true)
	// https://connect.build/docs/go/deployment
	return &http.Client{
		Transport: &http.Transport{
			Protocols: protocols,
		},
	}
}
