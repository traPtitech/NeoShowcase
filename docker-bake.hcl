group "default" {
  targets = [
    "ns-dashboard",
    "ns-sablier",
    "ns-migrate",
    "ns",
    "ns-auth-dev",
    "ns-builder",
    "ns-controller",
    "ns-gateway",
    "ns-gitea-integration",
    "ns-ssgen",
  ]
}

variable "APP_VERSION" {
  type    = string
  default = "dev"
}

variable "APP_REVISION" {
  type    = string
  default = "local"
}

# We accept comma-separated tags list here, which is generated in docker/metadata-action with sep-tags ",".
variable "TAGS" {
  type    = string
  default = "main"
}

function "tags" {
  params = [image]
  result = [for tag in split(",", TAGS) : "ghcr.io/traptitech/${image}:${tag}"]
}

# https://github.com/docker/metadata-action?tab=readme-ov-file#bake-definition
target "docker-metadata-action" {}

target "base" {
  inherits = ["docker-metadata-action"]
  args = {
    APP_VERSION  = APP_VERSION
    APP_REVISION = APP_REVISION
  }
  platforms = [
    "linux/amd64",
    "linux/arm64",
  ]
}

target "ns-dashboard" {
  inherits = ["base"]
  context  = "./dashboard"
  tags     = tags("ns-dashboard")
}

target "ns-sablier" {
  inherits = ["base"]
  context  = "./sablier"
  tags     = tags("ns-sablier")
}

target "ns-migrate" {
  inherits = ["base"]
  target   = "ns-migrate"
  tags     = tags("ns-migrate")
}

target "ns" {
  inherits = ["base"]
  target   = "ns"
  tags     = tags("ns")
}

target "ns-auth-dev" {
  inherits = ["base"]
  target   = "ns-auth-dev"
  tags     = tags("ns-auth-dev")
}

target "ns-builder" {
  inherits = ["base"]
  target   = "ns-builder"
  tags     = tags("ns-builder")
}

target "ns-controller" {
  inherits = ["base"]
  target   = "ns-controller"
  tags     = tags("ns-controller")
}

target "ns-gateway" {
  inherits = ["base"]
  target   = "ns-gateway"
  tags     = tags("ns-gateway")
}

target "ns-gitea-integration" {
  inherits = ["base"]
  target   = "ns-gitea-integration"
  tags     = tags("ns-gitea-integration")
}

target "ns-ssgen" {
  inherits = ["base"]
  target   = "ns-ssgen"
  tags     = tags("ns-ssgen")
}
