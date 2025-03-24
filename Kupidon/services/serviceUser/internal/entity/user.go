package entity

// @Description User structure
// @Param id path int true "User ID"
// @Param name formData string true "User's name"
// @Param age formData int true "User's age"
// @Param city formData string true "City of the user"
// @Param gender formData string true "Gender of the user"
// @Param description formData string true "Description of the user"
// @Param telegram_id formData int64 true "Telegram ID of the user"
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
