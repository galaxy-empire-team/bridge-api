package ratinghandlers

import "github.com/galaxy-empire-team/bridge-api/internal/models"

func toGetRatingsResponse(ratings models.Ratings) GetRatingsResponse {
	return GetRatingsResponse{
		User:  toRatingResponse(ratings.User, true),
		Fleet: toRatingResponse(ratings.Fleet, false),
	}
}

func toRatingResponse(rating models.Rating, isUserRating bool) RatingResponse {
	return RatingResponse{
		Top:  toRatingPlayersResponse(rating.Top, isUserRating),
		Near: toRatingPlayersResponse(rating.Near, isUserRating),
	}
}

func toRatingPlayersResponse(players []models.RatingPlayer, isUserRating bool) []RatingPlayerResponse {
	result := make([]RatingPlayerResponse, 0, len(players))
	for _, player := range players {
		var score uint64
		if isUserRating {
			score = player.SpentResources
		} else {
			score = player.FleetPower
		}

		result = append(result, RatingPlayerResponse{
			Login:        player.Login,
			Rank:         player.Rank,
			PreviousRank: player.PreviousRank,
			Score:        score,
		})
	}

	return result
}
