package entity

type User struct {
	ID          int    `json:"id"`
	TelegramID  int64  `json:"telegram_id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	City        string `json:"city,omitempty"`
	Gender      string `json:"gender,omitempty"`
	Description string `json:"description"`
	Photo       string `json:"photo"`
}
