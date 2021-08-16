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
	if session.Stopped == "1000-01-01 00:00:00" {
		l.Stopped = "" +
			"0000-00-00 00:00:00"
	} else {
		l.Stopped = session.Stopped
	}
	l.Name = session.Name
}
