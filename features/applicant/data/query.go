package data

import (
	"context"
	"errors"
	"strings"

	core "firlydev/backend-test/features/applicant"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("applicant not found")

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Create(ctx context.Context, item core.Applicant) (core.Applicant, error) {
	m := Applicant{Name: item.Name, Email: item.Email, Position: item.Position, Stage: string(item.Stage)}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return core.Applicant{}, err
	}
	return toEntity(m), nil
}
func (r *Repository) List(ctx context.Context, f core.ListFilter) (core.PageResult, error) {
	q := r.db.WithContext(ctx).Model(&Applicant{})
	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("name LIKE ? OR email LIKE ? OR position LIKE ?", like, like, like)
	}
	if f.Stage != "" {
		q = q.Where("stage = ?", strings.ToUpper(f.Stage))
	}
	if f.Position != "" {
		q = q.Where("LOWER(position) LIKE ?", "%"+strings.ToLower(f.Position)+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return core.PageResult{}, err
	}
	var models []Applicant
	offset := (f.Page - 1) * f.Limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(f.Limit).Find(&models).Error; err != nil {
		return core.PageResult{}, err
	}
	items := make([]core.Applicant, 0, len(models))
	for _, m := range models {
		items = append(items, toEntity(m))
	}
	return core.PageResult{Items: items, Page: int64(f.Page), Limit: int64(f.Limit), Total: total}, nil
}
func (r *Repository) GetByID(ctx context.Context, id uint) (core.Applicant, error) {
	var m Applicant
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.Applicant{}, ErrNotFound
	}
	if err != nil {
		return core.Applicant{}, err
	}
	return toEntity(m), nil
}
func (r *Repository) UpdateStage(ctx context.Context, id uint, stage core.Stage) (core.Applicant, error) {
	var m Applicant
	err := r.db.WithContext(ctx).First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.Applicant{}, ErrNotFound
	}
	if err != nil {
		return core.Applicant{}, err
	}
	if err := r.db.WithContext(ctx).Model(&m).Update("stage", string(stage)).Error; err != nil {
		return core.Applicant{}, err
	}
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return core.Applicant{}, err
	}
	return toEntity(m), nil
}
func (r *Repository) Summary(ctx context.Context) (core.Summary, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&Applicant{}).Count(&total).Error; err != nil {
		return core.Summary{}, err
	}
	type row struct {
		Stage string
		Count int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).Model(&Applicant{}).Select("stage, COUNT(*) as count").Group("stage").Order("stage ASC").Scan(&rows).Error; err != nil {
		return core.Summary{}, err
	}
	by := make([]core.StageSummary, 0, len(rows))
	for _, v := range rows {
		by = append(by, core.StageSummary{Stage: core.Stage(v.Stage), Count: v.Count})
	}
	return core.Summary{Total: total, ByStage: by}, nil
}
func toEntity(m Applicant) core.Applicant {
	return core.Applicant{ID: m.ID, Name: m.Name, Email: m.Email, Position: m.Position, Stage: core.Stage(m.Stage), CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}
