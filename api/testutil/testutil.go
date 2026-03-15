package testutil

import (
	"context"
	"encoding/json"
	"os"
	"time"

	api "codeberg.org/jessienyan/booruview"
	gb "codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

func Setup() {
	// Only log fatal messages during tests
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro}).
		Level(zerolog.FatalLevel)

	api.LoadEnv()
	if err := api.InitValkey(); err != nil {
		panic(err)
	}
}

func Flush() {
	vk := api.Valkey()
	err := vk.Do(context.Background(), vk.B().Flushall().Build()).Error()
	if err != nil {
		log.Fatal().Msgf("error flushing valkey: %v", err)
	}

	api.FlushRateLimits()
}

type MockGelbooruClient struct {
	mock.Mock
}

var _ gb.GelbooruClient = (*MockGelbooruClient)(nil)

func (m *MockGelbooruClient) ListPosts(tags string, page int) (*gb.PostList, error) {
	args := m.Called(tags, page)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*gb.PostList), args.Error(1)
}

func (m *MockGelbooruClient) ListTags(tags string) ([]api.TagResponse, error) {
	args := m.Called(tags)
	return args.Get(0).([]api.TagResponse), args.Error(1)
}

func (m *MockGelbooruClient) SearchTags(query string) ([]api.TagResponse, error) {
	args := m.Called(query)
	return args.Get(0).([]api.TagResponse), args.Error(1)
}

func MustUnmarshalJSON(data []byte, dst any) {
	if err := json.Unmarshal(data, dst); err != nil {
		log.Fatal().Msgf("failed to unmarshal JSON: %v", err)
	}
}

func MustMarshalJSON(val any) []byte {
	data, err := json.Marshal(val)
	if err != nil {
		log.Fatal().Msgf("failed to marshal JSON: %v", err)
	}
	return data
}
