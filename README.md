// ...existing code...
# 🌍 KWBsite – Çok Dilli İçerik Yönetimi ve RESTful API

Geliştirici: Hüseyin KOÇ  
Firma: AYSTEK MÜHENDİSLİK — https://aysteknolojileri.com

KWBsite, Gin (Go) ve MongoDB tabanlı, çok dilli içerik yönetimi, kullanıcı/rol yönetimi, medya ve bildirim altyapısı sağlayan yönetim paneli ve API projesidir. Hem REST hem de GraphQL endpoint'leri, otomatik Swagger dokümantasyonu ve geliştirici araçları içerir.

## İçindekiler
- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Mimari ve Klasör Yapısı](#mimari-ve-klasör-yapısı)
- [Gereksinimler](#gereksinimler)
- [Ortam Değişkenleri (.env)](#ortam-değişkenleri-env)
- [Yerel Çalıştırma & Geliştirme Akışı](#yerel-çalıştırma--geliştirme-akışı)
- [Docker / Docker Compose (Öneri)](#docker--docker-compose-öneri)
- [Önemli Endpoint'ler & Örnekler](#önemli-endpointler--örnekler)
- [Profiling & Debugging](#profiling--debugging)
- [Testler, Linting & CI](#testler-linting--ci)
- [Üretim Hazırlıkları & Güvenlik](#üretim-hazirliklari--guvenlik)
- [Sorun Giderme](#sorun-giderme)
- [Katkıda Bulunma](#katkida-bulunma)
- [İletişim & Lisans](#iletisim--lisans)

## Hızlı Başlangıç
1. Depoyu klonlayın:
```bash
git clone <repository-url>
cd /workspace/go/admin-panel
go mod tidy
```

2. Örnek .env dosyası oluşturun (aşağıya bakın).

3. Geliştirme:
```bash
# doğrudan
go run main.go

# veya derle
go build -o admin-panel ./...
./admin-panel
```

Not: Host makinede tarayıcı açmak için:
```bash
"$BROWSER" http://localhost:9090/swagger/index.html
```
(Dev container: Debian GNU/Linux 12 (bookworm) üzerinde çalışır.)

## Mimari ve Klasör Yapısı
- configs/     — MongoDB bağlantısı ve konfig yükleme
- routes/      — REST & GraphQL route grupları
- services/    — iş mantığı / servis init'leri
- middlewares/ — logger, CORS, hata yakalama vb.
- docs/        — Swagger/döküman örnekleri
- uploads/     — kullanıcı yüklemeleri (statik servis)
- main.go      — uygulama başlatma, servis init, router ayarları

Akış: main.go -> configs.Init() -> services.Init*() -> routes.*Routes(r) -> r.Run(PORT)

## Gereksinimler
- Go 1.19+
- MongoDB 5.0+
- Geliştirme: git, curl, make (isteğe bağlı)

## Ortam Değişkenleri (.env)
Örnek:
```env
MONGO_URI=mongodb://localhost:27017
DB_NAME=admin_panel
PORT=9090
JWT_SECRET=supersecret
EMAIL_HOST=smtp.example.com
EMAIL_PORT=587
EMAIL_USER=you@example.com
EMAIL_PASS=secret
```
- PORT yoksa main.go içindeki default :9090 kullanılır.
- Hassas verileri secrets manager veya ortam değişkenleri ile yönetin.

## Yerel Çalıştırma & Geliştirme Akışı
- Kod değişince `go run main.go` çalıştırın veya hot-reload için `air`/`reflex` kurun.
- Yeni rota eklediğinizde routes paketine ekleyin ve uygun init çağrısını main.go'da kontrol edin.
- Servisleri services/ içinde Init* fonksiyonları üzerinden başlatın.

## Docker / Docker Compose (Öneri)
Basit Dockerfile örneği:
```dockerfile
FROM golang:1.19-bookworm AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o admin-panel ./...

FROM gcr.io/distroless/static
COPY --from=build /app/admin-panel /admin-panel
EXPOSE 9090
ENTRYPOINT ["/admin-panel"]
```

docker-compose.yml önerisi (Mongo ile):
```yaml
version: "3.8"
services:
  mongo:
    image: mongo:6
    restart: always
    ports:
      - "27017:27017"
    volumes:
      - mongo-data:/data/db
  app:
    build: .
    environment:
      - MONGO_URI=mongodb://mongo:27017
      - DB_NAME=admin_panel
      - PORT=9090
    ports:
      - "9090:9090"
    depends_on:
      - mongo
volumes:
  mongo-data:
```

## Önemli Endpoint'ler & Örnekler
- Swagger UI: GET /swagger/index.html
- Statik:
  - Uploads: GET /uploads/<file>
  - Docs: GET /docs/<file>
- Dil/slug örneği:
```bash
curl http://localhost:9090/tr/anasayfa
# -> {"language":"tr","slug":"anasayfa"}
```
- Başlatma noktası: main.go (servis init ve r.Run(":9090"))

## Profiling & Debugging
- pprof aktif: localhost:6060 (pprof import edildi)
  - Profil almak:
    ```bash
    go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
    ```
- Loglama: hem standart log hem de middleware tabanlı request loglama mevcut.

## Testler, Linting & CI
- Birim testleri:
```bash
go test ./... -v
```
- Kod formatı:
```bash
gofmt -w .
```
- Önerilen linter: golangci-lint
- CI: PR'larda `go test`, `gofmt` ve linter çalıştırılmalı.

## Üretim Hazırlıkları & Güvenlik
- TLS: reverse proxy (nginx/Caddy) ile HTTPS sonlandırma önerilir.
- JWT_SECRET ve e-mail şifreleri secret manager ile saklanmalı.
- pprof sadece iç ağda veya kapalı tutulmalı.
- DB bağlantı sınırları ve connection pooling gözden geçirilmeli.

## Sorun Giderme
- Mongo bağlantı hatası: MONGO_URI ve Mongo servisinin çalıştığını kontrol edin.
- Swagger boşsa: swagger yorumları generate edilmemiş olabilir; proje kökünde `swag init` çalıştırın.
- Statik dosya görünmüyorsa: uploads klasörü izinlerini kontrol edin.

## Katkıda Bulunma
- Fork → feature branch → PR. PR açıklamasında test ve migration notları ekleyin.
- Büyük değişiklikler için issue açıp tasarım-migration planı tartışın.

## İletişim & Lisans
Geliştirici: Hüseyin KOÇ  
Firma: AYSTEK MÜHENDİSLİK — https://aysteknolojileri.com

Lisans: MIT — detaylar LICENSE dosyasında.

---  
Geliştirilmesini istediğiniz bölümleri (ör: ayrıntılı API örnekleri, CI/CD pipeline, kapsamlı Docker Compose, İngilizce sürüm) belirtin; ilgili eklemeleri yaparım.
// ...existing code...