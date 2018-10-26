package gossip

// A Signer can consume bytes and produce a signature for those bytes. This
// signature can be used by a Verifier to extract the signatory.
type Signer interface {
	Sign(data []byte) ([]byte, error)
}

// A Verifier can consume bytes and a signature for those bytes, and extract
// the signatory.
type Verifier interface {
	Verify(data []byte, signature []byte) error
}
