package service

import (
	"log"
	"time"

	"rbac/repository"
	"rbac/utils"
)

type EscalationService struct {
	repo   *repository.EscalationRepository
	mailer *utils.Mailer
}

func NewEscalationService(
	repo *repository.EscalationRepository,
	mailer *utils.Mailer,
) *EscalationService {
	return &EscalationService{
		repo:   repo,
		mailer: mailer,
	}
}

/* =====================
   RUN ESCALATION CHECK
===================== */
func (s *EscalationService) Run() {

	cutoff := time.Now().Add(-7 * 24 * time.Hour)

	tickets, err := s.repo.FindOverdueTickets(cutoff)
	if err != nil {
		log.Println("❌ escalation fetch failed:", err)
		return
	}

	for _, t := range tickets {

		body := `
			<h2>🚨 Ticket Escalation Alert</h2>
			<p><b>Ticket ID:</b> ` + t.ID.String() + `</p>
			<p><b>Title:</b> ` + t.Title + `</p>
			<p><b>Status:</b> ` + string(t.Status) + `</p>
			<p>This ticket has been open for more than 7 days.</p>
		`

		_ = s.mailer.Send(
			"emerd@gmail.com",
			"🚨 Ticket Escalation Alert",
			body,
		)

		_ = s.mailer.Send(
			"veemerd@gmail.com",
			"🚨 Ticket Escalation Alert",
			body,
		)

		_ = s.repo.MarkEscalated(t.ID)
	}
}
