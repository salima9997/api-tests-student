package advertisementCreation

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	advertisements "api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	base "api-tests-template/tests"
)

func TestCreateAdvertisementPositive(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")
	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
	token := loginData.Token

	title := fmt.Sprintf("QA Test %d", time.Now().UnixNano())
	description := "Тестовое объявление"
	price := 1999.0
	quantity := 2
	photoPath := filepath.Join("testdata", "photo.png")

	t.Log("Шаг 1. Создаём объявление со всеми параметрами")

	created := advertisements.CreateAdvertisement(
		t, token, title, description, price, quantity, photoPath,
	)

	require.NotEmpty(t, created.ID, "Не получен ID созданного объявления")

	t.Cleanup(func() {
		t.Log("Очистка: удаляем тестовое объявление")
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	t.Log("Шаг 2. Получаем объявление по ID")

	actual := advertisements.GetAdvertisement(t, token, created.ID)

	require.Equal(t, created.ID, actual.ID)
	require.Equal(t, title, actual.Title)
	require.Equal(t, description, actual.Description)
	require.Equal(t, price, actual.Price)
	require.Equal(t, quantity, actual.Quantity)

	t.Log("Поля объявления совпадают с переданными")

	t.Log("Шаг 3. Ищем объявление по полному названию")

	list := advertisements.GetAdvertisementsList(t, token, title)

	assertions.ContainsAdvertisementByTitle(t, list.Items, title)

	t.Log("Объявление найдено в поиске")

	t.Log("Шаг 4. Проверяем загруженную фотографию")

	require.NotEmpty(t, created.Photos,
		"Фотография отсутствует в ответе создания")

	photos := advertisements.GetAdvertisementPhotos(
		t, token, created.ID,
	)

	require.NotEmpty(t, photos,
		"GET /advertisements/{id}/photos вернул пустой список")

	require.Equal(t, created.Photos[0].ID, photos[0].ID,
		"ID загруженной фотографии не совпадает")

	require.NotEmpty(t, photos[0].URL,
		"У фотографии отсутствует URL")

	t.Log("Фотография успешно загружена и доступна через API")
}
