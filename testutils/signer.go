package testutils

type MockSinger struct {
}

func (signer MockSinger) Sign(data []byte) ([]byte, error) {
	return data, nil
}

// A mock verifier will always return true when verifying signature.
type MockVerifier struct {
}

func (verifier MockVerifier) Verify(data []byte, signature []byte) error {
	return nil
}
