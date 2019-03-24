package main

import (
	"strings"
)

type Result struct {
	Accounts []map[string]interface{} `json:"accounts"`
}

var disableProperties = map[string]bool{
	"limit":     true,
	"interests": true,
	"likes":     true,
	"email":     true,
}

func getAccountsByFilter(filters map[string]string) Result {
	query := Query{}
	properties := []string{}

	for name, value := range filters {
		predicate := strings.Split(name, "_")
		switch predicate[0] {
		case "sex":
			query.And("sex", "=", value)
		case "email":
			switch predicate[1] {
			case "domain":
				query.And("email", "domain", value)
			case "lt":
				query.And("email", "<", value)
			case "gt":
				query.And("email", ">", value)
			}

		case "status":
			status := "0"
			if value == StatusComplex {
				status = "1"
			}
			if value == StatusNotfree {
				status = "2"
			}

			switch predicate[1] {
			case "eq":
				query.And("status", "=", status)
			case "neq":
				query.And("status", "!=", status)
			}

		case "fname":
			switch predicate[1] {
			case "eq":
				query.And("fname", "=", value)
			case "any":
				query.And("fname", "in", value)
			case "null":
				if value == "0" {
					query.And("fname", "!=", "")
				} else {
					query.And("fname", "=", "")
				}
			}

		case "sname":
			switch predicate[1] {
			case "eq":
				query.And("sname", "=", value)
			case "starts":
				query.And("sname", "starts", value)
			case "null":
				if value == "0" {
					query.And("sname", "!=", "")
				} else {
					query.And("sname", "=", "")
				}
			}

		case "phone":
			switch predicate[1] {
			case "code":
				query.And("phone", "=", value)
			case "null":
				if value == "0" {
					query.And("phone", "!=", "")
				} else {
					query.And("phone", "=", "")
				}
			}

		case "country":
			switch predicate[1] {
			case "eq":
				query.And("country", "=", value)
			case "null":
				if value == "0" {
					query.And("country", "!=", "")
				} else {
					query.And("country", "=", "")
				}
			}

		case "city":
			switch predicate[1] {
			case "eq":
				query.And("city", "=", value)
			case "any":
				query.And("city", "in", value)
			case "null":
				if value == "0" {
					query.And("city", "!=", "")
				} else {
					query.And("city", "=", "")
				}
			}

		case "birth":
			switch predicate[1] {
			case "lt":
				query.And("birth", "<", value)
			case "gt":
				query.And("birth", ">", value)
			case "year":
				query.And("birth", "year", value)
			}

		case "premium":
			switch predicate[1] {
			case "null":
				if value == "0" {
					query.And("premium", "!=", "")
				} else {
					query.And("premium", "=", "")
				}
			case "now":
				query.And("premium", "=", "now")
			}
		case "interests":
			switch predicate[1] {
			case "contains":
				query.And("interests", "contains", value)
			case "any":
				query.And("interests", "any", value)
			}
		case "likes":
			switch predicate[1] {
			case "contains":
				query.And("likes", "contains", value)
			}
		}

		if !disableProperties[predicate[0]] {
			if !(predicate[1] == "null" && value == "1") {
				properties = append(properties, predicate[0])
			}
		}
	}

	query.And("limit", "=", filters["limit"])

	var result = Result{Accounts: query.Exec(properties)}

	return result
}
