package billsemester

import (
	"Edupay/model"
	"Edupay/repository"
	"Edupay/usecase/mail"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BillSemesterUseCase interface {
	CreateBillSemesterUseCase(payload *model.BillSemester) (*model.BillSemester, error)
	GetAllBillSemesterUseCase(page, limit int, semester, year string) ([]*model.BillSemester, error)
	GetBillSemesterByIdUseCase(billSemesterId string) (*model.BillSemester, error)
	UpdateBillSemesterByIdUseCase(billSemesterId string, payload *model.BillSemester) (*model.BillSemester, error)
	DeleteBillSemesterByIdUseCase(billSemesterId string) error
	GetBillsByStudentIdUseCase(studentID string) ([]*model.BillSemester, error)
	GetAllUnpaidBillsUseCase() ([]*model.BillSemester, error)
	BillInquirySemesterUseCase(userId string, payload *model.OyBillerApi) (*model.Transaction, error)
	PayBillSemesterUseCase(userId string, payload *model.OyBillerApi) (*model.Transaction, error)
	BillSemesterStatusUseCase(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error)
}

type billSemesterUseCase struct {
	billSemesterRepository repository.BillSemesterRepository
	userRepository         repository.UserRepository
	transactionRepository  repository.TransactionRepository
	billerOyApi            repository.BillerOyApiRepository
}

func NewBillSemesterUseCase(billSemesterRepository repository.BillSemesterRepository, userRepository repository.UserRepository, transactionRepository repository.TransactionRepository, billerOyApiRepository repository.BillerOyApiRepository) *billSemesterUseCase {
	return &billSemesterUseCase{
		billSemesterRepository: billSemesterRepository,
		userRepository:         userRepository,
		transactionRepository:  transactionRepository,
		billerOyApi:            billerOyApiRepository,
	}
}

// 📝 Membuat tagihan semester baru dengan jatuh tempo 7 hari
func (uc *billSemesterUseCase) CreateBillSemesterUseCase(payload *model.BillSemester) (*model.BillSemester, error) {
	payload.DueDate = time.Now().AddDate(0, 0, 7) // Set jatuh tempo 7 hari

	billSemester, err := uc.billSemesterRepository.CreateBillRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("error creating bill semester in database: %w", err)
	}
	return billSemester, nil
}

func (uc *billSemesterUseCase) GetAllBillSemesterUseCase(page, limit int, semester, year string) ([]*model.BillSemester, error) {
	bills, err := uc.billSemesterRepository.GetAllBillsRepository(page, limit, semester, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get all bills: %w", err)
	}
	return bills, nil
}

func (uc *billSemesterUseCase) GetBillSemesterByIdUseCase(billSemesterId string) (*model.BillSemester, error) {
	bill, err := uc.billSemesterRepository.GetBillByIdRepository(billSemesterId)
	if err != nil {
		return nil, fmt.Errorf("failed to get bill by id: %w", err)
	}
	return bill, nil
}

func (uc *billSemesterUseCase) UpdateBillSemesterByIdUseCase(billSemesterId string, payload *model.BillSemester) (*model.BillSemester, error) {
	bill, err := uc.billSemesterRepository.UpdateBillByIdRepository(billSemesterId, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update bill: %w", err)
	}
	return bill, nil
}

func (uc *billSemesterUseCase) DeleteBillSemesterByIdUseCase(billSemesterId string) error {
	if err := uc.billSemesterRepository.DeleteBillByIdRepository(billSemesterId); err != nil {
		return fmt.Errorf("failed to delete bill: %w", err)
	}
	return nil
}

func (uc *billSemesterUseCase) GetBillsByStudentIdUseCase(studentID string) ([]*model.BillSemester, error) {
	bills, err := uc.billSemesterRepository.GetBillsByStudentIdRepository(studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bills by student ID: %w", err)
	}
	return bills, nil
}

func (uc *billSemesterUseCase) GetAllUnpaidBillsUseCase() ([]*model.BillSemester, error) {
	bills, err := uc.billSemesterRepository.GetAllUnpaidBillsRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to get unpaid bills: %w", err)
	}
	return bills, nil
}

// 📌 Inquiry tagihan semester
func (uc *billSemesterUseCase) BillInquirySemesterUseCase(userId string, payload *model.OyBillerApi) (*model.Transaction, error) {
	user, err := uc.userRepository.GetUserByIdRepository(userId)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	// Gunakan `user.Name` dalam `Description`
	description := fmt.Sprintf("Tagihan Semester untuk %s - %s", user.Name, time.Now().Format("2006"))

	// Panggil API OY untuk inquiry
	oy, err := uc.billerOyApi.BillInquryRepository(payload)
	if err != nil {
		return nil, err
	}

	totalPrice := payload.Amount + oy.AdminFee

	// Simpan transaksi baru
	transaction := &model.Transaction{
		Id:            fmt.Sprintf("BILL-%s", uuid.New().String()),
		UserId:        userId,
		Status:        model.STATUS_UNPAID,
		ProductType:   "semester_fee",
		Description:   description,
		AdminFee:      model.ADMIN_FEE,
		Price:         payload.Amount,
		TotalPrice:    totalPrice,
		ProductDetail: payload,
	}

	_, err = uc.transactionRepository.CreateTransactionRepository(transaction)
	if err != nil {
		return nil, fmt.Errorf("error creating semester bill in database: %w", err)
	}
	return transaction, nil
}

// 💰 Membayar tagihan semester
func (uc *billSemesterUseCase) PayBillSemesterUseCase(userId string, payload *model.OyBillerApi) (*model.Transaction, error) {
	transaction, err := uc.transactionRepository.GetTransactionByIdRepository(payload.PartnerTxId)
	if err != nil {
		return nil, errors.New("transaction not found")
	} else if transaction.Status == model.STATUS_SUCCESSFUL {
		return nil, errors.New("this semester bill has been paid")
	}

	user, err := uc.userRepository.GetUserByIdRepository(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Amount < transaction.TotalPrice {
		transactionFail := &model.Transaction{
			Status:    model.STATUS_FAIL,
			UpdatedAt: time.Now(),
		}

		_, err := uc.transactionRepository.UpdateTransactionByIdRepository(payload.PartnerTxId, transactionFail)
		if err != nil {
			return nil, errors.New("your balance is not enough")
		}

		return nil, errors.New("your balance is not enough")
	}

	// Update saldo pengguna
	user.Amount -= transaction.TotalPrice
	err = uc.userRepository.UpdateUserBalanceRepository(userId, user.Amount)
	if err != nil {
		return nil, err
	}

	// Update status transaksi
	updateTransaction := &model.Transaction{
		Status:    model.STATUS_SUCCESSFUL,
		UpdatedAt: time.Now(),
	}

	resp, err := uc.transactionRepository.UpdateTransactionByIdRepository(payload.PartnerTxId, updateTransaction)
	if err != nil {
		return nil, fmt.Errorf("error updating semester bill in database: %w", err)
	}

	// Kirim email konfirmasi
	mailsend := model.PayloadMail{
		OrderId:       transaction.Id,
		CustomerName:  user.Name,
		Status:        resp.Status,
		RecipentEmail: user.Email,
		TransactionAt: resp.UpdatedAt,
		ProductType:   "SEMESTER",
		Description:   transaction.Description,
		AdminFee:      transaction.AdminFee,
		Price:         transaction.Price,
		TotalPrice:    transaction.TotalPrice,
	}
	mail.SendingMail(mailsend)

	return resp, nil
}

// 🔍 Cek status pembayaran tagihan semester
func (uc *billSemesterUseCase) BillSemesterStatusUseCase(payload *model.OyBillerApi) (*model.OyBillerApiResponse, error) {
	semesterBill, err := uc.billerOyApi.BillInquryRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve semester bill: %v", err)
	}
	return semesterBill, nil
}
