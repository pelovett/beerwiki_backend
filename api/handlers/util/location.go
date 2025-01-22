package util

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type City struct {
	Name       string `json:"name"`
	Country    string `json:"country"`
	ISO2       string `json:"iso2"`
	ISO3       string `json:"iso3"`
	AdminName  string `json:"admin_name"`
	Population int    `json:"population"`
	Flag       string `json:"flag"`
}

func GetCityList() (*[]City, error) {
	file, err := os.Open("assets/worldcities.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var cities []City
	for _, record := range records[1:] {
		var pop float64
		if record[5] != "" {
			pop, err = strconv.ParseFloat(record[5], 64)
			if err != nil {
				return nil, err
			}
		} else {
			pop = 1.0
		}
		cities = append(cities, City{
			Name:       record[0],
			Country:    record[1],
			ISO2:       record[2],
			ISO3:       record[3],
			AdminName:  record[4],
			Population: int(pop),
			Flag:       record[6],
		})
	}
	return &cities, nil
}

type SearchItem struct {
	city  *City
	score int
}

func FindCities(query string, cityList *[]City) ([]City, error) {
	// Hardcoded to get top 5
	numResults := 5
	var bestMatch []SearchItem
	lastPlaceDistance := 10000
	for _, city := range *cityList {
		matchScore := fuzzy.RankMatchNormalizedFold(query, city.Name)
		if matchScore == -1 {
			continue
		}

		if len(bestMatch) < numResults || lastPlaceDistance > matchScore {
			bestMatch = append(bestMatch, SearchItem{city: &city, score: matchScore})
			slices.SortFunc(bestMatch, func(a, b SearchItem) int {
				if a.score == b.score {
					// Favor bigger cities
					return b.city.Population - a.city.Population
				}
				return a.score - b.score
			})
			bestMatch = bestMatch[:min(len(bestMatch), numResults)]
			lastPlaceDistance = bestMatch[len(bestMatch)-1].score
		}
	}
	var results []City
	for _, item := range bestMatch {
		results = append(results, *item.city)
	}
	return results, nil
}

func SearchCities(cityList *[]City) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Request.URL.Query()
		queryList, queryPresent := params["q"]
		if !queryPresent || len(queryList) == 0 {
			c.Status(http.StatusBadRequest)
			return
		}
		query := queryList[0]

		cities, err := FindCities(query, cityList)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Printf("Error while searching cities: %v", err)
			return
		}

		if len(cities) > 0 {
			c.JSON(http.StatusOK, map[string]any{
				"results": cities,
			})
		} else {
			c.Status(http.StatusNotFound)
		}
	}

}
