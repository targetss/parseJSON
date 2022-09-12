package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type js struct {
	API         string `json:"API"`
	Description string `json:"Description"`
	Auth        string `json:"Auth"`
	HTTPS       bool   `json:"HTTPS"`
	Cors        string `json:"Cors"`
	Link        string `json:"Link"`
	Category    string `json:"Category"`
}

const times string = "2006-01-02"

func main() {
	rand.Seed(time.Now().UnixNano())

	var (
		result = new([]js)
		newstr []byte
		search string
		Year   int
		Month  int
		Day    int
	)
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

	resp, err := http.Get("https://api.publicapis.org/entries")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // возвращает []byte

	//fmt.Println(strings.Index(string(body), "[")) // находим индекс символа ""

	newstr = body[strings.Index(string(body), "[") : strings.Index(string(body), "]")+1] //для массива нужно включать в строку "[" и "]"

	if err := json.Unmarshal(newstr, result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println("Введите слово для поиска:")
	fmt.Scan(&search)

	err = PrintDataCategory(result, search)
	if err != nil {
		fmt.Println(err)
	}

	UniqueCategoryData(result)

}

func UniqueCategoryData(data *[]js) {
	fmt.Println("UniqueCategoryData")
	var (
		uniquedata = make([]string, 0)
		//count      int
	)
	fmt.Println("len:", len(uniquedata))
Loopb:
	for ind, val := range *data {
		if ind == 0 {
			uniquedata = append(uniquedata, val.Category)
		}
		for i := 0; i < len(uniquedata); i++ {
			if val.Category == uniquedata[i] {
				continue Loopb
			}
			uniquedata = append(uniquedata, val.Category)
		}
	}
	for _, val := range uniquedata {
		fmt.Println(val)
	}
}

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
