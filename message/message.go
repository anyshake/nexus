package message

func New(msg string) (Message, error) {
	messageObj := Message{RawMessage: msg}
	if err := messageObj.Parse(); err != nil {
		return Message{}, err
	}

	return messageObj, nil
}
