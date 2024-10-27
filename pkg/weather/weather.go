package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Структура для ответа от OpenWeatherMap
type WeatherResponse struct {
	Name string `json:"name"` // Название города
	Main struct {
		Temp     float64 `json:"temp"`     // Температура
		Humidity int     `json:"humidity"` // Влажность
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"` // Описание погоды
	} `json:"weather"`
}

// Client — структура для хранения API-ключа и выполнения запросов к OpenWeatherMap
type Client struct {
	APIKey string
}

// NewClient создает новый экземпляр клиента для OpenWeatherMap
func New(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

// GetCurrentWeather получает текущую погоду для указанного города
func (c *Client) GetCurrentWeather(city string) (*WeatherResponse, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric&lang=ru", encodedCity, c.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к OpenWeatherMap: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить данные о погоде, статус: %s", resp.Status)
	}

	var weatherData WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &weatherData, nil
}
