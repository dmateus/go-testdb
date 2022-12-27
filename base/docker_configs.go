package base

type DockerConfigs struct {
	Image   string
	Tag     string
	EnvVars []string
}

func NewDockerConfigs(image string, tag string, envVars []string) *DockerConfigs {
	return &DockerConfigs{Image: image, Tag: tag, EnvVars: envVars}
}
