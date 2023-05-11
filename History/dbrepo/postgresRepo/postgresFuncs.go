package postgresrepo

import (
	"context"
	"time"

	"github.com/sobhankazemi/GroupChat/History/models"
)

func (repo *Repository) SaveChatHistory(message models.Message) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into history (user_id,username,message,send_time , room_id) 
			 values($1,$2,$3,$4 , $5)`
	time, _ := time.Parse("2006-01-02 15:04", message.Time)
	_, err := repo.db.ExecContext(ctx, query,
		message.UserID,
		message.UserName,
		message.Message,
		time,
		message.Room_id,
	)

	return err == nil
}
func (repo *Repository) GetHistory(room_id, page int) ([]models.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	pageSize := 10
	query := `select user_id , username , message , send_time , room_id from history where room_id = $1 offset $2 rows fetch next $3 rows only`
	rows, err := repo.db.QueryContext(ctx, query, room_id, (page-1)*pageSize, pageSize)
	result := make([]models.Message, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.Message{}
		rows.Scan(&temp.UserID, &temp.UserName, &temp.Message, &temp.Time, &temp.Room_id)
		result = append(result, temp)
	}
	return result, nil
}
