package routes

import (
	"Edupay/controller"
	"Edupay/repository"
	authusecase "Edupay/usecase/auth"
	"Edupay/usecase/book"
	"Edupay/usecase/cart"
	"Edupay/usecase/cart_item"
	middleware "Edupay/usecase/middlewares"
	"Edupay/usecase/order"
	"Edupay/usecase/payment"
	"Edupay/usecase/payment_history"
	"Edupay/usecase/shirt"
	"Edupay/usecase/student"
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

	// Student
	studentRepository := repository.NewStudentRepository(db)
	studentUseCase := student.NewStudentUseCase(studentRepository)
	studentController := controller.NewStudentController(studentUseCase)

	// Cart
	cartRepository := repository.NewCartRepository(db)
	cartUseCase := cart.NewCartUseCase(cartRepository)
	cartController := controller.NewCartController(cartUseCase)

	// Cart Item
	cartItemRepository := repository.NewCartItemRepository(db)
	cartItemUseCase := cart_item.NewCartItemUseCase(cartItemRepository, cartRepository)
	cartItemController := controller.NewCartItemController(cartItemUseCase)

	// Order
	orderRepository := repository.NewOrderRepository(db)
	orderUseCase := order.NewOrderUseCase(orderRepository, cartItemRepository)
	orderController := controller.NewOrderController(orderUseCase)

	// Payment
	paymentRepository := repository.NewPaymentRepository(db)
	paymentUseCase := payment.NewPaymentUseCase(paymentRepository, orderRepository)
	paymentController := controller.NewPaymentController(paymentUseCase)

	// Payment History
	paymentHistoryRepository := repository.NewPaymentHistoryRepository(db)
	paymentHistoryUseCase := payment_history.NewPaymentHistoryUseCase(paymentHistoryRepository)
	paymentHistoryController := controller.NewPaymentHistoryController(paymentHistoryUseCase)
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

	// ====== Buku Admin =======
	admin.POST("/book", bookController.CreateBookController)
	admin.GET("/books", bookController.GetAllBooksController)
	admin.GET("/book/:id", bookController.GetBookByIDController)
	admin.GET("/book/code/:code", bookController.GetBookByCodeController)
	admin.PUT("/book/:id", bookController.UpdateBookByIDController)
	admin.DELETE("/book/:id", bookController.DeleteBookByIDController)

	// ====== Shirt Admin =======
	admin.POST("/shirt", shirtController.CreateShirtController)
	admin.GET("/shirts", shirtController.GetAllShirtsController)
	admin.GET("/shirt/:id", shirtController.GetShirtByIdController)
	admin.GET("/shirt/code/:code", shirtController.GetShirtByCodeController)
	admin.PUT("/shirt/:id", shirtController.UpdateShirtByIdController)
	admin.DELETE("/shirt/:id", shirtController.DeleteShirtByIdController)

	// ====== Cart Admin =======
	admin.POST("/cart", cartController.CreateCartController)
	admin.GET("/carts", cartController.GetAllCartsController)
	admin.GET("/cart/:id", cartController.GetCartByIDController)
	admin.GET("/cart/:user_id", cartController.GetCartByUserIDController)
	admin.PUT("/cart/:id", cartController.UpdateCartByIdController)
	admin.DELETE("/cart/:id", cartController.DeleteCartByIdController)

	// ====== CART ADMIN ======
	admin.POST("/cart-item", cartItemController.AddCartItemController)
	admin.GET("/cart-items", cartItemController.GetAllCartItemsController)
	admin.GET("/cart-items/:cart_id", cartItemController.GetCartItemsByCartIDController)
	admin.PUT("/cart-item/:id", cartItemController.UpdateCartItemByIDController)
	admin.DELETE("/cart-item/:id", cartItemController.DeleteCartItemByIDController)

	// ====== ORDER ADMIN ======
	admin.POST("/order", orderController.CreateOrderController)
	admin.GET("/orders", orderController.GetAllOrdersController)
	admin.GET("/order/:id", orderController.GetOrderByIdController)
	admin.PUT("/order/:id", orderController.UpdateOrderByIdController)
	admin.DELETE("/order/:id", orderController.DeleteOrderByIdController)

	// ====== PAYMENT ADMIN ======
	admin.POST("/payment", paymentController.CreatePaymentController)
	admin.GET("/payments", paymentController.GetAllPaymentsController)
	admin.GET("/payment/:id", paymentController.GetPaymentByOrderIDController)
	admin.PUT("/payment/:id", paymentController.UpdatePaymentByIDController)
	admin.DELETE("/payment/:id", paymentController.DeletePaymentByIDController)

	// ====== PAYMENT HISTORY ADMIN ======
	admin.GET("/payment-history/:user_id", paymentHistoryController.GetPaymentHistoryByUserIDController)

	// ====== STUDENT ADMIN =======
	admin.POST("/student", studentController.CreateStudentController)
	admin.GET("/students", studentController.GetAllStudentsController)
	admin.GET("/students/by-parent-name", studentController.GetStudentsByParentNameController)
	admin.GET("/student/:id", studentController.GetStudentByIdController)
	admin.PUT("/student/:id", studentController.UpdateStudentByIdController)
	admin.DELETE("/student/:id", studentController.DeleteStudentByIdController)

	// ====== USER ADMIN =======
	admin.GET("/users", userController.GetAllUsersController)
	admin.GET("/user/:id", userController.GetUserByIdController)
	admin.GET("/user/email/:email", userController.GetUserByEmailController)
	admin.GET("/user/phone/:phone", userController.GetUserByPhoneController)
	admin.GET("/user/query", userController.GetUserByQueryController)
	admin.PUT("/user/:id", userController.UpdateUserByIdController)
	admin.DELETE("/user/:id", userController.DeleteUserByIdController)

	// ====== USER ROLE =======
	user.GET("/profile", userController.GetUserByIdController)
	user.PUT("/user", userController.GetUserByIdController)
	user.DELETE("/user", userController.DeleteUserByIdController)

	// ====== USER CART =======
	user.POST("/cart", cartController.CreateCartController)
	user.GET("/cart/:id", cartController.GetCartByUserIDController)
	user.PUT("cart/:id", cartController.UpdateCartByIdController)
	user.DELETE("/cart/:id", cartController.DeleteCartByIdController)

	// ====== USER CART ITEM =======
	user.POST("/cart-item", cartItemController.AddCartItemController)
	user.GET("/cart-items/:cart_id", cartItemController.GetCartItemsByCartIDController)
	user.PUT("/cart-item/:id", cartItemController.UpdateCartItemByIDController)
	user.DELETE("/cart-item/:id", cartItemController.DeleteCartItemByIDController)

	// ====== USER STUDENT =======
	user.POST("/student", studentController.CreateStudentController)
	user.GET("/student/:id", studentController.GetStudentByIdController)
	user.GET("/students", studentController.GetStudentsByParentNameController) // Ambil siswa berdasarkan parent_name

	// ====== USER PAYMENT =======
	user.POST("/payment", paymentController.CreatePaymentController)
	user.GET("/payment/:id", paymentController.GetPaymentByOrderIDController)
	user.DELETE("payment/:id", paymentController.DeletePaymentByIDController)

	// ====== USER PAYMENT HISTORY ======
	user.GET("/payment-history/:id", paymentHistoryController.GetPaymentHistoryByUserIDController)

	// ====== USER BOOK ======
	user.GET("/books", bookController.GetAllBooksController)
	user.GET("/book/:id", bookController.GetBookByIDController)

	// ====== USER SHIRT ======
	user.GET("/shirts", shirtController.GetAllShirtsController)
	user.GET("/shirt/:id", shirtController.GetShirtByIdController)

}
