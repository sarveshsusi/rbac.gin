package repository

import "gorm.io/gorm"

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) FetchAdminStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	var usersCount int64
	var ticketsCount int64
	var pendingTickets int64
	var closedTickets int64

	if err := r.db.
		Table("users").
		Where("is_active = true").
		Count(&usersCount).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Table("tickets").
		Count(&ticketsCount).Error; err != nil {
		return nil, err
	}

	// Pending = customer created OR admin reviewed
	if err := r.db.
		Table("tickets").
		Where("status IN ?", []string{
			"customer_created",
			"admin_reviewed",
		}).
		Count(&pendingTickets).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Table("tickets").
		Where("status = ?", "closed_by_admin").
		Count(&closedTickets).Error; err != nil {
		return nil, err
	}

	stats["users"] = usersCount
	stats["tickets"] = ticketsCount
	stats["pending_tickets"] = pendingTickets
	stats["closed_tickets"] = closedTickets

	return stats, nil
}
