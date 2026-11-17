package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	favoritesClient "api-tests-template/internal/client/http/favorites"
	advertisementModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/favorites/models"
)

func GetFavoritesList(t *testing.T, token, search string) models.FavoritesList {
	raw := favoritesClient.HttpGetFavorites(t, token, search, http.StatusOK)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.FavoritesList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func AddFavorite(t *testing.T, token, id string) {
	raw := favoritesClient.HttpPostFavorite(t, token, id, http.StatusCreated)
	defer raw.Body.Close()
}

func DeleteFavorite(t *testing.T, token, id string) {
	raw := favoritesClient.HttpDeleteFavorite(t, token, id, http.StatusNoContent)
	defer raw.Body.Close()
}

func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) advertisementModels.ErrorResponse {
	raw := favoritesClient.HttpGetFavorites(t, token, "", expectedStatus)
	return parseFavoritesError(t, raw)
}

func AddFavoriteExpectError(t *testing.T, token, id string, expectedStatus int) advertisementModels.ErrorResponse {
	raw := favoritesClient.HttpPostFavorite(t, token, id, expectedStatus)
	return parseFavoritesError(t, raw)
}

func DeleteFavoriteExpectError(t *testing.T, token, id string, expectedStatus int) advertisementModels.ErrorResponse {
	raw := favoritesClient.HttpDeleteFavorite(t, token, id, expectedStatus)
	return parseFavoritesError(t, raw)
}

func parseFavoritesError(t *testing.T, raw *http.Response) advertisementModels.ErrorResponse {
	t.Helper()
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := advertisementModels.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
