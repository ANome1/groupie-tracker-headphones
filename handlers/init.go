package handlers

import "groupie-tracker/services"

// Variables globales pour accéder aux services depuis les handlers
var AuthService *services.AuthService
var RoomService *services.RoomService

// Init initialise les services pour les handlers
func Init(as *services.AuthService, rs *services.RoomService) {
	AuthService = as
	RoomService = rs
}
