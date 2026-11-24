package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpPostAdvertisement(t *testing.T, token, title, description string, price float64, quantity int, photoPath string) *http.Response {
	return HttpPostAdvertisementWithStatus(
		t, token, title, description, price, quantity,
		photoPath, http.StatusCreated,
	)
}

func HttpPostAdvertisementWithStatus(t *testing.T, token, title, description string, price float64, quantity int, photoPath string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.FormatFloat(price, 'f', -1, 64)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		Status(expectedStatus).
		End().Response
}

func HttpGetAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		Status(http.StatusNoContent).
		End().Response
}

func HttpGetAdvertisementPhotos(t *testing.T, token, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

func HttpPostAdvertisementRawWithStatus(
	t *testing.T,
	token, title, description, price, quantity string,
	photoPaths []string,
	expectedStatus int,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", price).
		MultipartFormData("quantity", quantity).
		MultipartFile("photos", photoPaths...).
		Expect(t).
		Status(expectedStatus).
		End().Response
}
