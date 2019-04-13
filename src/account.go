package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const StatusFree = "свободны"
const StatusComplex = "всё сложно"
const StatusNotfree = "заняты"

const PremiumNow = "1"
const PremiumExist = "2"
const PremiumNone = ""

var statuses = map[string]bool{
	StatusFree:    true,
	StatusComplex: true,
	StatusNotfree: true,
}

type Account struct {
	ID        uint32   `json:"id"`
	Email     string   `json:"email,omitempty"`
	Fname     string   `json:"fname,omitempty"`
	Sname     string   `json:"sname,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Sex       string   `json:"sex,omitempty"`
	Birth     int      `json:"birth,omitempty"`
	Country   string   `json:"country,omitempty"`
	City      string   `json:"city,omitempty"`
	Joined    int      `json:"joined,omitempty"`
	Status    string   `json:"status,omitempty"`
	Interests []string `json:"interests,omitempty"`
	Premium   Premium  `json:"premium,omitempty"`
	Likes     []Like   `json:"likes,omitempty"`
}

type Premium struct {
	Start  int `json:"start"`
	Finish int `json:"finish"`
}

type Like struct {
	ID uint32 `json:"id"`
	Ts int    `json:"ts"`
}

func (account *Account) Encode() string {
	status := 0
	if account.Status == StatusComplex {
		status = 1
	}
	if account.Status == StatusNotfree {
		status = 2
	}

	return fmt.Sprintf("%s,1,%s,2,%s,3,%s,%s,%d,5,%s,6,%s,7,%d,%d,%d",
		account.Email,
		account.Fname,
		account.Sname,
		account.Phone,
		account.Sex,
		account.Birth,
		account.Country,
		account.City,
		status,
		account.Premium.Start,
		account.Premium.Finish,
	)

}

func accountDecode(accountStr string) Account {
	accountSlice := strings.Split(accountStr, ",")
	switch accountSlice[14] {
	case "0":
		accountSlice[14] = StatusFree
	case "1":
		accountSlice[14] = StatusComplex
	case "2":
		accountSlice[14] = StatusNotfree
	}
	account := Account{}
	account.Email = accountSlice[0]
	account.Fname = accountSlice[2]
	account.Sname = accountSlice[4]
	account.Phone = accountSlice[6]
	account.Sex = accountSlice[7]
	account.Birth, _ = strconv.Atoi(accountSlice[8])
	account.Country = accountSlice[10]
	account.City = accountSlice[12]
	account.Status = accountSlice[14]
	start, _ := strconv.Atoi(accountSlice[15])
	finish, _ := strconv.Atoi(accountSlice[16])
	account.Premium = Premium{
		Start:  start,
		Finish: finish,
	}

	return account
}

func prepareAccount(id uint32, properties []string) map[string]interface{} {
	accountObj := accountDecode(accounts[id])
	tmpAccount := make(map[string]interface{})
	tmpAccount["id"] = id
	tmpAccount["email"] = accountObj.Email
	for _, property := range properties {
		switch property {
		case "sex":
			tmpAccount["sex"] = accountObj.Sex
		case "fname":
			tmpAccount["fname"] = accountObj.Fname
		case "sname":
			tmpAccount["sname"] = accountObj.Sname
		case "status":
			tmpAccount["status"] = accountObj.Status
		case "birth":
			tmpAccount["birth"] = accountObj.Birth
		case "phone":
			tmpAccount["phone"] = accountObj.Phone
		case "country":
			tmpAccount["country"] = accountObj.Country
		case "city":
			tmpAccount["city"] = accountObj.City
		case "premium":
			tmpAccount["premium"] = accountObj.Premium
		}
	}

	return tmpAccount
}

func (account *Account) GetEmailDomain() string {
	emailArr := strings.Split(account.Email, "@")
	return emailArr[1]
}

func (account *Account) GetPhoneCode() string {
	phoneArr := strings.Split(account.Phone, "(")
	if len(phoneArr) > 1 {
		phoneArr = strings.Split(phoneArr[1], ")")
		if len(phoneArr) > 1 {
			return phoneArr[1]
		}
	}

	return ""
}

func (account *Account) GetYear() string {
	birthTime := time.Unix(int64(account.Birth), 0)
	return strconv.Itoa(birthTime.Year())
}

func (account *Account) GetPremium() string {

	if int64(account.Premium.Start) <= timestamp && int64(account.Premium.Finish) >= timestamp {
		return PremiumNow
	}

	if account.Premium.Start > 0 && account.Premium.Finish > 0 {
		return PremiumExist
	}

	return PremiumNone
}
