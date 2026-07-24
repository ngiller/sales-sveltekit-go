package utils

import (
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SendEmailPayload struct {
	MailTo  []string `json:"email_to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func CheckAndSendPOEmail(db *gorm.DB, quotationID string) {
	// Retrieve the quotation with relations
	var quotation models.Quotation
	err := db.Preload("Customer").Preload("SalesPerson").Preload("Details").First(&quotation, "id = ?", quotationID).Error
	if err != nil {
		log.Printf("[Email Notification] Failed to find quotation %s: %v", quotationID, err)
		return
	}

	// If status is not P.O. (3), we reset the notif flag if it was true, and return
	if quotation.Status != 3 {
		if quotation.Notif {
			db.Model(&quotation).Update("notif", false)
			log.Printf("[Email Notification] Reset notif flag to false for quotation %s because status is no longer P.O.", quotation.QuotationID)
		}
		return
	}

	// If already notified, do nothing
	if quotation.Notif {
		return
	}

	// Determine recipient setting code based on quotation type
	// 2 (Retail) -> code 10 (Admin Sales Email)
	// 1 (Project), 3 (Maintenance) -> code 11 (Admin Project Email)
	var settingCode int
	switch quotation.QuotationType {
	case 2:
		settingCode = 10 // Admin Sales Email
	case 1, 3:
		settingCode = 11 // Admin Project Email
	default:
		// If type is others (e.g. 4), we do not send email based on instructions
		return
	}

	// Fetch setting value
	var setting models.Setting
	err = db.Where("code = ? AND property_id = ?", settingCode, quotation.PropertyID).First(&setting).Error
	if err != nil {
		log.Printf("[Email Notification] Failed to fetch setting code %d for property %d: %v", settingCode, quotation.PropertyID, err)
		return
	}

	// Parse emails
	var mailTo []string
	for _, line := range strings.Split(setting.Value, "\n") {
		email := strings.TrimSpace(line)
		if email != "" {
			mailTo = append(mailTo, email)
		}
	}

	if len(mailTo) == 0 {
		log.Printf("[Email Notification] No recipient emails found in setting code %d", settingCode)
		return
	}

	// Format values for the email body
	qDateStr := "-"
	if quotation.QuotationDate != nil {
		qDateStr = quotation.QuotationDate.Format("2006-01-02")
	}

	poDateStr := "-"
	if quotation.PoDate != nil {
		poDateStr = quotation.PoDate.Format("2006-01-02")
	}

	poNoStr := "-"
	if quotation.PoNo != nil {
		poNoStr = *quotation.PoNo
	}

	subjectStr := "-"
	if quotation.Subject != nil {
		subjectStr = *quotation.Subject
	}

	grandTotalStr := "Rp 0"
	if quotation.GrandTotal != nil {
		grandTotalStr = FormatRupiah(*quotation.GrandTotal)
	}

	salesPersonName := "-"
	if quotation.SalesPerson != nil {
		salesPersonName = quotation.SalesPerson.Name
	}

	// Build HTML Body
	var bodyBuilder strings.Builder
	bodyBuilder.WriteString("<h2>Quotation PO Notification</h2>")
	bodyBuilder.WriteString("<p>A quotation status has been updated to <strong>P.O.</strong></p>")
	bodyBuilder.WriteString("<table border=\"1\" cellpadding=\"8\" cellspacing=\"0\" style=\"border-collapse: collapse; font-family: sans-serif; min-width: 400px;\">")
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Quotation ID</strong></td><td>%s</td></tr>", quotation.QuotationID))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Date</strong></td><td>%s</td></tr>", qDateStr))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Customer</strong></td><td>%s</td></tr>", quotation.Customer.Name))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Subject</strong></td><td>%s</td></tr>", subjectStr))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Grand Total</strong></td><td><strong>%s</strong></td></tr>", grandTotalStr))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>PO Number</strong></td><td>%s</td></tr>", poNoStr))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>PO Date</strong></td><td>%s</td></tr>", poDateStr))
	bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"background-color: #f9f9f9;\"><strong>Sales Person</strong></td><td>%s</td></tr>", salesPersonName))
	bodyBuilder.WriteString("</table>")

	if len(quotation.Details) > 0 {
		bodyBuilder.WriteString("<br/><h3>Item Details</h3>")
		bodyBuilder.WriteString("<table border=\"1\" cellpadding=\"8\" cellspacing=\"0\" style=\"border-collapse: collapse; font-family: sans-serif; width: 100%; max-width: 800px;\">")
		bodyBuilder.WriteString("<thead><tr style=\"background-color: #f2f2f2;\"><th>No</th><th>Description</th><th>Qty</th><th>Price</th><th>Total</th></tr></thead>")
		bodyBuilder.WriteString("<tbody>")
		for _, detail := range quotation.Details {
			desc := "-"
			if detail.Description != nil {
				desc = *detail.Description
			}
			qty := 0.0
			if detail.Qty != nil {
				qty = *detail.Qty
			}
			priceStr := FormatRupiah(detail.Price)
			totalStr := FormatRupiah(detail.Total)

			noStr := ""
			if detail.No != nil {
				noStr = fmt.Sprintf("%d", *detail.No)
			} else {
				noStr = fmt.Sprintf("%d", detail.Line)
			}

			bodyBuilder.WriteString(fmt.Sprintf("<tr><td style=\"text-align: center;\">%s</td><td>%s</td><td style=\"text-align: center;\">%.0f</td><td style=\"text-align: right;\">%s</td><td style=\"text-align: right;\">%s</td></tr>", noStr, desc, qty, priceStr, totalStr))
		}
		bodyBuilder.WriteString("</tbody></table>")
	}

	bodyBuilder.WriteString("<br/><p style=\"font-size: 11px; color: #888;\">This is an automated notification. Please do not reply directly to this email.</p>")

	// Prepare API payload
	subject := fmt.Sprintf("[PO Notification] Quotation %s - %s", quotation.QuotationID, quotation.Customer.Name)
	payload := SendEmailPayload{
		MailTo:  mailTo,
		Subject: subject,
		Body:    bodyBuilder.String(),
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Email Notification] Failed to marshal payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "http://app.magnumsolusion.co.id:8025/api/sendemail", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("[Email Notification] Failed to create HTTP request: %v", err)
		return
	}

	// Basic Auth
	req.SetBasicAuth("magnumsolusion", "dA29LafX5qGeJKwfBJkdDuzP6GAhRLmA")
	// Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "access_token=O5s9stbwrO2BIkcqlwrs80aq4IOo0gZmrsHwI41sHsol0UIPHb")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Email Notification] HTTP request failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("[Email Notification] API returned non-OK status: %d", resp.StatusCode)
		return
	}

	// Update quotation Notif flag to true
	err = db.Model(&quotation).Update("notif", true).Error
	if err != nil {
		log.Printf("[Email Notification] Failed to update notif flag for quotation %s: %v", quotationID, err)
	} else {
		log.Printf("[Email Notification] Email notification successfully sent and flag updated for quotation %s", quotation.QuotationID)
	}
}

func FormatRupiah(val float64) string {
	// Format float to 0 decimal places
	str := fmt.Sprintf("%.0f", val)

	var result []string
	length := len(str)
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}
	return "Rp " + strings.Join(result, ".")
}
