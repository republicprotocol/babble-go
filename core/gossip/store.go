package gossip

// Store is used to read and write Messages that are being disseminated
// throughout the network.
type Store interface {

	// InsertMessage into the store. If there is an existing Message with the
	// same key, but a lower nonce, then the existing Message will be
	// overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	// It returns an empty message with zero nonce if there is no message with
	// the associated key in the store.
	Message(key []byte) (Message, error)
}
