package dto

// @Description A simple health check response.
type PongResponse struct {
	Message string `json:"message" example:"pong"`
}
