package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var err error

func init() {
	//Начало работы, читаем настроечный фаил
	if err = godotenv.Load(); err != nil {
		fmt.Println("Can't load enc file - ", err.Error())
	}
}

func main() {
	//Загружаем модуль логирования
	if err = Init(os.Getenv("logger_path")); err != nil {
		fmt.Println("Error opening logger subsystem ", err.Error())
		return
	}

	//Подключение к базе данных
	if err = ConnectDB(); err != nil {
		Info.Println("Error open DB", err.Error())
		fmt.Println("Error open DB", err.Error())
		return
	}
	defer GetDB().Close() // не забывает закрыть подключение

	Info.Println("Start work...")
	fmt.Println("Start work...")
	//----------------------------------------------------------------------

	tableReg, tableArea, err := GetRegionInfo()
	fmt.Println(tableArea)
	tableTL := getTrafficLights()
	for numReg, nameReg := range tableReg {
		var strTxt []string
		pathReg := os.Getenv("dir_path") + "//" + strconv.Itoa(numReg)
		os.Mkdir(pathReg, os.ModePerm)
		strTxt = append(strTxt, nameReg)
		for numArea, nameArea := range tableArea[nameReg] {
			pathArea := pathReg + "//" + strconv.Itoa(numArea)
			os.Mkdir(pathArea, os.ModePerm)
			straa := strconv.Itoa(numArea) + " " + nameArea
			strTxt = append(strTxt, straa)
			for _, TL := range tableTL {
				pathTL := pathArea + "//" + strconv.Itoa(TL.ID)
				if numReg == TL.Regin && TL.Area == numArea {
					os.Mkdir(pathTL, os.ModePerm)
					tempstr := strconv.Itoa(TL.ID) + "   " + TL.Description
					fmt.Println(TL.Description)
					err = makeBmp(TL, pathTL+"//")
					if err != nil {
						Info.Println(err.Error())
					}
					strTxt = append(strTxt, tempstr)
				}
			}
		}

		SaveFile(pathReg+"//Info.txt", strTxt)
	}
	fmt.Println("DONE!!!")
}

func SaveFile(FileName string, data []string) (err error) {
	file, err := os.Create(FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, str := range data {
		fmt.Fprint(file, str+"\n")
	}
	return nil
}
