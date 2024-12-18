package utils

import (
	"regexp"
	"strings"
)

// generateSlug oluşturulan bir başlık için SEO dostu slug oluşturur
func GenerateSlug(title string) string {
	// Büyük harfleri küçük harfe çevir
	slug := strings.ToLower(title)

	// Özel karakterleri ve noktalama işaretlerini kaldır
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	// Boşlukları tire (-) ile değiştir
	slug = strings.ReplaceAll(slug, " ", "-")

	// Birden fazla tireyi tek bir tireye indir
	re = regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	return slug
}
