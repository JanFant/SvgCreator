package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

//RegionInfo расшифровка региона
type RegionInfo struct {
	Num        int    `json:"num"`        //уникальный номер региона
	NameRegion string `json:"nameRegion"` //расшифровка номера
}

type AreaInfo struct {
	Num      int    `json:"num"`      //уникальный номер зоны
	NameArea string `json:"nameArea"` //расшифровка номера
}

type TrafficLights struct {
	ID          int
	Regin       int
	Area        int
	Description string
	Points      Point
}

//GetRegionInfo получить таблицу регионов
func GetRegionInfo() (region map[int]string, area map[string]map[int]string, err error) {
	region = make(map[int]string)
	area = make(map[string]map[int]string)
	sqlStr := fmt.Sprintf("select region, nameregion, area, namearea from %s", os.Getenv("region_table"))
	rows, err := GetDB().Raw(sqlStr).Rows()
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		var (
			tempReg  = &RegionInfo{}
			tempArea = &AreaInfo{}
		)
		err = rows.Scan(&tempReg.Num, &tempReg.NameRegion, &tempArea.Num, &tempArea.NameArea)
		if err != nil {
			return nil, nil, err
		}
		if _, ok := region[tempReg.Num]; !ok {
			region[tempReg.Num] = tempReg.NameRegion
		}

		if _, ok := area[tempReg.NameRegion][tempArea.Num]; !ok {
			if _, ok := area[tempReg.NameRegion]; !ok {
				area[tempReg.NameRegion] = make(map[int]string)
			}
			area[tempReg.NameRegion][tempArea.Num] = tempArea.NameArea
		}
	}
	return region, area, err
}

//func getRegion() (region []Region) {
//	temp := &Region{}
//	sqlquery := fmt.Sprintf("select region, area, nameregion, namearea from %s", os.Getenv("region_table"))
//	rows, _ := GetDB().Raw(sqlquery).Rows()
//	for rows.Next() {
//		rows.Scan(&temp.Region, &temp.Name)
//		region = append(region, *temp)
//	}
//	return
//}

func getTrafficLights() (trLight []TrafficLights) {
	var dgis string
	temp := &TrafficLights{}
	sqlquery := fmt.Sprintf("select region, id, area, dgis, describ from %s", os.Getenv("gis_table"))
	rows, _ := GetDB().Raw(sqlquery).Rows()
	for rows.Next() {
		rows.Scan(&temp.Regin, &temp.ID, &temp.Area, &dgis, &temp.Description)
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
	file, err := os.Create(filepath + "map.png")
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
