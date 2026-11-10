package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"

	"github.com/stretchr/testify/require"
)

func ContainsAdvertisementByTitle(
	t *testing.T,
	items []models.Advertisement,
	title string,
) {
	t.Helper()

	titles := make([]string, 0, len(items))

	for _, item := range items {
		titles = append(titles, item.Title)
	}

	require.Contains(
		t,
		titles,
		title,
		"Искомое объявление отсутствует среди найденных",
	)
}
