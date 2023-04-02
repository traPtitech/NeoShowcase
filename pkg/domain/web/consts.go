package web

const (
	HeaderNameAPIAuthorization  = "X-Forwarded-User"
	HeaderNameAuthorizationType = "X-NS-Auth-Type"
	HeaderNameShowcaseUser      = "X-Showcase-User"
	HeaderNameSSGenAppID        = "X-NS-App-Id"
)

const (
	TraefikHTTPEntrypoint     = "web"
	TraefikHTTPSEntrypoint    = "websecure"
	TraefikAuthSoftMiddleware = "nsapp_auth_soft@file"
	TraefikAuthHardMiddleware = "nsapp_auth_hard@file"
	TraefikAuthMiddleware     = "nsapp_auth@file"
	TraefikCertResolver       = "nsresolver@file"
)
