package gossip

// A Message is a unit of data that can be disseminated throughout the network.
// An outdated Message can be overwritten by disseminating a newer Message with
// the same `Key` but an incremented `Nonce`. Nodes in the network will discard
// the lower `Nonce` Message in favour of the higher `Nonce` Message. A
// `Signature` is used to verify the authenticity of the Message.
type Message struct {
	Nonce     uint64 `json:"nonce"`
	Key       []byte `json:"key"`
	Value     []byte `json:"value"`
	Signature []byte `json:"signature"`
}

// NewMessage returns a new Message with given nonce, key, value and signature.
func NewMessage(nonce uint64, key, value, signature []byte) Message {
	return Message{nonce, key, value, signature}
}

// Messages is used to read and write Messages to persistent storage.
type Messages interface {

	// InsertMessage into the MessageStore. If there is an existing Message
	// with the same key, but a lower nonce, then the existing Message will
	// be overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	// It returns an empty message with zero nonce if there is no message with
	// the associated key in the store.
	Message(key []byte) (Message, error)
}
