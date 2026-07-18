package ratinghandlers

type ErrorResponse struct {
	Err string `json:"err"`
}

type RatingPlayerResponse struct {
	Login        string `json:"login"`
	Rank         uint32 `json:"rank"`
	PreviousRank uint32 `json:"previousRank"`
	Score        uint64 `json:"score"`
}

type RatingResponse struct {
	Top  []RatingPlayerResponse `json:"top"`
	Near []RatingPlayerResponse `json:"near"`
}

type GetRatingsResponse struct {
	User  RatingResponse `json:"user"`
	Fleet RatingResponse `json:"fleet"`
}
