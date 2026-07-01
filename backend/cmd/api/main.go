package main

import (
	"log"
	"time"

	"nexusweb-market/backend/internal/config"
	"nexusweb-market/backend/internal/database"
	"nexusweb-market/backend/internal/middleware"
	"nexusweb-market/backend/internal/modules/auth"

	"github.com/gin-contrib/cors"

	userHandler "nexusweb-market/backend/internal/modules/user/handler"
	userRepository "nexusweb-market/backend/internal/modules/user/repository"
	userService "nexusweb-market/backend/internal/modules/user/service"

	categoryHandler "nexusweb-market/backend/internal/modules/category/handler"
	categoryRepository "nexusweb-market/backend/internal/modules/category/repository"
	categoryService "nexusweb-market/backend/internal/modules/category/service"

	serviceHandler "nexusweb-market/backend/internal/modules/service/handler"
	serviceRepository "nexusweb-market/backend/internal/modules/service/repository"
	serviceService "nexusweb-market/backend/internal/modules/service/service"

	packageHandler "nexusweb-market/backend/internal/modules/package/handler"
	packageRepository "nexusweb-market/backend/internal/modules/package/repository"
	packageService "nexusweb-market/backend/internal/modules/package/service"

	orderHandler "nexusweb-market/backend/internal/modules/order/handler"
	orderRepository "nexusweb-market/backend/internal/modules/order/repository"
	orderService "nexusweb-market/backend/internal/modules/order/service"

	requirementHandler "nexusweb-market/backend/internal/modules/orderrequirement/handler"
	requirementRepository "nexusweb-market/backend/internal/modules/orderrequirement/repository"
	requirementService "nexusweb-market/backend/internal/modules/orderrequirement/service"

	progressHandler "nexusweb-market/backend/internal/modules/orderprogress/handler"
	progressRepository "nexusweb-market/backend/internal/modules/orderprogress/repository"
	progressService "nexusweb-market/backend/internal/modules/orderprogress/service"

	invoiceHandler "nexusweb-market/backend/internal/modules/invoice/handler"
	invoiceRepository "nexusweb-market/backend/internal/modules/invoice/repository"
	invoiceService "nexusweb-market/backend/internal/modules/invoice/service"

	paymentHandler "nexusweb-market/backend/internal/modules/payment/handler"
	paymentRepository "nexusweb-market/backend/internal/modules/payment/repository"
	paymentService "nexusweb-market/backend/internal/modules/payment/service"

	dashboardHandler "nexusweb-market/backend/internal/modules/dashboard/handler"
	dashboardRepository "nexusweb-market/backend/internal/modules/dashboard/repository"
	dashboardService "nexusweb-market/backend/internal/modules/dashboard/service"

	orderFileHandler "nexusweb-market/backend/internal/modules/orderfile/handler"
	orderFileRepository "nexusweb-market/backend/internal/modules/orderfile/repository"
	orderFileService "nexusweb-market/backend/internal/modules/orderfile/service"

	activityLogHandler "nexusweb-market/backend/internal/modules/activitylog/handler"
	activityLogRepository "nexusweb-market/backend/internal/modules/activitylog/repository"
	activityLogService "nexusweb-market/backend/internal/modules/activitylog/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000", "http://127.0.0.1:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg, nil)
	authHandler := auth.NewHandler(authService)

	userRepo := userRepository.NewUserRepository(db)
	userSvc := userService.NewUserService(userRepo)
	userHdl := userHandler.NewUserHandler(userSvc)

	categoryRepo := categoryRepository.NewCategoryRepository(db)
	categorySvc := categoryService.NewCategoryService(categoryRepo)
	categoryHdl := categoryHandler.NewCategoryHandler(categorySvc)

	serviceRepo := serviceRepository.NewServiceRepository(db)
	serviceSvc := serviceService.NewServiceService(serviceRepo)
	serviceHdl := serviceHandler.NewServiceHandler(serviceSvc)

	packageRepo := packageRepository.NewPackageRepository(db)
	packageSvc := packageService.NewPackageService(packageRepo)
	packageHdl := packageHandler.NewPackageHandler(packageSvc)

	activityLogRepo := activityLogRepository.NewActivityLogRepository(db)
	activityLogSvc := activityLogService.NewActivityLogService(activityLogRepo)
	activityLogHdl := activityLogHandler.NewActivityLogHandler(activityLogSvc)

	orderRepo := orderRepository.NewOrderRepository(db)
	orderSvc := orderService.NewOrderService(orderRepo)
	orderHdl := orderHandler.NewOrderHandler(orderSvc, activityLogSvc)

	requirementRepo := requirementRepository.NewRequirementRepository(db)
	requirementSvc := requirementService.NewRequirementService(requirementRepo)
	requirementHdl := requirementHandler.NewRequirementHandler(requirementSvc)

	progressRepo := progressRepository.NewProgressRepository(db)
	progressSvc := progressService.NewProgressService(progressRepo)
	progressHdl := progressHandler.NewProgressHandler(progressSvc)

	invoiceRepo := invoiceRepository.NewInvoiceRepository(db)
	invoiceSvc := invoiceService.NewInvoiceService(invoiceRepo)
	invoiceHdl := invoiceHandler.NewInvoiceHandler(invoiceSvc, activityLogSvc)

	paymentRepo := paymentRepository.NewPaymentRepository(db)
	paymentSvc := paymentService.NewPaymentService(paymentRepo)
	paymentHdl := paymentHandler.NewPaymentHandler(paymentSvc, activityLogSvc)

	dashboardRepo := dashboardRepository.NewDashboardRepository(db)
	dashboardSvc := dashboardService.NewDashboardService(dashboardRepo)
	dashboardHdl := dashboardHandler.NewDashboardHandler(dashboardSvc)

	orderFileRepo := orderFileRepository.NewOrderFileRepository(db)
	orderFileSvc := orderFileService.NewOrderFileService(orderFileRepo)
	orderFileHdl := orderFileHandler.NewOrderFileHandler(orderFileSvc)

	authService = auth.NewService(authRepo, cfg, activityLogSvc)
	authHandler = auth.NewHandler(authService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "NexusWeb API running",
		})
	})

	api := r.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/protected/me", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"user_id": c.GetString("user_id"),
					"email":   c.GetString("email"),
					"role":    c.GetString("role"),
				})
			})

			users := protected.Group("/users")
			{
				users.GET("/me", userHdl.GetMe)
				users.GET("", middleware.AdminOnly(), userHdl.GetAllUsers)
				users.PATCH("/:id/status", middleware.AdminOnly(), userHdl.UpdateUserStatus)
			}
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHdl.GetAll)
				categories.GET("/:id", categoryHdl.GetByID)
				categories.POST("", middleware.AdminOnly(), categoryHdl.Create)
				categories.PUT("/:id", middleware.AdminOnly(), categoryHdl.Update)
				categories.DELETE("/:id", middleware.AdminOnly(), categoryHdl.Delete)
			}
			services := protected.Group("/services")
			{
				services.GET("", serviceHdl.GetAll)
				services.GET("/:id", serviceHdl.GetByID)
				services.POST("", middleware.AdminOnly(), serviceHdl.Create)
				services.PUT("/:id", middleware.AdminOnly(), serviceHdl.Update)
				services.DELETE("/:id", middleware.AdminOnly(), serviceHdl.Delete)
			}
			packages := protected.Group("/packages")
			{
				packages.GET("", packageHdl.GetAll)
				packages.GET("/:id", packageHdl.GetByID)
				packages.POST("", middleware.AdminOnly(), packageHdl.Create)
				packages.PUT("/:id", middleware.AdminOnly(), packageHdl.Update)
				packages.DELETE("/:id", middleware.AdminOnly(), packageHdl.Delete)
			}
			orders := protected.Group("/orders")
			{
				orders.GET("", middleware.AdminOnly(), orderHdl.GetAll)
				orders.GET("/:id", middleware.AdminOnly(), orderHdl.GetByID)
				orders.POST("", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), orderHdl.Create)
				orders.PATCH("/:id/status", middleware.AdminOnly(), orderHdl.UpdateStatus)
			}
			orderRequirements := protected.Group("/order-requirements")
			{
				orderRequirements.GET("/order/:orderId", middleware.AdminOnly(), requirementHdl.GetByOrderID)
				orderRequirements.POST("/order/:orderId", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), requirementHdl.Create)
				orderRequirements.GET("/:id", middleware.AdminOnly(), requirementHdl.GetByID)
				orderRequirements.PUT("/:id", middleware.AdminOnly(), requirementHdl.Update)
				orderRequirements.DELETE("/:id", middleware.AdminOnly(), requirementHdl.Delete)
			}
			orderProgress := protected.Group("/order-progress")
			{
				orderProgress.GET("/order/:orderId", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), progressHdl.GetByOrderID)
				orderProgress.POST("/order/:orderId", middleware.AdminOnly(), progressHdl.Create)
				orderProgress.GET("/:id", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), progressHdl.GetByID)
				orderProgress.PUT("/:id", middleware.AdminOnly(), progressHdl.Update)
				orderProgress.DELETE("/:id", middleware.AdminOnly(), progressHdl.Delete)
			}
			invoices := protected.Group("/invoices")
			{
				invoices.GET("", middleware.AdminOnly(), invoiceHdl.GetAll)
				invoices.GET("/order/:orderId", middleware.AdminOnly(), invoiceHdl.GetByOrderID)
				invoices.GET("/:id", middleware.AdminOnly(), invoiceHdl.GetByID)
				invoices.POST("", middleware.AdminOnly(), invoiceHdl.Create)
				invoices.PATCH("/:id/status", middleware.AdminOnly(), invoiceHdl.UpdateStatus)
			}
			payments := protected.Group("/payments")
			{
				payments.GET("", middleware.AdminOnly(), paymentHdl.GetAll)
				payments.GET("/invoice/:invoiceId", middleware.AdminOnly(), paymentHdl.GetByInvoiceID)
				payments.GET("/:id/whatsapp", middleware.AdminOnly(), paymentHdl.GetWhatsAppLink)
				payments.GET("/:id", middleware.AdminOnly(), paymentHdl.GetByID)
				payments.POST("", middleware.AdminOnly(), paymentHdl.Create)
				payments.PATCH("/:id/status", middleware.AdminOnly(), paymentHdl.UpdateStatus)
			}
			my := protected.Group("/my")
			{
				my.GET("/orders", orderHdl.GetByCustomerID)
				my.GET("/invoices", invoiceHdl.GetByCustomerID)
				my.GET("/payments", paymentHdl.GetByCustomerID)
			}
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", middleware.AdminOnly(), dashboardHdl.GetStats)
			}
			r.Static("/uploads", "./uploads")

			orderFiles := protected.Group("/order-files")
			{
				orderFiles.POST("/upload", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), orderFileHdl.Upload)
				orderFiles.GET("/order/:orderId", middleware.RequireRoles(middleware.RoleSuperAdmin, middleware.RoleAdmin, middleware.RoleCustomer), orderFileHdl.GetByOrderID)
				orderFiles.DELETE("/:id", middleware.AdminOnly(), orderFileHdl.Delete)
			}
			activityLogs := protected.Group("/activity-logs")
			{
				activityLogs.GET("", middleware.AdminOnly(), activityLogHdl.GetAll)
				activityLogs.GET("/user/:userId", middleware.AdminOnly(), activityLogHdl.GetByUserID)
				activityLogs.POST("", middleware.AdminOnly(), activityLogHdl.Create)
			}
		}
	}

	port := ":" + cfg.AppPort
	log.Println("server running on port", port)

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
