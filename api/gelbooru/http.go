package gelbooru

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"

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
			err = errors.Wrap(GelbooruError{Code: 503}, "gelbooru timeout")
		}

		return nil, err
	}

	if resp.StatusCode != 200 {
		body, _ := httputil.DumpResponse(resp, true)
		log.Error().Msgf("non-200 response: %s", string(body))
		err := errors.Wrap(GelbooruError{Code: resp.StatusCode}, "non-200 response")
		return nil, err
	}

	return resp, nil
}

func httpGetJson[T any](params url.Values, dst T) error {
	resp, err := httpGet(ApiUrl, params)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &dst); err != nil {
		err = errors.Wrap(err, "failed to parse json")
		log.Err(err).Msg("")
		return err
	}

	return nil
}
