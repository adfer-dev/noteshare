package services

import (
	"errors"
	"noteshare-api/database"
	"noteshare-api/models"
	"noteshare-api/storage"

	"github.com/google/uuid"
)

type RoomBody struct {
	Name    string `json:"name" validate:"required,max=20"`
	Creator uint   `json:"creator" validate:"required,number"`
}

type UpdateRoomBody struct {
	Name string `json:"name" validate:"max=20"`
}

var roomStorage storage.Storage = &storage.RoomStorage{}

func GetAllRooms() ([]*models.Room, error) {
	var rooms []*models.Room

	database := database.GetInstance().GetDB()
	results, err := database.Query("SELECT * FROM rooms;")

	if err != nil {
		return rooms, err
	}
	defer results.Close()

	for results.Next() {
		room, scanErr := roomStorage.Scan(results)

		if scanErr != nil {
			return nil, scanErr
		}

		rooms = append(rooms, room.(*models.Room))
	}

	return rooms, nil
}

func GetRoomById(id int) (*models.Room, error) {
	room, err := roomStorage.Get(id)

	if err != nil {
		return nil, err
	}

	return room.(*models.Room), nil
}

func GetRoomByInvite(inviteCode string) (*models.Room, error) {
	database := database.GetInstance().GetDB()

	result, err := database.Query("SELECT * FROM rooms WHERE invite LIKE ? ;", inviteCode)

	if err != nil {
		return nil, err
	}

	if result.Next() {
		room, scanErr := roomStorage.Scan(result)

		if scanErr != nil {
			return nil, scanErr
		}

		return room.(*models.Room), nil
	} else {
		return nil, errors.New(storage.RoomNotFoundErr)
	}
}

func GetUserRooms(userId int) ([]*models.Room, error) {
	var rooms []*models.Room

	database := database.GetInstance().GetDB()
	results, err := database.Query("SELECT rooms.* FROM rooms "+
		"JOIN users_rooms ON users_rooms.room_id = rooms.id WHERE users_rooms.user_id = ?;", userId)

	if err != nil {
		return rooms, err
	}
	defer results.Close()

	found := false
	for results.Next() {
		room, scanErr := roomStorage.Scan(results)

		if scanErr != nil {
			return nil, scanErr
		}

		rooms = append(rooms, room.(*models.Room))
		found = true
	}

	if !found {
		return nil, errors.New("user has no rooms")
	}

	return rooms, nil
}

func CreateRoom(roomBody RoomBody) (*models.Room, error) {

	room := &models.Room{
		Name:    roomBody.Name,
		Invite:  uuid.NewString(),
		Creator: roomBody.Creator,
	}

	if err := roomStorage.Create(room); err != nil {
		return nil, err
	}

	return room, nil
}

func UpdateRoom(id int, roomBody UpdateRoomBody) (*models.Room, error) {
	room, queryErr := GetRoomById(id)

	if queryErr != nil {
		return nil, queryErr
	}

	if roomBody.Name != "" {
		room.Name = roomBody.Name
	}

	updateErr := roomStorage.Update(room)

	return room, updateErr
}

func DeleteRoom(id int) error {
	room, queryErr := GetRoomById(id)

	if queryErr != nil {
		return queryErr
	}

	return roomStorage.Delete(room)
}
