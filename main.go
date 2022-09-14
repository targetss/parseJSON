// Rick And Morty API
package main

import (
	"encoding/json"
	"strings"

	//"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	//"strings"
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
		url       string     = "https://rickandmortyapi.com/api/character/?page="
		result    *[]JsonRaM = new([]JsonRaM)
		search    string
		uniquearr []string
		//newstr []byte
		//search string
		//Year   int
		//Month  int
		//Day    int
	)

	RequestData(url, result)

	for {
		fmt.Println("Введите поле для сортировки:")
		fmt.Scan(&search)

		uniquearr = UniqueData(result, strings.ToLower(search))

		for ind, val := range uniquearr {
			fmt.Printf("Index:%v\tValue:%v\n", ind, val)
		}
		fmt.Println("========================================================================================")
	}

	//fmt.Println(*result)
	/*
		fmt.Println("Введите дату 'От'\nВведите год:")
		fmt.Scan(&Year)
		fmt.Println("Введите месяц:")
		fmt.Scan(&Month)
		fmt.Println("Введите день:")
		fmt.Scan(&Day)

		if Month > 9 && Month < 13 {

		}
		date1 := fmt.Sprintf("%d-0%d-0%d", Year, Month, Day)
		date2, _ := time.Parse(times, date1)
		fmt.Println(date2)
	*/

	//fmt.Println("Введите слово для поиска:")
	//fmt.Scan(&search)

	//err = PrintDataCategory(result, search)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//UniqueCategoryData(result)

}

func RequestData(url string, datajson *[]JsonRaM) {
	var (
		rr JsonRaM
	)
	respT, err := http.Get(fmt.Sprintf("%v1", url))
	if err != nil {
		fmt.Println("No response from request")
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
	fmt.Println(sort)
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

/*
func CountCategoryData(data *[]js, strSort string) int {
	var (
		count int
	)

	for _, val := range *data {
		if val.Category == strSort {
			count++
		}
	}
	return count
}

func PrintDataCategory(str *[]js, strSearch string) error {
	var (
		count       uint16 = uint16(len(*str))
		count_range uint16
	)

	err := errors.New("Data is not found!")

	if count != 0 {
		fmt.Printf("%15s %100s\n", "Type", "Link")
		for _, val := range *str {
			switch searchstr := strings.Contains(val.Category, strSearch); searchstr {
			case true:
				fmt.Printf("%s %100s\n", val.API, val.Link)
			case false:
				count_range++
				break
			}
		}

		return nil
	} else {
		return err
	}
}

*/
