package advertisementCreation

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	advertisements "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementNegative(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	photoPath := filepath.Join("testdata", "photo.png")

	cases := []struct {
		name            string
		token           string
		title           string
		price           float64
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "Пустое название",
			token:           loginData.Token,
			title:           "",
			price:           100,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "title is required",
		},
		{
			name:            "Отрицательная цена",
			token:           loginData.Token,
			title:           "QA Negative Test",
			price:           -100,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "price must be between 0 and 1000000 inclusive",
		},
		{
			name:            "Слишком длинное название",
			token:           loginData.Token,
			title:           strings.Repeat("A", 1000),
			price:           100,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "title must not exceed 50 characters",
		},
		{
			name:            "Без токена",
			token:           "",
			title:           "QA Negative Test",
			price:           100,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "Authorization header required",
		},
		{
			name:            "Неверный токен",
			token:           "invalid-token",
			title:           "QA Negative Test",
			price:           100,
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "Invalid or expired token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Отправляем запрос на создание объявления")

			resp := advertisements.CreateAdvertisementExpectError(
				t,
				tc.token,
				tc.title,
				"Описание тестового объявления",
				tc.price,
				1,
				photoPath,
				tc.expectedStatus,
			)

			require.Contains(t, resp.Message, tc.expectedMessage, "Некорректное сообщение об ошибке")

			t.Logf("Ожидаемый статус: %d", tc.expectedStatus)
			t.Logf("Сообщение об ошибке: %s", resp.Message)
		})
	}
}
