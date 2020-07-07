package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type GooglePlaces struct {
	HTMLAttributions []interface{} `json:"html_attributions"`
	NextPageToken    string        `json:"next_page_token"`
	Results          []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
			Viewport struct {
				Northeast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"northeast"`
				Southwest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"southwest"`
			} `json:"viewport"`
		} `json:"geometry"`
		Icon         string `json:"icon"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		OpeningHours struct {
			OpenNow     bool          `json:"open_now"`
			WeekdayText []interface{} `json:"weekday_text"`
		} `json:"opening_hours,omitempty"`
		Photos []struct {
			Height           int      `json:"height"`
			HTMLAttributions []string `json:"html_attributions"`
			PhotoReference   string   `json:"photo_reference"`
			Width            int      `json:"width"`
		} `json:"photos,omitempty"`
		PlaceID   string   `json:"place_id"`
		Reference string   `json:"reference"`
		Scope     string   `json:"scope"`
		Types     []string `json:"types"`
		Vicinity  string   `json:"vicinity"`
		Rating    float64  `json:"rating,omitempty"`
	} `json:"results"`
	Status string `json:"status"`
}

func searchPlaces(page string) {
	apiKey := ""                    // Enter your API key here
	keyword := "mosque"             // Enter the type of location you want
	latLong := "5.636009,-0.234358" // Your position
	pageToken := page
	var buffer bytes.Buffer

	buffer.WriteString("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=")
	buffer.WriteString(latLong)
	buffer.WriteString("&radius=50000&keyword=")
	buffer.WriteString(keyword)
	buffer.WriteString("&key=")
	buffer.WriteString(apiKey)
	buffer.WriteString("&pagetoken=")
	buffer.WriteString(pageToken)

	query := buffer.String()

	// 1. Open the file
	csvfile, err := os.OpenFile("marks.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close()
	// 2. Initialize the writer
	writer := csv.NewWriter(csvfile)

	defer writer.Flush()

	// PRINT CURRENT SEARCH
	println("query is ", query)
	println("\n")

	// SEND REQUEST WITH QUERY
	resp, err := http.Get(query)
	if err != nil {
		log.Fatal(err)
	}
	// CLOSE THE PRECLOSER THATS RETURNED WITH HTTP RESPONSE
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	res := GooglePlaces{}
	json.Unmarshal([]byte(body), &res)

	var tache []string
	var listings bytes.Buffer
	for i := 0; i < len(res.Results); i++ {
		listings.WriteString(strconv.Itoa(i + 1))
		listings.WriteString("\nName: ")
		listings.WriteString(res.Results[i].Name)
		listings.WriteString("\nAddress: ")
		listings.WriteString(res.Results[i].Vicinity)
		listings.WriteString("\nPlace ID: ")
		listings.WriteString(res.Results[i].PlaceID)
		listings.WriteString("\nLatitude: ")
		listings.WriteString(strconv.FormatFloat(res.Results[i].Geometry.Location.Lat, 'E', -1, 64))
		listings.WriteString("\nLongitude: ")
		listings.WriteString(strconv.FormatFloat(res.Results[i].Geometry.Location.Lng, 'E', -1, 64))
		listings.WriteString("\n---------------------------------------------\n\n")

	}
	listings.WriteString("\npagetoken is now:\n")
	listings.WriteString(res.NextPageToken)
	tache = append(tache, listings.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(listings.String())
	fmt.Printf("\n\n\n")

	// 3. Write all the records
	err = writer.Write(tache)
	if err != nil {
		fmt.Println("An error encountered ::", err)
	}
	// LOOP BACK THROUGH FUNCTION

	if res.NextPageToken != "" {
		time.Sleep(5000 * time.Millisecond)
		searchPlaces(res.NextPageToken)
	} else {
		fmt.Println("No more pagetoken, we're done.")
	}

}

func main() {
	searchPlaces("")
}
