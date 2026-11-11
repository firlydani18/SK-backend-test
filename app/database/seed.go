package database

import (
	"firlydev/backend-test/features/applicant/data"

	"gorm.io/gorm"
)

// SeedDemoData inserts a small, deterministic set of mock applicants.
// FirstOrCreate makes seeding safe to run every time the application starts.
func SeedDemoData(db *gorm.DB) error {
	seed := []data.Applicant{
		{Name: "Alya Putri", Email: "alya.putri@example.com", Position: "Backend Engineer", Stage: "APPLIED"},
		{Name: "Bima Pratama", Email: "bima.pratama@example.com", Position: "Frontend Developer", Stage: "SCREENING"},
		{Name: "Citra Maharani", Email: "citra.maharani@example.com", Position: "QA Engineer", Stage: "INTERVIEW"},
		{Name: "Dimas Saputra", Email: "dimas.saputra@example.com", Position: "DevOps Engineer", Stage: "OFFER"},
		{Name: "Nadia Lestari", Email: "nadia.lestari@example.com", Position: "Product Designer", Stage: "REJECTED"},
	}

	for i := range seed {
		if err := db.Where("email = ?", seed[i].Email).FirstOrCreate(&seed[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
