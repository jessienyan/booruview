package routes

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"github.com/valkey-io/valkey-go"
)

const (
	settingShareMaxLen = 300_000 // KiB
)

func cacheShareKey(code string) string {
	return "share:" + code
}

// Generates a 12 digit code based on the hash of the data
func generateShareCode(data []byte) string {
	var num int64

	// using a hash gives us a random source for generating the code while also
	// ensuring we don't store the same data twice under two different codes
	hash := sha1.Sum(data)

	// convert the hash to a number with at least 12 digits
	for i := 0; i < len(hash) && num < 100_000_000_000; i++ {
		num = (num << 8) | int64(hash[i])
	}

	// format code as 0000-0000-0000
	code := fmt.Sprintf("%12d", num)
	code = code[:4] + "-" + code[4:8] + "-" + code[8:12]
	return code
}

type settingExportResponse struct {
	Code string `json:"code"`
}

func SettingExportHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, settingExportApiCost) {
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	if !json.Valid(body) {
		respondWithBadRequest(w, "must be valid json")
	}

	if len(body) > settingShareMaxLen {
		respondWithBadRequest(w, fmt.Sprintf("settings data is too large (max %d bytes)", settingShareMaxLen))
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
		respondWithInternalError(w, err)
		return
	}

	respondJson(w, http.StatusOK, settingExportResponse{Code: code})
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
		respondWithInternalError(w, err)
		return
	}

	var data SettingImportRequest
	if err := json.Unmarshal(body, &data); err != nil {
		respondWithInternalError(w, err)
		return
	}

	if len(data.Code) == 0 {
		respondWithBadRequest(w, "`code` is required")
		return
	}

	vc := api.Valkey()
	settings, err := vc.Do(req.Context(), vc.B().Get().Key(cacheShareKey(data.Code)).Build()).AsBytes()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			respondWithNotFound(w, "code is invalid or may have expired")
			return
		}

		respondWithInternalError(w, err)
		return
	}

	respondJson(w, http.StatusOK, settings)
}
