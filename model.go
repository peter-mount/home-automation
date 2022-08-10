package automation

import (
	"github.com/peter-mount/home-automation/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func (s *Service) LoadModel() error {
	f, err := os.Open(*s.modelFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	house := &model.House{}
	err = yaml.Unmarshal(b, house)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.house = house
	return nil
}

func (s *Service) GetModel() *model.House {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.house
}
