package repository

import (
	"context"
	"crypsis-backend/internal/entity"
	"crypsis-backend/internal/model"
	"errors"
	"fmt"
	"log"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"gorm.io/gorm"
)

type adminRepository struct {
	cache    *lru.Cache
	adminIDs map[string]string
	mu       sync.RWMutex
	db       *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	cache, err := lru.New(10) // Cache up to 100 admin
	if err != nil {
		log.Fatal(err)
	}
	return &adminRepository{db: db, cache: cache}
}

func (r *adminRepository) Create(ctx context.Context, admin *entity.Admins) error {
	var existing entity.Admins

	err := r.db.WithContext(ctx).
		Where("username = ?", admin.Username).
		First(&existing).Error

	if err == nil {
		return model.AdminErrAlreadyExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}
	_ = r.LoadAdminIDs(ctx)
	return nil
}

func (r *adminRepository) GetByUsername(ctx context.Context, username string) (*entity.Admins, error) {
	var admin *entity.Admins
	if err := r.db.WithContext(ctx).First(&admin, "username = ?", username).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}
	return admin, nil
}

func (r *adminRepository) GetByID(ctx context.Context, id string) (*entity.Admins, error) {
	var admin *entity.Admins
	if err := r.db.WithContext(ctx).First(&admin, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin by id: %w", err)
	}
	return admin, nil
}

func (r *adminRepository) GetByClientID(ctx context.Context, id string) (*entity.Admins, error) {
	var admin *entity.Admins
	if err := r.db.WithContext(ctx).First(&admin, "client_id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin by id: %w", err)
	}
	return admin, nil
}

func (r *adminRepository) GetList(ctx context.Context, offset, limit int, orderBy, sort string) (*[]entity.Admins, error) {
	if orderBy == "" {
		orderBy = "created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	allowedOrderFields := map[string]bool{
		"id":         true,
		"username":   true,
		"client_id":  true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedOrderFields[orderBy] {
		orderBy = "created_at"
	}

	var admins *[]entity.Admins
	if err := r.db.WithContext(ctx).
		Order(fmt.Sprintf("%s %s", orderBy, sort)).
		Offset(offset).
		Limit(limit).
		Find(&admins).Error; err != nil {
		return nil, fmt.Errorf("failed to get list of admins: %w", err)
	}
	return admins, nil
}

func (r *adminRepository) Update(ctx context.Context, admin *entity.Admins) error {
	if err := r.db.WithContext(ctx).Save(admin).Error; err != nil {
		return fmt.Errorf("failed to update admin: %w", err)
	}
	return nil
}

func (r *adminRepository) Delete(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Admins{}).Error; err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}
	return nil
}

func (r *adminRepository) LoadAdminIDs(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var admins []entity.Admins
	if err := r.db.WithContext(ctx).Find(&admins).Error; err != nil {
		return fmt.Errorf("failed to load admin IDs: %w", err)
	}
	r.adminIDs = make(map[string]string, len(admins))
	for _, admin := range admins {
		r.adminIDs[admin.ClientID] = admin.ID
	}
	return nil
}

func (r *adminRepository) IsAdmin(ctx context.Context, clientID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.adminIDs[clientID]
	return ok
}
