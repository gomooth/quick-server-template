package auth

type tokenEntity struct {
	AccessToken string `json:"accessToken"`
	ExpireTime  int64  `json:"expire"`

	Profile *profileEntity `json:"profile,omitempty"`
}

type profileEntity struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	AvatarURL   string   `json:"avatar"`
	CurrentRole string   `json:"currentRole"`
	Roles       []string `json:"roles"`
}
