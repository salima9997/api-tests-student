package favorites

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	advertisements "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	favoritesManager "api-tests-template/internal/managers/favorites"
	base "api-tests-template/tests"
)

func TestFavoritesNegative(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	token := loginData.Token
	nonexistentID := "11111111-1111-4111-8111-111111111111"

	cases := []struct {
		name            string
		expectedMessage string
		run             func(t *testing.T) string
	}{
		{
			name:            "Добавление несуществующего объявления",
			expectedMessage: "Advertisement not found",
			run: func(t *testing.T) string {
				resp := favoritesManager.AddFavoriteExpectError(
					t, token, nonexistentID, http.StatusNotFound,
				)
				return resp.Message
			},
		},
		{
			name:            "Удаление несуществующего объявления",
			expectedMessage: "Favorite entry not found",
			run: func(t *testing.T) string {
				resp := favoritesManager.DeleteFavoriteExpectError(
					t, token, nonexistentID, http.StatusNotFound,
				)
				return resp.Message
			},
		},
		{
			name:            "Получение избранного без токена",
			expectedMessage: "Authorization header required",
			run: func(t *testing.T) string {
				resp := favoritesManager.GetFavoritesExpectError(
					t, "", http.StatusUnauthorized,
				)
				return resp.Message
			},
		},
		{
			name:            "Добавление в избранное без токена",
			expectedMessage: "Authorization header required",
			run: func(t *testing.T) string {
				resp := favoritesManager.AddFavoriteExpectError(
					t, "", nonexistentID, http.StatusUnauthorized,
				)
				return resp.Message
			},
		},
		{
			name:            "Удаление из избранного без токена",
			expectedMessage: "Authorization header required",
			run: func(t *testing.T) string {
				resp := favoritesManager.DeleteFavoriteExpectError(
					t, "", nonexistentID, http.StatusUnauthorized,
				)
				return resp.Message
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Отправляем негативный запрос")

			message := tc.run(t)

			require.Equal(t, tc.expectedMessage, message, "Некорректное сообщение об ошибке")

			t.Logf("Сообщение об ошибке: %s", message)
		})
	}
}

func TestAddOwnAdvertisementToFavorites(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	token := loginData.Token

	photoPath := filepath.Join(
		"..", "advertisementCreation", "testdata", "photo.png",
	)

	t.Log("Подготовка: создаём собственное объявление")

	created := advertisements.CreateAdvertisement(
		t,
		token,
		fmt.Sprintf("QA Own %d", time.Now().UnixNano()),
		"Проверка добавления своего объявления в избранное",
		100,
		1,
		photoPath,
	)

	require.NotEmpty(t, created.ID, "Не получен ID объявления")

	t.Cleanup(func() {
		t.Log("Очистка: удаляем тестовое объявление")
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	t.Log("Пытаемся добавить своё объявление в избранное")

	resp := favoritesManager.AddFavoriteExpectError(
		t,
		token,
		created.ID,
		http.StatusForbidden,
	)

	require.Equal(t, "Cannot add own advertisement to favorites",
		resp.Message, "Некорректное сообщение об ошибке")

	t.Logf("Сообщение об ошибке: %s", resp.Message)
}

func TestAddFavoriteTwice(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	token := loginData.Token

	t.Log("Подготовка: ищем чужое объявление")

	advertisementID := createForeignTestAdvertisement(
		t, loginData.User.Id,
	)

	require.NotEmpty(t, advertisementID,
		"Не найдено подходящее объявление")

	t.Log("Шаг 1. Добавляем объявление в избранное")

	favoritesManager.AddFavorite(t, token, advertisementID)

	t.Cleanup(func() {
		t.Log("Очистка: удаляем объявление из избранного")
		favoritesManager.DeleteFavorite(t, token, advertisementID)
	})

	t.Log("Шаг 2. Повторно добавляем то же объявление")

	resp := favoritesManager.AddFavoriteExpectError(
		t,
		token,
		advertisementID,
		http.StatusConflict,
	)

	require.Equal(t, "Advertisement is already in favorites",
		resp.Message, "Некорректное сообщение об ошибке")

	t.Logf("Сообщение об ошибке: %s", resp.Message)
}
