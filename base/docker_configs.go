package base

type DockerConfigs struct {
	Image   string
	Tag     string
	EnvVars []string
	Cmd     []string
}
