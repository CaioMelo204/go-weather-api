package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

type apiConfig struct {
	ApiKey string `json:"api_key"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig() (*apiConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	key := os.Getenv("API_KEY")

	return &apiConfig{ApiKey: key}, nil
}

func main() {
	http.HandleFunc("/", hello)

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Weather APP!"))
}

func query(city string) (weatherData, error) {
	config, err := loadApiConfig()
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=" + config.ApiKey)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	return d, nil
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
