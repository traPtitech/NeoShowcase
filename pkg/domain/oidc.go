package domain

import "net/http"

func NewLoginHandler() http.Handler {
	// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのログイン用ハンドラを返す関数
	return http.NotFoundHandler()
}

func NewCallbackHandler() http.Handler {
	// TODO: Provider等の必要なものを受け取って、各プロバイダー向けのコールバック用ハンドラを返す関数
	return http.NotFoundHandler()
}
