package middlewears

import (
	"PlantCare/controllers"
	"PlantCare/utils"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gorilla/mux"
)

func PotMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		if !ok {
			utils.JsonError(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}

		potId := mux.Vars(r)["potId"]

		cropPot, err := controllers.FindCropPotById(potId)
		if err != nil {
			utils.JsonError(w, "Crop pot not found!", http.StatusNotFound)
			return
		}

		if cropPot.ClerkUserID == nil || *cropPot.ClerkUserID != claims.Subject {
			utils.JsonError(w, "Unauthorized! You do not own this pot.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
