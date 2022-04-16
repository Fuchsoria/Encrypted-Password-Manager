package main

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	UserID             int64
	UnlockedAt         time.Time
	HashKey            string
	HashKeyDescription string
	Passwords          []Password `gorm:"foreignkey:UserID;references:user_id"`
}

type Password struct {
	gorm.Model
	UserID   int64
	Name     string
	Password []byte
}

func (s *Storage) Connect() error {
	conn, err := gorm.Open(sqlite.Open("./storage.db"), &gorm.Config{})
	if err != nil {
		logger.Error(err)

		return err
	}

	db, err := conn.DB()
	if err != nil {
		logger.Error(err)

		return err
	}
	db.SetMaxOpenConns(1)

	err = conn.AutoMigrate(&User{})
	if err != nil {
		logger.Error(err)

		return err
	}
	err = conn.AutoMigrate(&Password{})
	if err != nil {
		logger.Error(err)

		return err
	}

	s.db = conn

	return nil
}

func (s *Storage) Disconnect() error {
	sqldb, err := s.db.DB()
	if err != nil {
		logger.Error(err)

		return err
	}

	return sqldb.Close()
}

func (s *Storage) Save(item interface{}) error {
	return s.db.Save(item).Error
}

func (s *Storage) GetUser(id int64) (*User, error) {
	var user User
	err := s.db.First(&user, "user_id = ?", id).Error
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	return &user, nil
}

func (s *Storage) IsUserExists(id int64) bool {
	err := s.db.First(&User{}, "user_id = ?", id).Error
	if err != nil {
		logger.Error(err)

		return false
	}

	return true
}

func (s *Storage) UpdateUnlocked(userID int64) error {
	u, err := s.GetUser(userID)
	if err != nil {
		logger.Error(err)

		return err
	}

	u.UnlockedAt = time.Now()

	return s.db.Save(&u).Error
}

func (s *Storage) GetPasswords(userID int64) ([]Password, error) {
	var passwords []Password
	err := s.db.Where("user_id = ?", userID).Find(&passwords).Error
	if err != nil {
		logger.Error(err)

		return nil, err
	}

	return passwords, nil
}

func (s *Storage) CleanMyPasswords(userID int64) error {
	err := s.db.Where("user_id = ?", userID).Delete(&Password{}).Error
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (s *Storage) CleanMyData(userID int64) error {
	err := s.db.Where("user_id = ?", userID).Delete(&User{}).Error
	if err != nil {
		logger.Error(err)

		return err
	}

	err = s.db.Where("user_id = ?", userID).Delete(&Password{}).Error
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (s *Storage) RemovePassword(userID int64, passwordID string) error {
	err := s.db.Where("user_id = ? AND id = ?", userID, passwordID).Delete(&Password{}).Error
	if err != nil {
		logger.Error(err)

		return err
	}

	return nil
}

func (s *Storage) SavePassword(userID int64, name, password string) error {
	key, err := locker.GetKey(userID)
	if err != nil {
		logger.Error(err)

		return err
	}

	b, err := Encrypt(key, password)
	if err != nil {
		logger.Error(err)

		return err
	}

	p := Password{
		UserID:   userID,
		Name:     name,
		Password: b,
	}

	return s.db.Create(&p).Error
}

func (s *Storage) GetPassword(userID int64, passwordID string) (string, error) {
	var password Password
	err := s.db.First(&password, "user_id = ? AND id = ?", userID, passwordID).Error
	if err != nil {
		logger.Error(err)

		return "", err
	}

	key, err := locker.GetKey(userID)
	if err != nil {
		logger.Error(err)

		return "", err
	}

	b, err := Decrypt(key, password.Password)
	if err != nil {
		logger.Error(err)

		return "", err
	}

	return string(b), nil
}

func (s *Storage) SaveUser(userID int64, hashKey string, hashKeyDescription string) error {
	u := User{
		UserID:             userID,
		HashKey:            hashKey,
		HashKeyDescription: hashKeyDescription,
	}

	return s.db.Create(&u).Error
}

func NewStorage() *Storage {
	return &Storage{}
}
