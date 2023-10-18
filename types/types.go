package types

import "podcast/database"

type Paginator = database.Paginator

type Sorter = database.Sorter

type StripeBalances struct {
	Available        int64 `json:"available"`
	Pending          int64 `json:"pending"`
	InstantAvailable int64 `json:"instant_available"`
}
