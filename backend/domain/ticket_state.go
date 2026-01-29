package domain

import "rbac/models"

var ValidTransitions = map[models.TicketStatus][]models.TicketStatus{
	models.TicketCustomerCreated: {
		models.TicketAdminReviewed,
	},
	models.TicketAdminReviewed: {
		models.TicketAssignedSupport,
	},
	models.TicketAssignedSupport: {
		models.TicketResolvedSupport,
	},
	models.TicketResolvedSupport: {
		models.TicketClosedByAdmin,
	},
	models.TicketClosedByAdmin: {
		models.TicketFeedbackGiven,
	},
}

func CanTransition(from, to models.TicketStatus) bool {
	for _, s := range ValidTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}
