package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadOTSMeta(t *testing.T) {
	var (
		//#nosec:G101 // Hardcoded credentials, just test-data
		secretData = "U2FsdGVkX1+7kNgAK57O/qdbsukK3OchMyMyE1tWzVJVlc9f9bkp8iaFHbwR7Q3b8tWhWmPAcfeOoBJH2zl1iNbIHWsmMKu3+pzE5wTE4wl31dOboV8LgsMChBFL5RQpda0iGku32BcB4tYEyb2VHcM/kkXNJh9lW1vRyiNx0iF8pe05JUkkmJJrnzIKC+/efZEfF2YX7fOaBC1+8AAhlg=="
		//#nosec:G101 // Hardcoded credentials, just test-data
		pass = "IKeiXsyGuVWdMUG8Fj3R"
		s    Secret
	)

	err := s.read([]byte(secretData), pass)
	require.NoError(t, err)

	assert.Equal(t, Secret{
		Secret: "I'm a secret!",
		Attachments: []SecretAttachment{{
			Name:    "secret.txt",
			Type:    "text/plain",
			Data:    "SSdtIGEgdmVyeSBzZWNyZXQgZmlsZS4K",
			Content: []byte("I'm a very secret file.\n"),
		}},
	}, s)
}

func TestReadSimpleSecret(t *testing.T) {
	var (
		//#nosec:G101 // Hardcoded credentials, just test-data
		secretData = "U2FsdGVkX18cvbYVRsD5cxMKKAHtMRmteu88tPwRtOk="
		//#nosec:G101 // Hardcoded credentials, just test-data
		pass = "YQHdft6hDnp575olczeq"
		s    Secret
	)

	err := s.read([]byte(secretData), pass)
	require.NoError(t, err)

	assert.Equal(t, Secret{
		Secret: "I'm a secret!",
	}, s)
}

func TestSerializeOTSMeta(t *testing.T) {
	// NOTE(kahlers): We're using an empty passphrase here to achieve
	// testability of the output. The data is not encrypted in this
	// case.
	data, err := Secret{
		Secret: "I'm a secret!",
		Attachments: []SecretAttachment{{
			Name:    "secret.txt",
			Type:    "text/plain",
			Content: []byte("I'm a very secret file.\n"),
		}},
	}.serialize("")
	require.NoError(t, err)

	assert.Equal(t, []byte(`OTSMeta{"secret":"I'm a secret!","attachments":[{"name":"secret.txt","type":"text/plain","data":"SSdtIGEgdmVyeSBzZWNyZXQgZmlsZS4K"}]}`), data)
}

func TestSerializeSimpleSecret(t *testing.T) {
	// NOTE(kahlers): We're using an empty passphrase here to achieve
	// testability of the output. The data is not encrypted in this
	// case.
	data, err := Secret{Secret: "I'm a secret!"}.serialize("")
	require.NoError(t, err)

	assert.Equal(t, []byte("I'm a secret!"), data)
}
