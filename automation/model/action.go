package model

type Action struct {
	Scene  string             `yaml:"scene"`  // Scene to activate
	Global *map[string]string `yaml:"global"` // Global variables to set
}
