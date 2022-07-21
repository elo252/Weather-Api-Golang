package main

import (
	"encoding/json"
	
	"net/http"
	
	
	"html/template"
	"github.com/joho/godotenv"
	"os"
	"log"
)



type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Coord struct {
		Long float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}`json:"coord"`

}

func goDotEnvVariable(key string) string {

	
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }





func query(city string) (weatherData, error) {


	key := goDotEnvVariable("API_KEY")

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + key + "&q=" + city)
	
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

func WeatherApiData (w http.ResponseWriter, r *http.Request) {

	city := r.FormValue("city")
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	t, err := template.ParseFiles("static/Result.html")
	t.Execute(w, data)


}


func main() {
	
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileServer)

	

	http.HandleFunc("/weather",WeatherApiData)

	http.ListenAndServe(":8080", nil)
	
}