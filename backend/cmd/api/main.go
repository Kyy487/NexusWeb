package main

import (
	"log"

	"nexusweb-market/backend/internal/config"
	"nexusweb-market/backend/internal/database"
	"nexusweb-market/backend/internal/middleware"
	"nexusweb-market/backend/internal/modules/auth"

	userHandler "nexusweb-market/backend/internal/modules/user/handler"
	userRepository "nexusweb-market/backend/internal/modules/user/repository"
	userService "nexusweb-market/backend/internal/modules/user/service"

	categoryHandler "nexusweb-market/backend/internal/modules/category/handler"
	categoryRepository "nexusweb-market/backend/internal/modules/category/repository"
	categoryService "nexusweb-market/backend/internal/modules/category/service"

	serviceHandler "nexusweb-market/backend/internal/modules/service/handler"
	serviceRepository "nexusweb-market/backend/internal/modules/service/repository"
	serviceService "nexusweb-market/backend/internal/modules/service/service"

	packageRepository "nexusweb-market/backend/internal/modules/package/repository"
	packageService "nexusweb-market/backend/internal/modules/package/service"
	packageHandler "nexusweb-market/backend/internal/modules/package/handler"

	orderRepository "nexusweb-market/backend/internal/modules/order/repository"
	orderService "nexusweb-market/backend/internal/modules/order/service"
	orderHandler "nexusweb-market/backend/internal/modules/order/handler"

	requirementRepository "nexusweb-market/backend/internal/modules/orderrequirement/repository"
	requirementService "nexusweb-market/backend/internal/modules/orderrequirement/service"
	requirementHandler "nexusweb-market/backend/internal/modules/orderrequirement/handler"

	progressRepository "nexusweb-market/backend/internal/modules/orderprogress/repository"
	progressService "nexusweb-market/backend/internal/modules/orderprogress/service"
	progressHandler "nexusweb-market/backend/internal/modules/orderprogress/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewPostgresConnection(cfg)
	defer db.Close()

	r := gin.Default()

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
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

	orderRepo := orderRepository.NewOrderRepository(db)
	orderSvc := orderService.NewOrderService(orderRepo)
	orderHdl := orderHandler.NewOrderHandler(orderSvc)	

	requirementRepo := requirementRepository.NewRequirementRepository(db)
	requirementSvc := requirementService.NewRequirementService(requirementRepo)
	requirementHdl := requirementHandler.NewRequirementHandler(requirementSvc)

	progressRepo := progressRepository.NewProgressRepository(db)
	progressSvc := progressService.NewProgressService(progressRepo)
	progressHdl := progressHandler.NewProgressHandler(progressSvc)

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
				users.GET("", userHdl.GetAllUsers)
				users.PATCH("/:id/status", userHdl.UpdateUserStatus)
			}
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHdl.GetAll)
				categories.GET("/:id", categoryHdl.GetByID)
				categories.POST("", categoryHdl.Create)
				categories.PUT("/:id", categoryHdl.Update)
				categories.DELETE("/:id", categoryHdl.Delete)
			}
			services := protected.Group("/services")
			{
				services.GET("", serviceHdl.GetAll)
				services.GET("/:id", serviceHdl.GetByID)
				services.POST("", serviceHdl.Create)
				services.PUT("/:id", serviceHdl.Update)
				services.DELETE("/:id", serviceHdl.Delete)
			}
			packages := protected.Group("/packages")
			{
				packages.GET("", packageHdl.GetAll)
				packages.GET("/:id", packageHdl.GetByID)
				packages.POST("", packageHdl.Create)
				packages.PUT("/:id", packageHdl.Update)
				packages.DELETE("/:id", packageHdl.Delete)
			}
			orders := protected.Group("/orders")
			{
				orders.GET("", orderHdl.GetAll)
				orders.GET("/:id", orderHdl.GetByID)
				orders.POST("", orderHdl.Create)
				orders.PATCH("/:id/status", orderHdl.UpdateStatus)
			}
			orderRequirements := protected.Group("/order-requirements")
			{
				orderRequirements.GET("/order/:orderId", requirementHdl.GetByOrderID)
				orderRequirements.POST("/order/:orderId", requirementHdl.Create)
				orderRequirements.GET("/:id", requirementHdl.GetByID)
				orderRequirements.PUT("/:id", requirementHdl.Update)
				orderRequirements.DELETE("/:id", requirementHdl.Delete)
			}
			orderProgress := protected.Group("/order-progress")
			{
				orderProgress.GET("/order/:orderId", progressHdl.GetByOrderID)
				orderProgress.POST("/order/:orderId", progressHdl.Create)
				orderProgress.GET("/:id", progressHdl.GetByID)
				orderProgress.PUT("/:id", progressHdl.Update)
				orderProgress.DELETE("/:id", progressHdl.Delete)
			}
		}
	}

	port := ":" + cfg.AppPort
	log.Println("server running on port", port)

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}