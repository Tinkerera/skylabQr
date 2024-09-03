package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	DB_NAME     = "urlshortener"
	DB_HOST     = "db"
)

func setupDB() *sql.DB {
	// PostgreSQL bağlantı dizesi
	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)

	// Veritabanına bağlanma
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Bağlantının düzgün çalışıp çalışmadığını kontrol etme
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Veritabanı bağlantısı başarılı!")
	return db
}

func shortenURL(originalURL string, db *sql.DB) (string, string, error) {
	// URL'yi hashleyerek kısa bir URL oluşturma
	hash := sha256.Sum256([]byte(originalURL + time.Now().String()))
	shortURL := base64.URLEncoding.EncodeToString(hash[:])[:8]

	// Veritabanına kaydetme
	expiration := time.Now().Add(24 * time.Hour) // 24 saat geçerlilik süresi
	_, err := db.Exec("INSERT INTO url_mappings (short_url, original_url, expiration_date) VALUES ($1, $2, $3)", shortURL, originalURL, expiration)
	if err != nil {
		return "", "", err
	}

	// QR kodunu oluşturma
	qrCode, err := qr.Encode(shortURL, qr.M, qr.Auto)
	if err != nil {
		return "", "", err
	}
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return "", "", err
	}

	var qrBuffer bytes.Buffer
	err = png.Encode(&qrBuffer, qrCode)
	if err != nil {
		return "", "", err
	}

	qrBase64 := base64.StdEncoding.EncodeToString(qrBuffer.Bytes())
	return shortURL, qrBase64, nil
}
func getOriginalURL(shortURL string, db *sql.DB) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM url_mappings WHERE short_url = $1 AND expiration_date > $2", shortURL, time.Now()).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("URL bulunamadı veya süresi doldu")
		}
		return "", err
	}
	return originalURL, nil
}

func generateQRCode(url string) error {
	// QR kodu oluşturma
	qrCode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		return err
	}

	// QR kodunu ölçeklendirme
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return err
	}

	// QR kodunu dosyaya kaydetme
	file, err := os.Create("qrcode.png")
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, qrCode)
	if err != nil {
		return err
	}

	fmt.Println("QR kodu başarıyla oluşturuldu: qrcode.png")
	return nil
}

func main() {
	db := setupDB()
	defer db.Close()

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		var request URLRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Geçersiz istek", http.StatusBadRequest)
			return
		}

		shortURL, qrBase64, err := shortenURL(request.URL, db)
		if err != nil {
			http.Error(w, "URL kısaltılamadı", http.StatusInternalServerError)
			return
		}

		// HTML yanıtı oluşturma
		htmlResponse := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>URL Kısaltma Sonuçları</title>
			</head>
			<body>
				<h1>Kısa URL Oluşturuldu</h1>
				<p>Orijinal URL: %s</p>
				<p>Kısa URL: %s</p>
				<p><img src="data:image/png;base64,%s" alt="QR Code" /></p>
			</body>
			</html>
		`, request.URL, shortURL, qrBase64)

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlResponse))
	})

	http.HandleFunc("/expand", func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Query().Get("short_url")
		if shortURL == "" {
			http.Error(w, "short_url parametresi gerekli", http.StatusBadRequest)
			return
		}

		originalURL, err := getOriginalURL(shortURL, db)
		if err != nil {
			http.Error(w, "URL bulunamadı", http.StatusNotFound)
			return
		}

		// Orijinal URL'yi JSON olarak döndür
		response := URLResponse{ShortURL: originalURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	fmt.Println("Sunucu 8080 portunda çalışıyor...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"short_url"`
}
