package model

import (
	"errors"

	"gorm.io/gorm"
)

type Model[T any] struct {
	db *gorm.DB
}

func NewModel[T any](db *gorm.DB) *Model[T] {
	return &Model[T]{db: db}
}

func (m *Model[T]) GetById(dest *T, id uint) (bool, error) {
	err := m.db.First(dest, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Model[T]) Find(dest *[]T, conds ...any) error {
	return m.db.Find(dest, conds...).Error
}

func (m *Model[T]) Create(value *T) error {
	return m.db.Create(value).Error
}

func (m *Model[T]) Save(value *T) error {
	return m.db.Save(value).Error
}

func (m *Model[T]) Update(column string, value any) error {
	return m.db.Model(new(T)).Update(column, value).Error
}

func (m *Model[T]) Updates(values any) error {
	return m.db.Model(new(T)).Updates(values).Error
}

func (m *Model[T]) Delete(conds ...any) error {
	return m.db.Delete(new(T), conds...).Error
}

func (m *Model[T]) DeleteByID(id any) error {
	return m.db.Delete(new(T), id).Error
}

func (m *Model[T]) Count(conds ...any) (int64, error) {
	var count int64
	query := m.db.Model(new(T))
	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}
	err := query.Count(&count).Error
	return count, err
}

func (m *Model[T]) Exists(conds ...any) (bool, error) {
	count, err := m.Count(conds...)
	return count > 0, err
}

func (m *Model[T]) FindPage(dest *[]T, page, pageSize int, conds ...any) (int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query := m.db.Model(new(T))
	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	if total == 0 {
		*dest = []T{}
		return 0, nil
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(dest).Error
	return total, err
}
