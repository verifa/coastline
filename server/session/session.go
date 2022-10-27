package session

type UserClaims struct {
	Sub     string   `json:"sub"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Groups  []string `json:"groups"`
	Picture string   `json:"picture"`
}
