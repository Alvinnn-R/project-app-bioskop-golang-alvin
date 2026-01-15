# 🎬 Cinema Booking System API

RESTful API untuk sistem pemesanan tiket bioskop menggunakan **Golang** dengan arsitektur Clean Architecture.

## 📋 Daftar Isi

- [Teknologi](#-teknologi)
- [Arsitektur](#-arsitektur)
- [Fitur Utama](#-fitur-utama)
  - [Token Authentication](#1-token-authentication)
  - [Email OTP Verification](#2-email-otp-verification)
  - [Booking Seats](#3-booking-seats)
  - [Login Required Routes](#4-login-required-routes)
  - [DRY Principle](#5-dry-principle)
  - [Goroutines Implementation](#6-goroutines-implementation)
- [Unit Testing](#-unit-testing)
- [API Endpoints](#-api-endpoints)
- [Instalasi & Menjalankan](#-instalasi--menjalankan)
- [Struktur Folder](#-struktur-folder)

---

## 🛠 Teknologi

| Teknologi | Deskripsi |
|-----------|-----------|
| **Go (Golang)** | Bahasa pemrograman utama |
| **go-chi/chi v5** | HTTP Router untuk RESTful API |
| **PostgreSQL** | Database relasional |
| **pgx/v5** | PostgreSQL driver dengan connection pooling |
| **pgxmock/v3** | Mock library untuk testing database |
| **testify** | Library untuk assertions dalam testing |
| **bcrypt** | Password hashing |
| **UUID** | Token generation |
| **Zap Logger** | Structured logging |
| **Viper** | Configuration management |

---

## 🏗 Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan layer:

```
┌─────────────────────────────────────────────────────────┐
│                      HTTP Request                        │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   Adaptor Layer                          │
│              (HTTP Handler / Controller)                 │
│         internal/adaptor/*.go                           │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                   UseCase Layer                          │
│                 (Business Logic)                         │
│          internal/usecase/*.go                          │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                 Repository Layer                         │
│              (Database Operations)                       │
│       internal/data/repository/*.go                     │
└─────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                     PostgreSQL                           │
└─────────────────────────────────────────────────────────┘
```

---

## 🚀 Fitur Utama

### 1. Token Authentication

Sistem autentikasi menggunakan **UUID Token** yang disimpan di tabel `sessions`.

#### Alur Login:
```
User Login → Validasi Password → Generate UUID Token → Simpan ke Session → Return Token
```

#### Implementasi:

**Generate Token (usecase/auth.go):**
```go
func (a *AuthUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
    // 1. Cari user berdasarkan email
    user, err := a.authRepo.GetUserByEmail(ctx, req.Email)
    
    // 2. Validasi password dengan bcrypt
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    
    // 3. Generate UUID token
    token := uuid.New().String()
    
    // 4. Simpan session ke database (expired 24 jam)
    session := entity.Session{
        UserID:    user.ID,
        Token:     token,
        ExpiredAt: time.Now().Add(24 * time.Hour),
    }
    a.authRepo.CreateSession(ctx, session)
    
    return &dto.LoginResponse{Token: token, User: user}
}
```

**Validasi Token (middleware/auth.go):**
```go
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Ambil token dari header Authorization
        authHeader := r.Header.Get("Authorization")
        token := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 2. Validasi token via usecase
        userID, err := m.authUseCase.ValidateToken(r.Context(), token)
        
        // 3. Inject userID ke context
        ctx := context.WithValue(r.Context(), "userID", userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

### 2. Email OTP Verification

Sistem verifikasi email menggunakan **6-digit OTP** dengan masa berlaku **5 menit**.

#### Alur OTP:
```
Register → Generate OTP → Simpan ke DB → Kirim Email (Async) → User Input OTP → Verify → Activate Account
```

#### Implementasi:

**Generate OTP (usecase/auth.go):**
```go
func generateOTP() string {
    rand.Seed(time.Now().UnixNano())
    otp := rand.Intn(900000) + 100000  // Range: 100000-999999
    return strconv.Itoa(otp)
}

func (a *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) error {
    // 1. Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    // 2. Generate OTP
    otpCode := generateOTP()
    
    // 3. Simpan OTP ke database (expired 5 menit)
    otp := entity.OTP{
        UserID:    userID,
        OTPCode:   otpCode,
        ExpiredAt: time.Now().Add(5 * time.Minute),
    }
    
    // 4. Kirim email secara async (goroutine)
    go func() {
        utils.SendOTP(req.Email, otpCode)
    }()
}
```

**Verify OTP (usecase/auth.go):**
```go
func (a *AuthUseCase) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) error {
    // 1. Ambil OTP dari database
    otp, err := a.authRepo.GetOTPByUserID(ctx, userID)
    
    // 2. Cek apakah OTP expired
    if time.Now().After(otp.ExpiredAt) {
        return errors.New("OTP has expired")
    }
    
    // 3. Validasi kode OTP
    if otp.OTPCode != req.OTPCode {
        return errors.New("invalid OTP code")
    }
    
    // 4. Aktivasi user
    a.authRepo.ActivateUser(ctx, userID)
}
```

---

### 3. Booking Seats

Pemesanan kursi menggunakan **database transaction** untuk memastikan atomicity.

#### Alur Booking:
```
User Request → Check Seat Availability → Begin Transaction → Insert Booking → Insert Booking_Seats → Commit
```

#### Implementasi:

**Repository dengan Transaction (repository/booking.go):**
```go
func (r *BookingRepository) CreateBooking(ctx context.Context, booking entity.Booking, seatIDs []int) (int, error) {
    // 1. Begin transaction
    tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
    if err != nil {
        return 0, err
    }
    defer tx.Rollback(ctx)  // Rollback jika ada error
    
    // 2. Insert booking
    var bookingID int
    err = tx.QueryRow(ctx, `
        INSERT INTO bookings (user_id, showtime_id, total_price, status)
        VALUES ($1, $2, $3, $4) RETURNING id
    `, booking.UserID, booking.ShowtimeID, booking.TotalPrice, booking.Status).Scan(&bookingID)
    
    // 3. Insert booking_seats (multiple rows)
    for _, seatID := range seatIDs {
        _, err = tx.Exec(ctx, `
            INSERT INTO booking_seats (booking_id, seat_id) VALUES ($1, $2)
        `, bookingID, seatID)
    }
    
    // 4. Commit transaction
    err = tx.Commit(ctx)
    return bookingID, nil
}
```

---

### 4. Login Required Routes

Endpoint yang membutuhkan autentikasi diproteksi menggunakan **Middleware Pattern**.

#### Implementasi (wire/wire.go):

```go
func SetupRouter(authMiddleware *middleware.AuthMiddleware, ...) *chi.Mux {
    r := chi.NewRouter()
    
    // ===== PUBLIC ROUTES (Tidak perlu login) =====
    r.Post("/register", authAdaptor.Register)
    r.Post("/login", authAdaptor.Login)
    r.Post("/verify-otp", authAdaptor.VerifyOTP)
    r.Get("/movies", movieAdaptor.GetMovies)
    r.Get("/cinemas", cinemaAdaptor.GetCinemas)
    
    // ===== PROTECTED ROUTES (Perlu login) =====
    r.Group(func(r chi.Router) {
        r.Use(authMiddleware.RequireAuth)  // Middleware validasi token
        
        r.Post("/logout", authAdaptor.Logout)
        r.Post("/booking", bookingAdaptor.CreateBooking)
        r.Get("/user/bookings", bookingAdaptor.GetUserBookings)
        r.Post("/pay", paymentAdaptor.ProcessPayment)
    })
    
    return r
}
```

#### Cara Kerja Middleware:

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│   Request   │────▶│  RequireAuth    │────▶│   Handler   │
│  + Token    │     │   Middleware    │     │  (Adaptor)  │
└─────────────┘     └─────────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────────┐
                    │ Validate Token  │
                    │ Inject UserID   │
                    │ to Context      │
                    └─────────────────┘
```

---

### 5. DRY Principle

Kode yang reusable untuk menghindari duplikasi.

#### 5.1 Response Helper (utils/response.go):

```go
// Reusable response functions
func ResponseSuccess(w http.ResponseWriter, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(dto.Response{
        Status:  "success",
        Message: message,
        Data:    data,
    })
}

func ResponseCreated(w http.ResponseWriter, message string, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(dto.Response{...})
}

func ResponseError(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(dto.Response{
        Status:  "error",
        Message: message,
    })
}

func ResponseUnauthorized(w http.ResponseWriter, message string) {
    ResponseError(w, http.StatusUnauthorized, message)
}

func ResponseBadRequest(w http.ResponseWriter, message string) {
    ResponseError(w, http.StatusBadRequest, message)
}
```

#### 5.2 Penggunaan di Adaptor:

```go
// Sebelum DRY (duplikasi kode)
func (a *AuthAdaptor) Login(w http.ResponseWriter, r *http.Request) {
    // ... logic ...
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{...})
}

// Setelah DRY (clean & reusable)
func (a *AuthAdaptor) Login(w http.ResponseWriter, r *http.Request) {
    // ... logic ...
    utils.ResponseSuccess(w, "Login successful", result)
}

func (a *AuthAdaptor) Register(w http.ResponseWriter, r *http.Request) {
    // ... logic ...
    utils.ResponseCreated(w, "User registered successfully", nil)
}
```

---

### 6. Goroutines Implementation

#### 6.1 Fire-and-Forget Pattern (Email Notification)

Mengirim email secara **asynchronous** tanpa menunggu response.

```go
// usecase/auth.go - Register dengan async email
func (a *AuthUseCase) Register(ctx context.Context, req dto.RegisterRequest) error {
    // ... create user logic ...
    
    // Fire-and-forget: kirim email di background
    go func() {
        if err := utils.SendOTP(req.Email, otpCode); err != nil {
            a.logger.Error("Failed to send OTP email", zap.Error(err))
        }
    }()
    
    return nil  // Return immediately, tidak menunggu email terkirim
}

// usecase/payment.go - Payment dengan async email notification
func (p *PaymentUseCase) ProcessPayment(ctx context.Context, req dto.PaymentRequest) error {
    // ... payment logic ...
    
    // Fire-and-forget: kirim email konfirmasi
    go func() {
        if err := utils.SendPaymentConfirmation(user.Email, booking); err != nil {
            p.logger.Error("Failed to send payment confirmation", zap.Error(err))
        }
    }()
    
    return nil
}
```

#### 6.2 Concurrent Pattern dengan Channels

Mengambil data secara **parallel** menggunakan goroutines dan channels.

```go
// usecase/booking.go - GetUserBookings dengan concurrent fetching
type bookingDetailResult struct {
    movie    *entity.Movie
    showtime *entity.Showtime
    seats    []entity.Seat
    err      error
}

func (b *BookingUseCase) GetUserBookings(ctx context.Context, userID int) ([]dto.BookingDetailResponse, error) {
    // 1. Ambil semua bookings user
    bookings, err := b.bookingRepo.GetBookingsByUserID(ctx, userID)
    
    var results []dto.BookingDetailResponse
    
    for _, booking := range bookings {
        // 2. Buat buffered channel untuk hasil
        resultChan := make(chan bookingDetailResult, 3)
        
        // 3. Goroutine 1: Fetch movie data
        go func(showtimeID int) {
            movie, err := b.movieRepo.GetMovieByShowtimeID(ctx, showtimeID)
            resultChan <- bookingDetailResult{movie: movie, err: err}
        }(booking.ShowtimeID)
        
        // 4. Goroutine 2: Fetch showtime data
        go func(showtimeID int) {
            showtime, err := b.showtimeRepo.GetShowtimeByID(ctx, showtimeID)
            resultChan <- bookingDetailResult{showtime: showtime, err: err}
        }(booking.ShowtimeID)
        
        // 5. Goroutine 3: Fetch seats data
        go func(bookingID int) {
            seats, err := b.seatRepo.GetSeatsByBookingID(ctx, bookingID)
            resultChan <- bookingDetailResult{seats: seats, err: err}
        }(booking.ID)
        
        // 6. Collect results dengan timeout
        var movie *entity.Movie
        var showtime *entity.Showtime
        var seats []entity.Seat
        
        for i := 0; i < 3; i++ {
            select {
            case result := <-resultChan:
                if result.err != nil {
                    return nil, result.err
                }
                if result.movie != nil {
                    movie = result.movie
                }
                if result.showtime != nil {
                    showtime = result.showtime
                }
                if result.seats != nil {
                    seats = result.seats
                }
            case <-ctx.Done():
                return nil, ctx.Err()  // Timeout handling
            }
        }
        
        // 7. Build response
        results = append(results, dto.BookingDetailResponse{
            Booking:  booking,
            Movie:    movie,
            Showtime: showtime,
            Seats:    seats,
        })
    }
    
    return results, nil
}
```

#### Diagram Concurrent Fetching:

```
                    ┌─────────────────┐
                    │  GetUserBookings │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
       ┌──────────┐   ┌──────────┐   ┌──────────┐
       │Goroutine1│   │Goroutine2│   │Goroutine3│
       │Get Movie │   │Get Show  │   │Get Seats │
       └────┬─────┘   └────┬─────┘   └────┬─────┘
            │              │              │
            └──────────────┼──────────────┘
                           ▼
                    ┌─────────────┐
                    │   Channel   │
                    │  (Buffer 3) │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Collect   │
                    │   Results   │
                    └─────────────┘
```

---

## 🧪 Unit Testing

Unit test dilakukan pada **Repository Layer** menggunakan **pgxmock** untuk mocking database.

### Test Coverage: **81.1%**

### Menjalankan Test:

```bash
# Run semua test di repository
go test ./internal/data/repository/... -v

# Run test dengan coverage
go test ./internal/data/repository/... -cover

# Generate coverage report
go test ./internal/data/repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Contoh Test (repository/auth_test.go):

```go
func TestAuthRepository_GetUserByEmail(t *testing.T) {
    // 1. Setup mock database
    mock, err := pgxmock.NewPool()
    require.NoError(t, err)
    defer mock.Close()
    
    repo := repository.NewAuthRepository(mock)
    
    // 2. Setup expected rows
    rows := pgxmock.NewRows([]string{"id", "name", "email", "password", "is_active", "created_at", "updated_at"}).
        AddRow(1, "John Doe", "john@example.com", "hashedpassword", true, time.Now(), time.Now())
    
    // 3. Setup expectation
    mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
        WithArgs("john@example.com").
        WillReturnRows(rows)
    
    // 4. Execute
    user, err := repo.GetUserByEmail(context.Background(), "john@example.com")
    
    // 5. Assert
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
    
    // 6. Verify expectations met
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### File Test yang Dibuat:

| File | Deskripsi |
|------|-----------|
| `auth_test.go` | Test untuk user & session operations |
| `movie_test.go` | Test untuk movie & showtime operations |
| `cinema_test.go` | Test untuk cinema operations |
| `booking_test.go` | Test untuk booking operations |
| `payment_test.go` | Test untuk payment operations |
| `seat_test.go` | Test untuk seat operations |

---

## 📡 API Endpoints

### Public Endpoints (Tanpa Login)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/register` | Registrasi user baru |
| POST | `/login` | Login user |
| POST | `/verify-otp` | Verifikasi OTP |
| POST | `/resend-otp` | Kirim ulang OTP |
| GET | `/movies` | Daftar semua film |
| GET | `/movies/{id}` | Detail film |
| GET | `/cinemas` | Daftar semua bioskop |
| GET | `/showtimes` | Daftar jadwal tayang |

### Protected Endpoints (Perlu Login)

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/logout` | Logout user |
| POST | `/booking` | Buat booking baru |
| GET | `/user/bookings` | Daftar booking user |
| POST | `/pay` | Proses pembayaran |
| GET | `/seats/{showtime_id}` | Daftar kursi tersedia |

### Contoh Request:

**Login:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

**Booking (dengan token):**
```bash
curl -X POST http://localhost:8080/booking \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{"showtime_id": 1, "seat_ids": [1, 2, 3]}'
```

---

## 🚀 Instalasi & Menjalankan

### Prerequisites

- Go 1.21+
- PostgreSQL 14+

### Setup

1. **Clone repository:**
   ```bash
   git clone <repository-url>
   cd project-app-bioskop-golang-alvin
   ```

2. **Setup database:**
   ```bash
   # Import database schema
   psql -U postgres -d cinema_booking < database_file/cinema_booking_system-backup.sql
   ```

3. **Konfigurasi environment:**
   ```bash
   # Buat file .env
   cp .env.example .env
   
   # Edit konfigurasi
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=cinema_booking
   
   EMAIL_API_KEY=your_email_api_key
   EMAIL_API_URL=https://api.emailservice.com/send
   ```

4. **Install dependencies:**
   ```bash
   go mod download
   ```

5. **Jalankan aplikasi:**
   ```bash
   go run main.go
   ```

6. **Jalankan test:**
   ```bash
   go test ./internal/data/repository/... -cover
   ```

---

## 📁 Struktur Folder

```
project-app-bioskop-golang-alvin/
├── main.go                 # Entry point
├── go.mod                  # Go modules
├── go.sum
├── README.md               # Dokumentasi
│
├── cmd/
│   └── server.go           # Server configuration
│
├── database_file/
│   ├── cinema_booking_system-backup.sql
│   ├── migration_otp.sql
│   └── System Ticket Bioskop.postman_collection.json
│
├── internal/
│   ├── adaptor/            # HTTP Handlers
│   │   ├── adaptor.go
│   │   ├── auth.go
│   │   ├── booking.go
│   │   ├── cinema.go
│   │   ├── movie.go
│   │   ├── payment.go
│   │   └── seat.go
│   │
│   ├── data/
│   │   ├── entity/         # Data models
│   │   │   ├── booking.go
│   │   │   ├── cinema.go
│   │   │   ├── model.go
│   │   │   ├── otp.go
│   │   │   ├── session.go
│   │   │   └── user.go
│   │   │
│   │   └── repository/     # Database operations
│   │       ├── auth.go
│   │       ├── auth_test.go
│   │       ├── booking.go
│   │       ├── booking_test.go
│   │       ├── cinema.go
│   │       ├── cinema_test.go
│   │       ├── db_interface.go
│   │       ├── movie.go
│   │       ├── movie_test.go
│   │       ├── payment.go
│   │       ├── payment_test.go
│   │       ├── repository.go
│   │       ├── seat.go
│   │       └── seat_test.go
│   │
│   ├── dto/                # Data Transfer Objects
│   │   ├── pagination.go
│   │   ├── request.go
│   │   └── response.go
│   │
│   ├── usecase/            # Business logic
│   │   ├── auth.go
│   │   ├── auth_test.go
│   │   ├── booking.go
│   │   ├── booking_test.go
│   │   ├── cinema.go
│   │   ├── cinema_test.go
│   │   ├── movie.go
│   │   ├── movie_test.go
│   │   ├── payment.go
│   │   ├── payment_seat_test.go
│   │   ├── seat.go
│   │   ├── usecase.go
│   │   └── user.go
│   │
│   └── wire/
│       └── wire.go         # Dependency injection & routing
│
├── logs/                   # Log files
│
└── pkg/
    ├── database/
    │   └── postgres.go     # Database connection
    │
    ├── middleware/
    │   ├── auth.go         # Authentication middleware
    │   ├── logging.go      # Logging middleware
    │   └── middleware.go
    │
    └── utils/
        ├── config.go       # Configuration loader
        ├── email.go        # Email service
        ├── logger.go       # Zap logger setup
        └── response.go     # Response helpers (DRY)
```

---

## 👤 Author

**Alvin** - Bootcamp Golang Mini Project

---

## 📄 License

This project is for educational purposes as part of Bootcamp Golang Mini Project.
