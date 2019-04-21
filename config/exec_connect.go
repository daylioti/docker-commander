package config

// ExecConnect docker execute configs (where execute the command).
type ExecConnect struct {
	FromImage     string `yaml:"container_image"` // The name of the image from which the container is made.
	ContainerName string `yaml:"container_name"`  // Container Name
	ContainerID   string `yaml:"container_id"`    // Container id
}
