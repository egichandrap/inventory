package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/example/jwt-ddd-clean/internal/domain/model"
)

// ExportService handles data export to various formats
type ExportService struct{}

// NewExportService creates a new ExportService
func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportTransactionsToCSV exports transactions to CSV format
func (s *ExportService) ExportTransactionsToCSV(transactions []*model.Transaction, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"Transaction No",
		"Date",
		"Cashier",
		"Customer",
		"Payment Method",
		"Subtotal",
		"Discount",
		"Tax",
		"Total",
		"Status",
		"Notes",
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, t := range transactions {
		record := []string{
			t.TransactionNo(),
			t.CreatedAt().Format("2006-01-02 15:04:05"),
			t.CashierName(),
			t.CustomerName(),
			string(t.PaymentMethod()),
			fmt.Sprintf("%.2f", t.Subtotal()),
			fmt.Sprintf("%.2f", t.DiscountAmount()),
			fmt.Sprintf("%.2f", t.TaxAmount()),
			fmt.Sprintf("%.2f", t.TotalAmount()),
			string(t.Status()),
			t.Notes(),
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// ExportInventoryToCSV exports inventory to CSV format
func (s *ExportService) ExportInventoryToCSV(inventories []InventoryItem, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"SKU",
		"Name",
		"Description",
		"Quantity",
		"Unit",
		"Price",
		"Location",
		"Min Stock",
		"Max Stock",
		"Status",
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, inv := range inventories {
		status := "OK"
		if inv.Quantity <= 0 {
			status = "Out of Stock"
		} else if inv.Quantity <= inv.MinStock {
			status = "Low Stock"
		}

		record := []string{
			inv.SKU,
			inv.Name,
			inv.Description,
			fmt.Sprintf("%d", inv.Quantity),
			inv.Unit,
			fmt.Sprintf("%.2f", inv.Price),
			inv.Location,
			fmt.Sprintf("%d", inv.MinStock),
			fmt.Sprintf("%d", inv.MaxStock),
			status,
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// ExportSalesReportToCSV exports sales report to CSV format
func (s *ExportService) ExportSalesReportToCSV(dailySales []model.DailySales, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"Date",
		"Total Sales",
		"Transactions",
		"Total Items",
		"Average Transaction",
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	for _, ds := range dailySales {
		avgTransaction := 0.0
		if ds.TransactionCount > 0 {
			avgTransaction = ds.TotalSales / float64(ds.TransactionCount)
		}

		record := []string{
			ds.Date.Format("2006-01-02"),
			fmt.Sprintf("%.2f", ds.TotalSales),
			fmt.Sprintf("%d", ds.TransactionCount),
			fmt.Sprintf("%d", ds.TotalItems),
			fmt.Sprintf("%.2f", avgTransaction),
		}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// GenerateCSVFilename generates a filename for CSV export
func (s *ExportService) GenerateCSVFilename(prefix string) string {
	return fmt.Sprintf("%s_%s.csv", prefix, time.Now().Format("20060102_150405"))
}

// InventoryItem represents inventory data for export
type InventoryItem struct {
	SKU         string
	Name        string
	Description string
	Quantity    int
	Unit        string
	Price       float64
	Location    string
	MinStock    int
	MaxStock    int
}

// ConvertToCSV converts data to CSV string
func (s *ExportService) ConvertToCSV(headers []string, records [][]string) string {
	var builder strings.Builder

	// Write headers
	builder.WriteString(strings.Join(headers, ","))
	builder.WriteString("\n")

	// Write records
	for _, record := range records {
		builder.WriteString(strings.Join(record, ","))
		builder.WriteString("\n")
	}

	return builder.String()
}
