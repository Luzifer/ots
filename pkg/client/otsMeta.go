package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/Luzifer/go-openssl/v4"
)

var metaMarker = []byte("OTSMeta")

type (
	// Secret represents a secret parsed from / prepared for
	// serialization to the OTS API
	Secret struct {
		Secret      string             `json:"secret"`
		Attachments []SecretAttachment `json:"attachments,omitempty"`
	}

	// SecretAttachment represents a file attached to a Secret. The Data
	// property must be the plain content (binary / text / ...) of the
	// file to attach. The base64 en-/decoding is done transparently.
	// The Name is the name of the file shown to the user (so ideally
	// should be the file-name on the source system). The Type should
	// contain the mime time of the file or an empty string.
	SecretAttachment struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Data    string `json:"data"`
		Content []byte `json:"-"`
	}
)

func (o *Secret) read(data []byte, passphrase string) (err error) {
	if passphrase != "" {
		if data, err = openssl.New().DecryptBytes(passphrase, data, KeyDerivationFunc); err != nil {
			return fmt.Errorf("decrypting data: %w", err)
		}
	}

	if !bytes.HasPrefix(data, metaMarker) {
		// We have a simple secret, makes less effort for us
		o.Secret = string(data)
		return nil
	}

	if err = json.Unmarshal(data[len(metaMarker):], o); err != nil {
		return fmt.Errorf("decoding JSON payload: %w", err)
	}

	for i := range o.Attachments {
		o.Attachments[i].Content, err = base64.StdEncoding.DecodeString(o.Attachments[i].Data)
		if err != nil {
			return fmt.Errorf("decoding attachment %d: %w", i, err)
		}
	}

	return nil
}

func (o Secret) serialize(passphrase string) ([]byte, error) {
	var data []byte

	if len(o.Attachments) == 0 {
		// No attachments? No problem, we create a classic simple secret
		data = []byte(o.Secret)
	} else {
		for i := range o.Attachments {
			o.Attachments[i].Data = base64.StdEncoding.EncodeToString(o.Attachments[i].Content)
		}

		j, err := json.Marshal(o)
		if err != nil {
			return nil, fmt.Errorf("encoding JSON payload: %w", err)
		}

		data = append(metaMarker, j...) //nolint:gocritic // :shrug:
	}

	if passphrase == "" {
		// No encryption requested
		return data, nil
	}

	out, err := openssl.New().EncryptBytes(passphrase, data, KeyDerivationFunc)
	if err != nil {
		return nil, fmt.Errorf("encrypting data: %w", err)
	}
	return out, nil
}
