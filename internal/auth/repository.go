package auth

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
}

type PgUserRepo struct {
	DB *gorm.DB
}

func (r *PgUserRepo) Create(user *User) error {
	return r.DB.Create(user).Error
}

func (r *PgUserRepo) FindByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (r *PgUserRepo) FindByID(id uint) (*User, error) {
	var user User
	result := r.DB.First(&user, id)
	return &user, result.Error
}
