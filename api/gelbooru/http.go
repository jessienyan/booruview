package gelbooru

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"syscall"

	api "github.com/jessienyan/booruview"
	"github.com/rs/zerolog/log"
)

func httpGet(theUrl string, params url.Values) (*http.Response, error) {
	_url, err := url.Parse(theUrl + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	req := &http.Request{URL: _url, Header: http.Header{}}
	req.Header.Set("Cookie", "fringeBenefits=yup") // Enable all content

	resp, err := api.DoRequest(req)
	if err != nil {
		resetByPeer := errors.Is(err, syscall.ECONNRESET)
		if resetByPeer || os.IsTimeout(err) {
			err = errors.Join(GelbooruError{Code: -1}, err)
		}

		return nil, err
	}

	if resp.StatusCode != 200 {
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
