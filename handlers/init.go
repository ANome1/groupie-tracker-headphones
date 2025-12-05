package handlers

import "groupie-tracker/services"

var AuthService *services.AuthService
var RoomService *services.RoomService

func Init(as *services.AuthService, rs *services.RoomService) {
	AuthService = as
	RoomService = rs
}
