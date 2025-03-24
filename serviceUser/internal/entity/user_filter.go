package entity

type UserFilter struct {
	MinAge *int   `json:"min_age,omitempty"`
	MaxAge *int   `json:"max_age,omitempty"`
	City   string `json:"city,omitempty"`
	Gender string `json:"gender,omitempty"`
}
