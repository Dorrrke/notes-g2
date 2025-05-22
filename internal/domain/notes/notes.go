package notes

import "time"

type Note struct {
	NID        string    `json:"nid"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Status     Status    `json:"status"`
	Created_at time.Time `json:"created_at"`
	UID        string    `json:"uid"`
}

type NoteResponseFormat struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	Created_at string `json:"created_at"`
	UID        string `json:"uid"`
}

func NoteResponse(note Note) NoteResponseFormat {
	return NoteResponseFormat{
		Title:      note.Title,
		Content:    note.Content,
		Status:     note.Status.String(),
		Created_at: note.Created_at.Format(time.RFC3339),
		UID:        note.UID,
	}
}
