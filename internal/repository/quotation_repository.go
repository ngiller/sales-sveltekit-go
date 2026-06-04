package repository

import (
	"backend/internal/models"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type QuotationRepository struct {
	db *gorm.DB
}

// sanitizeTimePtr sets *time.Time pointer to nil if it points to a zero time value
// or a year outside MySQL DATE range (1000-9999).
// This prevents GORM's Save() from writing '0000-00-00' for old DB records.
func sanitizeTimePtr(p **time.Time) {
	if *p != nil && ((*p).IsZero() || (*p).Year() < 1000 || (*p).Year() > 9999) {
		*p = nil
	}
}

func NewQuotationRepository(db *gorm.DB) *QuotationRepository {
	return &QuotationRepository{db: db}
}

func (r *QuotationRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *QuotationRepository) FindAll(search, fromDate, toDate string, qType, status, userCreated, progress string, page, limit int, sortBy, sortDir string) ([]models.Quotation, int64, error) {
	var items []models.Quotation
	var total int64

	query := r.db.Model(&models.Quotation{}).
		Preload("Customer").
		Preload("Contact").
		Preload("ProgressInfo").
		Preload("StatusInfo").
		Preload("PaymentTerm").
		Preload("Level").
		Preload("Priority").
		Preload("SalesPerson").
		Joins("left join customer on customer.id = quotation.customer_id").
		Joins("left join quotation_progress on quotation_progress.id = quotation.progress").
		Joins("left join quotation_status on quotation_status.id = quotation.status")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("quotation.quotation_id LIKE ? OR quotation.subject LIKE ? OR customer.name LIKE ?", searchTerm, searchTerm, searchTerm)
	} else {
		if fromDate != "" && toDate != "" {
			query = query.Where("quotation.quotation_date >= ? AND quotation.quotation_date <= ?", fromDate, toDate)
		} else if fromDate != "" {
			query = query.Where("quotation.quotation_date >= ?", fromDate)
		} else if toDate != "" {
			query = query.Where("quotation.quotation_date <= ?", toDate)
		}
	}

	if qType != "" && qType != "all" {
		if val, err := strconv.Atoi(qType); err == nil {
			query = query.Where("quotation.quotation_type = ?", val)
		}
	}
	if status != "" && status != "all" {
		if val, err := strconv.Atoi(status); err == nil {
			query = query.Where("quotation.status = ?", val)
		}
	}
	if userCreated != "" && userCreated != "all" {
		if val, err := strconv.Atoi(userCreated); err == nil {
			query = query.Where("quotation.sales_id = ?", val)
		}
	}
	if progress != "" && progress != "all" {
		if val, err := strconv.Atoi(progress); err == nil {
			query = query.Where("quotation.progress = ?", val)
		}
	}

	query.Count(&total)

	if sortBy == "" || sortBy == "quotation_date" {
		sortBy = "quotation.quotation_date"
	} else if sortBy == "quotation_id" {
		sortBy = "quotation.quotation_id"
	} else if sortBy == "customer" {
		sortBy = "customer.name"
	} else if sortBy == "subject" {
		sortBy = "quotation.subject"
	} else if sortBy == "grand_total" {
		sortBy = "quotation.grand_total"
	} else if sortBy == "progress" {
		sortBy = "quotation_progress.progress"
	} else if sortBy == "next_followup" {
		sortBy = "quotation.next_followup"
	} else if sortBy == "status" {
		sortBy = "quotation_status.name"
	}

	if sortDir == "" {
		sortDir = "desc"
	}
	orderClause := sortBy + " " + sortDir

	offset := (page - 1) * limit
	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *QuotationRepository) FindByID(id string) (*models.Quotation, error) {
	var item models.Quotation
	// Find the default revision
	var defaultMaster models.QuotationMaster
	if err := r.db.Where("id = ? AND default_quot = ?", id, true).First(&defaultMaster).Error; err != nil {
		// Fallback to rev 0 if no default set
		if err := r.db.Where("id = ? AND rev_id = ?", id, 0).First(&defaultMaster).Error; err != nil {
			return nil, err
		}
	}

	err := r.db.
		Preload("Customer").
		Preload("Contact").
		Preload("ProgressInfo").
		Preload("StatusInfo").
		Preload("PaymentTerm").
		Preload("Level").
		Preload("Priority").
		Preload("Details", "rev_id = ?", defaultMaster.RevID).
		Preload("Details.Unit").
		First(&item, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// Sanitize zero time values from old DB records ('0000-00-00' → nil)
	if item.QuotationDate != nil && item.QuotationDate.IsZero() {
		item.QuotationDate = nil
	}
	if item.FollowupDate != nil && item.FollowupDate.IsZero() {
		item.FollowupDate = nil
	}
	if item.NextFollowup != nil && item.NextFollowup.IsZero() {
		item.NextFollowup = nil
	}

	// Fetch subdetails for the default revision
	var subdetails []models.QuotationSubdetail
	r.db.Preload("Unit").Where("id = ? AND rev_id = ?", id, defaultMaster.RevID).Find(&subdetails)

	// Fetch all revisions
	var revisions []models.QuotationMaster
	r.db.Where("id = ?", id).Order("rev_id DESC").Find(&revisions)

	item.Subdetails = subdetails
	item.Revisions = revisions

	// Update item with master values for the default revision
	item.Total = defaultMaster.Total
	item.Disc = defaultMaster.Disc
	item.Tax = defaultMaster.Tax
	item.TaxValue = defaultMaster.TaxValue
	item.PPh = defaultMaster.PPh
	item.PPhValue = defaultMaster.PPhValue
	item.GrandTotal = defaultMaster.GrandTotal
	item.HppTotal = defaultMaster.HppTotal
	item.Profit = &defaultMaster.Profit
	item.ProfitValue = defaultMaster.ProfitValue
	item.PaymentTermID = defaultMaster.PaymentTermID
	item.ValidUntil = defaultMaster.ValidUntil
	if item.ValidUntil != nil && item.ValidUntil.IsZero() {
		item.ValidUntil = nil
	}
	item.Commision = defaultMaster.Commision
	item.Notes = defaultMaster.Notes
	item.LevelID = defaultMaster.LevelID
	item.PriorityID = defaultMaster.PriorityID
	item.ProjectStart = defaultMaster.ProjectStart
	if item.ProjectStart != nil && item.ProjectStart.IsZero() {
		item.ProjectStart = nil
	}
	item.ProjectEnd = defaultMaster.ProjectEnd
	if item.ProjectEnd != nil && item.ProjectEnd.IsZero() {
		item.ProjectEnd = nil
	}
	item.PoNo = defaultMaster.PoNo
	item.PoDate = defaultMaster.PoDate
	if item.PoDate != nil && item.PoDate.IsZero() {
		item.PoDate = nil
	}
	item.PoFile = defaultMaster.PoFile
	item.PoAssignTo = defaultMaster.PoAssignTo

	return &item, nil
}

func (r *QuotationRepository) FindDetails(id string, revID int) ([]models.QuotationDetail, error) {
	var items []models.QuotationDetail
	err := r.db.Preload("Unit").Where("id = ? AND rev_id = ?", id, revID).Order("line ASC").Find(&items).Error
	return items, err
}

func (r *QuotationRepository) FindSubdetails(id string, revID int) ([]models.QuotationSubdetail, error) {
	var items []models.QuotationSubdetail
	err := r.db.Preload("Unit").Where("id = ? AND rev_id = ?", id, revID).Order("line ASC, subline ASC").Find(&items).Error
	return items, err
}

func (r *QuotationRepository) FindMaster(id string) ([]models.QuotationMaster, error) {
	var items []models.QuotationMaster
	err := r.db.Where("id = ?", id).Order("rev_id DESC").Find(&items).Error
	return items, err
}

func (r *QuotationRepository) Create(q *models.Quotation, details []models.QuotationDetail, subdetails []models.QuotationSubdetail, master *models.QuotationMaster) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(q).Error; err != nil {
			return err
		}

		for i := range details {
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}

		for i := range subdetails {
			if err := tx.Create(&subdetails[i]).Error; err != nil {
				return err
			}
		}

		if master != nil {
			if err := tx.Create(master).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *QuotationRepository) Update(q *models.Quotation, details []models.QuotationDetail, subdetails []models.QuotationSubdetail, master *models.QuotationMaster) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Sanitize zero time values (e.g., from old DB records with '0000-00-00')
		sanitizeTimePtr(&q.QuotationDate)
		sanitizeTimePtr(&q.ValidUntil)
		sanitizeTimePtr(&q.FollowupDate)
		sanitizeTimePtr(&q.NextFollowup)
		sanitizeTimePtr(&q.ProjectStart)
		sanitizeTimePtr(&q.ProjectEnd)
		sanitizeTimePtr(&q.PoDate)

		if err := tx.Save(q).Error; err != nil {
			return err
		}

		if master != nil {
			// Sanitize zero time values on master struct as well
			sanitizeTimePtr(&master.ProjectStart)
			sanitizeTimePtr(&master.ProjectEnd)
			sanitizeTimePtr(&master.QuotationDate)
			sanitizeTimePtr(&master.ValidUntil)
			sanitizeTimePtr(&master.PoDate)

			if master.RevID == -1 {
				// Auto-increment logic
				var maxRev int
				tx.Model(&models.QuotationMaster{}).Where("id = ?", master.ID).Select("COALESCE(MAX(rev_id), -1)").Scan(&maxRev)
				master.RevID = maxRev + 1
			}

			// Update related items to use the same RevID
			for i := range details {
				details[i].RevID = master.RevID
			}
			for i := range subdetails {
				subdetails[i].RevID = master.RevID
			}

			// Reset all other revisions to NOT default
			tx.Model(&models.QuotationMaster{}).Where("id = ?", master.ID).Update("default_quot", false)
			master.DefaultQuot = true

			if err := tx.Create(master).Error; err != nil {
				return err
			}

			// Sync parent quotation table with the new default revision's totals
			updates := map[string]interface{}{
				"total":        master.Total,
				"disc":         master.Disc,
				"tax":          master.Tax,
				"tax_value":    master.TaxValue,
				"pph":          master.PPh,
				"pph_value":    master.PPhValue,
				"grand_total":  master.GrandTotal,
				"hpp_total":    master.HppTotal,
				"profit":       master.Profit,
				"profit_value": master.ProfitValue,
				"notes":        master.Notes,
			}

			if err := tx.Model(&models.Quotation{}).Where("id = ?", master.ID).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Create details and subdetails with correct RevID
		for i := range details {
			if err := tx.Create(&details[i]).Error; err != nil {
				return err
			}
		}
		for i := range subdetails {
			if err := tx.Create(&subdetails[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *QuotationRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.Quotation{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&models.QuotationDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&models.QuotationSubdetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", id).Delete(&models.QuotationMaster{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *QuotationRepository) CreateRevision(oldID string, newID string, newQuotationID string, newDate time.Time, userID uint, userInisial string) error {
	oldQ := &models.Quotation{}
	if err := r.db.First(oldQ, "id = ?", oldID).Error; err != nil {
		return err
	}

	var defaultMaster models.QuotationMaster
	if err := r.db.Where("id = ? AND default_quot = ?", oldID, true).First(&defaultMaster).Error; err != nil {
		if err := r.db.Where("id = ? AND rev_id = ?", oldID, 0).First(&defaultMaster).Error; err != nil {
			return err
		}
	}

	details, _ := r.FindDetails(oldID, defaultMaster.RevID)
	subdetails, _ := r.FindSubdetails(oldID, defaultMaster.RevID)
	masterList, _ := r.FindMaster(oldID)

	maxRevID := 0
	for _, m := range masterList {
		if m.RevID > maxRevID {
			maxRevID = m.RevID
		}
	}

	newQ := *oldQ
	newQ.ID = newID
	newQ.QuotationID = newQuotationID
	newQ.QuotationDate = &newDate
	newQ.UserCreated = &userID
	newQ.UserUpdate = &userID
	newQ.CreatedAt = time.Now()
	newQ.UpdatedAt = time.Now()

	if newQ.Subject != nil {
		newQ.Subject = strPtr("[REVISION] " + *newQ.Subject)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newQ).Error; err != nil {
			return err
		}

		for i := range details {
			d := details[i]
			d.ID = newID
			d.RevID = 0
			if err := tx.Create(&d).Error; err != nil {
				return err
			}
		}

		for i := range subdetails {
			s := subdetails[i]
			s.ID = newID
			s.RevID = 0
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		}

		profitVal := 0.0
		if newQ.Profit != nil { profitVal = *newQ.Profit }

		newMaster := models.QuotationMaster{
			ID:            newID,
			RevID:         0,
			QuotationDate: newQ.QuotationDate,
			Subject:       newQ.Subject,
			Total:         newQ.Total,
			Disc:         newQ.Disc,
			Tax:           newQ.Tax,
			TaxValue:      newQ.TaxValue,
			PPh:          newQ.PPh,
			PPhValue:     newQ.PPhValue,
			GrandTotal:   newQ.GrandTotal,
			HppTotal:     newQ.HppTotal,
			Profit:       profitVal,
			ProfitValue:  newQ.ProfitValue,
			PaymentTermID: newQ.PaymentTermID,
			ValidUntil:  newQ.ValidUntil,
			Commision:   newQ.Commision,
			Notes:       newQ.Notes,
			SalesID:      newQ.SalesID,
			UserCreated:  &userID,
			UserUpdate: &userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DefaultQuot: true,
			LevelID:     newQ.LevelID,
			PriorityID: newQ.PriorityID,
		}
		if err := tx.Create(&newMaster).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *QuotationRepository) GetCounter(propertyID int, counterName string, quotType int, ym string) (int, error) {
	var counter models.CounterID
	err := r.db.Where("property_id = ? AND counter_name = ? AND type = ? AND ym = ?", propertyID, counterName, quotType, ym).First(&counter).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			counter = models.CounterID{
				PropertyID:   propertyID,
				CounterName:  counterName,
				Type:        quotType,
				Ym:          ym,
				Counter:     1,
			}
			if err := r.db.Create(&counter).Error; err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}

	counter.Counter++
	if err := r.db.Save(&counter).Error; err != nil {
		return 0, err
	}
	return counter.Counter, nil
}

func GenerateQuotationID(quotType, counter int, month int, year int, userInisial string) string {
	romanMonth := []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	typePrefix := map[int]string{1: "NR", 2: "R", 3: "M", 4: "O"}
	prefix := typePrefix[quotType]
	if prefix == "" {
		prefix = "R"
	}
	return fmt.Sprintf("%s%03d/%s/%d/%s", prefix, counter, romanMonth[month-1], year, userInisial)
}

func GenerateID() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()

	randomNum, _ := rand.Int(rand.Reader, big.NewInt(10000000))
	randomPart := fmt.Sprintf("%07d", randomNum.Int64())

	return fmt.Sprintf("%04d%02d%02d%s", year, month, day, randomPart)
}

func strPtr(s string) *string {
	return &s
}

func (r *QuotationRepository) SetDefault(id string, revID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Reset all to non-default
		if err := tx.Model(&models.QuotationMaster{}).Where("id = ?", id).Update("default_quot", false).Error; err != nil {
			return err
		}
		// Set selected as default
		if err := tx.Model(&models.QuotationMaster{}).Where("id = ? AND rev_id = ?", id, revID).Update("default_quot", true).Error; err != nil {
			return err
		}

		// Sync parent quotation table with the new default revision's totals
		var master models.QuotationMaster
		if err := tx.Where("id = ? AND rev_id = ?", id, revID).First(&master).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"total":        master.Total,
			"disc":         master.Disc,
			"tax":          master.Tax,
			"tax_value":    master.TaxValue,
			"pph":          master.PPh,
			"pph_value":    master.PPhValue,
			"grand_total":  master.GrandTotal,
			"hpp_total":    master.HppTotal,
			"profit":       master.Profit,
			"profit_value": master.ProfitValue,
			"notes":        master.Notes,
		}

		if err := tx.Model(&models.Quotation{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}
