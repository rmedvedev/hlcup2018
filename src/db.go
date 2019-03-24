package main

import (
	"sort"
	"strconv"
	"strings"
)

type Query struct {
	Predicates []Predicate
}

type Predicate struct {
	Field     string
	Condition string
	Value     string
}

func (q *Query) And(field string, condition string, value string) *Query {
	q.Predicates = append(q.Predicates, Predicate{
		Field:     field,
		Condition: condition,
		Value:     value,
	})
	return q
}

func (q *Query) Exec(properties []string) []map[string]interface{} {
	var accountsResult = []map[string]interface{}{}
	var ids = []uint32{}
	var indexOperations = []IndexOperation{}
	var seqScanFields = []Predicate{}

	limit := 100

	for _, predicate := range q.Predicates {
		switch predicate.Field {
		case "interests":
			interests := strings.Split(predicate.Value, ",")
			switch predicate.Condition {
			case "contains":
				index := IndexOperation{Condition: "AND"}
				for _, interest := range interests {
					index.Indexes = append(index.Indexes, indexes.GetIds("interests", interest))
				}
				indexOperations = append(indexOperations, index)
			case "any":
				index := IndexOperation{Condition: "OR"}
				for _, interest := range interests {
					index.Indexes = append(index.Indexes, indexes.GetIds("interests", interest))
				}
				indexOperations = append(indexOperations, index)
			}
		case "email_domain":
			index := IndexOperation{Condition: "AND"}
			index.Indexes = append(index.Indexes, indexes.GetIds("email_domain", predicate.Value))
			indexOperations = append(indexOperations, index)
		case "fname":
			switch predicate.Condition {
			case "=":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("fname", predicate.Value))
				indexOperations = append(indexOperations, index)
			case "in":
				index := IndexOperation{Condition: "OR"}
				names := strings.Split(predicate.Value, ",")
				for _, name := range names {
					index.Indexes = append(index.Indexes, indexes.GetIds("fname", name))
				}
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "sname":
			switch predicate.Condition {
			case "=":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("sname", predicate.Value))
				indexOperations = append(indexOperations, index)
			case "starts":
				index := IndexOperation{Condition: "OR"}
				for name, ids := range indexes["sname"] {
					if strings.Contains(name, predicate.Value) {
						index.Indexes = append(index.Indexes, ids)
					}
				}
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "country":
			switch predicate.Condition {
			case "=":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("country", predicate.Value))
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "birth":
			switch predicate.Condition {
			case "year":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("birth_year", predicate.Value))
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "city":
			switch predicate.Condition {
			case "=":
				if predicate.Value == "" {
					seqScanFields = append(seqScanFields, predicate)
				} else {
					index := IndexOperation{Condition: "AND"}
					index.Indexes = append(index.Indexes, indexes.GetIds("city", predicate.Value))
					indexOperations = append(indexOperations, index)
				}
			case "in":
				index := IndexOperation{Condition: "OR"}
				cities := strings.Split(predicate.Value, ",")
				for _, city := range cities {
					index.Indexes = append(index.Indexes, indexes.GetIds("city", city))
				}
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "premium":
			switch predicate.Condition {
			case "=":
				switch predicate.Value {
				case "":
					seqScanFields = append(seqScanFields, predicate)
				case "now":
					index := IndexOperation{Condition: "AND"}
					index.Indexes = append(index.Indexes, indexes.GetIds("premium", "1"))
					indexOperations = append(indexOperations, index)
				}
			case "!=":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("premium", "2"))
				indexOperations = append(indexOperations, index)
			}
		case "phone":
			switch predicate.Condition {
			case "=":
				index := IndexOperation{Condition: "AND"}
				index.Indexes = append(index.Indexes, indexes.GetIds("phone", predicate.Value))
				indexOperations = append(indexOperations, index)
			default:
				seqScanFields = append(seqScanFields, predicate)
			}
		case "likes":
			switch predicate.Condition {
			case "contains":
				index := IndexOperation{Condition: "AND"}
				accIds := strings.Split(predicate.Value, ",")
				for _, accId := range accIds {
					accIdInt, _ := strconv.Atoi(accId)
					if _, ok := accounts[uint32(accIdInt)]; !ok {
						return accountsResult
					}
					index.Indexes = append(index.Indexes, indexes.GetIds("likes", accId))
				}
				indexOperations = append(indexOperations, index)
			}
		case "limit":
			limit, _ = strconv.Atoi(predicate.Value)
		default:
			seqScanFields = append(seqScanFields, predicate)
		}
	}

	ids = intersection(limit, seqScanFields, indexOperations)

	for _, id := range ids {
		accountsResult = append(accountsResult, prepareAccount(id, properties))
	}

	return accountsResult
}

func intersection(limit int, seqScanFields []Predicate, indexes []IndexOperation) (result []uint32) {

	if len(indexes) > 0 {
		sort.Slice(indexes, func(i, j int) bool {
			if indexes[i].Condition == indexes[j].Condition {
				return len(indexes[i].Indexes[0]) < len(indexes[j].Indexes[0])
			}
			return indexes[i].Condition < indexes[j].Condition
		})

		keys := make(map[int]map[int]int)
		for id := rangeIndexes(&keys, indexes); id > 0; id = rangeIndexes(&keys, indexes) {
			if limit == 0 {
				break
			}

			if seqScan(id, seqScanFields) {
				result = append(result, id)
				limit--
			}
		}

	} else {
		for _, id := range accountsIds {
			if limit == 0 {
				break
			}

			if seqScan(id, seqScanFields) {
				result = append(result, id)
				limit--
			}
		}
	}

	return result
}

func seqScan(id uint32, predicates []Predicate) bool {
	var accountTmp = Account{}
	var accountStr = accounts[id]

	for _, predicate := range predicates {
		if predicate.Field == "sex" {
			if !strings.Contains(accountStr, ","+predicate.Value+",") {
				return false
			}
		} else if predicate.Field == "email" {
			if predicate.Condition == ">" {
				accountTmp = accountEncodeToDest(accounts[id], accountTmp)
				if !(accountTmp.Email > predicate.Value) {
					return false
				}
			} else if predicate.Condition == "<" {
				accountTmp = accountEncodeToDest(accounts[id], accountTmp)
				if !(accountTmp.Email < predicate.Value) {
					return false
				}
			} else if predicate.Condition == "domain" {
				if !(strings.Contains(accounts[id], predicate.Value)) {
					return false
				}
			}
		} else if predicate.Field == "status" {
			if predicate.Condition == "=" {
				if !strings.Contains(accountStr, "7,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accountStr, "7,"+predicate.Value+",") {
					return false
				}
			}
		} else if predicate.Field == "fname" {
			if predicate.Condition == "=" {
				if !strings.Contains(accounts[id], "1,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accounts[id], "1,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "in" {
				if accountTmp.Fname == "" || !strings.Contains(predicate.Value, accountTmp.Fname) {
					return false
				}
			}
		} else if predicate.Field == "sname" {
			if predicate.Condition == "=" {
				if !strings.Contains(accounts[id], "2,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accounts[id], "2,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "starts" {
				if !strings.Contains(accountStr, ","+predicate.Value) {
					return false
				}
			}
		} else if predicate.Field == "phone" {
			if predicate.Condition == "=" {
				if predicate.Value == "" {
					if !strings.Contains(accounts[id], "3,"+predicate.Value+",") {
						return false
					}
				} else {
					if !strings.Contains(accounts[id], "("+predicate.Value+")") {
						return false
					}
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accounts[id], "3,"+predicate.Value+",") {
					return false
				}
			}
		} else if predicate.Field == "country" {
			if predicate.Condition == "=" {
				if !strings.Contains(accounts[id], "5,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accounts[id], "5,"+predicate.Value+",") {
					return false
				}
			}
		} else if predicate.Field == "city" {
			if predicate.Condition == "=" {
				if !strings.Contains(accounts[id], "6,"+predicate.Value+",") {
					return false
				}
			} else if predicate.Condition == "!=" {
				if strings.Contains(accounts[id], "6,"+predicate.Value+",") {
					return false
				}
			}
		} else if predicate.Field == "birth" {
			val, _ := strconv.Atoi(predicate.Value)
			accountTmp = accountEncodeToDest(accounts[id], accountTmp)
			if predicate.Condition == ">" {
				if !(accountTmp.Birth > val) {
					return false
				}
			} else if predicate.Condition == "<" {
				if !(accountTmp.Birth < val) {
					return false
				}
			}
		} else if predicate.Field == "premium" {
			accountTmp = accountEncodeToDest(accounts[id], accountTmp)
			if predicate.Condition == "=" {
				if !(accountTmp.Premium.Start == 0) {
					return false
				}
			}
		}
	}
	return true
}

func accountEncodeToDest(accountStr string, account Account) Account {
	if account.ID == 0 {
		account = accountDecode(accountStr)
	}

	return account
}
