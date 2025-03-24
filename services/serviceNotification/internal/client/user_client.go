package clientsUser

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"serviceNotification/internal/entity"
)

type HTTPUserServiseClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPUserServiseClient(baseURL string) *HTTPUserServiseClient {
	return &HTTPUserServiseClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *HTTPUserServiseClient) GetUserByID(userID int64) (*entity.User, error) {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	// Создаем HTTP-запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Отправляем запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Обрабатываем статус ответа
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("User with ID %d not found", userID)
		return nil, nil // Возвращаем nil без ошибки
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Декодируем JSON-ответ
	var user entity.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &user, nil
}
