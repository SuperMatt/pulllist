package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
	//"reflect"
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

type allJSON struct {
	Results map[string]book
}

type book struct {
	ID     string `json:"id"`
	Fields bookFields
}

type bookFields struct {
	Title string `json:"dc_solr_sortable_title"`
	Type  string `json:"type"`
}

func midnight(t time.Time) (midnight time.Time) {

	yy := t.Year()
	mm := t.Month()
	dd := t.Day()

	return time.Date(yy, mm, dd, 0, 0, 0, 0, time.UTC)
}

func main() {
	advance := flag.Int("advance", 0, "How many weeks in advance to show")
	jsonOnly := flag.Bool("json", false, "Show the JSON results only")
	single := flag.Bool("1", false, "Single page mode")
	flag.Parse()

	date := time.Now()
	date = date.Add(time.Duration(*advance*7*24) * time.Hour)

	day := date.Weekday()

	daysUntilWednesday := 0
	if day < 3 {
		daysUntilWednesday = 3 - int(day)
	} else if day > 3 {
		daysUntilWednesday = 10 - int(day)
	}

	nextWednesday := date.Add(time.Duration(24*daysUntilWednesday) * time.Hour)
	nextThursday := date.Add(time.Duration(24*(daysUntilWednesday+1)) * time.Hour)

	nextWednesdayStart := midnight(nextWednesday)
	nextThursdayStart := midnight(nextThursday)

	start := strconv.FormatInt(nextWednesdayStart.Unix(), 10)
	end := strconv.FormatInt(nextThursdayStart.Unix(), 10)

	booktype := "comic%7Cgraphic_novel"

	url := "http://www.dccomics.com/proxy/search?type=" + booktype + "&startdate=" + start + "&enddate=" + end
	//fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if *jsonOnly {
		fmt.Println(string(body))
		os.Exit(0)
	}

	var alljson allJSON
	err = json.Unmarshal(body, &alljson)
	if err != nil {
		fmt.Println(err.Error())
	}

	var comics []string
	var graphicNovels []string

	for _, v := range alljson.Results {
		if v.Fields.Type == "comic" {
			comics = append(comics, strings.Title(strings.ToLower(v.Fields.Title)))
		} else if v.Fields.Type == "graphic_novel" {
			graphicNovels = append(graphicNovels, strings.Title(strings.ToLower(v.Fields.Title)))
		}
	}

	sort.Strings(comics)
	sort.Strings(graphicNovels)

	_, height, err := terminal.GetSize(0)
	maxlines := height - 1

	var list []string

	list = append(list, fmt.Sprintf("Comics:\n"))
	for _, v := range comics {
		list = append(list, fmt.Sprintf("%v\n", v))
	}
	list = append(list, fmt.Sprintf("\n"))
	list = append(list, fmt.Sprintf("Graphic Novels:\n"))
	for _, v := range graphicNovels {
		list = append(list, fmt.Sprintf("%v\n", v))
	}

	for k, v := range list {
		if k == 0 || k%maxlines != 0 {
			fmt.Printf(v)
			continue
		}
		if !*single {
			var input string
			fmt.Fprintf(os.Stderr, "MORE")
			fmt.Scanln(&input)
		}
	}
}
