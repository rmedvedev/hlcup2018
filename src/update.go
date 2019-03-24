package main

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/valyala/fasthttp"
)

var removeIndexes = map[string]map[string][]uint32{
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
var addIndexes = map[string]map[string][]uint32{
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

func Update(ctx *fasthttp.RequestCtx, id int) {

	account := accountDecode(accounts[uint32(id)])
	account.ID = uint32(id)
	accountTmp := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewReader(ctx.PostBody()))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&account); err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	json.Unmarshal(ctx.PostBody(), &accountTmp)

	if accountTmp["email"] != nil {
		if _, ok := emails[accountTmp["email"].(string)]; ok {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	if accountTmp["phone"] != nil {
		if _, ok := phones[accountTmp["phone"].(string)]; ok {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	for name, val := range accountTmp {
		switch name {
		case "email":
			if !strings.Contains(val.(string), "@") {
				ctx.Error("", fasthttp.StatusBadRequest)
				return
			}
		}
	}

	if _, ok := statuses[account.Status]; !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	updateAccount(account, accountTmp)

	ctx.Response.SetBody([]byte("{}"))
	ctx.SetStatusCode(fasthttp.StatusAccepted)
}

func updateAccount(account Account, changedFields map[string]interface{}) {
	mutexChangeAccount.Lock()
	currentAccount := accountDecode(accounts[account.ID])
	currentAccount.ID = account.ID
	for name := range changedFields {
		switch name {
		case "interests":
			for _, interest := range currentAccount.Interests {
				removeIndexes["interests"][interest] = append(removeIndexes["interests"][interest], account.ID)
			}
			for _, interest := range account.Interests {
				addIndexes["interests"][interest] = append(addIndexes["interests"][interest], account.ID)
			}
		case "likes":
		case "premium":
			oldPremium := currentAccount.GetPremium()
			premium := currentAccount.GetPremium()
			if oldPremium == "1" {
				removeIndexes["premium"]["2"] = append(removeIndexes["premium"]["2"], account.ID)
			}
			removeIndexes["premium"][oldPremium] = append(removeIndexes["premium"][oldPremium], account.ID)

			if premium == "1" {
				addIndexes["premium"]["2"] = append(addIndexes["premium"]["2"], account.ID)
			}
			addIndexes["premium"][premium] = append(addIndexes["premium"][premium], account.ID)

		case "sname":
			removeIndexes["sname"][currentAccount.Sname] = append(removeIndexes["sname"][currentAccount.Sname], account.ID)
			addIndexes["sname"][account.Sname] = append(addIndexes["sname"][account.Sname], account.ID)
		case "fname":
			removeIndexes["fname"][currentAccount.Fname] = append(removeIndexes["fname"][currentAccount.Fname], account.ID)
			addIndexes["fname"][account.Fname] = append(addIndexes["fname"][account.Fname], account.ID)
		case "city":
			removeIndexes["city"][currentAccount.City] = append(removeIndexes["city"][currentAccount.City], account.ID)
			addIndexes["city"][account.City] = append(addIndexes["city"][account.City], account.ID)
		case "country":
			removeIndexes["country"][currentAccount.Country] = append(removeIndexes["country"][currentAccount.Country], account.ID)
			addIndexes["country"][account.Country] = append(addIndexes["country"][account.Country], account.ID)
		case "email":
			oldEmailDomain := currentAccount.GetEmailDomain()
			emailDomain := currentAccount.GetEmailDomain()
			removeIndexes["email_domain"][oldEmailDomain] = append(removeIndexes["email_domain"][oldEmailDomain], account.ID)
			addIndexes["email_domain"][emailDomain] = append(addIndexes["email_domain"][emailDomain], account.ID)
			delete(emails, account.Email)
		case "phone":
			oldPhoneCode := currentAccount.GetPhoneCode()
			phoneCode := account.GetPhoneCode()
			if oldPhoneCode != "" {
				removeIndexes["phone"][oldPhoneCode] = append(removeIndexes["phone"][oldPhoneCode], account.ID)
			}
			if phoneCode != "" {
				addIndexes["phone"][phoneCode] = append(addIndexes["phone"][phoneCode], account.ID)
			}
			delete(phones, account.Phone)
		case "birth":
			oldBirthYear := currentAccount.GetYear()
			birthYear := account.GetYear()
			removeIndexes["birth_year"][oldBirthYear] = append(removeIndexes["birth_year"][oldBirthYear], account.ID)
			addIndexes["birth_year"][birthYear] = append(addIndexes["birth_year"][birthYear], account.ID)
		}
	}

	accounts[account.ID] = account.Encode()
	defer mutexChangeAccount.Unlock()
}
