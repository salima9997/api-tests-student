package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

type FavoriteItem struct {
	AddedAt       string                            `json:"added_at"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

type FavoritesList struct {
	Items []FavoriteItem `json:"items"`
}
