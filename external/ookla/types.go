package ookla

import (
	"encoding/xml"
	"strconv"
)

type ClientConfig struct {
	XMLName    xml.Name `xml:"settings"`
	Client     Client   `xml:"client"`
	LicenseKey string   `xml:"licensekey"`
}
type Client struct {
	XMLName xml.Name `xml:"client"`
	IP      string   `xml:"ip,attr"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
}

type ServersList []ServerSummary

type ServerSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Host     string `json:"host"`
	Lat      string `json:"lat"`
	Lon      string `json:"lon"`
	Distance uint64 `json:"distance"`
	Country  string `json:"country"`
	CC       string `json:"CC"`
	Sponsor  string `json:"sponsor"`
}

func (ss ServerSummary) GetID() uint64 {
	ui, err := strconv.ParseUint(ss.ID, 10, 64)
	_ = err // suppress because it can not happen
	return ui
}

func (ss ServerSummary) GetLat() float64 {
	f, err := strconv.ParseFloat(ss.Lat, 64)
	_ = err // suppress because it can not happen
	return f
}

func (ss ServerSummary) GetLon() float64 {
	f, err := strconv.ParseFloat(ss.Lon, 64)
	_ = err // suppress because it can not happen
	return f
}
