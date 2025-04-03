package repository

import (
	"Edupay/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BillSemesterRepository adalah interface untuk operasi CRUD pada tagihan semester
type BillSemesterRepository interface {
	GetAllBillsRepository(page, limit int, semester, year string) ([]*model.BillSemester, error)
	GetBillByIdRepository(id string) (*model.BillSemester, error)
	CreateBillRepository(bill *model.BillSemester) (*model.BillSemester, error)
	UpdateBillByIdRepository(id string, bill *model.BillSemester) (*model.BillSemester, error)
	DeleteBillByIdRepository(id string) error
	GetBillsByStudentIdRepository(studentID string) ([]*model.BillSemester, error)
	GetAllUnpaidBillsRepository() ([]*model.BillSemester, error)
	GetDueBillsRepository() ([]*model.BillSemester, error)
	UpdateBillStatusRepository(billID string, status string) error
}

// Struct repository yang mengimplementasikan BillSemesterRepository
type billSemesterRepository struct {
	db *gorm.DB
}

// NewBillSemesterRepository membuat instance baru dari billSemesterRepository
func NewBillSemesterRepository(db *gorm.DB) *billSemesterRepository {
	return &billSemesterRepository{db}
}

// GetAllBillsRepository mengambil semua tagihan semester dengan pagination dan pencarian berdasarkan semester dan tahun
func (r *billSemesterRepository) GetAllBillsRepository(page, limit int, semester, year string) ([]*model.BillSemester, error) {
	var bills []*model.BillSemester
	offset := (page - 1) * limit

	query := r.db.Offset(offset).Limit(limit)
	if semester != "" {
		query = query.Where("semester LIKE ?", "%"+semester+"%")
	}
	if year != "" {
		query = query.Where("year = ?", year)
	}

	result := query.Preload("Transaction").Order("created_at DESC").Find(&bills)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting bills: %s", result.Error)
	}
	return bills, nil
}

// GetBillByIDRepository mengambil tagihan semester berdasarkan ID dengan relasi transaksi
func (r *billSemesterRepository) GetBillByIdRepository(id string) (*model.BillSemester, error) {
	var bill model.BillSemester
	result := r.db.Preload("Transaction").First(&bill, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("bill with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting bill with ID %s: %s", id, result.Error)
	}
	return &bill, nil
}

// CreateBillRepository membuat tagihan semester baru
func (r *billSemesterRepository) CreateBillRepository(bill *model.BillSemester) (*model.BillSemester, error) {
	// Atur jatuh tempo 7 hari setelah pembuatan
	bill.DueDate = time.Now().AddDate(0, 0, 7)

	// Simpan tagihan ke database
	result := r.db.Create(bill)
	if result.Error != nil {
		return nil, result.Error
	}

	// Buat transaksi terkait untuk tagihan ini
	transaction := &model.Transaction{
		Id:            fmt.Sprintf("BILL-%s", bill.ID),
		UserId:        bill.StudentId,
		Status:        model.STATUS_UNPAID,
		ProductType:   "semester_fee",
		Description:   fmt.Sprintf("Tagihan Semester %s - %s", bill.Semester, bill.Year),
		AdminFee:      model.ADMIN_FEE,
		Price:         bill.Amount,
		TotalPrice:    bill.Amount + model.ADMIN_FEE,
		ProductDetail: bill,
	}

	if err := r.db.Create(transaction).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction for bill: %w", err)
	}

	return bill, nil
}

// UpdateBillByIDRepository memperbarui tagihan semester berdasarkan ID
func (r *billSemesterRepository) UpdateBillByIdRepository(id string, bill *model.BillSemester) (*model.BillSemester, error) {
	result := r.db.Model(&model.BillSemester{}).Where("id = ?", id).Updates(bill)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("bill not found")
	}
	return bill, nil
}

// DeleteBillByIdRepository menghapus tagihan semester berdasarkan ID
func (r *billSemesterRepository) DeleteBillByIdRepository(id string) error {
	result := r.db.Delete(&model.BillSemester{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("bill not found")
	}
	return nil
}

// GetBillsByStudentIdRepository mengambil semua tagihan semester berdasarkan StudentID
func (r *billSemesterRepository) GetBillsByStudentIdRepository(studentID string) ([]*model.BillSemester, error) {
	var bills []*model.BillSemester
	result := r.db.Preload("Transaction").Where("student_id = ?", studentID).Find(&bills)
	if result.Error != nil {
		return nil, result.Error
	}
	return bills, nil
}

// GetAllUnpaidBillsRepository mengambil semua tagihan semester yang belum dibayar
func (r *billSemesterRepository) GetAllUnpaidBillsRepository() ([]*model.BillSemester, error) {
	var bills []*model.BillSemester
	result := r.db.Where("status = ?", model.STATUS_UNPAID).Find(&bills)
	if result.Error != nil {
		return nil, result.Error
	}
	return bills, nil
}

// GetDueBillsRepository mengambil semua tagihan semester yang sudah jatuh tempo
func (r *billSemesterRepository) GetDueBillsRepository() ([]*model.BillSemester, error) {
	var bills []*model.BillSemester
	result := r.db.Where("status = ? AND due_date < ?", model.STATUS_UNPAID, time.Now()).Find(&bills)
	if result.Error != nil {
		return nil, result.Error
	}
	return bills, nil
}

// UpdateBillStatusRepository memperbarui status tagihan semester setelah pembayaran
func (r *billSemesterRepository) UpdateBillStatusRepository(billID string, status string) error {
	result := r.db.Model(&model.BillSemester{}).Where("id = ?", billID).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("bill not found")
	}
	return nil
}
