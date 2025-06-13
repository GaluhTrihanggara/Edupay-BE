package balance

// import (
// 	"Edupay/model"
// 	"Edupay/repository"
// 	"fmt"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type BalanceUsecase interface {
// 	GenerateVaUseCase(userId uuid.UUID, payload model.GenerateVirtualAgregator) (*model.VaNumber, error)
// 	CreateBalanceUseCase(payload *model.PartnerCallbackVirtualAggregator) (*model.Payment, error)
// 	GetPayBalanceStatusUseCase(userId uuid.UUID, vaId string) (*model.VaNumber, error)
// }

// type balanceUsecase struct {
// 	balanceRepository     repository.BalanceRepository
// 	userRepository        repository.UserRepository
// 	paymentRepository     repository.PaymentRepository
// 	virtualAgregatorOyApi repository.VirtualAgregatorOyApi
// }

// func NewBalanceUsecase(balanceRepository repository.BalanceRepository, userRepository repository.UserRepository, paymentRepository repository.PaymentRepository, virtualAgregatorOyApi repository.VirtualAgregatorOyApi) *balanceUsecase {
// 	return &balanceUsecase{
// 		balanceRepository:     balanceRepository,
// 		userRepository:        userRepository,
// 		paymentRepository:     paymentRepository,
// 		virtualAgregatorOyApi: virtualAgregatorOyApi,
// 	}
// }

// func (uc *balanceUsecase) GenerateVaUseCase(userId uuid.UUID, payload model.GenerateVirtualAgregator) (*model.VaNumber, error) {
// 	user, err := uc.userRepository.GetUserByIdRepository(userId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	payload.PatnerUserId = userId.String()
// 	payload.Username = user.Name
// 	payload.IsOpen = true
// 	payload.SingleUse = false

// 	resp, err := uc.virtualAgregatorOyApi.GenerateVaApi(payload)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate VA: %v", err)
// 	}

// 	insert := &model.VaNumber{
// 		UserId:         user.ID,
// 		VaNumber:       resp.VaNumber,
// 		VaStatus:       resp.VaStatus,
// 		Amount:         resp.Amount,
// 		ExpirationTime: resp.ExpirationTime,
// 		Name:           resp.Name,
// 	}

// 	existingVa, err := uc.balanceRepository.InsertVaRepository(insert)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if existingVa != nil {
// 		return &model.VaNumber{
// 			UserId:         user.ID,
// 			VaNumber:       existingVa.VaNumber,
// 			VaStatus:       existingVa.VaStatus,
// 			BankCode:       payload.BankCode,
// 			Amount:         existingVa.Amount,
// 			ExpirationTime: existingVa.ExpirationTime,
// 			Name:           existingVa.Name,
// 		}, nil
// 	}

// 	return resp, nil
// }

// func (uc *balanceUsecase) CreateBalanceUseCase(payload *model.PartnerCallbackVirtualAggregator) (*model.Payment, error) {
// 	userId, err := uuid.Parse(payload.PartnerUserId)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid user ID format: %v", err)
// 	}

// 	user, err := uc.userRepository.GetUserByIdRepository(userId)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal mengambil data user: %v", err)
// 	}

// 	user.Amount += payload.Amount
// 	user.UpdatedAt = time.Now()

// 	_, err = uc.userRepository.UpdateUserAmountByIdRepository(userId, user)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal memperbarui saldo user: %v", err)
// 	}

// 	// Total harga termasuk biaya admin
// 	total := payload.Amount + model.ADMIN_FEE

// 	// Membuat payment yang baru
// 	createPayment := &model.Payment{
// 		ID:          uuid.New(),
// 		Amount:      payload.Amount,
// 		Method:      "topup", // Metode pembayaran, misalnya "topup"
// 		Status:      model.STATUS_SUCCESSFUL,
// 		OrderID:     uuid.New(), // ID Order jika ada
// 		AdminFee:    model.ADMIN_FEE,
// 		TotalPrice:  total,
// 		PaymentDate: time.Now(),
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 		DeletedAt:   gorm.DeletedAt{}, // Timestamp penghapusan (set kosong jika tidak dihapus)
// 	}

// 	// Simpan payment ke dalam repository
// 	resp, err := uc.paymentRepository.CreatePaymentRepository(createPayment) // Menggunakan paymentRepository
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }

// func (uc *balanceUsecase) GetPayBalanceStatusUseCase(userId, vaId string) (*model.VaNumber, error) {
// 	resp, err := uc.virtualAgregatorOyApi.GetVaIdStatusVaApi(vaId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp, nil
// }
