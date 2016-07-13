package main

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
)

var upgrader = websocket.Upgrader{}
var countConnections = 0

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "index.html")
    })

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        var conn, _ = upgrader.Upgrade(w, r, nil)
        countConnections++
        println("Connected. Count ",countConnections)

        go func(conn *websocket.Conn) {
            defer func(){
                countConnections--
                println("Diconnected. Count ",countConnections)
                conn.Close()
            } ()

            for {
                mType, msg, err := conn.ReadMessage()
                if err != nil {
                    return
                }
                conn.WriteMessage(mType, msg)
                println(string(msg))
            }
        }(conn)
    })

    err := http.ListenAndServe(":3000", nil) // setting listening port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}