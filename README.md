# URL Kısaltma ve QR Kod Oluşturma Sistemi

Bu proje, uzun URL'leri kısaltarak kısa URL'ler oluşturur ve bu kısa URL'ler için QR kodları üretir. Proje, Go dilinde geliştirilmiş olup, PostgreSQL ve Redis kullanarak yüksek performanslı ve ölçeklenebilir bir sistem sağlar. Docker ve Docker Compose kullanarak kolayca dağıtılabilir.

## Özellikler

- **URL Kısaltma**: Uzun URL'leri kısaltarak daha kısa URL'ler oluşturur.
- **QR Kod Oluşturma**: Kısa URL'ler için QR kodları üretir ve bu QR kodlarını web sayfasında gösterir.
- **PostgreSQL Entegrasyonu**: URL verilerini saklamak için PostgreSQL kullanır.
- **Docker ile Dağıtım**: Docker ve Docker Compose kullanarak uygulamanızı hızlıca başlatabilir ve dağıtabilirsiniz.

## Başlangıç

### Gereksinimler

- Docker ve Docker Compose
- Go 1.18 veya daha yeni bir sürüm
- PostgreSQL

### Kurulum

1. **Proje Klonlama**

   Projeyi yerel makinenize klonlayın:

   ```bash
   git clone <repo-url>
   cd url-shortener
   ```

2. **Docker ve Docker Compose ile Başlatma**

   Docker Compose kullanarak PostgreSQL veritabanı ve uygulama servisini başlatın:

   ```bash
   docker-compose up --build
   ```

   Bu komut, PostgreSQL veritabanını başlatır, gerekli tabloları oluşturur ve uygulama servisini çalıştırır.

3. **Veritabanı Şemasını Otomatik Uygulama**

   `docker-compose` başlatıldığında, veritabanı şeması otomatik olarak uygulanacaktır. Ancak, `init.sql` dosyası içindeki şemayı manuel olarak uygulamanız gerekebilir.

### Kullanım

#### URL Kısaltma

- **Yöntem:** POST
- **URL:** `http://localhost:8080/shorten`
- **Başlıklar:**
  - `Content-Type: application/json`
- **Gövde (Body):**

  ```json
  {
    "url": "https://www.example.com"
  }
  ```

- **Yanıt:**

  HTML formatında bir yanıt alırsınız. Yanıt HTML sayfasında QR kodunu ve kısa URL'yi içerir.

#### URL Genişletme

- **Yöntem:** GET
- **URL:** `http://localhost:8080/expand?short_url=<short_url>`
- **Başlıklar:**
  - `Content-Type: application/json`

- **Yanıt:**

  ```json
  {
    "short_url": "https://www.example.com"
  }
  ```
