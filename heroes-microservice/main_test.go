package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	hero "github.com/tgl-dogg/golang-microservice-play/heroes-data"
)

func Test_GetRaces_OK(t *testing.T) {
	r := gin.New()
	r.GET("/", getRaces)
	resp := emulateRequest(r, "/")

	if resp.Code != http.StatusOK {
		t.Error("HTTP request status code error.")
	}

	var races []hero.Race
	err := json.NewDecoder(resp.Body).Decode(&races)
	if err != nil {
		t.Fatal(err)
	}

	if len(races) < 1 {
		t.Error("Invalid records found:", races)
	}
}

func Test_GetRaceByID_OK(t *testing.T) {
	r := gin.New()
	r.GET("/:id", getRaceByID)
	resp := emulateRequest(r, "/1")

	if resp.Code != http.StatusOK {
		t.Error("HTTP request status code error.")
	}

	var race hero.Race
	err := json.NewDecoder(resp.Body).Decode(&race)
	if err != nil {
		t.Fatal(err)
	}

	if race.Name != "Human" {
		t.Error("Invalid record found:", race)
	}
}

func Test_GetRaceByID_NOTFOUND(t *testing.T) {
	r := gin.New()
	r.GET("/:id", getRaceByID)
	resp := emulateRequest(r, "/-1")

	if resp.Code != http.StatusNotFound {
		t.Error("HTTP request status code error.")
	}
}

func Test_GetRaceByID_IDINVALID(t *testing.T) {
	invalidID := "98a11010-d019-11ec-9d64-0242ac120002"

	r := gin.New()
	r.GET("/:id", getRaceByID)
	resp := emulateRequest(r, "/"+invalidID)

	if resp.Code != http.StatusBadRequest {
		t.Error("HTTP request status code error")
	}

	body := resp.Body.String()

	if !strings.Contains(body, invalidID) {
		t.Error("Invalid response error:", body)
	}
}

func Test_GetRaceByRecommendedClasses_OK(t *testing.T) {
	r := gin.New()
	r.GET("/mock", getRacesByRecommendedClasses)
	resp := emulateRequest(r, "/mock?classes=thi&classes=war")

	var races []hero.Race
	err := json.NewDecoder(resp.Body).Decode(&races)
	if err != nil {
		t.Fatal(err, resp.Body)
	}

	if resp.Code != http.StatusOK {
		t.Error("HTTP request status code error.")
	}

	if !isRacePresent(races, "Human") {
		t.Error("Expected Human race to be found.")
	} else if !isRacePresent(races, "Elf") {
		t.Error("Expected Elf race to be found.")
	} else if !isRacePresent(races, "Dwarf") {
		t.Error("Expected Dwarf race to be found.")
	}
}

func Test_GetRaceByRecommendedClasses_EMPTY(t *testing.T) {
	r := gin.New()
	r.GET("/mock", getRacesByRecommendedClasses)
	resp := emulateRequest(r, "/mock?classes=biruleibes")

	var races []hero.Race
	err := json.NewDecoder(resp.Body).Decode(&races)
	if err != nil {
		t.Fatal(err, resp.Body)
	}

	if resp.Code != http.StatusOK {
		t.Error("HTTP request status code error.")
	}

	if len(races) > 0 {
		t.Error("Expected to found empty races set:", races)
	}
}

func isRacePresent(races []hero.Race, name string) bool {
	for i := range races {
		if races[i].Name == name {
			return true
		}
	}

	return false
}

func emulateRequest(r *gin.Engine, url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
