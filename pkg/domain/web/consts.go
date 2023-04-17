package web

const (
	HeaderNameAPIAuthorization  = "X-Forwarded-User"
	HeaderNameAuthorizationType = "X-Controller-Auth-Type"
	HeaderNameShowcaseUser      = "X-Showcase-User"
	HeaderNameSSGenAppID        = "X-Controller-App-Id"
)

const (
	TraefikHTTPEntrypoint     = "web"
	TraefikHTTPSEntrypoint    = "websecure"
	TraefikAuthSoftMiddleware = "nsapp_auth_soft@file"
	TraefikAuthHardMiddleware = "nsapp_auth_hard@file"
	TraefikAuthMiddleware     = "nsapp_auth@file"
)
