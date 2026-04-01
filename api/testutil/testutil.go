package testutil

import (
	"context"
	"encoding/json"
	"os"
	"time"

	api "codeberg.org/jessienyan/booruview"
	gb "codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/models"
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

func ResetUserData(userID int64) {
	db := models.New(api.UserDB())
	err := db.UpdateUserData(context.Background(), models.UpdateUserDataParams{
		UserID: userID,
		Data:   "",
	})
	if err != nil {
		log.Fatal().Msgf("error resetting user data: %v", err)
	}
}

func UpdateUserData(userID int64, data models.UserData) {
	db := models.New(api.UserDB())
	err := db.UpdateUserData(context.Background(), models.UpdateUserDataParams{
		UserID: userID,
		Data:   data.Data,
	})
	if err != nil {
		log.Fatal().Msgf("error updating user data: %v", err)
	}
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

func (m *MockGelbooruClient) ListTags(tags string) (api.TagList, error) {
	args := m.Called(tags)
	return args.Get(0).(api.TagList), args.Error(1)
}

func (m *MockGelbooruClient) SearchTags(query string) (api.TagList, error) {
	args := m.Called(query)
	return args.Get(0).(api.TagList), args.Error(1)
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

func CreateUser(username, password string) (models.Users, models.UserData) {
	salt := api.GenerateSalt()
	passHash := api.HashPassword(password, salt)

	db := models.New(api.UserDB())
	user, err := db.CreateUser(context.Background(), models.CreateUserParams{
		Username:     username,
		Password:     passHash,
		PasswordSalt: salt,
	})
	if err != nil {
		log.Fatal().Msgf("failed to create user: %v", err)
	}

	data, err := db.CreateUserData(context.Background(), models.CreateUserDataParams{
		UserID: user.ID,
		Data:   "",
	})
	if err != nil {
		log.Fatal().Msgf("failed to create user data: %v", err)
	}

	return user, data
}

func Time() time.Time {
	testTime, _ := time.Parse(time.RFC3339, "2026-04-01T01:23:45Z")
	return testTime
}
