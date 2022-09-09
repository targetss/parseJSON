package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func main() {
	var (
		newstr []byte
		search string
	)
	resp, err := http.Get("https://api.publicapis.org/entries")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // возвращает []byte

	//fmt.Println(strings.Index(string(body), "[")) // находим индекс символа ""

	newstr = body[strings.Index(string(body), "[") : strings.Index(string(body), "]")+1] //для массива нужно включать в строку "[" и "]"

	//fmt.Println(string(newstr))

	var result []js
	if err := json.Unmarshal(newstr, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	fmt.Println("Введите слово для поиска:")
	fmt.Scan(&search)

	fmt.Printf("%15s %100s", "Type", "Link")
	for rec := range result {
		switch searchn := strings.Contains(result[rec].Category, search); searchn {
		case true:
			fmt.Printf("%s %100s\n", result[rec].API, result[rec].Link)
		case false:
			break
		}
	}
}
