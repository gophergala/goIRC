package main

import (
	"net"
)

type ConnectionStatus int

const (
	UserConnected ConnectionStatus = iota
	UserRegistered
)

type Client struct {
	Status     ConnectionStatus
	connection net.Conn
}
