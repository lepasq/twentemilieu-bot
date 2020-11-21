package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"twentemilieu-whatsapp-bot/config"
)

const (
	calendarUrl = "https://twentemilieuapi.ximmio.com/api/GetCalendar"
	addressUrl  = "https://twentemilieuapi.ximmio.com/api/FetchAdress"
	companyCode = "8d97bb56-5afd-4cbc-a651-b4f7314264b4"
)

/*
 * Struct for making an API request to get one's calendar.
 */
type Calendar struct {
	CompanyCode     string
	UniqueAddressID string
	StartDate       string
	EndDate         string
}

func createCalendar(id string) *Calendar {
	calendar := Calendar{
		CompanyCode:     companyCode,
		UniqueAddressID: id,
		StartDate:       getDateTomorrow(),
		EndDate:         getDateTomorrow(),
	}
	return &calendar
}

func getDateTomorrow() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02") // For testing, use the line below
	//return "2020-10-26" // For real date, use the line above
}

func getAddressRequest(address *config.Address) (*string, error) {
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
	return nil, errors.New("Couldn't get Address.")
}

func getCalendar(address *config.Address) (*string, error) {
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
	return nil, errors.New("Couldn't get deposits for today.")
}

func GetMessage(config *config.Config) (*string, error) {
	pickup, ok := getCalendar(config.Api)
	if ok != nil {
		return nil, errors.New("There are no deposits for today.")
	}
	pickupMessage := fmt.Sprint("Tomorrow, on the ", getDateTomorrow(), " you need to deposit a ", *pickup, " container.")
	return &pickupMessage, nil
}
