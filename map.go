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

var (
	str1 = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
			<svg
			xmlns:svg="http://www.w3.org/2000/svg"
			xmlns="http://www.w3.org/2000/svg"
			width="450mm"
			height="450mm"
			viewBox="0 0 450 450">
			<foreignObject x="5" y="5" width="100" height="450">
			<div xmlns="http://www.w3.org/1999/xhtml"
		style="font-size:8px;font-family:sans-serif">`
	str2 = `</div>
			</foreignObject>
 			</svg>`
)

func makeBmp(TL TrafficLights, filepath string) (err error) {
	url := fmt.Sprintf("https://static-maps.yandex.ru/1.x/?ll=%3.15f,%3.15f&z=19&l=map&size=450,450", TL.Points.Y, TL.Points.X)
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
	str3 := fmt.Sprintf("%s", TL.Description)
	fmt.Fprintln(file1, str1, str3, str2)

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
