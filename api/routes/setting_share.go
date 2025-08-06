package routes

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "github.com/jessienyan/booruview"
	"github.com/valkey-io/valkey-go"
)

const (
	settingShareMaxLen = 150_000 // KiB

	settingExportApiCost = 3
	settingImportApiCost = 1
)

func cacheShareKey(code string) string {
	return "share:" + code
}

// Generates a 12 digit code based on the hash of the data
func generateShareCode(data []byte) string {
	hash := sha1.Sum(data)

	var num int64
	for i := 0; i < len(hash) && num < 100_000_000_000; i++ {
		num = (num << 8) | int64(hash[i])
	}

	code := fmt.Sprintf("%12d", num)
	code = code[:4] + "-" + code[4:8] + "-" + code[8:12]
	return code
}

func SettingExportHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, settingExportApiCost) {
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	if !json.Valid(body) {
		handle400Error(w, "must be valid json")
	}

	if len(body) > settingShareMaxLen {
		handle400Error(w, fmt.Sprintf("settings data is too large (max %d bytes)", settingShareMaxLen))
		return
	}

	code := generateShareCode(body)
	vc := api.Valkey()

	// Write to redis
	err = vc.Do(
		req.Context(),
		vc.B().Setex().
			Key(cacheShareKey(code)).
			Seconds(api.SettingShareTtl).
			Value(string(body)).
			Build()).Error()

	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"code":"%s"}`, code)))
}

type SettingImportRequest struct {
	Code string `json:"code"`
}

func SettingImportHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, settingImportApiCost) {
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	var data SettingImportRequest
	if err := json.Unmarshal(body, &data); err != nil {
		handleError(w, err)
		return
	}

	if len(data.Code) == 0 {
		handle400Error(w, "`code` is required")
		return
	}

	vc := api.Valkey()
	stored, err := vc.Do(req.Context(), vc.B().Get().Key(cacheShareKey(data.Code)).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			handle4xxError(w, 404, "code is invalid or may have expired")
			return
		}

		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(stored))
}
