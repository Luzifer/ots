package client

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePassword(t *testing.T) {
	pass, err := genPass()
	require.NoError(t, err)

	assert.Len(t, pass, PasswordLength)
	assert.Regexp(t, regexp.MustCompile(`^[0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ]+$`), pass)
}

func TestIntegration(t *testing.T) {
	s := Secret{
		Secret: "I'm a secret!",
		Attachments: []SecretAttachment{{
			Name:    "secret.txt",
			Type:    "text/plain",
			Content: []byte("I'm a very secret file.\n"),
		}},
	}

	secretURL, _, err := Create("https://ots.fyi/", s, time.Minute)
	require.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile(`^https://ots.fyi/#[0-9a-f-]+%7C[0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ]+$`), secretURL)

	apiSecret, err := Fetch(secretURL)
	require.NoError(t, err)

	assert.Equal(t, s, apiSecret)
}
