package testutil

import (
	"context"

	api "codeberg.org/jessienyan/booruview"
	gb "codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/stretchr/testify/mock"
)

func Setup() {
	api.LoadEnv()
	if err := api.InitValkey(); err != nil {
		panic(err)
	}
}

func FlushCache() {
	vk := api.Valkey()
	err := vk.Do(context.Background(), vk.B().Flushall().Build()).Error()
	if err != nil {
		panic(err)
	}
}

type MockGelbooruClient struct {
	mock.Mock
}

var _ gb.GelbooruClient = (*MockGelbooruClient)(nil)

func (m *MockGelbooruClient) ListPosts(tags string, page int) (*gb.PostList, error) {
	args := m.Called(tags, page)
	return args.Get(0).(*gb.PostList), args.Error(1)
}

func (m *MockGelbooruClient) ListTags(tags string) ([]api.TagResponse, error) {
	args := m.Called(tags)
	return args.Get(0).([]api.TagResponse), args.Error(1)
}

func (m *MockGelbooruClient) SearchTags(query string) ([]api.TagResponse, error) {
	args := m.Called(query)
	return args.Get(0).([]api.TagResponse), args.Error(1)
}
