group "default" {
  targets = [
    "ns-dashboard",
    "ns-sablier",
    "ns",
    "ns-components",
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

target "ns" {
  inherits = ["base"]
  target   = "ns"
  tags     = tags("ns")
}

target "ns-components" {
  inherits = ["base"]

  matrix = {
    component = [
      "migrate",
      "auth-dev", 
      "builder",
      "controller",
      "gateway",
      "gitea-integration",
      "ssgen",
    ]
  }
  
  name   = "ns-${component}"
  target = "ns-${component}"
  tags   = tags("ns-${component}")
}
