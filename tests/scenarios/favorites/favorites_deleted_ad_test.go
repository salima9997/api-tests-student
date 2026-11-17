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

func TestFavoriteRemovedAfterAdvertisementDeletion(t *testing.T) {
	base.SetupSuite()
	t.Cleanup(base.TearDownSuite)

	base.Precondition("Авторизуем двух пользователей")

	login2 := os.Getenv("TEST_LOGIN_2")
	password2 := os.Getenv("TEST_PASSWORD_2")

	require.NotEmpty(t, login2, "Не задан TEST_LOGIN_2")
	require.NotEmpty(t, password2, "Не задан TEST_PASSWORD_2")

	owner := auth.Login(
		t,
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	viewer := auth.Login(t, login2, password2)

	require.NotEqual(t, owner.User.Id, viewer.User.Id,
		"Для теста необходимы два разных пользователя")

	t.Log("Шаг 1. Первый пользователь создаёт объявление")

	photoPath := filepath.Join(
		"..", "advertisementCreation", "testdata", "photo.png",
	)

	created := advertisements.CreateAdvertisement(
		t,
		owner.Token,
		fmt.Sprintf("QA Favorite Delete %d", time.Now().UnixNano()),
		"Тест удаления объявления из избранного",
		100,
		1,
		photoPath,
	)

	require.NotEmpty(t, created.ID, "Не получен ID объявления")

	adDeleted := false
	favoriteAdded := false

	t.Cleanup(func() {
		if !adDeleted {
			if favoriteAdded {
				favoritesManager.DeleteFavorite(t, viewer.Token, created.ID)
			}
			advertisements.DeleteAdvertisement(t, owner.Token, created.ID)
		}
	})

	t.Log("Шаг 2. Второй пользователь добавляет объявление в избранное")

	favoritesManager.AddFavorite(t, viewer.Token, created.ID)
	favoriteAdded = true

	list := favoritesManager.GetFavoritesList(t, viewer.Token, "")

	found := false
	for _, item := range list.Items {
		if item.Advertisement.ID == created.ID {
			found = true
			break
		}
	}

	require.True(t, found,
		"Объявление не появилось в избранном второго пользователя")

	t.Log("Шаг 3. Первый пользователь удаляет объявление")

	advertisements.DeleteAdvertisement(t, owner.Token, created.ID)
	adDeleted = true

	t.Log("Шаг 4. Проверяем избранное второго пользователя")

	list = favoritesManager.GetFavoritesList(t, viewer.Token, "")

	for _, item := range list.Items {
		require.NotEqual(t, created.ID, item.Advertisement.ID,
			"Удалённое объявление осталось в избранном")
	}

	t.Log("Удалённое объявление отсутствует в избранном")
}
