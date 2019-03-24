package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var timeStart2Phase = time.Now()
var isRating = false
var indexesUpdated = false

var availableFilters = map[string]bool{
	"sex_eq":             true,
	"email_domain":       true,
	"email_lt":           true,
	"email_gt":           true,
	"status_eq":          true,
	"status_neq":         true,
	"fname_eq":           true,
	"fname_any":          true,
	"fname_null":         true,
	"sname_eq":           true,
	"sname_starts":       true,
	"sname_null":         true,
	"phone_code":         true,
	"phone_null":         true,
	"country_eq":         true,
	"country_null":       true,
	"city_eq":            true,
	"city_any":           true,
	"city_null":          true,
	"birth_lt":           true,
	"birth_gt":           true,
	"birth_year":         true,
	"interests_contains": true,
	"interests_any":      true,
	"likes_contains":     true,
	"premium_now":        true,
	"premium_null":       true,
	"limit":              true,
}

var availableGroupArgs = map[string]bool{
	"keys":      true,
	"limit":     true,
	"order":     true,
	"sex":       true,
	"status":    true,
	"fname":     true,
	"sname":     true,
	"phone":     true,
	"country":   true,
	"city":      true,
	"birth":     true,
	"interests": true,
	"likes":     true,
	"premium":   true,
	"joined":    true,
}

var availableKeys = map[string]bool{
	"sex":       true,
	"status":    true,
	"interests": true,
	"country":   true,
	"city":      true,
}

func Group(ctx *fasthttp.RequestCtx) {

	args := make(map[string]string)
	valid := true
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "query_id" {
			return
		}
		//валидация полей фильтра
		_, ok := availableGroupArgs[string(key)]
		if !ok {
			valid = false
			return
		}
		args[string(key)] = string(value)
	})

	limit, ok := args["limit"]
	if !ok {
		valid = false
	}

	_, err := strconv.Atoi(limit)
	if err != nil {
		valid = false
	}

	for _, key := range strings.Split(args["keys"], ",") {
		if _, ok = availableKeys[key]; !ok {
			valid = false
			break
		}
	}

	if !valid {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	result := getGroups(args)

	json, _ := json.Marshal(result)

	ctx.Response.SetBody(json)
}

func Filter(ctx *fasthttp.RequestCtx) {

	args := make(map[string]string)
	valid := true
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "query_id" {
			return
		}
		//валидация полей фильтра
		_, ok := availableFilters[string(key)]
		if !ok {
			valid = false
			return
		}
		args[string(key)] = string(value)
	})

	limit, ok := args["limit"]
	if !ok {
		valid = false
	}

	_, err := strconv.Atoi(limit)
	if err != nil {
		valid = false
	}

	if !valid {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	result := getAccountsByFilter(args)

	json, _ := json.Marshal(result)

	ctx.Response.SetBody(json)
}

func main() {

	fmt.Println(time.Now().UTC(), " Start loading data.")
	err := PrepareData("data/data/data.zip")
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Now().UTC(), " Data loaded success.")

	if len(accountsIds) > 100000 {
		isRating = true
	}

	ticker := time.NewTicker(2000 * time.Millisecond)
	countNewAccounts := 0
	go func() {
		for t := range ticker.C {
			sinceStart := t.Sub(timeStart2Phase).Seconds()

			if ((isRating && sinceStart > 1430) || (!isRating && sinceStart > 210)) && !indexesUpdated && countNewAccounts == len(newAccountIds) {
				fmt.Println(time.Now().UTC(), " Update indexes.")
				updateIndexes()
				fmt.Println(time.Now().UTC(), " End of update indexes.")
			}

			if countNewAccounts != len(newAccountIds) {
				countNewAccounts = len(newAccountIds)
			}

			PrintMemUsage()
		}
	}()

	// the corresponding fasthttp request handler
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetContentType("application/json; charset=utf-8")
		switch string(ctx.Method()) {
		case "GET":
			switch string(ctx.Path()) {
			case "/accounts/filter/":
				Filter(ctx)
			case "/accounts/group/":
				Group(ctx)
			default:
				pathParts := strings.Split(string(ctx.Path())[1:], "/")
				if len(pathParts) > 2 && pathParts[0] == "accounts" {
					if id, ok := strconv.Atoi(pathParts[1]); ok == nil {
						if _, exist := accounts[uint32(id)]; exist {
							ctx.Response.SetBody([]byte("{\"accounts\":[]}"))
						} else {
							ctx.Error("Unsupported path", fasthttp.StatusNotFound)
						}
					} else {
						ctx.Error("Unsupported path", fasthttp.StatusNotFound)
					}
				} else {
					ctx.Error("Unsupported path", fasthttp.StatusNotFound)
				}
			}
		case "POST":
			switch string(ctx.Path()) {
			case "/accounts/new/":
				NewAccount(ctx)
			default:
				pathParts := strings.Split(string(ctx.Path())[1:], "/")
				if len(pathParts) > 1 && pathParts[0] == "accounts" {
					if pathParts[1] == "likes" {
						ctx.Response.SetBody([]byte("{}"))
						ctx.SetStatusCode(fasthttp.StatusAccepted)
					} else if id, ok := strconv.Atoi(pathParts[1]); ok == nil {
						mutexChangeAccount.RLock()
						if _, exist := accounts[uint32(id)]; exist {
							mutexChangeAccount.RUnlock()
							Update(ctx, id)
						} else {
							mutexChangeAccount.RUnlock()
							ctx.Error("Unsupported path", fasthttp.StatusNotFound)
						}
					} else {
						ctx.Error("Unsupported path", fasthttp.StatusNotFound)
					}
				}
			}
		}

	}

	log.Fatal(fasthttp.ListenAndServe(":8080", requestHandler))
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
