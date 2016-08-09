package main

import (
    "sync"
    "errors"
)

type Client struct {
    id string
    in chan string
    out chan string
}

type ActiveClients struct {
    mutex *sync.RWMutex
    data  map[string]Client
}

var (
    ErrorAlreadyExist = errors.New("User already exists!")
    ErrorNotFound     = errors.New("There is no such user in the storage!")
)

func InitActiveClients() ActiveClients {
    return ActiveClients{data: make(map[string]Client), mutex: new(sync.RWMutex)}
}

func (m ActiveClients) GetClient(id string) (Client, error) {
    m.mutex.RLock()
    val, ok := m.data[id]
    m.mutex.RUnlock()
    if ok {
        return val, nil
    }
    return val, ErrorNotFound

}
func (m ActiveClients) AddClient(id string, client Client) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    _, exists := m.data[id]
    if exists {
        return ErrorAlreadyExist
    }
    m.data[id] = client
    return nil

}

func (m ActiveClients) RemoveClient(id string) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    _, exists := m.data[id]
    if !exists {
        return ErrorNotFound
    }
    delete(m.data, id)
    return nil

}
