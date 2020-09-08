package updator

import (
	"encoding/base64"
	"github.com/jamesruan/sodium"
)

func (gc *GithubClient) getRawPublicKey() ([]byte, error) {
	base64PubKey := gc.PublicKey.GetKey()
	pk := make([]byte, 32)

	_, err := base64.StdEncoding.Decode(pk, []byte(base64PubKey))
	if err != nil {
		return []byte(""), err
	}

	return pk, nil
}

func encryptSodium(data string, pk []byte) string {
	enc := sodium.Bytes(data).SealedBox(sodium.BoxPublicKey{Bytes: pk})
	return base64.StdEncoding.EncodeToString(enc)
}
