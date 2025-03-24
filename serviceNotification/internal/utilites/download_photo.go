package utilites

import (
	"io"
	"net/http"
)

func DownloadImageAsBytes(url string) ([]byte, error) {
	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем все байты из ответа
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}
