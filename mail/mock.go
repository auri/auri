package mail

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/stanislas-m/mocksmtp"
)

// getMock returns the current instance converted to the mock type. Panics if anything goes wrong
func getMock() *mocksmtp.MockSMTP {
	sender, err := GetSender()
	if err != nil {
		panic(err)
	}
	ms, ok := sender.(*mocksmtp.MockSMTP)
	if ok != true {
		panic("Sender is expected to be MockSMTP")
	}
	return ms
}

//NewMock returns initialized Mail mock
func NewMock(_, _, _, _ string, _, _ bool) (mail.Sender, error) {
	return mocksmtp.New(), nil
}

//ClearMock clears the mock with all messages
func ClearMock() {
	getMock().Clear()
}

// LastMessageFromMock returns the last message from mock
func LastMessageFromMock() mail.Message {
	ms, err := getMock().LastMessage()
	if err != nil {
		panic(err)
	}
	return ms
}

// MessagesFromMock returns all messages from the mock
func MessagesFromMock() []mail.Message {
	return getMock().Messages()
}

// MessageExistsInMock returns true if mail mock has some messages
func MessageExistsInMock() bool {
	return len(getMock().Messages()) > 0
}
