package gossip

// MessageStore is used to read and write Message disseminated throughout the
// network.
type MessageStore interface {

	// InsertMessage into the MessageStore. If there is an existing Message
	// with the same key, then the existing Message will be overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	// It returns an empty message with zero nonce if there is no message with
	// the associated key in the store.
	Message(key []byte) (Message, error)
}
