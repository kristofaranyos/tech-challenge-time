package dto

import "github.com/kristofaranyos/tech-challenge-time/repository"

type StopRequest struct {
	SessionId string `json:"sessionId"`
}

type NameUpdateRequest struct {
	SessionId string `json:"sessionId"`
	Name      string `json:"name"`
}

type ListResponse struct {
	Id      string `json:"sessionId"`
	Created string `json:"createdAt"`
	Stopped string `json:"stoppedAt"`
	Name    string `json:"name"`
}

func (l *ListResponse) Fill(session *repository.Session) {
	l.Id = session.Id
	l.Created = session.Created
	l.Stopped = session.Stopped
	l.Name = session.Name
}
