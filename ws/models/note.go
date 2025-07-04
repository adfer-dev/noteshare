package models

type Note struct {
	ID        uint   `json:"id" validate:"required"`
	Content   string `json:"content"`
	Color     string `json:"color" validate:"required"`
	Creator   string `json:"creator" validate:"required"`
	RoomID    uint   `json:"room_id" validate:"required"`
	CreatedAt string `json:"created_at" validate:"required"`
	EditedAt  string `json:"edited_at" validate:"required"`
}

type NoteIDMessage struct {
	NoteID uint `json:"note_id" validate:"required"`
	RoomID uint `json:"room_id" validate:"required"`
}
