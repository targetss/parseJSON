// Rick And Morty API
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/jpeg"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type InfoData struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type Origin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type JsonRaM struct {
	Info    InfoData     `json:"info"`
	Results []PersonInfo `json:"results"`
}

type PersonInfo struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Species  string   `json:"species"`
	Type     string   `json:"type"`
	Gender   string   `json:"gender"`
	Origin   Origin   `json:"origin"`
	Location Location `json:"location"`
	Image    string   `json:"image"`
	Episode  []string `json:"episode"`
	Url      string   `json:"url"`
	Created  string   `json:"created"`
}

type PersonUni struct {
	ID       int
	Name     string
	Status   string
	Species  string
	Type     string
	Gender   string
	Origin   string
	Location string
	Image    string
	Created  string
}

type JsonRaMInfo []JsonRaM

func (p *JsonRaMInfo) ImportData(url string) {
	if len(*p) > 0 {
		p = new(JsonRaMInfo)
	}

	var (
		rr JsonRaM
	)
	respT, err := http.Get(fmt.Sprintf("%v1", url))
	if err != nil {
		fmt.Errorf("Произошла ошибка: %v", err)
	}
	defer respT.Body.Close()
	body, err := io.ReadAll(respT.Body) // возвращает []byte

	if err := json.Unmarshal(body, &rr); err != nil {
		fmt.Println(err)
	}

	pathCfg := filepath.Join(pathConfig, fileNameCache) //Полный путь до файла
	errchdir := os.Chdir(pathConfig)                    //изменяет текущий рабочий каталог на именованный каталог
	if errchdir != nil {
		err := os.Mkdir(pathConfig, 0755)
		if err == nil {
			log.Println("Создание директории для конфигурационных файлов")
		}
		_, errfle := os.Create(pathCfg)
		if errfle != nil {
			panic(fmt.Sprintf("Ошибка создания файла данный Json, ошибка: %v", err))
		}
	} else {
		if _, err := os.Stat(pathCfg); os.IsNotExist(err) {
			log.Println("Создание файла конфигурации")
			os.Create(pathCfg)
		}
	}

	bytefile, _ := os.ReadFile(pathCfg)
	if len(bytefile) != 0 {
		var tempJson []JsonRaM
		json.Unmarshal([]byte(bytefile), &tempJson)
		if tempJson[0].Info.Count == rr.Info.Count {
			fmt.Println("Считывание данных с файла конфигурации Json..")
			*p = tempJson
			return
		}
	}

	*p = append(*p, rr) // записываем в массив структур данные первой страницы, отсюда вычисляем общее кол-во страниц

	if (*p)[0].Info.Pages > 1 {
		for i := 2; i <= (*p)[0].Info.Pages; i++ {
			var (
				rr JsonRaM
			)
			response, err := http.Get(fmt.Sprintf("%v%d", url, i))
			if err != nil {
				return
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err := json.Unmarshal(body, &rr); err != nil {
				fmt.Println(err)
			}
			*p = append(*p, rr)
		}
	}

	filewrite, err := os.Create(pathCfg)
	bt, _ := json.Marshal(*p)
	countbt, errbt := filewrite.Write(bt)
	if errbt != nil {
		log.Println("Ошибка записи в файл json")
	}
	log.Println(fmt.Sprintf("Записано байт: %v\n", countbt))
	filewrite.Close()
}

func (p *JsonRaMInfo) AddDataInFile() bool {
	cacheDir, _ := os.UserCacheDir()
	os.Mkdir(filepath.Join(cacheDir, "ConfigRaM"), 0755)
	if _, fifeinfo := os.Stat(filepath.Join(cacheDir, "ConfigRaM", "JsonInfo.dat")); os.IsExist(fifeinfo) {
		return false
	}
	fl, _ := os.Create(filepath.Join(cacheDir, "ConfigRaM", "JsonInfo.dat"))

	//file, erropen := os.Open(filepath.Join(cacheDir, "ConfigRaM", "JsonInfo.dat"))
	//if erropen != nil {
	//	log.Println(erropen)
	//}
	//bt := make([]byte, len(*p))
	bt, rtr := json.Marshal(*p)
	fmt.Println(string(bt))
	if rtr != nil {
		log.Println("Ошибка Marshal json")
	}
	fl.Write(bt)
	fl.Close()
	return true
}

func (p *JsonRaMInfo) DownloadImageCharacter() {
	cacheDir, _ := os.UserCacheDir()
	fullpath := filepath.Join(cacheDir, "RaMImg") //join для того, чтобы корректно расставить разделители как в вашей ОС
	os.Mkdir(fullpath, 0755)
	for ind, _ := range *p {
		for _, path2 := range (*p)[ind].Results {
			go func(id int, url, path string) {
				if _, errfile := os.Stat(filepath.Join(fullpath, fmt.Sprintf("%v.jpeg", id))); os.IsExist(errfile) {
					log.Printf("Файл с именем %v.jped существует\n", id)
					return
				}

				//time.Sleep(100 * time.Millisecond) //срабатывает защита от ддос атак
				//fmt.Println("Gorutine go id:", id, "URL: ", url)
				response, err := http.Get(url)
				if err != nil {
					log.Print(err)
				}
				if response.StatusCode == 429 {
					//fmt.Println("Response code = 429, ID=", id, "\n", "Status=", response.Status, "\nStatusCode=", response.StatusCode)
					response.Body.Close()
					time.Sleep(8 * time.Second)

					response, _ := http.Get(url)
					fileCreate, _ := os.Create(filepath.Join(fullpath, fmt.Sprintf("%v.jpeg", id)))
					io.Copy(fileCreate, response.Body)
					fileCreate.Close()
					response.Body.Close()
					return
				}
				fileCreate, err := os.Create(filepath.Join(fullpath, fmt.Sprintf("%v.jpeg", id)))
				io.Copy(fileCreate, response.Body)

				errEOF := errors.New("unexpected EOF")
				loadImage, errdec := jpeg.Decode(fileCreate)
				if !errors.Is(errdec, errEOF) {
					resp, _ := http.Get(url)
					fileCreate.Seek(0, 0)
					io.Copy(fileCreate, resp.Body)
				}
				_ = loadImage

				response.Body.Close()
				fileCreate.Close()
			}(path2.ID, path2.Image, cacheDir)
		}
	}
}

const (
	times string = "2006-01-02"
	url   string = "https://rickandmortyapi.com/api/character/?page="
)

var (
	userCachePath, _        = os.UserCacheDir()
	nameDirCache     string = "RaM"
	pathConfig              = filepath.Join(userCachePath, nameDirCache)
	fileNameCache    string = "JsonData.dat"
)

type CardLayout struct {
}

func (d *CardLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	weight, height := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		weight += childSize.Width
		height += childSize.Height
	}
	return fyne.NewSize(weight, height)
}

func (d *CardLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	fmt.Printf("Type: %T", objects)
	pos := fyne.NewPos(0, containerSize.Height-d.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width, size.Height))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var (
		/*result *[]JsonRaM = new([]JsonRaM)*/
		result      = new(JsonRaMInfo)
		cacheDir, _ = os.UserCacheDir()
	)
	result.ImportData(url)
	//_ = result.AddDataInFile()
	result.DownloadImageCharacter()

	//RequestData(url, result)

	uniqueelemstruct := struct { //анонимная структура с уникальными полями, для полей выбора сортировки
		//ID []int
		Name         []string
		Status       []string
		Species      []string
		Type         []string
		Gender       []string
		OriginName   []string
		LocationName []string
	}{
		Name:         UniqueData(*result, "Name"),
		Status:       UniqueData(*result, "Status"),
		Species:      UniqueData(*result, "Species"),
		Type:         UniqueData(*result, "Type"),
		Gender:       UniqueData(*result, "Gender"),
		OriginName:   UniqueData(*result, "Origin"),
		LocationName: UniqueData(*result, "Location"),
	}

	_ = uniqueelemstruct

	test := PersonID(*result)

	a := app.New()
	w := a.NewWindow("Rick And Morty")
	w.Resize(fyne.NewSize(1000, 900))

	file_item1 := fyne.NewMenuItem("Обновить", func() {
		result.ImportData(url)
	})
	menu1 := fyne.NewMenu("Файл", file_item1)
	main_menu := fyne.NewMainMenu(menu1)
	w.SetMainMenu(main_menu)

	// =============== Поле карточки персонажа ================//
	image := canvas.NewImageFromFile(fmt.Sprintf("%v\\RaMImg\\1.jpeg", cacheDir))
	image.Resize(fyne.Size{Height: 250, Width: 250})

	labelNameField := widget.NewLabel("Имя:")
	labelName := widget.NewLabel("Unknown")

	labelStatusField := widget.NewLabel("Статус")
	labelStatus := widget.NewLabel("Unknown")

	labelSpeciesField := widget.NewLabel("Разновидность:")
	labelSpecies := widget.NewLabel("Unknown")

	labelTypeField := widget.NewLabel("Тип:")
	labelType := widget.NewLabel("Unknown")

	labelGenderField := widget.NewLabel("Пол:")
	labelGender := widget.NewLabel("Unknown")

	tableCard := container.NewVBox(container.NewHBox(labelNameField, labelName), container.NewHBox(labelStatusField, labelStatus),
		container.NewHBox(labelSpeciesField, labelSpecies), container.NewHBox(labelTypeField, labelType), container.NewHBox(labelGenderField, labelGender))
	// =============== Поле карточки персонажа ================//

	listID := widget.NewList(
		func() int {
			return len(test)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(fmt.Sprintf("ID: %v, Name: %v", i+1, test[i].Name))
		})

	var contN *fyne.Container

	rr := container.NewWithoutLayout(image)

	rr2 := container.NewVBox(rr, tableCard)

	listID.OnSelected = func(id widget.ListItemID) {
		img := canvas.NewImageFromFile(fmt.Sprintf("%v\\RaMImg\\%v.jpeg", cacheDir, id+1))
		img.FillMode = canvas.ImageFillOriginal
		//img.Resize(fyne.Size{Height: 250, Width: 250})
		//img.Move(fyne.NewPos(250, 250))
		//ttt := container.NewWithoutLayout(img)
		//contN = container.New(&CardLayout{}, img, tableCard)
		contN = container.NewVBox(img, tableCard)

		labelName.SetText(test[id].Name)
		labelStatus.SetText(test[id].Status)
		labelSpecies.SetText(test[id].Species)
		labelType.SetText(test[id].Type)
		labelGender.SetText(test[id].Gender)
		w.SetContent(container.NewHSplit(listID, contN))
		w.Show()
	}

	//res, _ := fyne.LoadResourceFromURLString("https://rickandmortyapi.com/api/character/avatar/21.jpeg")
	//img := canvas.NewImageFromResource(res)
	//l := container.New(layout.NewGridLayout(3), listStatus, listSpecies, img)
	w.SetContent(container.NewHSplit(listID, rr2))
	w.ShowAndRun()

}

func DownloadImage(id int, url string) {
	_, err := os.Stat(fmt.Sprintf("./img/%v.jpeg", id))
	if err == nil {
		//fmt.Println("File exists")
		return
	}
	resp, _ := http.Get(url)

	defer resp.Body.Close()

	filecrt, err := os.Create(fmt.Sprintf("./img/%v.jpeg", id))
	if err != nil {
		fmt.Println("Ошибка создания файла!")
	}
	io.Copy(filecrt, resp.Body)
	filecrt.Write([]byte{1, 4})
	defer filecrt.Close()

	errEOF := errors.New("unexpected EOF")
	loadImage, errdec := jpeg.Decode(filecrt)
	if !errors.Is(errdec, errEOF) {
		resp, _ := http.Get(url)
		filecrt.Seek(0, io.SeekStart)
		io.Copy(filecrt, resp.Body)
	}
	_ = loadImage
	//fmt.Println("Завершение создания файла")
}

func UniqueData(data JsonRaMInfo, sort string) []string {
	//fmt.Println(sort)
	sort = strings.ToLower(sort)
	//fmt.Println("UniqueCategoryData")
	var (
		arrdata    = make([]string, 0)
		uniquedata = make([]string, 0)
		//nametable  = make([]string, 0)
		//count      int
	)
	//fmt.Println((*data)[0].Results[0].Name)

	for _, vl := range data[0].Results {
		//nametable = append(nametable, string(vl)) //дописать
		fmt.Sprint(vl)
	}

	for _, val := range data {
		for _, valn := range val.Results {
			switch sort {
			case "name":
				arrdata = append(arrdata, valn.Name)
			case "status":
				arrdata = append(arrdata, valn.Status)
			case "species":
				arrdata = append(arrdata, valn.Species)
			case "type":
				arrdata = append(arrdata, valn.Type)
			case "gender":
				arrdata = append(arrdata, valn.Gender)
			case "origin":
				arrdata = append(arrdata, valn.Origin.Name)
			case "location":
				arrdata = append(arrdata, valn.Location.Name)
			default:
				break
			}
		}
	}

	uniquedata = append(uniquedata, arrdata[0])
Loop:
	for ind, val := range arrdata {
		if ind == 0 {
			continue Loop
		}
		for ind2, val2 := range uniquedata {
			if val == val2 {
				continue Loop
			}
			if ind2+1 == len(uniquedata) && val != val2 {
				uniquedata = append(uniquedata, val)
			}
		}
	}
	return uniquedata
}

func PersonID(data JsonRaMInfo) []PersonUni {
	var (
		temp = make([]PersonUni, 0)
	)

	for _, val := range data {
		for _, val2 := range val.Results {
			var dt PersonUni = PersonUni{
				ID:       val2.ID,
				Name:     val2.Name,
				Status:   val2.Status,
				Species:  val2.Species,
				Type:     val2.Type,
				Gender:   val2.Gender,
				Origin:   val2.Origin.Name,
				Location: val2.Location.Name,
				Image:    val2.Image,
				Created:  val2.Created,
			}
			temp = append(temp, dt)
		}
	}
	return temp
}
