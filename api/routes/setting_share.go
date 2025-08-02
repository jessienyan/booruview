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
	settingShareMaxLen = 50_000 // KiB
)

func cacheShareKey(key string) string {
	return "share:" + key
}

// Generates a 12 digit key based on the hash of the data
func generateShareKey(data SettingImportRequest) string {
	hash := sha1.Sum([]byte(data.Data))

	var num int64
	for i := 0; i < len(hash) && num < 100_000_000_000; i++ {
		num = (num << 8) | int64(hash[i])
	}

	key := fmt.Sprintf("%12d", num)
	key = key[:4] + "-" + key[4:8] + "-" + key[8:12]
	return key
}

type SettingImportRequest struct {
	Data string `json:"data"`
}

func SettingImportHandler(w http.ResponseWriter, req *http.Request) {
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

	if len(data.Data) == 0 {
		handle400Error(w, "`data` is required")
		return
	}

	if len(data.Data) > settingShareMaxLen {
		handle400Error(w, fmt.Sprintf("settings data is too large (max %d bytes)", settingShareMaxLen))
		return
	}

	key := generateShareKey(data)
	vc := api.Valkey()

	// Write to redis
	err = vc.Do(
		req.Context(),
		vc.B().Setex().
			Key(cacheShareKey(key)).
			Seconds(api.SettingShareTtl).
			Value(data.Data).
			Build()).Error()

	if err != nil {
		handleError(w, err)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"key":"%s"}`, key)))
}

type SettingExportRequest struct {
	Key string `json:"key"`
}

func SettingExportHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(w, err)
		return
	}

	var data SettingExportRequest
	if err := json.Unmarshal(body, &data); err != nil {
		handleError(w, err)
		return
	}

	if len(data.Key) == 0 {
		handle400Error(w, "`key` is required")
		return
	}

	vc := api.Valkey()
	stored, err := vc.Do(req.Context(), vc.B().Get().Key(cacheShareKey(data.Key)).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			handle4xxError(w, 404, "key is invalid or may have expired")
			return
		}

		handleError(w, err)
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"data":"%s"}`, stored)))
}
