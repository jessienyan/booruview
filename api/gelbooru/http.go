package gelbooru

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	stderrors "errors"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

var (
	defaultHTTPClient = &http.Client{Timeout: 4 * time.Second}
)

func doRequest(httpClient *http.Client, req *http.Request) (*http.Response, error) {
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

func httpGet(httpClient *http.Client, theUrl string, params url.Values) (*http.Response, error) {
	_url, err := url.Parse(theUrl + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	req := &http.Request{URL: _url, Header: http.Header{}}
	req.Header.Set("Cookie", "fringeBenefits=yup") // Enable all content

	resp, err := doRequest(httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", stderrors.Join(GelbooruError{Code: 503}, err))
	}

	if resp.StatusCode != 200 {
		body, _ := httputil.DumpResponse(resp, true)
		log.Error().Msgf("non-200 response: %s", string(body))
		err := errors.Wrap(GelbooruError{Code: resp.StatusCode}, "non-200 response")
		return nil, err
	}

	return resp, nil
}

func httpGetJson[T any](httpClient *http.Client, apiUrl string, params url.Values, dst T) error {
	resp, err := httpGet(httpClient, apiUrl, params)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return GelbooruError{Code: resp.StatusCode}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return transformTimeoutError(err)
	}

	if len(body) == 0 {
		return errors.Wrap(GelbooruError{Code: resp.StatusCode}, "empty response body, expected JSON")
	}

	if err := json.Unmarshal(body, &dst); err != nil {
		err = errors.Wrap(err, "failed to parse json")
		log.Err(err).Msg("")
		return err
	}

	return nil
}
