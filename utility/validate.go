package utility

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations(v *validator.Validate) {

	_ = v.RegisterValidation("timezone", func(fl validator.FieldLevel) bool {
		tz := fl.Field().String()
		_, err := time.LoadLocation(tz)
		return err == nil
	})
}

func ValidateImageryFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := []string{".tif", ".tiff", ".geotiff"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("invalid imagery format: %s. Expected .tif or .tiff", ext)
}

func ValidateDocumentFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := []string{".pdf", ".docx", ".doc"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("invalid document format: %s. Expected .pdf or .docx", ext)
}
