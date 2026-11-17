package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetFavorites(t *testing.T, token, search string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Query("search", search).
		Query("page_size", "50").
		Expect(t).
		Status(expectedStatus).
		End().Response
}

func HttpPostFavorite(t *testing.T, token, id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		JSON(map[string]string{"advertisement_id": id}).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token, id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoritesPath + "/" + id).
		Expect(t).
		Status(expectedStatus).
		End().Response
}
