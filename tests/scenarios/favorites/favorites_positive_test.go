package favorites

import (
	"fmt"
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

func TestFavoritesAddAndRemove(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизация пользователя")

	loginData := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
	token := loginData.Token

	require.NotEmpty(t, loginData.User.Id, "Не получен ID пользователя")

	t.Log("Подготовка: ищем чужое объявление, которого нет в избранном")

	advertisementID := createForeignTestAdvertisement(
		t, loginData.User.Id,
	)

	require.NotEmpty(t, advertisementID,
		"Не найдено чужое объявление, отсутствующее в избранном")

	t.Log("Шаг 1. Добавляем объявление в избранное")

	favoritesManager.AddFavorite(t, token, advertisementID)

	needsCleanup := true

	t.Cleanup(func() {
		if needsCleanup {
			t.Log("Очистка: удаляем объявление из избранного")
			favoritesManager.DeleteFavorite(t, token, advertisementID)
		}
	})

	t.Log("Шаг 2. Проверяем наличие объявления в избранном")

	list := favoritesManager.GetFavoritesList(t, token, "")

	found := false

	for _, item := range list.Items {
		if item.Advertisement.ID == advertisementID {
			found = true
			break
		}
	}

	require.True(t, found,
		"Добавленное объявление отсутствует в избранном")

	t.Log("Объявление успешно добавлено в избранное")

	t.Log("Шаг 3. Удаляем объявление из избранного")

	favoritesManager.DeleteFavorite(t, token, advertisementID)
	needsCleanup = false

	t.Log("Шаг 4. Проверяем, что объявление исчезло из избранного")

	list = favoritesManager.GetFavoritesList(t, token, "")

	for _, item := range list.Items {
		require.NotEqual(t, advertisementID, item.Advertisement.ID,
			"Удалённое объявление осталось в избранном")
	}

	t.Log("Объявление успешно удалено из избранного")
}

func createForeignTestAdvertisement(t *testing.T, viewerUserID string) string {
	t.Helper()

	login := os.Getenv("TEST_LOGIN_2")
	password := os.Getenv("TEST_PASSWORD_2")

	require.NotEmpty(t, login, "Не задан TEST_LOGIN_2")
	require.NotEmpty(t, password, "Не задан TEST_PASSWORD_2")

	owner := auth.Login(t, login, password)

	require.NotEqual(t, viewerUserID, owner.User.Id,
		"Для создания чужого объявления нужны два разных пользователя")

	photoPath := filepath.Join(
		"..", "advertisementCreation", "testdata", "photo.png",
	)

	created := advertisements.CreateAdvertisement(
		t,
		owner.Token,
		fmt.Sprintf("QA Favorite %d", time.Now().UnixNano()),
		"Объявление для тестирования избранного",
		100,
		1,
		photoPath,
	)

	require.NotEmpty(t, created.ID, "Не получен ID объявления")

	t.Cleanup(func() {
		t.Log("Очистка: удаляем тестовое объявление")
		advertisements.DeleteAdvertisement(t, owner.Token, created.ID)
	})

	return created.ID
}
