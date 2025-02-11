package repository

import (
	"Edupay/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// **ItemRepository Interface**
type ItemRepository interface {
	GetAllItemsRepository(page, limit int, name, code string) ([]*model.Item, error)
	GetItemByIdRepository(id string) (*model.Item, error)
	CreateItemRepository(item *model.Item) (*model.Item, error)
	UpdateItemByIdRepository(id string, item *model.Item) (*model.Item, error)
	DeleteItemByIdRepository(id string) error
	CalculateTotalPrice(item *model.Item) (float64, error)
}

// **Struct untuk Implementasi ItemRepository**
type itemRepository struct {
	db              *gorm.DB
	bookRepository  BookRepository
	shirtRepository ShirtRepository
	userRepository  UserRepository
}

// **Membuat Instance Baru dari ItemRepository**
func NewItemRepository(db *gorm.DB, bookRepository BookRepository, shirtRepository ShirtRepository, userRepository UserRepository) *itemRepository {
	return &itemRepository{
		db:              db,
		bookRepository:  bookRepository,
		shirtRepository: shirtRepository,
		userRepository:  userRepository,
	}
}

// **CreateItemRepository - Simpan item baru ke database**
func (r *itemRepository) CreateItemRepository(item *model.Item) (*model.Item, error) {
	// **Debugging log**
	fmt.Printf("✅ Creating item with Code: %s, UserID: %s\n", item.Code, item.UserId)

	// **1. Generate UUID jika ID kosong**
	if item.ID == "" {
		item.ID = uuid.New().String()
		fmt.Printf("🔄 Generated new UUID for item: %s\n", item.ID)
	}

	// **2. Simpan Item Utama ke Database**
	result := r.db.Create(item)
	if result.Error != nil {
		return nil, result.Error
	}

	// **3. Simpan item_details ke database**
	for _, itemDetail := range item.ItemDetails {
		itemDetail.ItemID = item.ID

		// **Pastikan `type` tidak kosong**
		if itemDetail.Type == "" {
			fmt.Println("❌ Error: Product type is missing in ItemDetail!")
			return nil, fmt.Errorf("product type is required for all items")
		}

		// **Simpan ke database**
		if err := r.db.Create(&itemDetail).Error; err != nil {
			return nil, fmt.Errorf("error saving item details: %w", err)
		}
	}

	// **4. Hitung Total Price**
	totalPrice, err := r.CalculateTotalPrice(item)
	if err != nil {
		return nil, fmt.Errorf("error calculating total price: %w", err)
	}
	item.TotalPrice = totalPrice

	// **5. Update total harga di database setelah insert**
	if err := r.db.Model(&model.Item{}).Where("id = ?", item.ID).Update("total_price", totalPrice).Error; err != nil {
		return nil, err
	}

	fmt.Printf("✅ Successfully created item with ID: %s, Total Price: %.2f\n", item.ID, item.TotalPrice)
	return item, nil
}

// **GetAllItemsRepository - Mengambil semua item dengan pagination & filter**
func (r *itemRepository) GetAllItemsRepository(page, limit int, name, code string) ([]*model.Item, error) {
	var items []*model.Item
	offset := (page - 1) * limit

	query := r.db.Preload("ItemDetails.Product").Offset(offset).Limit(limit)
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	result := query.Order("created_at DESC").Find(&items)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting items: %s", result.Error)
	}

	// **Hitung ulang TotalPrice untuk setiap item**
	for _, item := range items {
		totalPrice, err := r.CalculateTotalPrice(item)
		if err != nil {
			return nil, fmt.Errorf("error calculating total price for item %s: %s", item.Code, err)
		}
		item.TotalPrice = totalPrice
	}

	return items, nil
}

// **GetItemByIdRepository - Mengambil item berdasarkan ID**
func (r *itemRepository) GetItemByIdRepository(id string) (*model.Item, error) {
	var item model.Item
	result := r.db.Preload("ItemDetails.Product").First(&item, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item with ID %s not found", id)
		}
		return nil, fmt.Errorf("error getting item with ID %s: %s", id, result.Error)
	}

	// **Hitung ulang TotalPrice**
	totalPrice, err := r.CalculateTotalPrice(&item)
	if err != nil {
		return nil, fmt.Errorf("error calculating total price for item %s: %s", item.Code, err)
	}
	item.TotalPrice = totalPrice

	return &item, nil
}

// **CalculateTotalPrice - Menghitung total harga berdasarkan banyak buku & baju**
func (r *itemRepository) CalculateTotalPrice(item *model.Item) (float64, error) {
	var totalPrice float64

	// **Hitung total harga untuk semua produk dalam item**
	for _, itemDetail := range item.ItemDetails {
		if itemDetail.Type == "book" {
			book, err := r.bookRepository.GetBookByIdRepository(itemDetail.ProductID)
			if err != nil {
				return 0, fmt.Errorf("book not found: %w", err)
			}
			totalPrice += book.Price * float64(itemDetail.Quantity)
		} else if itemDetail.Type == "shirt" {
			shirt, err := r.shirtRepository.GetShirtByIdRepository(itemDetail.ProductID)
			if err != nil {
				return 0, fmt.Errorf("shirt not found: %w", err)
			}
			totalPrice += shirt.Price * float64(itemDetail.Quantity)
		}
	}

	return totalPrice, nil
}

// **UpdateItemByIdRepository - Memperbarui data item berdasarkan ID**
func (r *itemRepository) UpdateItemByIdRepository(id string, item *model.Item) (*model.Item, error) {
	// **Hapus item lama**
	if err := r.db.Where("item_id = ?", id).Delete(&model.ItemDetail{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old item details: %w", err)
	}

	// **Update item utama**
	result := r.db.Model(&model.Item{}).Where("id = ?", id).Updates(item)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("item not found")
	}

	// **Tambahkan item baru ke database**
	for _, itemDetail := range item.ItemDetails {
		itemDetail.ItemID = id
		if err := r.db.Create(&itemDetail).Error; err != nil {
			return nil, fmt.Errorf("error saving item details: %w", err)
		}
	}

	return item, nil
}

// **DeleteItemByIdRepository - Menghapus item berdasarkan ID**
func (r *itemRepository) DeleteItemByIdRepository(id string) error {
	// **Hapus item details terlebih dahulu**
	if err := r.db.Where("item_id = ?", id).Delete(&model.ItemDetail{}).Error; err != nil {
		return err
	}

	// **Hapus item utama**
	result := r.db.Delete(&model.Item{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}
