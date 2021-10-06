package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Price struct {
	Price string `json:"price"`
	Name  string `json:"name"`
}

func main() {
	var lowestPrice int64 = 0
	lowestHotel := ""
	allHotels := make([]Price, 0)

	collector := colly.NewCollector(
		colly.AllowedDomains("expedia.co.in", "www.expedia.co.in"),
	)
	collector.OnHTML(".results li", func(element *colly.HTMLElement) {

		//		hotelPrice, err := strconv.Atoi(element.ChildAttr(".loyalty-display-price", "span"))
		//		if err != nil {
		//			log.Println("Could not get id")
		//		}
		hotelPrice := element.ChildText(".uitk-cell .loyalty-display-price .all-cell-shrink")
		hotelPrice = strings.Replace(hotelPrice, "Rs", "", -1)
		hotelPrice = strings.Replace(hotelPrice, ",", "", -1)

		//hotelName := element.ChildText("h3")
		hotelName := element.ChildText("h3")
		//Hack to fix the name issue
		hotelName = hotelName[:len(hotelName)-len(hotelName)/2]

		price := Price{
			Price: hotelPrice,
			Name:  hotelName,
		}
		if len(hotelPrice) > 0 {
			allHotels = append(allHotels, price)
		}
		//fmt.Println("Name", hotelName)
		if i, err := strconv.ParseInt(hotelPrice, 10, 0); err == nil {
			if lowestPrice == 0 || i < lowestPrice {
				lowestPrice = i
				lowestHotel = hotelName
			}
		}

	})

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	url := "https://www.expedia.co.in/Hotel-Search?adults=2&d1=2021-12-31&d2=2022-01-01&destination=Bengaluru%20%28and%20vicinity%29%2C%20Karnataka%2C%20India&endDate=2022-01-01&latLong=12.999551%2C77.587685&regionId=6053307&rooms=1&semdtl=&sort=RECOMMENDED&startDate=2021-12-31&theme=&useRewards=false&userIntent="
	collector.Visit(url)
	fmt.Println("Lowest Hotel : ", lowestHotel, " Price : ", lowestPrice)
	writeJSON(allHotels)

}

func writeJSON(data []Price) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("hotel.json", file, 0644)
}
