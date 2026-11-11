package handler

type CreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Position string `json:"position"`
}
type UpdateStageRequest struct {
	Stage string `json:"stage"`
}
