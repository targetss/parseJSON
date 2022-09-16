// Rick And Morty API
package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
<<<<<<< HEAD
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
=======
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
>>>>>>> origin/master
	"fyne.io/fyne/v2/widget"
	"io"
	"math/rand"
	"net/http"
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

const times string = "2006-01-02"

func main() {
	rand.Seed(time.Now().UnixNano())

	var (
		url    string     = "https://rickandmortyapi.com/api/character/?page="
		result *[]JsonRaM = new([]JsonRaM)
<<<<<<< HEAD
		i      int        = 1
=======
>>>>>>> origin/master
		//search    string
		//uniquearr []string
		//newstr []byte
		//search string
		//Year   int
		//Month  int
		//Day    int
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
<<<<<<< HEAD

	_ = uniqueelemstruct

	a := app.New()
	w := a.NewWindow("Rick And Morty")
	w.Resize(fyne.NewSize(1000, 900))

	file_item1 := fyne.NewMenuItem("Обновить", func() {
		RequestData(url, result)
	})

	menu1 := fyne.NewMenu("Файл", file_item1)

	main_menu := fyne.NewMainMenu(menu1)
	w.SetMainMenu(main_menu)

	str_test := binding.NewString()
	str_test.Set("Test Name")
	txt := widget.NewLabelWithData(str_test)

	btn_next := widget.NewButton("Далее", func() {
		i *= 2
		str_test.Set(string(i))
	})

	menu := container.NewVBox(txt, btn_next)

	/*card_form := widget.NewForm(
		widget.NewFormItem("Имя: ", canvas.NewText(uniqueelemstruct.Name[1], color.Black)),
		widget.NewFormItem("Статус: ", canvas.NewText(uniqueelemstruct.Status[0], color.Black)),
	)*/

	/*
		listStatus := widget.NewList(
			func() int {
				return len(uniqueelemstruct.Status)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(uniqueelemstruct.Status[i])
			})

		listSpecies := widget.NewList(
			func() int {
				return len(uniqueelemstruct.Species)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(uniqueelemstruct.Species[i])
			})
	*/
	//res, _ := fyne.LoadResourceFromURLString("https://rickandmortyapi.com/api/character/avatar/21.jpeg")
	//img := canvas.NewImageFromResource(res)
	//l := container.New(layout.NewGridLayout(3), listStatus, listSpecies, img)
	w.SetContent(menu)
	w.ShowAndRun()

=======

	a := app.New()
	w := a.NewWindow("Rick And Morty")
	w.Resize(fyne.NewSize(1000, 900))

	listStatus := widget.NewList(
		func() int {
			return len(uniqueelemstruct.Status)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(uniqueelemstruct.Status[i])
		})

	listSpecies := widget.NewList(
		func() int {
			return len(uniqueelemstruct.Species)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(uniqueelemstruct.Species[i])
		})

	res, _ := fyne.LoadResourceFromURLString("https://rickandmortyapi.com/api/character/avatar/21.jpeg")
	img := canvas.NewImageFromResource(res)
	l := container.New(layout.NewGridLayout(3), listStatus, listSpecies, img)
	w.SetContent(l)
	w.ShowAndRun()

	/*
		for {
			fmt.Println("Введите поле для сортировки:")
			fmt.Scan(&search)

			uniquearr = UniqueData(result, strings.ToLower(search))

			for ind, val := range uniquearr {
				fmt.Printf("Index:%v\tValue:%v\n", ind, val)
			}
			fmt.Println("========================================================================================")
		}
	*/
>>>>>>> origin/master
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
		}
	}
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
