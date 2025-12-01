package gelbooru

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	httpClient = &http.Client{Timeout: 4 * time.Second}
)

func doRequest(req *http.Request) (*http.Response, error) {
	earlier := time.Now()
	resp, err := httpClient.Do(req)

	method := req.Method
	if method == "" {
		method = "GET"
	}

	if err != nil {
		log.Err(err).Str("method", method).Str("url", req.URL.String()).Msg("http request failed")
	} else {
		log.Info().
			Dur("duration", time.Since(earlier)).
			Int("status", resp.StatusCode).
			Str("method", method).
			Str("url", req.URL.String()).
			Send()
	}

	return resp, err
}

func httpGet(theUrl string, params url.Values) (*http.Response, error) {
	_url, err := url.Parse(theUrl + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	req := &http.Request{URL: _url, Header: http.Header{}}
	req.Header.Set("Cookie", "fringeBenefits=yup") // Enable all content

	resp, err := doRequest(req)
	if err != nil {
		resetByPeer := errors.Is(err, syscall.ECONNRESET)
		isTimeout := os.IsTimeout(err)
		isCtxDeadline := errors.Is(err, context.DeadlineExceeded)

		// Timeouts or closed connections generally mean Gelbooru isn't available
		if resetByPeer || isTimeout || isCtxDeadline {
			err = errors.Join(GelbooruError{Code: -1}, err)
		}

		return nil, err
	}

	// Gelbooru is down (cloudflare error)
	if resp.StatusCode == 521 {
		err = GelbooruError{Code: -1}
		return nil, err
	} else if resp.StatusCode != 200 {
		body, _ := httputil.DumpResponse(resp, true)
		log.Error().Msgf("non-200 response: %s", string(body))
		return nil, GelbooruError{Code: resp.StatusCode}
	}

	return resp, nil
}

func httpGetJson[T any](params url.Values, dst T) error {
	resp, err := httpGet(ApiUrl, params)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return GelbooruError{Code: resp.StatusCode}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// This gets received sometimes with status=200, very cool
	if bytes.HasPrefix(body, []byte("Too deep!")) {
		return GelbooruError{Code: 429}
	}

	if err := json.Unmarshal(body, &dst); err != nil {
		log.Err(err).Str("body", string(body)).Msg("failed to parse json")
		return GelbooruError{Code: 500}
	}

	return nil
}
