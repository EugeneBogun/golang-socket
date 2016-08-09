package main

import (
    "net/http"
    "golang.org/x/net/websocket"
    "strconv"
    "time"
    "fmt"
)

type Message struct {
    Type string `json:"type"`
    Data string `json:"data"`
    SenderId string `json:"sender_id"`
    ReceiverId string `json:"receiver_id"`
}

var clients ActiveClients

func main() {
    clients = InitActiveClients()
    http.HandleFunc("/", renderMainPage)
    http.Handle("/ws", websocket.Handler(handlerClient))

    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        panic("ListenAndServe: " + err.Error())
    }
}

func renderMainPage(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func handlerClient(ws *websocket.Conn) {
    var uid string
    defer func() {
        fmt.Println(uid+" :closed")
        ws.Close()
    }()

    uid = strconv.FormatInt(time.Now().Unix(), 10)
    client := Client{id:uid, in:make(chan string), out:make(chan string)}
    err := clients.AddClient(uid, client)
    if(err != nil) {
        fmt.Println(err)
        return
    }

    fmt.Println(uid+" :opened")

    go Read(ws)
    go Write(ws, client.out)

    for msg := range client.in {
        println("in msg:"+msg)
        client.out <- msg
    }
}

func Read(ws *websocket.Conn)  {
    var data Message
    for {
        err := websocket.JSON.Receive(ws, &data)
        if err != nil {
            return
        }
        client, err  := clients.GetClient(data.ReceiverId)
        if err != nil {
            println("Chanel not found " + data.ReceiverId)
            return
        }
        client.in <- data.Data
    }
}

func Write(ws *websocket.Conn, out chan string) {
    for msg := range out {
        println("out msg:"+msg)
        _, err := ws.Write([]byte(msg))
        if err != nil {
            fmt.Println(err)
            return
        }
    }
}