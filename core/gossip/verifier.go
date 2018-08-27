package gossip

// A Signer can consume bytes and produce a unique signature for those bytes.
// This is usually done using the private key of an asymmetrical cryptography
// scheme.
type Signer interface {
	Sign(data []byte) ([]byte, error)
}

// A Verifier can consume bytes and a signature for those bytes, and extract
// signatory.
type Verifier interface {
	Verify(data []byte, signature []byte) error
}
