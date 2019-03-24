package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"sort"
	"sync"

	"github.com/valyala/fasthttp"
)

var mutexChangeAccount = &sync.RWMutex{}
var newAccountIds = []uint32{}

var newIndexes = map[string]map[string][]uint32{
	"email_domain": make(map[string][]uint32),
	"fname":        make(map[string][]uint32),
	"sname":        make(map[string][]uint32),
	"country":      make(map[string][]uint32),
	"city":         make(map[string][]uint32),
	"birth_year":   make(map[string][]uint32),
	"premium":      make(map[string][]uint32),
	"phone":        make(map[string][]uint32),
	"interests":    make(map[string][]uint32),
}

var newLikes = make(map[uint32][]uint32)

func NewAccount(ctx *fasthttp.RequestCtx) {

	account := Account{}

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&account); err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if account.ID == 0 {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if _, ok := emails[account.Email]; ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if account.Phone != "" {
		if _, ok := phones[account.Phone]; ok {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	createAccount(account)

	ctx.Response.SetBody([]byte("{}"))
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func createAccount(account Account) {
	mutexChangeAccount.Lock()

	email_domain := account.GetEmailDomain()
	newIndexes["email_domain"][email_domain] = append(newIndexes["email_domain"][email_domain], account.ID)
	emails[account.Email] = true
	phones[account.Phone] = true
	newIndexes["fname"][account.Fname] = append(newIndexes["fname"][account.Fname], account.ID)
	newIndexes["sname"][account.Sname] = append(newIndexes["sname"][account.Sname], account.ID)
	newIndexes["country"][account.Country] = append(newIndexes["country"][account.Country], account.ID)
	newIndexes["city"][account.City] = append(newIndexes["city"][account.City], account.ID)
	birth_year := account.GetYear()
	newIndexes["birth_year"][birth_year] = append(newIndexes["birth_year"][birth_year], account.ID)
	if account.Phone != "" {
		phone_code := account.GetPhoneCode()
		newIndexes["phone"][phone_code] = append(newIndexes["phone"][phone_code], account.ID)
	}
	for _, interest := range account.Interests {
		newIndexes["interests"][interest] = append(newIndexes["interests"][interest], account.ID)
	}

	premium := account.GetPremium()
	newIndexes["premium"][premium] = append(newIndexes["premium"][premium], account.ID)

	for _, like := range account.Likes {
		newLikes[like.ID] = append(newLikes[like.ID], account.ID)
	}

	accounts[account.ID] = account.Encode()
	newAccountIds = append(newAccountIds, account.ID)
	defer mutexChangeAccount.Unlock()
}

func updateIndexes() {
	indexesUpdated = true
	mutexChangeAccount.Lock()
	sort.Slice(newAccountIds, func(i int, j int) bool {
		return newAccountIds[i] > newAccountIds[j]
	})
	accountsIds = append(newAccountIds, accountsIds...)

	for name, newIndex := range newIndexes {
		if _, ok := indexes[name]; ok {
			for value, newIds := range newIndex {
				sort.Slice(newIds, func(i int, j int) bool {
					return newIds[i] > newIds[j]
				})
				indexes[name][value] = append(newIds, indexes[name][value]...)
			}
		}
	}

	for name, index := range addIndexes {
		for val, ids := range index {
			indexes[name][val] = append(indexes[name][val], ids...)
			sort.Slice(indexes[name][val], func(i int, j int) bool {
				return indexes[name][val][i] > indexes[name][val][j]
			})
		}
	}

	for name, index := range removeIndexes {
		for val, ids := range index {
			for _, id := range ids {
				i := sort.Search(len(indexes[name][val]), func(i int) bool {
					return indexes[name][val][i] <= id
				})

				if i != len(indexes[name][val]) && indexes[name][val][i] == id {
					indexes[name][val] = append(indexes[name][val][0:i], indexes[name][val][i+1:]...)
				}
			}
		}
	}

	newIndexes = nil
	addIndexes = nil
	newLikes = nil
	removeIndexes = nil
	indexes["likes"] = make(map[string][]uint32)

	defer mutexChangeAccount.Unlock()
	fmt.Println(accountsIds[0:10])
	runtime.GC()
}

func appendSorted(ints []uint32, id uint32) []uint32 {
	l := len(ints)
	if l == 0 {
		ints = []uint32{id}
	}

	i := sort.Search(l, func(i int) bool { return ints[i] < id })

	if i == l {
		ints = append([]uint32{id}, ints...)
		return ints
	}

	if i == 0 {
		ints = append([]uint32{id}, ints...)
		return ints
	}

	return append(ints[0:i], append([]uint32{id}, ints[i:]...)...)
}
