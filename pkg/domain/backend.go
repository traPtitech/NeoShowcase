package domain

type ContainerCreateArgs struct {
	ApplicationID string
	EnvironmentID string
	ImageName     string
	ImageTag      string
	Labels        map[string]string
	Envs          map[string]string
	HTTPProxy     *ContainerHTTPProxy
	Recreate      bool
}

type ContainerHTTPProxy struct {
	Domain string
	Port   int
}

type Container struct {
	ApplicationID string
	EnvironmentID string
	State         ContainerState
}

type ContainerState int

const (
	ContainerStateRunning ContainerState = iota
	ContainerStateRestarting
	ContainerStateStopped
	ContainerStateOther
)
