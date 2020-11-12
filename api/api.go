package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	calendarUrl = "https://twentemilieuapi.ximmio.com/api/GetCalendar"
	addressUrl  = "https://twentemilieuapi.ximmio.com/api/FetchAdress"
	companyCode = "8d97bb56-5afd-4cbc-a651-b4f7314264b4"
	postCode    = "7545MR"
	houseNumber = 4
	houseLetter = "a"
)

/*
 * Struct for making an API request to get one's address.
 */
type Address struct {
	CompanyCode string
	PostCode    string
	HouseNumber int
	HouseLetter string
}

/*
 * Struct for making an API request to get one's calendar.
 */
type Calendar struct {
	CompanyCode     string
	UniqueAddressID string
	StartDate       string
	EndDate         string
}

func main() {
	message, ok := getMessage(createAddress())
	if ok != nil {
		fmt.Println("Something went wrong.")
		return
	}
	fmt.Println(*message)
}

func createAddress() *Address {
	address := Address{
		CompanyCode: companyCode,
		PostCode:    postCode,
		HouseNumber: houseNumber,
		HouseLetter: houseLetter,
	}
	return &address
}

func createCalendar(id string) *Calendar {
	calendar := Calendar{
		CompanyCode:     companyCode,
		UniqueAddressID: id,
		StartDate:       getDateToday(),
		EndDate:         getDateToday(),
	}
	return &calendar
}

func getDateToday() string {
	return time.Now().Format("2006-01-02") // For testing, use the line below
	//return "2020-10-26" // For real date, use the line above
}

func getAddressRequest(address *Address) (*string, error) {
	bytearray, err := json.Marshal(*address)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(addressUrl, "application/json", bytes.NewBuffer(bytearray))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	for _, item := range m["dataList"].([]interface{}) {
		a := fmt.Sprintf("%v", item.(map[string]interface{})["UniqueId"])
		return &a, nil
	}
	return nil, errors.New("Couldn't get Calendar.")
}

func getCalendar(address *Address) (*string, error) {
	id, ok := getAddressRequest(address)
	if ok != nil {
		return nil, errors.New("Couldn't get Address.")
	}
	calendar := createCalendar(*id)

	bytearray, err := json.Marshal(calendar)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(calendarUrl, "application/json", bytes.NewBuffer(bytearray))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}

	for _, item := range m["dataList"].([]interface{}) {
		pickupType := fmt.Sprintf("%v", item.(map[string]interface{})["_pickupTypeText"])
		return &pickupType, nil
	}
	return nil, errors.New("Couldn't get Calendar.")
}

func getMessage(address *Address) (*string, error) {
	pickup, ok := getCalendar(address)
	if ok != nil {
		return nil, errors.New("Couldn't get Calendar.")
	}
	pickupMessage := fmt.Sprint("Today, on the ", getDateToday(), " you need to deposit a ", *pickup, " container.")
	return &pickupMessage, nil
}
