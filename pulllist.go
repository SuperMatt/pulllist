package main

import(
    "net/http"
    "fmt"
    "time"
    //"reflect"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "sort"
    "strings"
)

type AllJson struct {
    Results map[string]Book
}

type Book struct {
    Id string `json:"id"`
    Fields BookFields
}

type BookFields struct {
    Title string `json:"dc_solr_sortable_title"`
    Type string `json:"type"`
}

func start_of_day(t time.Time) time.Time {
    t_hour := t.Hour()
    t_minute := t.Minute()
    t_second := t.Second()
    t_nanosecond := t.Nanosecond()
    return t.Add(-time.Duration(t_hour)*time.Hour).Add(-time.Duration(t_minute)*time.Minute).Add(-time.Duration(t_second)*time.Second).Add(-time.Duration(t_nanosecond))
}

func ToCamel(s string) string {
	s = strings.ToLower(s)
	ss := strings.Split(s, " ")
	var cs []string
	for _,v := range(ss) {
		l := string(v[0])
		ns := strings.Replace(v, l, strings.ToUpper(l), 1)
		cs = append(cs, ns)
	}
	rs := strings.Join(cs, " ")
	return rs
}

func main() {
    now := time.Now()
    day := now.Weekday()

    days_until_wednesday := 0
    if day < 3 {
        days_until_wednesday = 3 - int(day)
    } else if day > 3 {
        days_until_wednesday = 10 - int(day)
    }

    next_wednesday := now.Add(time.Duration(24*days_until_wednesday)*time.Hour)
    next_thursday := now.Add(time.Duration(24*(days_until_wednesday + 1))*time.Hour)

    next_wednesday_start := start_of_day(next_wednesday)
    next_thursday_start := start_of_day(next_thursday)

    start := strconv.FormatInt(next_wednesday_start.Unix(), 10)
    end := strconv.FormatInt(next_thursday_start.Unix(), 10)

    booktype := "comic%7Cgraphic_novel"

    url := "http://www.dccomics.com/proxy/search?type=" + booktype + "&startdate=" + start + "&enddate=" + end
    //fmt.Println(url)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println(err.Error())
    }

    body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println(string(body))

    var alljson AllJson
    err = json.Unmarshal(body, &alljson)
    if err != nil {
        fmt.Println(err.Error())
    }

    var comics []string
    var graphic_novels []string

    for _, v := range alljson.Results {
        if v.Fields.Type == "comic" {
            comics = append(comics, ToCamel(v.Fields.Title))
        } else if v.Fields.Type == "graphic_novel" {
            graphic_novels = append(graphic_novels, ToCamel(v.Fields.Title))
        }
    }

    sort.Strings(comics)
    sort.Strings(graphic_novels)
    fmt.Println("Comics:")
    for _, v := range(comics) {
        fmt.Println(v)
    }
	fmt.Println("")
    fmt.Println("Graphic Novels:")
    for _, v := range(graphic_novels) {
        fmt.Println(v)
    }
}
