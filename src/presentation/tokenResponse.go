package presentation

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
}

func NewTokenResponse(accessToken string) Token {
	return Token{AccessToken: accessToken}
}
