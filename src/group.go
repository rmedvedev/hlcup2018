package main

import (
	"strings"
)

type ResultGroups struct {
	Groups []map[string]interface{} `json:"groups"`
}

func getGroups(filters map[string]string) ResultGroups {
	query := Query{}

	for name, _ := range filters {
		predicate := strings.Split(name, "_")
		switch predicate[0] {
		case "sex":
		case "email":

		case "status":

		case "fname":

		case "sname":

		case "phone":

		case "country":

		case "city":

		case "birth":

		case "premium":

		case "interests":

		case "likes":

		}
	}

	query.And("limit", "=", filters["limit"])

	var result = ResultGroups{Groups: []map[string]interface{}{}}
	return result
}
