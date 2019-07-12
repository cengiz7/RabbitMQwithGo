package HttpJobs

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"../MQJobs"
)

var QueueName string
var Connection string
var boolCh chan bool
var strCh chan string

func SetHandlers(queueName, connection string, indexChannel chan bool, messageCh chan string){
	QueueName  = queueName
	Connection = connection
	boolCh     = indexChannel
	strCh  	   = messageCh
	http.HandleFunc(`/messages/send`,messageHandler)
	http.HandleFunc(`/messages`,index)
	err := http.ListenAndServe( port(),nil); if err != nil {
		log.Fatal("Failed to listen and serve from port "+port())
	}
}
func messageHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}
		fmt.Println(string(body))
		priority , err := strconv.Atoi(r.URL.Query()["priority"][0]); if err != nil{
			log.Println("Couldn't convert priority string to integer")
		}
		MQJobs.SentToQueue(body,QueueName,Connection,uint8(priority))
		fmt.Fprint(w, "ok")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
type QueueStatus struct {
	LastMessage string
}

func index(w http.ResponseWriter, r *http.Request) {
	boolCh <- true

	queueStatus := QueueStatus{<-strCh}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, queueStatus); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func Post() {
	for count := 0; count <= 500 ; count++{
		priority := count % 3
		message := "Hello mello otello "+ strconv.Itoa(count) + "\nPriority is "+ strconv.Itoa(priority)

		req, _ := http.NewRequest("POST", "http://localhost"+ port() +"/messages/send?priority="+strconv.Itoa(priority),
			bytes.NewBuffer([]byte(message)))

		client := &http.Client{}
		_, err := client.Do(req)
		if err != nil {
			panic(err)
		}
	}
}

func port() string{
	port := os.Getenv("PORT")
	if len(port) == 0{
		port = "8080"
	}
	return ":" +port
}