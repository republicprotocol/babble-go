package foundation

type Message struct {
	Nonce     uint64 `json:"nonce"`
	Key       []byte `json:"key"`
	Value     []byte `json:"value"`
	Signature []byte `json:"signature"`
}
