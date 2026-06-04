package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"backend/config"
	"backend/internal/handlers"
	"backend/internal/models"
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

	// Add CORS middleware
	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Set("Access-Control-Allow-Origin", origin)
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Set("Access-Control-Allow-Credentials", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
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
	configureDatabase(db)

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
	routes.SetupRoutes(app, db, authHandler, deptHandler, roleHandler, userHandler, custCatHandler, custHandler, custContactHandler, policyHandler, menuAccessHandler, paymentTermHandler, projectLevelHandler, projectPriorityHandler, quotationProgressHandler, quotationStatusHandler, unitHandler, quotationHandler, quotationFollowupHandler, stockHandler, productCategoryHandler, brandHandler, kanbanBoardHandler, kanbanListHandler, kanbanCardHandler, kanbanLabelHandler, kanbanChecklistHandler, kanbanAttachmentHandler, kanbanCommentHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5500"
	}

	// certDir := "./cert"
	// log.Printf("Server starting on port %s (HTTPS) for app.magnumsolusion.co.id", port)
	// log.Fatal(app.ListenTLS(":"+port, filepath.Join(certDir, "cert.pem"), filepath.Join(certDir, "key.pem")))

	log.Printf("Server starting on port %s (HTTP)", port)
	log.Fatal(app.Listen(":" + port))
}

func configureDatabase(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.MasterTableAccess{},
		&models.UserGroupPolicy{},
		&models.GroupPolicy{},
		&models.Departement{},
		&models.UserGroup{},
		&models.CustomerCategory{},
		&models.Customer{},
		&models.CustomerContact{},
		&models.PaymentTerm{},
		&models.ProjectLevel{},
		&models.ProjectPriority{},
		&models.QuotationProgress{},
		&models.QuotationStatus{},
		&models.Quotation{},
		&models.QuotationDetail{},
		&models.QuotationSubdetail{},
		&models.QuotationMaster{},
		&models.QuotationFollowup{},
		// Kanban models
		&models.KanbanBoard{},
		&models.KanbanList{},
		&models.KanbanCard{},
		&models.KanbanLabel{},
		&models.KanbanChecklist{},
		&models.KanbanChecklistItem{},
		&models.KanbanAttachment{},
		&models.KanbanComment{},
	)

	// Manual migration for composite primary keys and ID size
	// GORM AutoMigrate often fails to change existing PKs in MySQL
	db.Exec("ALTER TABLE quotation MODIFY id VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	db.Exec("ALTER TABLE quotation_detail MODIFY id VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	db.Exec("ALTER TABLE quotation_subdetail MODIFY id VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	db.Exec("ALTER TABLE quotation_master MODIFY id VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")
	db.Exec("ALTER TABLE quotation_followup MODIFY id VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci")

	// We use IGNORE or check existence to prevent errors if already migrated
	db.Exec("ALTER TABLE quotation_detail DROP PRIMARY KEY, ADD PRIMARY KEY (id, rev_id, line)")
	db.Exec("ALTER TABLE quotation_subdetail DROP PRIMARY KEY, ADD PRIMARY KEY (id, rev_id, line, subline)")
	db.Exec("ALTER TABLE quotation_master DROP PRIMARY KEY, ADD PRIMARY KEY (id, rev_id)")

	// Recreate stored procedure to ensure the 'qid' parameter uses the exact same collation and sizing (VARCHAR 20)
	// This prevents the "Illegal mix of collations" error inside the trigger when inserting/updating/deleting follow-ups.
	db.Exec("DROP PROCEDURE IF EXISTS update_quotation_status")
	db.Exec(`
		CREATE PROCEDURE update_quotation_status(IN qid VARCHAR(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)
		BEGIN
		  DECLARE q_status, q_progress, foll_by, count INT DEFAULT 0;
		  DECLARE foll_date date DEFAULT null;
		  DECLARE next_foll date DEFAULT null;
		  DECLARE f_po_date varchar(10) DEFAULT null;
		  DECLARE f_po_no varchar(255) DEFAULT null;
		  DECLARE f_po_file varchar(255) DEFAULT null;
		  
		  SELECT count(*) INTO count FROM quotation_followup WHERE id=qid;
		  IF count > 0 THEN
			  SELECT status, progress, followup_date, followup_by, next_followup, po_no, po_date, po_file 
			  INTO q_status, q_progress, foll_date, foll_by, next_foll, f_po_no, f_po_date, f_po_file 
			  FROM quotation_followup 
			  WHERE id=qid AND line_id IN (SELECT MAX(line_id) FROM quotation_followup WHERE id=qid AND followup_date IN (SELECT MAX(followup_date) FROM quotation_followup WHERE id=qid));
			  
			  IF f_po_date = "0000-00-00" THEN 
			   SET f_po_date = NULL;
			  END IF;
			 
			  UPDATE quotation SET status=q_status, progress=q_progress, followup_date=foll_date, followup_by=foll_by, next_followup=next_foll, po_no=f_po_no, po_date=f_po_date, po_file=f_po_file WHERE id=qid;
		  ELSE
			   SELECT status, progress, quotation_date, followup_by, next_followup, po_no, po_date, po_file 
			   INTO q_status, q_progress, foll_date, foll_by, next_foll, f_po_no, f_po_date, f_po_file 
			   FROM quotation WHERE id=qid;
			   
			   IF f_po_date = "0000-00-00" THEN 
				  SET f_po_date = NULL;
			   END IF;
			   
			   UPDATE quotation SET status=q_status, progress=q_progress, followup_date=foll_date, followup_by=foll_by, next_followup=next_foll, po_no=f_po_no, po_date=f_po_date, po_file=f_po_file WHERE id=qid;
		  END IF;
		END
	`)

	log.Println("Manual database migrations completed successfully")

	// Seed Kanban menu
	var kanbanMenu models.MasterTableAccess
	kanbanResult := db.Where("name = ?", "Kanban Board").First(&kanbanMenu)
	if kanbanResult.Error == gorm.ErrRecordNotFound {
		kanbanIcon := "columns-3"
		kanbanSort := 45
		kanbanPath := "/kanban"
		kanbanEp := "kanban-boards"
		db.Create(&models.MasterTableAccess{
			Name:      "Kanban Board",
			MenuName:  "Kanban Board",
			Path:      &kanbanPath,
			Endpoint:  &kanbanEp,
			Icon:      &kanbanIcon,
			SortOrder: kanbanSort,
			IsActive:  true,
		})
	}

	// Fix existing '0000-00-00' date values that cause MySQL error 1292
	// Temporarily disable strict mode to allow comparison with '0000-00-00' on DATE columns
	db.Exec("SET SESSION sql_mode = ''")
	db.Exec("UPDATE quotation SET project_start = NULL WHERE project_start = '0000-00-00'")
	db.Exec("UPDATE quotation SET project_end = NULL WHERE project_end = '0000-00-00'")
	db.Exec("UPDATE quotation SET quotation_date = NULL WHERE quotation_date = '0000-00-00'")
	db.Exec("UPDATE quotation SET valid_until = NULL WHERE valid_until = '0000-00-00'")
	db.Exec("UPDATE quotation SET followup_date = NULL WHERE followup_date = '0000-00-00'")
	db.Exec("UPDATE quotation SET next_followup = NULL WHERE next_followup = '0000-00-00'")
	db.Exec("UPDATE quotation SET po_date = NULL WHERE po_date = '0000-00-00'")
	db.Exec("UPDATE quotation_master SET project_start = NULL WHERE project_start = '0000-00-00'")
	db.Exec("UPDATE quotation_master SET project_end = NULL WHERE project_end = '0000-00-00'")
	db.Exec("UPDATE quotation_master SET quotation_date = NULL WHERE quotation_date = '0000-00-00'")
	db.Exec("UPDATE quotation_master SET valid_until = NULL WHERE valid_until = '0000-00-00'")
	db.Exec("UPDATE quotation_master SET po_date = NULL WHERE po_date = '0000-00-00'")
	db.Exec("UPDATE quotation_followup SET followup_date = NULL WHERE followup_date = '0000-00-00'")
	db.Exec("UPDATE quotation_followup SET next_followup = NULL WHERE next_followup = '0000-00-00'")
	db.Exec("UPDATE quotation_followup SET po_date = NULL WHERE po_date = '0000-00-00'")

	// Seed MasterTableAccess with endpoints
	apiEndpoints := map[string]string{
		"departements":        "departements",
		"users":               "users",
		"usergroups":          "roles",
		"property":            "property",
		"usergroupspolicies":  "policies",
		"customers":           "customers",
		"customer category":   "customer-categories",
		"payment term":        "payment-terms",
		"project level":       "project-levels",
		"project priority":    "project-priorities",
		"quotation progress":  "quotation-progress",
		"quotation status":    "quotation-statuses",
		"units":               "units",
		"quotations":          "quotations",
		"quotation followups": "quotation-followups",
	}

	for name, endpoint := range apiEndpoints {
		var table models.MasterTableAccess
		result := db.Where("name = ?", name).First(&table)

		ep := endpoint // copy to addressable variable
		if result.Error == gorm.ErrRecordNotFound {
			db.Create(&models.MasterTableAccess{Name: name, MenuName: name, Endpoint: &ep, IsActive: true})
		} else {
			if table.Endpoint == nil || *table.Endpoint != endpoint {
				table.Endpoint = &ep
				db.Save(&table)
			}
		}
	}

	// Seed Reports parent menu
	var reportsMenu models.MasterTableAccess
	reportsResult := db.Where("name = ?", "Reports").First(&reportsMenu)
	if reportsResult.Error == gorm.ErrRecordNotFound {
		icon := "bar-chart-3"
		sortOrder := 50
		db.Create(&models.MasterTableAccess{
			Name:      "Reports",
			MenuName:  "Reports",
			ParentID:  nil,
			Path:      nil,
			Endpoint:  nil,
			Icon:      &icon,
			SortOrder: sortOrder,
			IsActive:  true,
		})
		db.Where("name = ?", "Reports").First(&reportsMenu)
	}

	// Seed report child menus
	type reportChild struct {
		Name     string
		Path     string
		Endpoint string
		Icon     string
		Sort     int
	}

	reportChildren := []reportChild{
		{"Quotations Report", "/reports/quotations-report", "quotations", "file-text", 1},
		{"Top Customers", "/reports/top-customers", "quotations", "award", 2},
		{"Followups Report", "/reports/followups-report", "quotation-followups", "clipboard-list", 3},
		{"Sales Summary By Customer", "/reports/sales-summary-by-customer", "quotations", "bar-chart-3", 4},
		{"Sales Summary By Sales Person", "/reports/sales-summary-by-sales-person", "quotations", "bar-chart-3", 5},
		{"Sales Charts By Customer", "/reports/sales-charts-by-customer", "quotations", "bar-chart-3", 6},
		{"Sales Charts By Sales Person", "/reports/sales-charts-by-sales-person", "quotations", "bar-chart-3", 7},
		{"Sales Detail By Customer", "/reports/sales-detail-by-customer", "quotations", "list", 8},
	}

	for _, rc := range reportChildren {
		var child models.MasterTableAccess
		err := db.Where("name = ?", rc.Name).First(&child).Error
		if err == gorm.ErrRecordNotFound {
			path := rc.Path
			ep := rc.Endpoint
			icon := rc.Icon
			db.Create(&models.MasterTableAccess{
				Name:      rc.Name,
				MenuName:  rc.Name,
				ParentID:  &reportsMenu.ID,
				Path:      &path,
				Endpoint:  &ep,
				Icon:      &icon,
				SortOrder: rc.Sort,
				IsActive:  true,
			})
		}
	}
}
