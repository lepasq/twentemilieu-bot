package main

import (
	"bytes"
	"encoding/json"
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
	addressId
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
	var address *Address = createAddress()
	getCalendar(address)
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

func getTimeNow() {
	currentTime := time.Now()
	fmt.Println("Today is: ", currentTime.Format("2006-01-02"))
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

	for _,item:=range m["dataList"].([]interface{}) {
		a := fmt.Sprintf("%v", item.(map[string]interface{})["UniqueId"])
		return &a, nil
	}
	return nil, nil
}


func getCalendar(address *Address) {
	id, ok := getAddressRequest(address)
	if ok != nil {
		fmt.Println("Error: ", id)
	}
	addressId := *id
	fmt.Println(addressId)
}

