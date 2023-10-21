// Package client implements a client library for OTS supporting the
// OTSMeta content format for file upload support
package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Luzifer/go-openssl/v4"
)

type (
	// HTTPClientIntf describes a minimal interface to be fulfilled
	// by the given HTTP client. This can be used for mocking and to
	// pass in authenticated clients
	HTTPClientIntf interface {
		Do(*http.Request) (*http.Response, error)
	}
)

// HTTPClient defines the client to use for create and fetch requests
// and can be overwritten to provide authentication
var HTTPClient HTTPClientIntf = http.DefaultClient

// KeyDerivationFunc defines the key derivation algorithm used in OTS
// to derive the key / iv from the password for encryption. You only
// should change this if you are running an OTS instance with modified
// parameters.
//
// The corresponding settings are found in `/src/crypto.js` in the OTS
// source code.
var KeyDerivationFunc = openssl.NewPBKDF2Generator(sha512.New, 300000) //nolint:gomnd // that's the definition

// PasswordLength defines the length of the generated encryption password
var PasswordLength = 20

// RequestTimeout defines how long the request to the OTS instance for
// create and fetch may take
var RequestTimeout = 5 * time.Second

// UserAgent defines the user-agent to send when interacting with an
// OTS instance. When using this library please set this to something
// the operator of the instance can determine your client from and
// provide an URL to useful information about your tool.
var UserAgent = "ots-client/1.x +https://github.com/Luzifer/ots"

// Create serializes the secret and creates a new secret on the
// instance given by its URL.
//
// The given URL should point to the frontend of the instance. Do not
// include the API paths, they are added automatically. For the
// expireIn parameter zero value can be used to use server-default.
//
// So for OTS.fyi you'd use `New("https://ots.fyi/")`
func Create(instanceURL string, secret Secret, expireIn time.Duration) (string, time.Time, error) {
	u, err := url.Parse(instanceURL)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("parsing instance URL: %w", err)
	}

	pass, err := genPass()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("generating password: %w", err)
	}

	data, err := secret.serialize(pass)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("serializing data: %w", err)
	}

	body := new(bytes.Buffer)
	if err = json.NewEncoder(body).Encode(struct {
		Secret string `json:"secret"`
	}{Secret: string(data)}); err != nil {
		return "", time.Time{}, fmt.Errorf("encoding request payload: %w", err)
	}

	createURL := u.JoinPath(strings.Join([]string{".", "api", "create"}, "/"))
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	if expireIn > time.Second {
		createURL.RawQuery = url.Values{
			"expire": []string{strconv.Itoa(int(expireIn / time.Second))},
		}.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, createURL.String(), body)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // possible leaked-fd, lib should not log, potential short-lived leak

	if resp.StatusCode != http.StatusCreated {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", time.Time{}, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
		}
		return "", time.Time{}, fmt.Errorf("unexpected HTTP status %d (%s)", resp.StatusCode, respBody)
	}

	var payload struct {
		ExpiresAt time.Time `json:"expires_at"`
		SecretID  string    `json:"secret_id"`
		Success   bool      `json:"success"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", time.Time{}, fmt.Errorf("decoding response: %w", err)
	}

	u.Fragment = strings.Join([]string{payload.SecretID, pass}, "|")

	return u.String(), payload.ExpiresAt, nil
}

// Fetch retrieves a secret by its given URL. The URL given must
// include the fragment (part after the `#`) with the secret ID and
// the encryption passphrase.
//
// The object returned will always be an OTSMeta object even in case
// the secret is a plain secret without attachments.
func Fetch(secretURL string) (s Secret, err error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		return s, fmt.Errorf("parsing secret URL: %w", err)
	}

	fragment, err := url.QueryUnescape(u.Fragment)
	if err != nil {
		return s, fmt.Errorf("unescaping fragment: %w", err)
	}
	fragmentParts := strings.SplitN(fragment, "|", 2) //nolint:gomnd

	fetchURL := u.JoinPath(strings.Join([]string{".", "api", "get", fragmentParts[0]}, "/")).String()
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fetchURL, nil)
	if err != nil {
		return s, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return s, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck // possible leaked-fd, lib should not log, potential short-lived leak

	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)
	}

	var payload struct {
		Secret string `json:"secret"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return s, fmt.Errorf("decoding response body: %w", err)
	}

	if err = s.read([]byte(payload.Secret), fragmentParts[1]); err != nil {
		return s, fmt.Errorf("decoding secret: %w", err)
	}

	return s, nil
}

func genPass() (string, error) {
	var (
		charSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		pass    = make([]byte, PasswordLength)

		n   int
		err error
	)

	for n < PasswordLength {
		n, err = rand.Read(pass)
		if err != nil {
			return "", fmt.Errorf("reading random data: %w", err)
		}
	}

	for i := 0; i < PasswordLength; i++ {
		pass[i] = charSet[int(pass[i])%len(charSet)]
	}

	return string(pass), nil
}
