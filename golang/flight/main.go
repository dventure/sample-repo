package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type Price struct {
	Price         string `json:"price"`
	Flight        string `json:"flight"`
	Arrival       string `json:"arrival"`
	Departure     string `json:"departure"`
	ArrTime       string `json:"arrTime"`
	DepTimeFlight string `json:"depTime"`
	Duration      string `json:"duration"`
}

func main() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	const (
		chromeDriverPath = "/Users/gijoy/Downloads/selenium/chromedriver"
		port             = 9515
	)
	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(), // Start an X frame buffer for the browser to run in.
		selenium.ChromeDriver(chromeDriverPath),
		selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	//selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)

	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		panic(err)
	}
	defer wd.Quit()
	url := "https://www.expedia.co.in/Flights-Search?leg1=from%3AChennai%20%28MAA-Chennai%20Intl.%29%2Cto%3ABengaluru%20%28BLR-Kempegowda%20Intl.%29%2Cdeparture%3A24%2F11%2F2021TANYT&mode=search&options=carrier%3A%2A%2Ccabinclass%3A%2Cmaxhops%3A1%2Cnopenalty%3AN&pageId=0&passengers=adults%3A1%2Cchildren%3A0%2Cinfantinlap%3AN&trip=oneway"

	if err := wd.Get(url); err != nil {
		panic(err)
	}

	time.Sleep(time.Millisecond * 20000)

	// Wait for the program to finish running and get the output.
	outputDiv, err := wd.FindElements(selenium.ByCSSSelector, "span[class=is-visually-hidden]")
	if err != nil {
		panic(err)
	}
	allFlights := make([]Price, 0)
	lowestFlight := ""
	lowestPrice := ""
	for i := 0; i < len(outputDiv); i++ {
		output, err := outputDiv[i].Text()
		if err != nil {
			panic(err)
		}
		if strings.Index(output, "Select and show fare") == 0 {
			str := strings.Replace(output, "Select and show fare information for ", "", 1)
			index := strings.Index(str, ",")
			flightName := str[:index]
			str = str[index+15:]
			//fmt.Println(i, str)
			fmt.Println(flightName)

			//str = strings.Replace(str, "departing at ","",1)
			index = strings.Index(str, "from")
			depTime := str[:index]
			str = str[index+5:]
			fmt.Println(depTime)

			index = strings.Index(str, ",")
			depPlace := str[:index]
			str = str[index+14:]
			fmt.Println(depPlace)

			index = strings.Index(str, " in")
			arrTime := str[:index]
			str = str[index+4:]
			fmt.Println(arrTime)

			index = strings.Index(str, ",")
			arrPlace := str[:index]
			str = str[index+14:]
			fmt.Println(arrPlace)

			index = strings.Index(str, " ")
			flightPrice := str[:index]
			str = str[index+24:]
			fmt.Println(flightPrice)

			index = strings.Index(str, " total travel time")
			dur := str[:index]
			fmt.Println(dur)

			price := Price{
				Price:         flightPrice,
				Flight:        flightName,
				Arrival:       arrPlace,
				Departure:     depPlace,
				ArrTime:       arrTime,
				DepTimeFlight: depTime,
				Duration:      dur,
			}
			if len(flightPrice) > 0 {
				allFlights = append(allFlights, price)
			}
			if len(lowestPrice) == 0 || lowestPrice < flightPrice {
				lowestPrice = flightPrice
				lowestFlight = flightName
			}
		}

	}
	writeJSON(allFlights)
	fmt.Println("*******************************************************")
	fmt.Println(lowestFlight, "is having lowest fare of", lowestPrice)
	fmt.Println("*******************************************************")

}

func writeJSON(data []Price) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("flight.json", file, 0644)
}
