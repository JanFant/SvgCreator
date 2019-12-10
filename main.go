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

	tableReg := getRegion()
	tableTL := getTrafficLights()
	for _, region := range tableReg {
		var strTxt []string
		pathReg := os.Getenv("dir_path") + "//" + strconv.Itoa(region.Region)
		os.Mkdir(pathReg, os.ModePerm)
		strTxt = append(strTxt, region.Name)
		for _, TL := range tableTL {
			pathTL := os.Getenv("dir_path") + "//" + strconv.Itoa(region.Region) + "//" + strconv.Itoa(TL.ID)
			if region.Region == TL.Regin {
				os.Mkdir(pathTL, os.ModePerm)
				tempstr := strconv.Itoa(TL.ID) + "   " + TL.Description
				fmt.Println(TL.Description)
				err = makeBmp(TL.Points, pathTL+"//")
				if err != nil {
					Info.Println(err)
				}
				strTxt = append(strTxt, tempstr)
			}
		}
		SaveFile(pathReg+"//Info.txt", strTxt)
	}

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
