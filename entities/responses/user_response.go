package responses

type UserResponse struct {
	Id          string `json:"id,omitempty"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserPublicResponse struct {
	Username string `json:"username"`
}
