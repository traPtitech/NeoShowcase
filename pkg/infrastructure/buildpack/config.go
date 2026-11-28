package buildpack

import "github.com/traPtitech/neoshowcase/pkg/domain/builder"

// DefaultImage is the builder image of Cloud Native Buildpacks used when none is configured.
// Its tag is pinned so that every step of a build, and every builder, runs the same image.
const DefaultImage = "paketobuildpacks/builder-jammy-full:0.3.724"

// DefaultRunAs is the build user of DefaultImage, used when none is configured.
var DefaultRunAs = builder.UIDGID{UID: 1001, GID: 1000}

type Config struct {
	// Image is the builder image of Cloud Native Buildpacks.
	Image string
	// RunAs is the build user of Image.
	RunAs       builder.UIDGID
	PlatformAPI string
}
