package routes

import (
	"Edupay/controller"
	"Edupay/repository"
	authusecase "Edupay/usecase/auth"
	billsemester "Edupay/usecase/bill_semester"
	"Edupay/usecase/book"
	"Edupay/usecase/item"
	middleware "Edupay/usecase/middlewares"
	"Edupay/usecase/shirt"
	"Edupay/usecase/student"
	"Edupay/usecase/transaction"
	transactionhistory "Edupay/usecase/transaction_history"
	user "Edupay/usecase/users"

	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, db *gorm.DB) {

	// User
	userRepository := repository.NewUserRepository(db)
	userUseCase := user.NewUserUseCase(userRepository)
	userController := controller.NewUserController(userUseCase)

	// Bill Semester
	billSemesterRepository := repository.NewBillSemesterRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)
	billSemesterUseCase := billsemester.NewBillSemesterUseCase(billSemesterRepository, transactionRepository)
	transactionUseCase := transaction.NewTransactionUseCase(transactionRepository)
	billSemesterController := controller.NewBillSemesterController(billSemesterUseCase, transactionUseCase)

	// Auth
	authRepository := repository.NewAuthRepository(db)
	authUseCase := authusecase.NewAuthUsecase(authRepository)
	authController := controller.NewAuthController(authUseCase)

	// Book
	bookRepository := repository.NewBookRepository(db)
	bookUseCase := book.NewBookUseCase(bookRepository)
	bookController := controller.NewBookController(bookUseCase)

	// Shirt
	shirtRepository := repository.NewShirtRepository(db)
	shirtUseCase := shirt.NewShirtUseCase(shirtRepository)
	shirtController := controller.NewShirtController(shirtUseCase)

	// Item
	itemRepository := repository.NewItemRepository(db, bookRepository, shirtRepository, userRepository)
	itemUseCase := item.NewItemUseCase(itemRepository, bookRepository, shirtRepository, userRepository)
	itemController := controller.NewItemController(itemUseCase)

	// Student
	studentRepository := repository.NewStudentRepository(db)
	studentUseCase := student.NewStudentUseCase(studentRepository)
	studentController := controller.NewStudentController(studentUseCase)

	// Transaction
	transactionRepository = repository.NewTransactionRepository(db)
	transactionUseCase = transaction.NewTransactionUseCase(transactionRepository)
	transactionController := controller.NewTransactionController(transactionUseCase)

	// Transaction History
	transactionHistoryRepository := repository.NewTransactionHistoryRepository(db)
	transactionHistoryUseCase := transactionhistory.NewTransactionHistoryUseCase(transactionHistoryRepository)
	transactionHistoryController := controller.NewTransactionHistoryController(transactionHistoryUseCase)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
			<h1>Welcome to Edupay API</h1>
		`)
	})

	// Routes Grouping
	api := e.Group("/api/v1")
	admin := api.Group("/admin", middleware.AdminAuthMiddleware) // Tambahkan middleware
	user := api.Group("/user", middleware.AuthMiddleware)        // Tambahkan middleware

	// ====== AUTH ROUTES =======
	api.POST("/login", authController.LoginController)
	api.POST("/register", authController.RegisterController)
	api.POST("/register/admin", authController.RegisterAdminController)

	// ====== ADMIN ROLE =======
	admin.POST("/bill-semester", billSemesterController.CreateBillSemesterController)
	admin.GET("/bill-semesters", billSemesterController.GetAllBillSemesterController)
	admin.GET("/bill-semester/:id", billSemesterController.GetBillSemesterByIdController)
	admin.PUT("/bill-semester/:id", billSemesterController.UpdateBillSemesterByIdController)
	admin.DELETE("/bill-semester/:id", billSemesterController.DeleteBillSemesterByIdController)

	// ====== Buku Admin =======
	admin.POST("/book", bookController.CreateBookController)
	admin.GET("/books", bookController.GetAllBooksController)
	admin.GET("/book/:id", bookController.GetBookByIdController)
	admin.PUT("/book/:id", bookController.UpdateBookByIdController)
	admin.DELETE("/book/:id", bookController.DeleteBookByIdController)

	// ====== Item Admin =======
	admin.POST("/item", itemController.CreateItemController)
	admin.GET("/items", itemController.GetAllItemsController)
	admin.GET("/item/:id", itemController.GetItemByIdController)
	admin.PUT("/item/:id", itemController.UpdateItemByIdController)
	admin.DELETE("/item/:id", itemController.DeleteItemByIdController)

	// ====== Shirt Admin =======
	admin.POST("/shirt", shirtController.CreateShirtController)
	admin.GET("/shirts", shirtController.GetAllShirtsController)
	admin.GET("/shirt/:id", shirtController.GetShirtByIdController)
	admin.PUT("/shirt/:id", shirtController.UpdateShirtByIdController)
	admin.DELETE("/shirt/:id", shirtController.DeleteShirtByIdController)

	// ====== Student Admin =======
	admin.POST("/student", studentController.CreateStudentController)
	admin.GET("/students", studentController.GetAllStudentsController)
	admin.GET("/students/by-parent-name", studentController.GetStudentsByParentNameController)
	admin.GET("/student/:id", studentController.GetStudentByIdController)
	admin.PUT("/student/:id", studentController.UpdateStudentByIdController)
	admin.DELETE("/student/:id", studentController.DeleteStudentByIdController)

	// ====== Transaction Admin =======
	admin.POST("/transaction", transactionController.CreateTransactionController)
	admin.GET("/transactions", transactionController.GetAllTransactionsController)
	admin.GET("/transaction/:id", transactionController.GetTransactionByIdController)

	// ====== Transaction History Admin =======
	admin.GET("/transaction-histories", transactionHistoryController.GetAllHistoriesController)
	admin.GET("/transaction-history/:id", transactionHistoryController.GetHistoryByIdController)
	admin.GET("/transaction-history/user/:id", transactionController.GetTransactionByUserIdController)
	admin.DELETE("/transaction-history/:id", transactionHistoryController.DeleteHistoryByIdController)

	// ====== User Admin =======
	admin.GET("/users", userController.GetAllUsersController)
	admin.GET("/user/:id", userController.GetUserByIdController)
	admin.GET("/user/:email", userController.GetUserByEmailController)
	admin.GET("/user/:phone", userController.GetUserByPhoneController)
	admin.GET("/user/query", userController.GetUserByQueryController)
	admin.PUT("/user/:id", userController.UpdateUserByIdController)
	admin.DELETE("/user/:id", userController.DeleteUserByIdController)

	// ====== USER ROLE =======
	user.GET("/profile", userController.GetUserByIdController)
	user.PUT("/user", userController.GetUserByIdController)
	user.DELETE("/user", userController.DeleteUserByIdController)

	// ====== USER TRANSACTION =======
	user.POST("/item", itemController.CreateItemController)
	user.POST("/item", itemController.CreateItemController)
	user.GET("/item/:id", itemController.GetItemByIdController)
	user.PUT("/item/:id", itemController.UpdateItemByIdController)
	user.DELETE("/item/:id", itemController.DeleteItemByIdController)

	// ====== USER TRANSACTION =======
	user.GET("/transactions", transactionController.GetTransactionByUserIdController)
	user.GET("/transaction-histories", transactionHistoryController.GetHistoriesByUserIdController)

	// ====== USER STUDENT =======
	user.POST("/student", studentController.CreateStudentController)
	user.GET("/student/:id", studentController.GetStudentByIdController)
	user.GET("/students", studentController.GetStudentsByParentNameController) // Ambil siswa berdasarkan parent_name

	// ====== USER TRANSACTION =======
	user.GET("/transactions", transactionController.GetTransactionByUserIdController) // Ambil semua transaksi user
	user.GET("/transaction/:id", transactionController.GetTransactionByIdController)  // Ambil transaksi berdasarkan ID
	user.POST("/transaction", transactionController.CreateTransactionController)      // Buat transaksi baru

	// ====== USER TRANSACTION HISTORY =======
	user.GET("/transaction-histories", transactionHistoryController.GetHistoriesByUserIdController) // Ambil semua riwayat transaksi user
	user.GET("/transaction-history/:id", transactionHistoryController.GetHistoryByIdController)     // Ambil detail riwayat transaksi berdasarkan ID

}
