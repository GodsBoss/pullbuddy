package pullbuddy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type schedulerAPI interface {
	list() []image
	schedule(id string) error
}

func newHandler(api schedulerAPI) http.Handler {
	router := chi.NewRouter()
	router.Get("/list", listHandler(api))
	router.Post("/schedule", scheduleHandler(api))
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

func scheduleHandler(api schedulerAPI) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		sr := scheduleRequest{}
		resp := scheduleResponse{}
		err := json.NewDecoder(request.Body).Decode(&sr)
		if err != nil {
			resp.Message = fmt.Sprintf("could not parse request body: %s", err)
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(resp)
			return
		}
		err = api.schedule(sr.ImageID)
		if isValidationFailedError(err) {
			resp.Message = fmt.Sprintf("validation failed: %s", err)
			response.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(response).Encode(resp)
			return
		}
		if err != nil {
			resp.Message = "a server error occured"
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(resp)
			return
		}
		resp.Message = "OK"
		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(resp)
	}
}
