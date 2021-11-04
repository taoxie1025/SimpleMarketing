package models

import "time"

const (
	TicketTypeGeneral        = iota
	TicketTypePayment        = iota
	TicketTypeBugReport      = iota
	TicketTypeApiSupport     = iota
	TicketTypeFeatureRequest = iota
)

const (
	TicketStatusNone = iota
	TicketStatusOpen
	TicketStatusClose
	TicketStatusResolved
	TicketStatusDelete
)

type Ticket struct {
	Email        string    `json:"email"`    // primary key
	TicketId     string    `json:"ticketId"` // sort key
	ProjectId    string    `json:"projectId"`
	ProjectName  string    `json:"projectName"`
	Name         string    `json:"name"`
	CreatedAt    int64     `json:"createdAt"`
	TicketStatus int       `json:"ticketStatus"`
	TicketType   int       `json:"ticketType"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Comments     []Comment `json:"comments"`
}

func NewTicket(email, projectId, projectName, name, title, body string, ticketType int) *Ticket {
	return &Ticket{
		Email:        email,
		ProjectId:    projectId,
		ProjectName:  projectName,
		Name:         name,
		Title:        title,
		Body:         body,
		Comments:     []Comment{},
		TicketStatus: TicketStatusOpen,
		TicketType:   ticketType,
		CreatedAt:    time.Now().UnixNano() / int64(time.Millisecond),
		TicketId:     GenerateTicketUUID(),
	}
}

type Comment struct {
	CommendId string `json:"commendId"`
	TicketId  string `json:"ticketId"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"createdAt"`
	Body      string `json:"body"`
}

func NewComment(email, ticketId, name, body string) *Comment {
	return &Comment{
		CommendId: GenerateCommentUUID(),
		TicketId:  ticketId,
		Email:     email,
		Name:      name,
		CreatedAt: time.Now().UnixNano() / int64(time.Millisecond),
		Body:      body,
	}
}
