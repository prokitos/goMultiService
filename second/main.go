package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

var secureCode string = "gpq74gpq"

func main() {

	log.SetLevel(log.DebugLevel)
	migrateStart()
	serverStart()

}

// запуск сервера
func serverStart() {

	router := mux.NewRouter()
	router.HandleFunc("/getter", getterRoute).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8112",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	srv.ListenAndServe()
}

// подключение пути для получения данных с 1 сервера
func getterRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var user showUser
	user.Name = r.FormValue("name")
	user.Age, _ = strconv.Atoi(r.FormValue("age"))
	user.Gender = r.FormValue("gender")
	user.Secure = r.FormValue("security")

	if user.Secure != secureCode {
		log.Info("wrong secure code !!!")
	} else {
		log.Info(user.Name, " ", user.Gender, " ", user.Age)
	}

	// передача на сервер
	InsertDocument(&user)
	log.Info("success add to server")

}

// миграция
func migrateStart() {

	duration := time.Second * 5
	time.Sleep(duration)

	db := ConnectToDb("postgress.env")
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Error("migration connection error")
		panic(err)
	}
}

// соединение с базой данных
func ConnectToDb(path string) *sql.DB {

	log.Info("connecting to the database")

	godotenv.Load(path)

	envUser := os.Getenv("User")
	envPass := os.Getenv("Pass")
	envHost := os.Getenv("Host")
	envPort := os.Getenv("Port")
	envName := os.Getenv("Name")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", envUser, envPass, envHost, envPort, envName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("database connection error")
		log.Debug("there is not connection with database")
	}

	db.Begin()

	return db
}

// добавление данных
func InsertDocument(DBmodel *showUser) {

	db := ConnectToDb("postgress.env")
	defer db.Close()

	_, err := db.Exec("insert into microInteract (name, age, gender) values ($1, $2, $3)", DBmodel.Name, DBmodel.Age, DBmodel.Gender)
	if err != nil {
		panic(err)
	}

}

// структура, с помощью которой идёт обмен внутри сервисов, и потом в базу
type showUser struct {
	Age    int    `json:"age" example:"25" format:"int64"`
	Name   string `json:"name" example:"ivan"`
	Gender string `json:"gender" example:"male"`
	Secure string `json:"secure" example:"gp46gp"`
}
