package database

import (
	"log"
	"rbac/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.TwoFAOTP{},
		&models.RememberedDevice{},
		&models.Customer{},
		&models.SupportEngineer{},
		&models.Brand{},
		&models.Category{},
		&models.Model{},
		&models.Product{},
		&models.CustomerProduct{},
		&models.AMCContract{},
		&models.AMCSchedule{},
		&models.Ticket{},
		&models.TicketAssignment{},
		&models.TicketStatusHistory{},
		&models.TicketComment{},
		&models.TicketAttachment{},
		&models.TicketFeedback{},
		&models.ServiceVisit{},
		&models.GPSLog{},
		&models.DigitalSignature{},
		&models.AuditLog{},
		&models.EscalationRule{},
		&models.TicketEscalation{},
	)

	if err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	log.Println("✅ Database migration completed successfully")
}
