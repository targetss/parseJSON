// Rick And Morty API
package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"math/rand"
	"net/http"
	"os"
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

const times string = "2006-01-02"
const url string = "https://rickandmortyapi.com/api/character/?page="

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
		//url    string     = "https://rickandmortyapi.com/api/character/?page="
		result *[]JsonRaM = new([]JsonRaM)
	)

	RequestData(url, result)

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
		Name:         UniqueData(result, "Name"),
		Status:       UniqueData(result, "Status"),
		Species:      UniqueData(result, "Species"),
		Type:         UniqueData(result, "Type"),
		Gender:       UniqueData(result, "Gender"),
		OriginName:   UniqueData(result, "Origin"),
		LocationName: UniqueData(result, "Location"),
	}

	_ = uniqueelemstruct

	test := PersonID(result)

	a := app.New()
	w := a.NewWindow("Rick And Morty")
	w.Resize(fyne.NewSize(1000, 900))

	file_item1 := fyne.NewMenuItem("Обновить", func() {
		RequestData(url, result)
	})
	menu1 := fyne.NewMenu("Файл", file_item1)
	main_menu := fyne.NewMainMenu(menu1)
	w.SetMainMenu(main_menu)

	// =============== Поле карточки персонажа ================//
	image := canvas.NewImageFromFile("./img/1.jpeg")
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
			o.(*widget.Label).SetText(test[i].Name)
		})

	var contN *fyne.Container

	rr := container.NewWithoutLayout(image)

	rr2 := container.NewVBox(rr, tableCard)

	listID.OnSelected = func(id widget.ListItemID) {
		img := canvas.NewImageFromFile(fmt.Sprintf("./img/%v.jpeg", id+1))
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

func RequestData(url string, datajson *[]JsonRaM) {
	var (
		rr JsonRaM
	)
	respT, err := http.Get(fmt.Sprintf("%v1", url))
	if err != nil {
		fmt.Errorf("Произошла ошибка: %v", err)
	}
	defer respT.Body.Close()
	body, err := io.ReadAll(respT.Body) // возвращает []byte

	//fmt.Println(strings.Index(string(body), "[")) // находим индекс символа ""

	//newstr = body[strings.Index(string(body), "[") : strings.Index(string(body), "]")+1] //для массива нужно включать в строку "[" и "]"

	//fmt.Println(string(body))

	if err := json.Unmarshal(body, &rr); err != nil { // Parse []byte to the go struct pointer
		fmt.Println(err)
	}

	*datajson = append((*datajson), rr) // записываем в массив структур данные первой страницы, отсюда вычисляем общее кол-во страниц

	for _, val := range *datajson {
		for _, val2 := range val.Results {
			fmt.Println(val2.ID, val2.Image)
			go DownloadImage(val2.ID, val2.Image)
		}
	}

	fmt.Println((*datajson)[0].Info.Pages) //сначала разыменовываем указатель, а после обращаемся по индексу

	if (*datajson)[0].Info.Pages > 1 {
		for i := 2; i <= (*datajson)[0].Info.Pages; i++ {
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
			(*datajson) = append((*datajson), rr)

			for _, val2 := range (*datajson)[i-1].Results {
				//fmt.Println(val2.ID, val2.Image)
				go DownloadImage(val2.ID, val2.Image)
			}
		}
	}
}

func DownloadImage(id int, url string) {
	_, err := os.Stat(fmt.Sprintf("./img/%v.jpeg", id))
	if err == nil {
		fmt.Println("File exists")
		return
	}

	fmt.Println(id, " | ", url)
	resp, _ := http.Get(url)

	defer resp.Body.Close()

	filecrt, err := os.Create(fmt.Sprintf("./img/%v.jpeg", id))
	if err != nil {
		fmt.Println("Ошибка создания файла!")
	}
	io.Copy(filecrt, resp.Body)
	defer filecrt.Close()
	fmt.Println("Завершение создания файла")
}

func UniqueData(data *[]JsonRaM, sort string) []string {
	//fmt.Println(sort)
	sort = strings.ToLower(sort)
	fmt.Println("UniqueCategoryData")
	var (
		arrdata    = make([]string, 0)
		uniquedata = make([]string, 0)
		//nametable  = make([]string, 0)
		//count      int
	)
	//fmt.Println((*data)[0].Results[0].Name)

	for _, vl := range (*data)[0].Results {
		//nametable = append(nametable, string(vl)) //дописать
		fmt.Sprint(vl)
	}

	for _, val := range *data {
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

func PersonID(data *[]JsonRaM) []PersonUni {
	var (
		temp = make([]PersonUni, 0)
	)

	for _, val := range *data {
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
