package pullbuddy

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type schedulerAPI interface {
	list() []image
	schedule(id string)
}

func newHandler(api schedulerAPI) http.Handler {
	router := chi.NewRouter()
	router.Get("list", listHandler(api))
	router.Post("schedule", scheduleHandler(api))
	return router
}

func listHandler(api schedulerAPI) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		images := api.list()
		lr := listResponse{
			Images: make([]listResponseImage, len(images)),
		}
		for i := range images {
			lr.Images[i] = listResponseImage{
				ID:     string(images[i].id),
				Status: string(images[i].status),
			}
			if images[i].err != nil {
				lr.Images[i].Error = images[i].err.Error()
			}
		}
		encoder := json.NewEncoder(response)
		encoder.Encode(lr)
	}
}

type listResponse struct {
	Images []listResponseImage `json:"images"`
}

type listResponseImage struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func scheduleHandler(api schedulerAPI) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sr := scheduleRequest{}
		err := json.NewDecoder(request.Body).Decode(&sr)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		api.schedule(sr.ImageID)
	}
}

type scheduleRequest struct {
	ImageID string `json:"image_id"`
}
