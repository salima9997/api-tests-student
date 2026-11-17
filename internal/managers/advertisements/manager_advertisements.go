package advertisements

import (
	"encoding/json"
	"io"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func CreateAdvertisement(t *testing.T, token, title, description string, price float64, quantity int, photoPath string) models.Advertisement {
	raw := advertisements.HttpPostAdvertisement(t, token, title, description, price, quantity, photoPath)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func GetAdvertisement(t *testing.T, token, id string) models.Advertisement {
	raw := advertisements.HttpGetAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.Advertisement{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()
}

func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := []models.Photo{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}

func CreateAdvertisementExpectError(t *testing.T, token, title, description string, price float64, quantity int, photoPath string, expectedStatus int) models.ErrorResponse {
	raw := advertisements.HttpPostAdvertisementWithStatus(
		t, token, title, description, price, quantity,
		photoPath, expectedStatus,
	)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	resp := models.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &resp))

	return resp
}
