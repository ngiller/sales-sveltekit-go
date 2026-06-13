package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"backend/config"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env vars")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: config.ErrorHandler,
	})

	app.Use(recover.New())

	// Serve static files
	app.Static("/uploads", "./uploads")

	// Add CORS middleware with origin whitelist
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")

		allowedOrigins := map[string]bool{}
		corsOrigins := os.Getenv("CORS_ORIGINS")
		if corsOrigins == "" {
			corsOrigins = "http://localhost:5173,http://localhost:5500,https://app.magnumsolusion.co.id,https://app.magnumsolusion.co.id:5500,https://sales.magnumsolusion.co.id"
		}
		for _, o := range strings.Split(corsOrigins, ",") {
			allowedOrigins[strings.TrimSpace(o)] = true
		}

		if origin != "" && allowedOrigins[origin] {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
		} else if origin == "" {
			c.Set("Access-Control-Allow-Origin", "")
		} else {
			c.Set("Access-Control-Allow-Origin", "null")
		}

		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Set("Access-Control-Max-Age", "86400")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// Security headers middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "0")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Set("Cross-Origin-Resource-Policy", "same-origin")
		return c.Next()
	})

	// Serve frontend static files from public/ directory (SPA)
	publicDir := "./public"
	app.Use(func(c *fiber.Ctx) error {
		// Skip API and uploads routes
		if strings.HasPrefix(c.Path(), "/api/") || strings.HasPrefix(c.Path(), "/uploads/") {
			return c.Next()
		}

		// Get absolute path to public directory
		absPublicDir, _ := filepath.Abs(publicDir)

		// Try to serve the exact file
		reqPath := c.Path()
		if reqPath == "/" {
			reqPath = "/index.html"
		}
		filePath := filepath.Join(absPublicDir, reqPath)

		// Security check: ensure the path is within public directory
		cleanPath, _ := filepath.Abs(filePath)
		if !strings.HasPrefix(cleanPath, absPublicDir) {
			return c.Status(400).SendString("Bad request")
		}

		// Check if file exists
		if info, err := os.Stat(cleanPath); err == nil {
			if !info.IsDir() {
				return c.SendFile(cleanPath)
			}

			// If it's a directory, check for index.html inside it
			indexPath := filepath.Join(cleanPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				return c.SendFile(indexPath)
			}
		}

		// SPA fallback: serve root index.html for all other routes
		return c.SendFile(filepath.Join(absPublicDir, "index.html"))
	})

	db := config.InitDB()

	userRepo := repository.NewUserRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	menuAccessRepo := repository.NewMenuAccessRepository(db)
	deptRepo := repository.NewDepartementRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	custCatRepo := repository.NewCustomerCategoryRepository(db)
	custRepo := repository.NewCustomerRepository(db)
	custContactRepo := repository.NewCustomerContactRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	paymentTermRepo := repository.NewPaymentTermRepository(db)
	projectLevelRepo := repository.NewProjectLevelRepository(db)
	projectPriorityRepo := repository.NewProjectPriorityRepository(db)
	quotationProgressRepo := repository.NewQuotationProgressRepository(db)
	quotationStatusRepo := repository.NewQuotationStatusRepository(db)
	quotationRepo := repository.NewQuotationRepository(db)
	quotationFollowupRepo := repository.NewQuotationFollowupRepository(db)

	// Dual database for stock
	stockDb := config.InitStockDB()
	stockRepo := repository.NewStockRepository(stockDb)
	productCategoryRepo := repository.NewProductCategoryRepository(stockDb)
	brandRepo := repository.NewBrandRepository(stockDb)
	unitRepo := repository.NewUnitRepository(stockDb)

	authHandler := handlers.NewAuthHandler(userRepo, menuRepo)
	deptHandler := handlers.NewDepartementHandler(deptRepo)
	roleHandler := handlers.NewRoleHandler(roleRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	custCatHandler := handlers.NewCustomerCategoryHandler(custCatRepo)
	custHandler := handlers.NewCustomerHandler(custRepo)
	custContactHandler := handlers.NewCustomerContactHandler(custContactRepo)
	policyHandler := handlers.NewPolicyHandler(policyRepo)
	menuAccessHandler := handlers.NewMenuAccessHandler(menuAccessRepo, menuRepo, db)
	paymentTermHandler := handlers.NewPaymentTermHandler(paymentTermRepo)
	projectLevelHandler := handlers.NewProjectLevelHandler(projectLevelRepo)
	projectPriorityHandler := handlers.NewProjectPriorityHandler(projectPriorityRepo)
	quotationProgressHandler := handlers.NewQuotationProgressHandler(quotationProgressRepo)
	quotationStatusHandler := handlers.NewQuotationStatusHandler(quotationStatusRepo)
	unitHandler := handlers.NewUnitHandler(unitRepo)
	quotationHandler := handlers.NewQuotationHandler(quotationRepo)
	quotationFollowupHandler := handlers.NewQuotationFollowupHandler(quotationFollowupRepo)
	stockHandler := handlers.NewStockHandler(stockRepo)
	productCategoryHandler := handlers.NewProductCategoryHandler(productCategoryRepo)
	brandHandler := handlers.NewBrandHandler(brandRepo)

	settingRepo := repository.NewSettingRepository(db)
	settingHandler := handlers.NewSettingHandler(settingRepo)

	// Kanban repositories
	kanbanBoardRepo := repository.NewKanbanBoardRepository(db)
	kanbanListRepo := repository.NewKanbanListRepository(db)
	kanbanCardRepo := repository.NewKanbanCardRepository(db)
	kanbanLabelRepo := repository.NewKanbanLabelRepository(db)
	kanbanChecklistRepo := repository.NewKanbanChecklistRepository(db)
	kanbanAttachmentRepo := repository.NewKanbanAttachmentRepository(db)
	kanbanCommentRepo := repository.NewKanbanCommentRepository(db)

	// Kanban handlers
	kanbanBoardHandler := handlers.NewKanbanBoardHandler(kanbanBoardRepo, kanbanListRepo, kanbanCardRepo)
	kanbanListHandler := handlers.NewKanbanListHandler(kanbanListRepo, kanbanBoardRepo)
	kanbanCardHandler := handlers.NewKanbanCardHandler(kanbanCardRepo, kanbanListRepo)
	kanbanLabelHandler := handlers.NewKanbanLabelHandler(kanbanLabelRepo)
	kanbanChecklistHandler := handlers.NewKanbanChecklistHandler(kanbanChecklistRepo)
	kanbanAttachmentHandler := handlers.NewKanbanAttachmentHandler(kanbanAttachmentRepo)
	kanbanCommentHandler := handlers.NewKanbanCommentHandler(kanbanCommentRepo)

	// Setup all routing
	routes.SetupRoutes(app, db, authHandler, deptHandler, roleHandler, userHandler, custCatHandler, custHandler, custContactHandler, policyHandler, menuAccessHandler, paymentTermHandler, projectLevelHandler, projectPriorityHandler, quotationProgressHandler, quotationStatusHandler, unitHandler, quotationHandler, quotationFollowupHandler, settingHandler, stockHandler, productCategoryHandler, brandHandler, kanbanBoardHandler, kanbanListHandler, kanbanCardHandler, kanbanLabelHandler, kanbanChecklistHandler, kanbanAttachmentHandler, kanbanCommentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	// Disable TLS for development by setting TLS_DISABLE=true
	if os.Getenv("TLS_DISABLE") == "true" {
		log.Printf("Server starting on port %s (HTTP - development mode)", port)
		log.Fatal(app.Listen(":" + port))
	} else {
		certDir := "./cert"
		log.Printf("Server starting on port %s (HTTPS) for app.magnumsolusion.co.id", port)
		log.Fatal(app.ListenTLS(":"+port, filepath.Join(certDir, "cert.pem"), filepath.Join(certDir, "key.pem")))
	}
}
