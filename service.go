package automation

import (
	"github.com/peter-mount/go-kernel/util/task"
	"github.com/peter-mount/home-automation/model"
	"github.com/peter-mount/home-automation/mq"
	"github.com/peter-mount/home-automation/state"
	"sync"
)

type Service struct {
	mq        *mq.MQ         `kernel:"inject"`
	states    *state.Service `kernel:"inject"`
	queueName *mq.Queue      `kernel:"config,automationQueue"`
	publisher *mq.Publisher  `kernel:"config,automationPublisher"`
	modelFile *string        `kernel:"config,modelFile"`
	worker    task.Queue     `kernel:"worker"`
	mutex     sync.Mutex
	house     *model.House
}

func (s *Service) Start() error {
	err := s.mq.AttachPublisher(s.publisher)
	if err != nil {
		return err
	}

	err = s.LoadModel()
	if err != nil {
		return err
	}

	return s.mq.ConsumeTask(s.queueName, "automation", mq.Guard(s.processMqMessage))
}
