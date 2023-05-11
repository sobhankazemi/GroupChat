package dbrepo

import "github.com/sobhankazemi/GroupChat/History/models"

type Repository interface {
	SaveChatHistory(models.Message) bool
	GetHistory (room_id int ,page int) ([]models.Message, error)
}
