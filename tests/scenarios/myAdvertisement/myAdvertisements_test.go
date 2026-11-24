package myAdvertisement

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	advList "api-tests-template/internal/managers/advertisements"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	title := fmt.Sprintf("QA My Ads %d", time.Now().UnixNano())
	photoPath := filepath.Join("..", "advertisementCreation", "testdata", "photo.png")

	created := advList.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		"Тестовое объявление",
		100,
		1,
		photoPath,
	)

	require.NotEmpty(s.T(), created.ID, "Не получен ID созданного объявления")

	s.T().Cleanup(func() {
		advList.DeleteAdvertisement(s.T(), s.loginData.Token, created.ID)
	})

	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(
			s.T(), s.loginData.Token, http.StatusOK, 50,
		)
	})

	var items gjson.Result
	s.Run("Проверяем наличие объявлений", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	})

	s.Run("Проверяем объявления и наличие созданного", func() {
		found := false

		for _, item := range items.Array() {
			require.Equal(s.T(), s.loginData.User.Id,
				gjson.Get(item.Raw, "user_id").String(), "user_id не совпадает")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "id").String(), "id пустое")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "updated_at").String(), "updated_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "title").String(), "Название пустое")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "description").String(), "Описание пустое")
			require.NotEmpty(s.T(), gjson.Get(item.Raw, "photos").Array(), "Фотографии отсутствуют")

			if gjson.Get(item.Raw, "id").String() == created.ID {
				found = true
				require.Equal(s.T(), title,
					gjson.Get(item.Raw, "title").String(), "Название не совпадает")
			}
		}

		require.True(s.T(), found,
			"Созданное объявление с ID %s не найдено в списке", created.ID)
	})
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления без токена
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), "incorrect_token", http.StatusUnauthorized, 50)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), "unauthorized", gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), "Invalid or expired token", gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
