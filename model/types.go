package model

type Service struct {
	ExposedPort       string
	ExposedIP         string
	ContainerID       string
	ContainerHostName string
	Domain            string
	Endpoint          string
}
