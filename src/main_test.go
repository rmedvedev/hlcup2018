package main

import (
	"testing"
)

func BenchmarkStart(b *testing.B) {
	// PrepareData("../../data.zip")
	PrepareData("../data/data/data.zip")
	PrintMemUsage()
}

func BenchmarkFnameEqSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("fname", "=", "Алексей")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"id", "email", "sex", "fname"})
	}
}

func BenchmarkFnameAnySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("fname", "any", "Алексей,Виктор")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"id", "email", "sex", "fname"})
	}
}

func BenchmarkFnameNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("fname", "=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"id", "email", "sex"})
	}
}

func BenchmarkFnameNotNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("fname", "!=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkInterestsAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("interests", "any", "Металлика,Туфли,Горы,Знакомство")
		query.And("sex", "=", "m")
		query.And("status", "=", "1")
		query.And("limit", "=", "50")
		query.Exec([]string{"id", "email", "sex", "status"})
	}
}

func BenchmarkInterestsContains(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("interests", "contains", "Целоваться,Честность,Боевые искусства")
		query.And("sex", "=", "m")
		query.And("status", "=", "1")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkSnameStarts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("sname", "starts", "Дана")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkSnameEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("sname", "=", "Луклентина")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkPhoneNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("phone", "=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkPhoneCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("phone", "=", "953")
		query.And("sex", "=", "f")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCountryEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("country", "=", "Гератрис")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCountryNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("country", "=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCountryNotNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("country", "!=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCityEq(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("city", "=", "Великодам")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCityAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("city", "in", "Великодам,Роттератск")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCityNullSnameStarts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("city", "=", "")
		query.And("sname", "starts", "Кол")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkCityNotNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("city", "!=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkBirthYear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("birth", "year", "1979")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkPremiumNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("premium", "=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkPremiumNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("premium", "=", "now")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex", "premium"})
	}
}

func BenchmarkPremiumNotNull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("premium", "!=", "")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkLikesContainsOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("likes", "contains", "9121")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}

func BenchmarkLikesContainsMany(b *testing.B) {
	for i := 0; i < b.N; i++ {
		query := Query{}
		query.And("likes", "contains", "10181,14951,17055")
		query.And("sex", "=", "m")
		query.And("limit", "=", "50")
		query.Exec([]string{"sex"})
	}
}
