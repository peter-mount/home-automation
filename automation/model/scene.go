package model

type Scene struct {
	Id          string            `yaml:"id"`
	Description string            `yaml:"description"`
	Devices     map[string]string `yaml:"devices"`
}
