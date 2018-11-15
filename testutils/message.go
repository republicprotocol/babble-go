package testutils

import (
	"sync"

	"github.com/republicprotocol/babble-go/core/gossip"
)

type MockMessages struct {
	messageMu *sync.Mutex
	messages  map[string]gossip.Message
}

func NewMockMessages() MockMessages {
	return MockMessages{
		messageMu: new(sync.Mutex),
		messages:  map[string]gossip.Message{},
	}
}

func (messages MockMessages) InsertMessage(message gossip.Message) error {
	messages.messageMu.Lock()
	defer messages.messageMu.Unlock()
	messages.messages[string(message.Key)] = message

	return nil
}

func (messages MockMessages) Message(key []byte) (gossip.Message, error) {
	messages.messageMu.Lock()
	defer messages.messageMu.Unlock()

	return messages.messages[string(key)], nil
}
