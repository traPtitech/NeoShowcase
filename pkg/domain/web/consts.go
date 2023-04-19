package web

const (
	HeaderNameAPIAuthorization  = "X-Forwarded-User"
	HeaderNameAuthorizationType = "X-Controller-Auth-Type"
	HeaderNameShowcaseUser      = "X-Showcase-User"
	HeaderNameSSGenAppID        = "X-Controller-App-Id"
)

const (
	TraefikHTTPEntrypoint  = "web"
	TraefikHTTPSEntrypoint = "websecure"
)
