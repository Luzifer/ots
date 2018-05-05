package http

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Luzifer/go_helpers/str"
)

func GetDigestAuth(resp *http.Response, method, requestPath, user, password string) string {
	params := map[string]string{}
	for _, part := range strings.Split(resp.Header.Get("Www-Authenticate"), " ") {
		if !strings.Contains(part, `="`) {
			continue
		}
		spl := strings.Split(strings.Trim(part, " ,"), "=")
		if !str.StringInSlice(spl[0], []string{"nonce", "realm", "qop"}) {
			continue
		}
		params[spl[0]] = strings.Trim(spl[1], `"`)
	}

	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)

	params["cnonce"] = fmt.Sprintf("%x", b)
	params["nc"] = "1"
	params["uri"] = requestPath
	params["username"] = user
	params["response"] = getMD5([]string{
		getMD5([]string{params["username"], params["realm"], password}),
		params["nonce"],
		params["nc"],
		params["cnonce"],
		params["qop"],
		getMD5([]string{method, requestPath}),
	})

	authParts := []string{}
	for k, v := range params {
		authParts = append(authParts, fmt.Sprintf("%s=%q", k, v))
	}
	return "Digest " + strings.Join(authParts, ", ")
}

func getMD5(in []string) string {
	h := md5.New()
	h.Write([]byte(strings.Join(in, ":")))
	return hex.EncodeToString(h.Sum(nil))
}
