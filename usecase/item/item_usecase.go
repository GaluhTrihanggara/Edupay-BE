package item

import (
	"Edupay/model"
	"Edupay/repository"
	"fmt"

	"github.com/google/uuid"
)

type ItemUseCase interface {
	CreateItemUseCase(payload *model.Item) (*model.Item, error)
	GetAllItemsUseCase(page, limit int, name, code string) ([]*model.Item, error)
	GetItemByIdUseCase(id string) (*model.Item, error)
	UpdateItemByIdUseCase(id string, payload *model.Item) (*model.Item, error)
	DeleteItemByIdUseCase(id string) error
}

type itemUseCase struct {
	itemRepository  repository.ItemRepository
	bookRepository  repository.BookRepository
	shirtRepository repository.ShirtRepository
	userRepository  repository.UserRepository
}

func NewItemUseCase(itemRepository repository.ItemRepository, bookRepository repository.BookRepository, shirtRepository repository.ShirtRepository, userRepository repository.UserRepository) *itemUseCase {
	return &itemUseCase{
		itemRepository:  itemRepository,
		bookRepository:  bookRepository,
		shirtRepository: shirtRepository,
		userRepository:  userRepository,
	}
}

// CreateItemUseCase membuat item baru dengan multiple produk dalam satu transaksi
func (uc *itemUseCase) CreateItemUseCase(payload *model.Item) (*model.Item, error) {
	// Debugging log untuk melihat data yang diterima
	fmt.Printf("✅ Received request to create item: %+v\n", payload)

	for i := range payload.ItemDetails {
		if payload.ItemDetails[i].ID == "" {
			payload.ItemDetails[i].ID = uuid.New().String()
		}
	}

	// **1. Validasi UserId tidak kosong**
	if payload.UserId == "" {
		fmt.Println("❌ User ID is empty!")
		return nil, fmt.Errorf("user ID is required")
	}

	// **2. Validasi user_id ada di database**
	user, err := uc.userRepository.GetUserByIdRepository(payload.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// **3. Format nama dan deskripsi berdasarkan user**
	payload.Name = fmt.Sprintf("%s - Pesanan", user.Name)
	payload.Description = fmt.Sprintf("Pesanan milik %s:", user.Name)

	// **4. Validasi item_details tidak kosong**
	if len(payload.ItemDetails) == 0 {
		return nil, fmt.Errorf("at least one product is required in item_details")
	}

	// **5. Hitung total harga & Validasi setiap item dalam `item_details`**
	var totalPrice float64
	for i, itemDetail := range payload.ItemDetails {
		// **Validasi product_id tidak kosong**
		if itemDetail.ProductID == "" {
			return nil, fmt.Errorf("product ID is required for all items")
		}

		// **Validasi Type (book atau shirt)**
		if itemDetail.Type != "book" && itemDetail.Type != "shirt" {
			return nil, fmt.Errorf("invalid product type: %s (must be 'book' or 'shirt')", itemDetail.Type)
		}

		// **Ambil harga produk berdasarkan Type**
		if itemDetail.Type == "book" {
			book, err := uc.bookRepository.GetBookByIdRepository(itemDetail.ProductID)
			if err != nil {
				return nil, fmt.Errorf("book not found: %w", err)
			}
			if book.Stock < itemDetail.Quantity {
				return nil, fmt.Errorf("insufficient stock for book: %s (requested: %d, available: %d)", book.Title, itemDetail.Quantity, book.Stock)
			}
			totalPrice += book.Price * float64(itemDetail.Quantity)

			// **Tambahkan Product ke item_detail**
			payload.ItemDetails[i].Product = model.Product{
				ID:    book.ID,
				Name:  book.Title,
				Price: book.Price,
			}
			payload.ItemDetails[i].Type = "book" // ✅ **Pastikan `type` diisi!**
			payload.Description += fmt.Sprintf(" %dx Buku %s,", itemDetail.Quantity, book.Title)

		} else if itemDetail.Type == "shirt" {
			shirt, err := uc.shirtRepository.GetShirtByIdRepository(itemDetail.ProductID)
			if err != nil {
				return nil, fmt.Errorf("shirt not found: %w", err)
			}
			if shirt.Stock < itemDetail.Quantity {
				return nil, fmt.Errorf("insufficient stock for shirt: %s (requested: %d, available: %d)", shirt.Name, itemDetail.Quantity, shirt.Stock)
			}
			totalPrice += shirt.Price * float64(itemDetail.Quantity)

			// **Tambahkan Product ke item_detail**
			payload.ItemDetails[i].Product = model.Product{
				ID:    shirt.ID,
				Name:  string(shirt.Name),
				Price: shirt.Price,
			}
			payload.ItemDetails[i].Type = "shirt" // ✅ **Pastikan `type` diisi!**
			payload.Description += fmt.Sprintf(" %dx %s,", itemDetail.Quantity, shirt.Name)
		}
	}

	// **6. Set total price**
	payload.TotalPrice = totalPrice

	// **7. Generate UUID jika `id` kosong**
	if payload.ID == "" {
		payload.ID = uuid.New().String()
		fmt.Printf("🔄 Generated new UUID for item: %s\n", payload.ID)
	}

	// **8. Simpan item ke repository**
	item, err := uc.itemRepository.CreateItemRepository(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	fmt.Printf("✅ Successfully created item with ID: %s\n", item.ID)
	return item, nil
}

// GetAllItemsUseCase mengambil semua item dengan filter nama dan kode
func (uc *itemUseCase) GetAllItemsUseCase(page, limit int, name, code string) ([]*model.Item, error) {
	return uc.itemRepository.GetAllItemsRepository(page, limit, name, code)
}

// GetItemByIdUseCase mengambil item berdasarkan ID
func (uc *itemUseCase) GetItemByIdUseCase(id string) (*model.Item, error) {
	return uc.itemRepository.GetItemByIdRepository(id)
}

// UpdateItemByIdUseCase memperbarui item berdasarkan ID
func (uc *itemUseCase) UpdateItemByIdUseCase(id string, payload *model.Item) (*model.Item, error) {
	// Validasi keberadaan item
	_, err := uc.itemRepository.GetItemByIdRepository(id)
	if err != nil {
		return nil, fmt.Errorf("item with ID %s not found: %v", id, err)
	}

	// Validasi stok dan perhitungan harga
	for i, itemDetail := range payload.ItemDetails {
		if itemDetail.Type == "book" {
			book, err := uc.bookRepository.GetBookByIdRepository(itemDetail.ProductID)
			if err != nil {
				return nil, fmt.Errorf("book not found: %w", err)
			}
			if book.Stock < itemDetail.Quantity {
				return nil, fmt.Errorf("insufficient stock for book: %s", book.Title)
			}
			payload.ItemDetails[i].Product = model.Product{
				ID:    book.ID,
				Name:  book.Title,
				Price: book.Price,
			}
		} else if itemDetail.Type == "shirt" {
			shirt, err := uc.shirtRepository.GetShirtByIdRepository(itemDetail.ProductID)
			if err != nil {
				return nil, fmt.Errorf("shirt not found: %w", err)
			}
			if shirt.Stock < itemDetail.Quantity {
				return nil, fmt.Errorf("insufficient stock for shirt: %s", shirt.Name)
			}
			payload.ItemDetails[i].Product = model.Product{
				ID:    shirt.ID,
				Name:  string(shirt.Name),
				Price: shirt.Price,
			}
		} else {
			return nil, fmt.Errorf("invalid product type: %s", itemDetail.Type)
		}
	}

	// Update item di repository
	updatedItem, err := uc.itemRepository.UpdateItemByIdRepository(id, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return updatedItem, nil
}

// DeleteItemByIdUseCase menghapus item berdasarkan ID
func (uc *itemUseCase) DeleteItemByIdUseCase(id string) error {
	// Validasi keberadaan item
	_, err := uc.itemRepository.GetItemByIdRepository(id)
	if err != nil {
		return fmt.Errorf("item with ID %s not found: %v", id, err)
	}

	// Hapus item
	return uc.itemRepository.DeleteItemByIdRepository(id)
}
