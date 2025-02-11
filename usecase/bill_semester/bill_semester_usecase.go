package billsemester

import (
	"Edupay/model"
	"Edupay/repository"
	"Edupay/usecase/service"
	"errors"
	"fmt"
	"time"
)

type BillSemesterUseCase interface {
	CreateBillSemesterUseCase(payload *model.BillSemester) (*model.BillSemester, error)
	GetAllBillSemesterUseCase(page, limit int, semester, year string) ([]*model.BillSemester, error)
	GetBillSemesterByIdUseCase(billSemesterId string) (*model.BillSemester, error)
	UpdateBillSemesterByIdUseCase(billSemesterId string, payload *model.BillSemester) (*model.BillSemester, error)
	DeleteBillSemesterByIdUseCase(billSemesterId string) error
}

type billSemesterUseCase struct {
	billSemesterRepository repository.BillSemesterRepository
	transactionRepository  repository.TransactionRepository
}

func NewBillSemesterUseCase(billSemesterRepository repository.BillSemesterRepository, transactionRepository repository.TransactionRepository) *billSemesterUseCase {
	return &billSemesterUseCase{
		billSemesterRepository: billSemesterRepository,
		transactionRepository:  transactionRepository,
	}
}

// CreateBillSemesterUseCase membuat tagihan semester baru
func (uc *billSemesterUseCase) CreateBillSemesterUseCase(payload *model.BillSemester) (*model.BillSemester, error) {
	// Buat tagihan semester baru
	billSemester, err := uc.billSemesterRepository.CreateBillRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating bill semester in database: %w", err)
	}
	emailBody := fmt.Sprintf(` 
	<h1> Tagihan Semester Baru </h1>
	<p> Anda memiliki tagihan semester baru untuk siswa %s.</p>
	<p> Jumlah yang harus dibaray: Rp %.2f</p>
	<p> Silakan login ke aplikasi Edupay untuk melakukan pembayaran.</p>`, billSemester.StudentId, billSemester.Amount)

	err = service.SendEmail("hanggaragaluh22@gmail.com", "Tagihan Semester Baru", emailBody)
	if err != nil {
		return nil, fmt.Errorf("error sending email: %w", err)
	}

	return billSemester, nil
}

// GetAllBillSemesterUseCase mengambil semua tagihan semester
func (uc *billSemesterUseCase) GetAllBillSemesterUseCase(page, limit int, semester, year string) ([]*model.BillSemester, error) {
	// Ambil semua tagihan dengan data relasi transaksi
	billSemesters, err := uc.billSemesterRepository.GetAllBillsRepository(page, limit, semester, year)
	if err != nil {
		return nil, err
	}
	return billSemesters, nil
}

// GetBillSemesterByIdUseCase mengambil tagihan semester berdasarkan ID
func (uc *billSemesterUseCase) GetBillSemesterByIdUseCase(billSemesterId string) (*model.BillSemester, error) {
	// Ambil tagihan dengan relasi ke transaksi
	billSemester, err := uc.billSemesterRepository.GetBillByIdRepository(billSemesterId)
	if err != nil {
		return nil, errors.New("bill semester not found")
	}
	return billSemester, nil
}

// UpdateBillSemesterByIdUseCase memperbarui tagihan semester berdasarkan ID
func (uc *billSemesterUseCase) UpdateBillSemesterByIdUseCase(billSemesterId string, payload *model.BillSemester) (*model.BillSemester, error) {
	// Validasi keberadaan tagihan semester
	existingBillSemester, err := uc.billSemesterRepository.GetBillByIdRepository(billSemesterId)
	if err != nil {
		return nil, fmt.Errorf("bill semester with ID %s not found: %v", billSemesterId, err)
	}

	// Perbarui data tagihan semester
	existingBillSemester.Semester = payload.Semester
	existingBillSemester.Year = payload.Year
	existingBillSemester.Amount = payload.Amount
	existingBillSemester.Status = payload.Status
	existingBillSemester.DueDate = payload.DueDate
	existingBillSemester.UpdatedAt = time.Now()

	updatedBillSemester, err := uc.billSemesterRepository.UpdateBillByIdRepository(billSemesterId, existingBillSemester)
	if err != nil {
		return nil, fmt.Errorf("failed to update bill semester: %v", err)
	}

	return updatedBillSemester, nil
}

// DeleteBillSemesterByIdUseCase menghapus tagihan semester berdasarkan ID
func (uc *billSemesterUseCase) DeleteBillSemesterByIdUseCase(billSemesterId string) error {
	// Validasi apakah tagihan semester memiliki transaksi terkait
	billSemester, err := uc.billSemesterRepository.GetBillByIdRepository(billSemesterId)
	if err != nil {
		return fmt.Errorf("bill semester with ID %s not found: %v", billSemesterId, err)
	}

	if billSemester.TransactionId != nil {
		return fmt.Errorf("cannot delete bill semester with ID %s because it is linked to a transaction", billSemesterId)
	}

	// Hapus tagihan semester
	err = uc.billSemesterRepository.DeleteBillByIdRepository(billSemesterId)
	if err != nil {
		return fmt.Errorf("failed to delete bill semester: %v", err)
	}
	return nil
}
