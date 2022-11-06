package bot

import "github.com/gempir/gempbot/internal/dto"

type Mockbot struct {
}

func NewMockbot() *Mockbot {
	return &Mockbot{}
}

func (mb *Mockbot) RegisterCommand(command string, handler func(dto.CommandPayload)) {
}

func (mb *Mockbot) Say(channel string, message string) {
}

func (mb *Mockbot) Reply(channel string, parentMsgId string, message string) {
}
