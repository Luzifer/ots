package client

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luzifer/ots/pkg/customization"
	"github.com/stretchr/testify/assert"
)

type custMockClient struct {
	Response *customization.Customize
}

func (c custMockClient) Do(r *http.Request) (*http.Response, error) {
	m := http.NewServeMux()
	m.HandleFunc(r.URL.Path, func(w http.ResponseWriter, _ *http.Request) {
		if c.Response == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		d, _ := c.Response.ToJSON()
		_, _ = w.Write([]byte(d))
	})

	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)

	return w.Result(), nil
}

func TestSanityCheck(t *testing.T) {
	var (
		err error
		m   = custMockClient{&customization.Customize{
			AcceptedFileTypes:      "text/*,image/png,.gif",
			DisableFileAttachment:  true,
			MaxAttachmentSizeTotal: 64,
		}}
		u = "http://localhost/"
	)

	HTTPClient = &m
	defer func() { HTTPClient = http.DefaultClient }()

	s := Secret{Secret: "ohai"}

	// no attachments & attachments disabled
	err = SanityCheck(u, s)
	assert.NoError(t, err)

	// attachments & attachmetns disabled
	s.Attachments = []SecretAttachment{
		{Name: "myfile.webm", Type: "video/webm", Content: []byte{0x0}},
	}

	err = SanityCheck(u, s)
	assert.ErrorIs(t, err, ErrAttachmentsDisabled)

	// disallowed attachment
	m.Response.DisableFileAttachment = false
	err = SanityCheck(u, s)
	assert.ErrorIs(t, err, ErrAttachmentTypeNotAllowed)

	// attachment allowed by extension
	s.Attachments = []SecretAttachment{
		{Name: "doesthiswork.gif", Type: "image/gif", Content: []byte{0x0}},
	}
	err = SanityCheck(u, s)
	assert.NoError(t, err)

	// attachment allowed by mime type
	s.Attachments = []SecretAttachment{
		{Name: "doesthiswork.png", Type: "image/png", Content: []byte{0x0}},
	}
	err = SanityCheck(u, s)
	assert.NoError(t, err)

	// attachment allowed by mime type wildcard
	s.Attachments = []SecretAttachment{
		{Name: "doesthiswork.md", Type: "text/markdown", Content: []byte{0x0}},
	}
	err = SanityCheck(u, s)
	assert.NoError(t, err)

	// attachment too large
	s.Attachments = []SecretAttachment{
		{Name: "doesthiswork.md", Type: "text/markdown", Content: bytes.Repeat([]byte{0x0}, 128)},
	}
	err = SanityCheck(u, s)
	assert.ErrorIs(t, err, ErrAttachmentsTooLarge)

	// check without settings API on instance
	m.Response = nil
	err = SanityCheck(u, s)
	assert.NoError(t, err)
}
