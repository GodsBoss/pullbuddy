package pullbuddy

type listResponse struct {
	Images []listResponseImage `json:"images"`
}

type listResponseImage struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type scheduleRequest struct {
	ImageID string `json:"image_id"`
}

type scheduleResponse struct {
	Message string `json:"messge"`
}
