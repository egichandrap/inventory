package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/skip2/go-qrcode"
)

// QRCodeConfig holds QR code generation configuration
type QRCodeConfig struct {
	BaseURL        string
	MerchantName   string
	MerchantLogo   string
	Size           int
	ErrorCorrection qrcode.RecoveryLevel
	ForegroundColor color.Color
	BackgroundColor color.Color
}

// QRCodeService handles QR code generation with custom design
type QRCodeService struct {
	config QRCodeConfig
	logo   image.Image
}

// NewQRCodeService creates a new QR code service
func NewQRCodeService(config QRCodeConfig) *QRCodeService {
	svc := &QRCodeService{
		config: config,
	}

	// Load logo if specified
	if config.MerchantLogo != "" {
		logo, err := loadLogo(config.MerchantLogo)
		if err == nil {
			svc.logo = logo
		}
	}

	return svc
}

// GenerateTableQR generates a QR code for a table
func (s *QRCodeService) GenerateTableQR(tableNumber int, tableID string) (image.Image, error) {
	// Generate URL for the table
	url := fmt.Sprintf("%s/order?table=%d&id=%s", s.config.BaseURL, tableNumber, tableID)

	return s.generateQR(url)
}

// GenerateQR generates a QR code from any URL
func (s *QRCodeService) generateQR(url string) (image.Image, error) {
	// Create QR code
	qr, err := qrcode.New(url, s.config.ErrorCorrection)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat QR code: %w", err)
	}

	// Set colors
	if s.config.ForegroundColor != nil {
		qr.ForegroundColor = s.config.ForegroundColor
	} else {
		qr.ForegroundColor = color.Black
	}
	
	if s.config.BackgroundColor != nil {
		qr.BackgroundColor = s.config.BackgroundColor
	} else {
		qr.BackgroundColor = color.White
	}

	// Generate image
	qrImage := qr.Image(s.config.Size)

	// Add logo if available
	if s.logo != nil {
		qrImage = s.addLogo(qrImage)
	}

	return qrImage, nil
}

// GenerateQRBytes returns QR code as PNG bytes
func (s *QRCodeService) GenerateQRBytes(tableNumber int, tableID string) ([]byte, error) {
	img, err := s.GenerateTableQR(tableNumber, tableID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("gagal encode QR code: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateQRFile generates QR code and saves to file
func (s *QRCodeService) GenerateQRFile(tableNumber int, tableID, outputPath string) error {
	bytes, err := s.GenerateQRBytes(tableNumber, tableID)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, bytes, 0644)
}

// GenerateBulkQR generates QR codes for multiple tables
func (s *QRCodeService) GenerateBulkQR(tables []struct {
	Number int
	ID     string
}) (map[int]image.Image, error) {
	result := make(map[int]image.Image)

	for _, table := range tables {
		qr, err := s.GenerateTableQR(table.Number, table.ID)
		if err != nil {
			return nil, err
		}
		result[table.Number] = qr
	}

	return result, nil
}

// GenerateQRString returns QR code as base64 string
func (s *QRCodeService) GenerateQRString(tableNumber int, tableID string) (string, error) {
	bytes, err := s.GenerateQRBytes(tableNumber, tableID)
	if err != nil {
		return "", err
	}

	// Return as base64 data URL
	return fmt.Sprintf("data:image/png;base64,%x", bytes), nil
}

// addLogo overlays the logo onto the QR code
func (s *QRCodeService) addLogo(qrImage image.Image) image.Image {
	bounds := qrImage.Bounds()
	logoSize := bounds.Dx() / 5 // Logo is 20% of QR code size

	// Create new image
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, qrImage, image.Point{}, draw.Over)

	// Resize logo if needed
	logo := s.resizeLogo(s.logo, logoSize)

	// Center logo
	logoBounds := logo.Bounds()
	offsetX := (bounds.Dx() - logoSize) / 2
	offsetY := (bounds.Dy() - logoSize) / 2

	// Draw white background for logo
	whiteBg := image.Rect(offsetX-5, offsetY-5, offsetX+logoSize+5, offsetY+logoSize+5)
	draw.Draw(result, whiteBg, &image.Uniform{color.White}, image.Point{}, draw.Over)

	// Draw logo
	draw.Draw(result, logoBounds.Add(image.Pt(offsetX, offsetY)), logo, image.Point{}, draw.Over)

	return result
}

// resizeLogo resizes the logo to fit size
func (s *QRCodeService) resizeLogo(logo image.Image, size int) image.Image {
	// Simple resizing (in production, use image/draw or image resize library)
	bounds := logo.Bounds()
	if bounds.Dx() <= size && bounds.Dy() <= size {
		return logo
	}

	// Create resized image
	resized := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(resized, resized.Bounds(), logo, bounds.Min, draw.Src)

	return resized
}

// loadLogo loads a logo from file
func loadLogo(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("gagal buka logo: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("gagal decode logo: %w", err)
	}

	return img, nil
}

// GenerateQRCodeAsFile is a helper to generate and save QR code
func (s *QRCodeService) GenerateQRCodeAsFile(tableNumber int, tableID, dir string) (string, error) {
	filename := fmt.Sprintf("%s/table-%d.png", dir, tableNumber)
	
	if err := s.GenerateQRFile(tableNumber, tableID, filename); err != nil {
		return "", err
	}

	return filename, nil
}
