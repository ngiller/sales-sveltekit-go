package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes registers all API endpoints for the application
func SetupRoutes(app *fiber.App, db *gorm.DB, authHandler *handlers.AuthHandler, deptHandler *handlers.DepartementHandler, roleHandler *handlers.RoleHandler, userHandler *handlers.UserHandler, customerCategoryHandler *handlers.CustomerCategoryHandler, customerHandler *handlers.CustomerHandler, customerContactHandler *handlers.CustomerContactHandler, policyHandler *handlers.PolicyHandler, menuAccessHandler *handlers.MenuAccessHandler, paymentTermHandler *handlers.PaymentTermHandler, projectLevelHandler *handlers.ProjectLevelHandler, projectPriorityHandler *handlers.ProjectPriorityHandler, quotationProgressHandler *handlers.QuotationProgressHandler, quotationStatusHandler *handlers.QuotationStatusHandler, unitHandler *handlers.UnitHandler, quotationHandler *handlers.QuotationHandler, quotationFollowupHandler *handlers.QuotationFollowupHandler, stockHandler *handlers.StockHandler, productCategoryHandler *handlers.ProductCategoryHandler, brandHandler *handlers.BrandHandler, kanbanBoardHandler *handlers.KanbanBoardHandler, kanbanListHandler *handlers.KanbanListHandler, kanbanCardHandler *handlers.KanbanCardHandler, kanbanLabelHandler *handlers.KanbanLabelHandler, kanbanChecklistHandler *handlers.KanbanChecklistHandler, kanbanAttachmentHandler *handlers.KanbanAttachmentHandler, kanbanCommentHandler *handlers.KanbanCommentHandler) {

	// API Group
	api := app.Group("/api")

	// Live Stocks Routes
	api.Get("/live-stocks", middleware.AuthMiddleware(), stockHandler.GetAllProducts)
	api.Get("/live-stocks/export", middleware.AuthMiddleware(), stockHandler.ExportToExcel)

	// Product Categories Routes (from stock DB)
	api.Get("/product-categories", middleware.AuthMiddleware(), productCategoryHandler.FindAll)

	// Brands Routes (from stock DB)
	api.Get("/brands", middleware.AuthMiddleware(), brandHandler.FindAll)

	// Auth Routes
	api.Post("/login", authHandler.Login)
	api.Post("/logout", authHandler.Logout)
	api.Get("/profile", middleware.AuthMiddleware(), authHandler.Profile)
	api.Put("/profile/:id", middleware.AuthMiddleware(), authHandler.UpdateProfile)
	api.Put("/change-password/:id", middleware.AuthMiddleware(), authHandler.ChangePassword)

	// Users Routes
	userGrp := api.Group("/users", middleware.AuthMiddleware())
	userGrp.Get("/", middleware.RequirePolicy(db, "read"), userHandler.FindAll)
	userGrp.Get("/:id", middleware.RequirePolicy(db, "read"), userHandler.FindByID)
	userGrp.Post("/", middleware.RequirePolicy(db, "create"), userHandler.Create)
	userGrp.Put("/:id", middleware.RequirePolicy(db, "update"), userHandler.Update)
	userGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), userHandler.Delete)
	userGrp.Post("/:id/avatar", middleware.RequirePolicy(db, "update"), userHandler.UploadAvatar)
	userGrp.Post("/:id/signature", middleware.RequirePolicy(db, "update"), userHandler.UploadSignature)

	// Customer Categories Routes
	custCatGrp := api.Group("/customer-categories", middleware.AuthMiddleware())
	custCatGrp.Get("/", middleware.RequirePolicy(db, "read"), customerCategoryHandler.FindAll)
	custCatGrp.Get("/:id", middleware.RequirePolicy(db, "read"), customerCategoryHandler.FindByID)
	custCatGrp.Post("/", middleware.RequirePolicy(db, "create"), customerCategoryHandler.Create)
	custCatGrp.Put("/:id", middleware.RequirePolicy(db, "update"), customerCategoryHandler.Update)
	custCatGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), customerCategoryHandler.Delete)

	// Customers Routes
	custGrp := api.Group("/customers", middleware.AuthMiddleware())
	custGrp.Get("/", middleware.RequirePolicy(db, "read"), customerHandler.FindAll)
	custGrp.Get("/:id", middleware.RequirePolicy(db, "read"), customerHandler.FindByID)
	custGrp.Post("/", middleware.RequirePolicy(db, "create"), customerHandler.Create)
	custGrp.Put("/:id", middleware.RequirePolicy(db, "update"), customerHandler.Update)
	custGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), customerHandler.Delete)

	// Customer Contacts Routes
	custContactGrp := api.Group("/customer-contacts", middleware.AuthMiddleware())
	custContactGrp.Get("/", middleware.RequirePolicy(db, "read"), customerContactHandler.FindAllByCustomer)
	custContactGrp.Get("/:id", middleware.RequirePolicy(db, "read"), customerContactHandler.FindByID)
	custContactGrp.Post("/", middleware.RequirePolicy(db, "create"), customerContactHandler.Create)
	custContactGrp.Put("/:id", middleware.RequirePolicy(db, "update"), customerContactHandler.Update)
	custContactGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), customerContactHandler.Delete)

	// Policies Routes
	policyGrp := api.Group("/policies", middleware.AuthMiddleware())
	policyGrp.Get("/", middleware.RequirePolicy(db, "read"), policyHandler.FindAll)
	policyGrp.Get("/:id", middleware.RequirePolicy(db, "read"), policyHandler.FindByID)
	policyGrp.Get("/group/:groupID", middleware.RequirePolicy(db, "read"), policyHandler.FindByGroupID)
	policyGrp.Post("/", middleware.RequirePolicy(db, "create"), policyHandler.Create)
	policyGrp.Put("/:id", middleware.RequirePolicy(db, "update"), policyHandler.Update)
	policyGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), policyHandler.Delete)

	// Departements Routes with Role Policies
	deptGrp := api.Group("/departements", middleware.AuthMiddleware())
	deptGrp.Get("/", middleware.RequirePolicy(db, "read"), deptHandler.FindAll)
	deptGrp.Get("/:id", middleware.RequirePolicy(db, "read"), deptHandler.FindByID)
	deptGrp.Post("/", middleware.RequirePolicy(db, "create"), deptHandler.Create)
	deptGrp.Put("/:id", middleware.RequirePolicy(db, "update"), deptHandler.Update)
	deptGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), deptHandler.Delete)

	// Roles (User Groups) Routes
	roleGrp := api.Group("/roles", middleware.AuthMiddleware())
	roleGrp.Get("/", middleware.RequirePolicy(db, "read"), roleHandler.FindAll)
	roleGrp.Get("/:id", middleware.RequirePolicy(db, "read"), roleHandler.FindByID)
	roleGrp.Post("/", middleware.RequirePolicy(db, "create"), roleHandler.Create)
	roleGrp.Put("/:id", middleware.RequirePolicy(db, "update"), roleHandler.Update)
	roleGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), roleHandler.Delete)

	// Menu Navigation (Master Table Access) Routes
	menuAccessGrp := api.Group("/menu-navigation", middleware.AuthMiddleware())
	menuAccessGrp.Get("/", middleware.RequirePolicy(db, "read"), menuAccessHandler.FindAll)
	menuAccessGrp.Get("/:id", middleware.RequirePolicy(db, "read"), menuAccessHandler.FindByID)
	menuAccessGrp.Post("/", middleware.RequirePolicy(db, "create"), menuAccessHandler.Create)
	menuAccessGrp.Put("/:id", middleware.RequirePolicy(db, "update"), menuAccessHandler.Update)
	menuAccessGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), menuAccessHandler.Delete)

	// Payment Term Routes
	paymentTermGrp := api.Group("/payment-terms", middleware.AuthMiddleware())
	paymentTermGrp.Get("/", middleware.RequirePolicy(db, "read"), paymentTermHandler.FindAll)
	paymentTermGrp.Get("/:id", middleware.RequirePolicy(db, "read"), paymentTermHandler.FindByID)
	paymentTermGrp.Post("/", middleware.RequirePolicy(db, "create"), paymentTermHandler.Create)
	paymentTermGrp.Put("/:id", middleware.RequirePolicy(db, "update"), paymentTermHandler.Update)
	paymentTermGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), paymentTermHandler.Delete)

	// Project Level Routes
	projectLevelGrp := api.Group("/project-levels", middleware.AuthMiddleware())
	projectLevelGrp.Get("/", middleware.RequirePolicy(db, "read"), projectLevelHandler.FindAll)
	projectLevelGrp.Get("/:id", middleware.RequirePolicy(db, "read"), projectLevelHandler.FindByID)
	projectLevelGrp.Post("/", middleware.RequirePolicy(db, "create"), projectLevelHandler.Create)
	projectLevelGrp.Put("/:id", middleware.RequirePolicy(db, "update"), projectLevelHandler.Update)
	projectLevelGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), projectLevelHandler.Delete)

	// Project Priority Routes
	projectPriorityGrp := api.Group("/project-priorities", middleware.AuthMiddleware())
	projectPriorityGrp.Get("/", middleware.RequirePolicy(db, "read"), projectPriorityHandler.FindAll)
	projectPriorityGrp.Get("/:id", middleware.RequirePolicy(db, "read"), projectPriorityHandler.FindByID)
	projectPriorityGrp.Post("/", middleware.RequirePolicy(db, "create"), projectPriorityHandler.Create)
	projectPriorityGrp.Put("/:id", middleware.RequirePolicy(db, "update"), projectPriorityHandler.Update)
	projectPriorityGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), projectPriorityHandler.Delete)

	// Quotation Progress Routes
	quotationProgressGrp := api.Group("/quotation-progress", middleware.AuthMiddleware())
	quotationProgressGrp.Get("/", quotationProgressHandler.FindAll)
	quotationProgressGrp.Get("/:id", middleware.RequirePolicy(db, "read"), quotationProgressHandler.FindByID)
	quotationProgressGrp.Post("/", middleware.RequirePolicy(db, "create"), quotationProgressHandler.Create)
	quotationProgressGrp.Put("/:id", middleware.RequirePolicy(db, "update"), quotationProgressHandler.Update)
	quotationProgressGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), quotationProgressHandler.Delete)

	// Quotation Status Routes
	quotationStatusGrp := api.Group("/quotation-statuses", middleware.AuthMiddleware())
	quotationStatusGrp.Get("/", quotationStatusHandler.FindAll)
	quotationStatusGrp.Get("/:id", middleware.RequirePolicy(db, "read"), quotationStatusHandler.FindByID)
	quotationStatusGrp.Post("/", middleware.RequirePolicy(db, "create"), quotationStatusHandler.Create)
	quotationStatusGrp.Put("/:id", middleware.RequirePolicy(db, "update"), quotationStatusHandler.Update)
	quotationStatusGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), quotationStatusHandler.Delete)

	// Units Routes
	unitGrp := api.Group("/units", middleware.AuthMiddleware())
	unitGrp.Get("/", middleware.RequirePolicy(db, "read"), unitHandler.FindAll)
	unitGrp.Get("/:id", middleware.RequirePolicy(db, "read"), unitHandler.FindByID)
	unitGrp.Post("/", middleware.RequirePolicy(db, "create"), unitHandler.Create)
	unitGrp.Put("/:id", middleware.RequirePolicy(db, "update"), unitHandler.Update)
	unitGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), unitHandler.Delete)

	// Quotations Routes
	quotationsGrp := api.Group("/quotations", middleware.AuthMiddleware())
	quotationsGrp.Get("/", middleware.RequirePolicy(db, "read"), quotationHandler.FindAll)
	quotationsGrp.Get("/export", middleware.RequirePolicy(db, "read"), quotationHandler.ExportExcel)
	quotationsGrp.Get("/top-customers", middleware.RequirePolicy(db, "read"), quotationHandler.TopCustomers)
	quotationsGrp.Get("/top-customers/export", middleware.RequirePolicy(db, "read"), quotationHandler.TopCustomersExport)
	quotationsGrp.Get("/sales-summary", middleware.RequirePolicy(db, "read"), quotationHandler.SalesSummaryByCustomer)
	quotationsGrp.Get("/sales-summary/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesSummaryByCustomerExport)
	quotationsGrp.Get("/sales-summary-by-sales-person", middleware.RequirePolicy(db, "read"), quotationHandler.SalesSummaryBySalesPerson)
	quotationsGrp.Get("/sales-summary-by-sales-person/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesSummaryBySalesPersonExport)
	quotationsGrp.Get("/sales-charts-by-customer", middleware.RequirePolicy(db, "read"), quotationHandler.SalesChartsByCustomer)
	quotationsGrp.Get("/sales-charts-by-sales-person", middleware.RequirePolicy(db, "read"), quotationHandler.SalesChartsBySalesPerson)
	quotationsGrp.Get("/sales-detail", middleware.RequirePolicy(db, "read"), quotationHandler.SalesDetailByCustomer)
	quotationsGrp.Get("/sales-detail/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesDetailByCustomerExport)
	quotationsGrp.Get("/sales-detail-by-sales-person", middleware.RequirePolicy(db, "read"), quotationHandler.SalesDetailBySalesPerson)
	quotationsGrp.Get("/sales-detail-by-sales-person/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesDetailBySalesPersonExport)
	quotationsGrp.Get("/sales-item-by-customer", middleware.RequirePolicy(db, "read"), quotationHandler.SalesItemByCustomer)
	quotationsGrp.Get("/sales-item-by-customer/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesItemByCustomerExport)
	quotationsGrp.Get("/sales-item-by-sales-person", middleware.RequirePolicy(db, "read"), quotationHandler.SalesItemBySalesPerson)
	quotationsGrp.Get("/sales-item-by-sales-person/export", middleware.RequirePolicy(db, "read"), quotationHandler.SalesItemBySalesPersonExport)
	quotationsGrp.Get("/report", middleware.RequirePolicy(db, "read"), quotationHandler.QuotationsReport)
	quotationsGrp.Get("/report/section", middleware.RequirePolicy(db, "read"), quotationHandler.ReportSection)
	quotationsGrp.Get("/report/export", middleware.RequirePolicy(db, "read"), quotationHandler.ReportExportExcel)
	quotationsGrp.Get("/chart-progress", middleware.RequirePolicy(db, "read"), quotationHandler.ChartProgressByYear)
	quotationsGrp.Get("/chart-po-analysis", middleware.RequirePolicy(db, "read"), quotationHandler.ChartPOAnalysis)
	quotationsGrp.Get("/stats", middleware.RequirePolicy(db, "read"), quotationHandler.QuotationStats)
	quotationsGrp.Get("/need-followup", middleware.RequirePolicy(db, "read"), quotationHandler.NeedFollowup)
	quotationsGrp.Get("/:id", middleware.RequirePolicy(db, "read"), quotationHandler.FindByID)
	quotationsGrp.Post("/", middleware.RequirePolicy(db, "create"), quotationHandler.Create)
	quotationsGrp.Put("/:id", middleware.RequirePolicy(db, "update"), quotationHandler.Update)
	quotationsGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), quotationHandler.Delete)
	quotationsGrp.Post("/:id/revision", middleware.RequirePolicy(db, "create"), quotationHandler.CreateRevision)
	quotationsGrp.Post("/:id/duplicate", middleware.RequirePolicy(db, "create"), quotationHandler.Duplicate)
	quotationsGrp.Post("/:id/set-default/:rev_id", middleware.RequirePolicy(db, "update"), quotationHandler.SetDefault)

	// Quotation Follow-ups Routes
	followupsGrp := api.Group("/quotation-followups", middleware.AuthMiddleware())
	followupsGrp.Get("/", middleware.RequirePolicy(db, "read"), quotationFollowupHandler.FindAllByQuotation)
	followupsGrp.Get("/report", middleware.RequirePolicy(db, "read"), quotationFollowupHandler.FollowupsReport)
	followupsGrp.Get("/report/export", middleware.RequirePolicy(db, "read"), quotationFollowupHandler.FollowupsReportExport)
	followupsGrp.Get("/:id", middleware.RequirePolicy(db, "read"), quotationFollowupHandler.FindByID)
	followupsGrp.Post("/", middleware.RequirePolicy(db, "create"), quotationFollowupHandler.Create)
	followupsGrp.Put("/:id", middleware.RequirePolicy(db, "update"), quotationFollowupHandler.Update)
	followupsGrp.Delete("/:id", middleware.RequirePolicy(db, "delete"), quotationFollowupHandler.Delete)
	followupsGrp.Post("/upload", middleware.RequirePolicy(db, "create"), quotationFollowupHandler.UploadPO)

	// === Kanban Board Routes ===
	kanbanBoardGrp := api.Group("/kanban-boards", middleware.AuthMiddleware())
	kanbanBoardGrp.Get("/", kanbanBoardHandler.FindAll)
	kanbanBoardGrp.Get("/:id", kanbanBoardHandler.FindByID)
	kanbanBoardGrp.Post("/", kanbanBoardHandler.Create)
	kanbanBoardGrp.Put("/:id", kanbanBoardHandler.Update)
	kanbanBoardGrp.Delete("/:id", kanbanBoardHandler.Delete)

	// === Kanban List Routes ===
	kanbanListGrp := api.Group("/kanban-lists", middleware.AuthMiddleware())
	kanbanListGrp.Get("/", kanbanListHandler.FindByBoardID)
	kanbanListGrp.Post("/", kanbanListHandler.Create)
	kanbanListGrp.Put("/:id", kanbanListHandler.Update)
	kanbanListGrp.Delete("/:id", kanbanListHandler.Delete)
	kanbanListGrp.Put("/reorder", kanbanListHandler.Reorder)

	// === Kanban Card Routes ===
	kanbanCardGrp := api.Group("/kanban-cards", middleware.AuthMiddleware())
	kanbanCardGrp.Get("/", kanbanCardHandler.FindByListID)
	kanbanCardGrp.Post("/", kanbanCardHandler.Create)
	kanbanCardGrp.Put("/move", kanbanCardHandler.Move)
	kanbanCardGrp.Put("/reorder", kanbanCardHandler.Reorder)
	kanbanCardGrp.Get("/:id", kanbanCardHandler.FindByID)
	kanbanCardGrp.Put("/:id", kanbanCardHandler.Update)
	kanbanCardGrp.Delete("/:id", kanbanCardHandler.Delete)
	kanbanCardGrp.Put("/:id/members", kanbanCardHandler.SyncMembers)
	kanbanCardGrp.Post("/:id/labels", kanbanCardHandler.SyncLabels)
	kanbanCardGrp.Delete("/:id/labels/:label_id", kanbanCardHandler.SyncLabels)

	// === Kanban Label Routes ===
	kanbanLabelGrp := api.Group("/kanban-labels", middleware.AuthMiddleware())
	kanbanLabelGrp.Get("/", kanbanLabelHandler.FindByBoardID)
	kanbanLabelGrp.Post("/", kanbanLabelHandler.Create)
	kanbanLabelGrp.Put("/:id", kanbanLabelHandler.Update)
	kanbanLabelGrp.Delete("/:id", kanbanLabelHandler.Delete)

	// === Kanban Checklist Routes ===
	kanbanChecklistGrp := api.Group("/kanban-checklists", middleware.AuthMiddleware())
	kanbanChecklistGrp.Get("/", kanbanChecklistHandler.FindByCardID)
	kanbanChecklistGrp.Post("/", kanbanChecklistHandler.CreateChecklist)
	kanbanChecklistGrp.Put("/:id", kanbanChecklistHandler.UpdateChecklist)
	kanbanChecklistGrp.Delete("/:id", kanbanChecklistHandler.DeleteChecklist)

	kanbanChecklistItemGrp := api.Group("/kanban-checklist-items", middleware.AuthMiddleware())
	kanbanChecklistItemGrp.Post("/", kanbanChecklistHandler.CreateItem)
	kanbanChecklistItemGrp.Put("/:id", kanbanChecklistHandler.UpdateItem)
	kanbanChecklistItemGrp.Delete("/:id", kanbanChecklistHandler.DeleteItem)

	// === Kanban Attachment Routes ===
	kanbanAttachmentGrp := api.Group("/kanban-attachments", middleware.AuthMiddleware())
	kanbanAttachmentGrp.Get("/", kanbanAttachmentHandler.FindByCardID)
	kanbanAttachmentGrp.Post("/", kanbanAttachmentHandler.Create)
	kanbanAttachmentGrp.Delete("/:id", kanbanAttachmentHandler.Delete)

	// === Kanban Comment Routes ===
	kanbanCommentGrp := api.Group("/kanban-comments", middleware.AuthMiddleware())
	kanbanCommentGrp.Get("/", kanbanCommentHandler.FindByCardID)
	kanbanCommentGrp.Post("/", kanbanCommentHandler.Create)
	kanbanCommentGrp.Put("/:id", kanbanCommentHandler.Update)
	kanbanCommentGrp.Delete("/:id", kanbanCommentHandler.Delete)

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return utils.SuccessResponse(c, fiber.StatusOK, "ok")
	})
}
