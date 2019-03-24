package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var timestamp int64 = time.Now().Unix()
var accounts = make(map[uint32]string, 1400000)
var accountsIds = make([]uint32, 0, 1400000)
var emails = make(map[string]bool)
var phones = make(map[string]bool)

var likes = make(map[uint32][]uint32)

type FileAccounts struct {
	Accounts []Account
}

func PrepareData(src string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	indexes.NewIndex("email_domain")
	indexes.NewIndex("fname")
	indexes.NewIndex("sname")
	indexes.NewIndex("country")
	indexes.NewIndex("city")
	indexes.NewIndex("birth_year")
	indexes.NewIndex("premium")
	indexes.NewIndex("likes")
	indexes.NewIndex("phone")
	indexes.NewIndex("interests")

	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		matched, err := regexp.MatchString("accounts.+\\.json", f.Name)
		buf := new(bytes.Buffer)
		if matched {
			buf.ReadFrom(rc)
			f := FileAccounts{}
			json.Unmarshal(buf.Bytes(), &f)

			sort.Slice(f.Accounts, func(i, j int) bool {
				return f.Accounts[i].ID > f.Accounts[j].ID
			})

			for _, account := range f.Accounts {
				accountsIds = append(accountsIds, uint32(account.ID))
				accounts[uint32(account.ID)] = account.Encode()
				emails[account.Email] = true
				phones[account.Phone] = true
				for _, interest := range account.Interests {
					indexes.UpdateIndex("interests", uint32(account.ID), interest)
				}
				indexes.UpdateIndex("email_domain", uint32(account.ID), account.GetEmailDomain())
				indexes.UpdateIndex("fname", uint32(account.ID), account.Fname)
				indexes.UpdateIndex("sname", account.ID, account.Sname)
				indexes.UpdateIndex("country", uint32(account.ID), account.Country)
				indexes.UpdateIndex("city", uint32(account.ID), account.City)
				indexes.UpdateIndex("birth_year", uint32(account.ID), account.GetYear())
				premium := account.GetPremium()
				if premium != "" {
					if premium == "1" {
						indexes.UpdateIndex("premium", uint32(account.ID), "2")
					}
					indexes.UpdateIndex("premium", uint32(account.ID), premium)
				}

				for _, like := range account.Likes {
					indexes.UpdateIndex("likes", uint32(account.ID), string(like.ID))
				}

				if account.Phone != "" {
					indexes.UpdateIndex("phone", uint32(account.ID), account.GetPhoneCode())
				} else {
					indexes.UpdateIndex("phone", uint32(account.ID), "")
				}
			}
		}

		return nil
	}

	dat, err := ioutil.ReadFile("/tmp/data/options.txt")
	if err != nil {
		// panic(err)
	}
	timeString := strings.Split(string(dat), "\n")
	tInt32, _ := strconv.Atoi(timeString[0])
	timestamp = int64(tInt32)

	sort.Slice(r.File, func(i, j int) bool {
		name1 := strings.Split(strings.Split(r.File[i].Name, ".")[0], "_")
		name2 := strings.Split(strings.Split(r.File[j].Name, ".")[0], "_")

		number1, _ := strconv.Atoi(name1[1])
		number2, _ := strconv.Atoi(name2[1])
		return number1 > number2
	})

	fmt.Println("Start import files")

	for _, f := range r.File {
		fmt.Println(f.Name)
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	fmt.Println(likes[1])

	runtime.GC()

	printIndexStat()

	return nil
}

func printIndexStat() {
	for key, index := range indexes {
		fmt.Println(key, len(index))
	}
}
