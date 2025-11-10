package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type key string

var NonceKey key = "nonces"

type Nonces struct {
	Htmx            string
	Alpine          string
	ResponseTargets string
	Tw              string
	HtmxCSSHash     string
	AlpineCSSHash   string
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			Alpine:          generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
			AlpineCSSHash:   "sha256-faU7yAF8NxuMTNEwVmBz+VcYeIoBQ2EMHW3WaVxCvnk=",
		}

		// set nonces in context
		ctx := context.WithValue(r.Context(), NonceKey, nonceSet)

		cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' 'nonce-%s'; style-src 'nonce-%s' '%s' '%s'; frame-src 'self';",
			nonceSet.Htmx,
			nonceSet.Alpine,
			//			nonceSet.ResponseTargets,
			nonceSet.Tw,
			nonceSet.HtmxCSSHash,
			nonceSet.AlpineCSSHash,
		)
		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func GetNonces(ctx context.Context) Nonces {
	nonceSet := ctx.Value(NonceKey)
	if nonceSet == nil {
		log.Fatal("error getting nonce set - is nil")
	}

	nonces, ok := nonceSet.(Nonces)
	if !ok {
		log.Fatal("error getting nonce set - not ok")
	}

	return nonces
}

func GetHtmxNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Htmx
}

func GetAlpineNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Alpine
}

func GetResponseTargetsNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.ResponseTargets
}

func GetTwNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Tw
}

func GetHtmxCSSHashNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.HtmxCSSHash
}
