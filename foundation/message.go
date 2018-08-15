package foundation

import "encoding"

// Message is an interface for data that can be sent over the gossip network.
type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	Nonce() uint64
	Data() []byte
	Signature() []byte
}
