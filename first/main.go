package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var secureCode string = "gpq74gpq"

func main() {

	log.SetLevel(log.DebugLevel)
	serverStart()

}

// запуск сервера
func serverStart() {

	router := mux.NewRouter()
	router.HandleFunc("/send", sendRoute).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8111",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	srv.ListenAndServe()
}

// подключение пути для обращения к серверу. GET, принимает имя
func sendRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var user showUser
	user.Name = r.FormValue("name")
	dataEnrich(&user)

	// передача данных на второй сервер.
	log.Info("send data to second server")
	sendToSecond(&user)

}

// обогащение полученых данных
func dataEnrich(curInstance *showUser) {
	curInstance.Secure = secureCode
	curInstance.Gender = "male"

	rand.Seed(time.Now().UnixNano())
	curInstance.Age = rand.Intn(100) + 1
}

// отправка данных на второй сервер
func sendToSecond(curInstance *showUser) {

	baseURL, _ := url.Parse("http://localhost:8112/getter")

	data := url.Values{
		"name":     {curInstance.Name},
		"age":      {strconv.Itoa(curInstance.Age)},
		"gender":   {curInstance.Gender},
		"security": {curInstance.Secure},
	}

	resp, _ := http.PostForm(baseURL.String(), data)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

// структура, с помощью которой идёт обмен внутри сервисов, и потом в базу
type showUser struct {
	Age    int    `json:"age" example:"25" format:"int64"`
	Name   string `json:"name" example:"ivan"`
	Gender string `json:"gender" example:"male"`
	Secure string `json:"secure" example:"gp46gp"`
}
