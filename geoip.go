package geoip

import (
	"context"
	"encoding/json"
	"strings"
	"fmt"
	"io"
	"net/http")


// holds the config to be pass to the plugin
type Config struct {
}

func CreateConfig() *Config {
	return &Config{}
}

type UIDdemo struct {
	next http.Handler
	name string
}

type GeoIP struct {
	ip			string	 `json:"ip"`
	countryCode	string	 `json:"country_code"`
	countryName	string	 `json:"country_name"`
	regionCode	string	 `json:"region_code"`
	regionName	string	 `json:"region_name"`
	city		string	 `json:"city"`
	zipCode		string	 `json:"zip_code"`
	timeZone	string	 `json:"time_zone"`
	latitude	float64	 `json:"latitude"`
	longitude	float64	 `json:"longitude"`
	metroCode	int		 `json:"metro_code"`
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &UIDdemo{
		next: next,
		name: name,
	}, nil
}

func (u *UIDdemo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ip := req.Header.Get("X-Forwarded-For")

	geoLoc, err := MakeRequest(ip)

	req.Header.Set("X-Country-Code", geoLoc.countryCode)
	req.Header.Set("X-Country-Name", geoLoc.countryName)

	u.next.ServeHTTP(res, req)
}

func MakeRequest(ip string) (*GeoIP, error) {
	resp, err := http.Get("https://freegeoip.live/json/"+ ip)

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	geoLocation := GeoIP{}
	err = json.Unmarshal(body, &geoLocation)

 	return &geoLocation, nil
}