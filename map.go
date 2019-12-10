package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type Region struct {
	Region int
	Name   string
}

type TrafficLights struct {
	ID          int
	Regin       int
	Description string
	Points      Point
}

func getRegion() (region []Region) {
	temp := &Region{}
	sqlquery := fmt.Sprintf("select region, name from %s", os.Getenv("region_table"))
	rows, _ := GetDB().Raw(sqlquery).Rows()
	for rows.Next() {
		rows.Scan(&temp.Region, &temp.Name)
		region = append(region, *temp)
	}
	return
}

func getTrafficLights() (trLight []TrafficLights) {
	var dgis string
	temp := &TrafficLights{}
	sqlquery := fmt.Sprintf("select region, id, dgis, describ from %s", os.Getenv("gis_table"))
	rows, _ := GetDB().Raw(sqlquery).Rows()
	for rows.Next() {
		rows.Scan(&temp.Regin, &temp.ID, &dgis, &temp.Description)
		temp.Points.StrToFloat(dgis)
		trLight = append(trLight, *temp)
	}
	return

}

func makeBmp(Points Point, filepath string) (err error) {
	url := fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%3.15f,%3.15f&z=19&l=map&size=450,450", Points.Y, Points.X)
	// don't worry about errors
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filepath + "map.bmp")
	if err != nil {
		return err
	}
	defer file.Close()

	file1, err := os.Create(filepath + "cross.svg")
	if err != nil {
		return err
	}
	defer file1.Close()
	str1 := " <svg  version=\"1.1\" width=\"1280\" height=\"1024\"> </svg>"
	fmt.Fprint(file1, str1)

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
