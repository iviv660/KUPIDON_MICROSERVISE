package clientsUser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"serviceBot/internal/entity"
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

func (c *HTTPUserServiseClient) CreateUser(name, city, gender, description string, age int, telegramID int64, file []byte, filename string) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// --- 1. Создаем JSON-объект ---
	data := map[string]interface{}{
		"name":        name,
		"city":        city,
		"gender":      gender,
		"description": description,
		"age":         age,
		"telegram_id": telegramID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// --- 2. Добавляем JSON в form-data ---
	if err := writer.WriteField("json", string(jsonData)); err != nil {
		return fmt.Errorf("failed to write JSON field: %w", err)
	}

	// --- 3. Добавляем файл ---
	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create file part: %w", err)
	}

	if _, err := io.Copy(filePart, bytes.NewReader(file)); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	// Закрываем writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// --- 4. Создаем HTTP-запрос ---
	url := fmt.Sprintf("%s/users", c.baseURL)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// --- 5. Отправляем запрос ---
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// --- 6. Обрабатываем ответ ---
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Received response: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
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

func (c *HTTPUserServiseClient) Delete(id int64) error {
	url := fmt.Sprintf("%s/users/%d", c.baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *HTTPUserServiseClient) SearchUser(MinAge, MaxAge int, City, Gender string) ([]entity.User, error) {
	url := fmt.Sprintf("%s/users/search?min_age=%d&max_age=%d&city=%s&gender=%s", c.baseURL, MinAge, MaxAge, City, Gender)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var users []entity.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return users, nil
}
