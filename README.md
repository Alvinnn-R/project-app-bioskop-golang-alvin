# 🎬 Cinema Booking System API

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Chi Router](https://img.shields.io/badge/Chi_Router-000000?style=for-the-badge&logo=go&logoColor=white)
![UUID](https://img.shields.io/badge/UUID-Token-orange?style=for-the-badge)

RESTful API untuk sistem pemesanan tiket bioskop dengan fitur autentikasi token, verifikasi email OTP, dan booking kursi. Dibuat menggunakan Go (Golang), Chi Router, dan PostgreSQL sebagai project **Golang Advance - Lumoshive Academy Bootcamp**.

## Video Demo

[![Watch Demo](https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white)](https://youtu.be/-NIZbD5uNIk)

**[Tonton Video Penjelasan Project Cinema Booking System](https://youtu.be/-NIZbD5uNIk)**

---

## Fitur Utama

- **Token Authentication** - Login system dengan UUID Token & Session Table
- **Email OTP Verification** - Verifikasi email dengan 6-digit OTP (expired 5 menit)
- **Booking Seats** - Pemesanan kursi dengan database transaction
- **Protected Routes** - Middleware pattern untuk route yang membutuhkan login
- **Movie & Cinema Management** - CRUD film dan bioskop
- **Showtime Management** - Jadwal tayang film
- **Payment Processing** - Proses pembayaran booking
- **Email Notification** - Notifikasi email async menggunakan goroutine
- **Concurrent Data Fetching** - Goroutine dengan channels untuk fetch data parallel
- **Logging System** - Zap Logger untuk structured logging
- **Unit Testing** - Testing repository layer dengan pgxmock (81.1% coverage)

---

## Konsep Programming yang Diimplementasikan

### 1. Clean Architecture Pattern

Arsitektur 3-layer: **Adaptor (Handler) → UseCase (Business Logic) → Repository (Data Access)**

- **Adaptor Layer**: HTTP request/response handling
- **UseCase Layer**: Business logic & validasi
- **Repository Layer**: Data access dengan PostgreSQL
- **Dependency Injection**: Interface-based design

### 2. Token Authentication

- Generate UUID Token saat login
- Simpan session ke database dengan expired time (24 jam)
- Validasi token menggunakan middleware

### 3. Email OTP Verification

- Generate random 6-digit OTP
- Expired time 5 menit
- Kirim OTP via email API secara async

### 4. Middleware Pattern

- Protected routes menggunakan middleware `RequireAuth`
- Token validation dari header `Authorization`
- Inject userID ke context untuk handler

### 5. Database Transaction

- Menggunakan `BeginTx()` untuk atomicity
- Insert booking dan booking_seats dalam satu transaction
- Rollback otomatis jika ada error

### 6. Goroutines Implementation

**Fire-and-Forget:**

- Kirim email notification async tanpa blocking

**Concurrent dengan Channels:**

- Fetch data parallel (movie, showtime, seats)
- Collect results dengan buffered channel
- Timeout handling dengan context

### 7. DRY Principle

- Reusable response functions (`ResponseSuccess`, `ResponseError`, dll)
- Centralized email service
- Repository pattern untuk data access

### 8. Unit Testing

- Repository layer testing dengan pgxmock
- Mock database connection
- Test coverage 81.1%

---

## Struktur Project

```
project-app-bioskop-golang-alvin/
├── cmd/
│   └── server.go              # Server configuration
├── database_file/
│   ├── cinema_booking_system-backup.sql  # Database schema
│   ├── migration_otp.sql      # OTP table migration
│   └── System Ticket Bioskop.postman_collection.json
├── internal/
│   ├── adaptor/               # HTTP Handlers
│   │   ├── adaptor.go         # Adaptor aggregator
│   │   ├── auth.go            # Auth handlers (register, login, logout, OTP)
│   │   ├── booking.go         # Booking handlers
│   │   ├── cinema.go          # Cinema handlers
│   │   ├── movie.go           # Movie handlers
│   │   ├── payment.go         # Payment handlers
│   │   └── seat.go            # Seat handlers
│   ├── data/
│   │   ├── entity/            # Domain models
│   │   │   ├── booking.go     # Booking & Payment entities
│   │   │   ├── cinema.go      # Cinema & Showtime entities
│   │   │   ├── otp.go         # OTP entity
│   │   │   ├── session.go     # Session entity
│   │   │   └── user.go        # User entity
│   │   └── repository/        # Data access layer
│   │       ├── auth.go        # User & session repository
│   │       ├── auth_test.go   # Auth repository tests
│   │       ├── booking.go     # Booking repository
│   │       ├── booking_test.go
│   │       ├── cinema.go      # Cinema repository
│   │       ├── cinema_test.go
│   │       ├── db_interface.go # DBPool interface for mocking
│   │       ├── movie.go       # Movie repository
│   │       ├── movie_test.go
│   │       ├── payment.go     # Payment repository
│   │       ├── payment_test.go
│   │       ├── repository.go  # Repository aggregator
│   │       ├── seat.go        # Seat repository
│   │       └── seat_test.go
│   ├── dto/                   # Data Transfer Objects
│   │   ├── pagination.go      # Pagination DTO
│   │   ├── request.go         # Request DTOs
│   │   └── response.go        # Response DTOs
│   ├── usecase/               # Business logic layer
│   │   ├── auth.go            # Auth logic (register, login, OTP)
│   │   ├── auth_test.go
│   │   ├── booking.go         # Booking logic with goroutines
│   │   ├── booking_test.go
│   │   ├── cinema.go          # Cinema logic
│   │   ├── cinema_test.go
│   │   ├── movie.go           # Movie logic
│   │   ├── movie_test.go
│   │   ├── payment.go         # Payment logic with email notification
│   │   ├── seat.go            # Seat logic
│   │   ├── usecase.go         # UseCase aggregator
│   │   └── user.go            # User logic
│   └── wire/
│       └── wire.go            # Dependency injection & routing
├── logs/                      # Application logs
├── pkg/
│   ├── database/
│   │   └── postgres.go        # PostgreSQL connection
│   ├── middleware/
│   │   ├── auth.go            # Authentication middleware
│   │   ├── logging.go         # Logging middleware
│   │   └── middleware.go      # Middleware aggregator
│   └── utils/
│       ├── config.go          # Configuration loader (Viper)
│       ├── email.go           # Email service (OTP & notifications)
│       ├── logger.go          # Zap logger setup
│       └── response.go        # Response helpers (DRY)
├── main.go                    # Entry point
├── go.mod                     # Go modules
└── README.md                  # Documentation
```

---

## Cara Menggunakan

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Git

### Instalasi

1. **Clone repository**

   ```bash
   git clone https://github.com/Alvinnn-R/project-app-bioskop-golang-alvin.git
   cd project-app-bioskop-golang-alvin
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Setup database**

   - Buat database: `createdb cinema_booking`
   - Import schema: `psql -U postgres -d cinema_booking -f database_file/cinema_booking_system-backup.sql`
   - Import OTP migration: `psql -U postgres -d cinema_booking -f database_file/migration_otp.sql`

4. **Konfigurasi .env**

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=cinema_booking
   SERVER_PORT=8080
   ```

5. **Jalankan aplikasi**

   ```bash
   go run .
   ```

6. **Akses API**: `http://localhost:8080`

---

## API Endpoints

### Public Endpoints (Tanpa Login)

| Method | Endpoint       | Deskripsi            |
| ------ | -------------- | -------------------- |
| POST   | `/register`    | Registrasi user baru |
| POST   | `/login`       | Login user           |
| POST   | `/verify-otp`  | Verifikasi OTP email |
| POST   | `/resend-otp`  | Kirim ulang OTP      |
| GET    | `/movies`      | Daftar semua film    |
| GET    | `/movies/{id}` | Detail film          |
| GET    | `/cinemas`     | Daftar semua bioskop |
| GET    | `/showtimes`   | Daftar jadwal tayang |

### Protected Endpoints (Perlu Login)

| Method | Endpoint               | Deskripsi             |
| ------ | ---------------------- | --------------------- |
| POST   | `/logout`              | Logout user           |
| POST   | `/booking`             | Buat booking baru     |
| GET    | `/user/bookings`       | Daftar booking user   |
| POST   | `/pay`                 | Proses pembayaran     |
| GET    | `/seats/{showtime_id}` | Daftar kursi tersedia |

## Database Schema

### Tables

- **users** - Data user (username, email, password_hash, is_verified)
- **otps** - OTP verification codes
- **sessions** - User sessions dengan token
- **movies** - Data film
- **cinemas** - Data bioskop
- **studios** - Studio dalam bioskop
- **showtimes** - Jadwal tayang film
- **seats** - Kursi dalam studio
- **bookings** - Data pemesanan
- **booking_seats** - Kursi yang dipesan
- **payment_methods** - Metode pembayaran
- **payments** - Data pembayaran

### ERD

Lihat file `database_file/cinema_booking_system-backup.sql` untuk schema lengkap.

---

## Author

**Alvin**

- GitHub: [@Alvinnn-R](https://github.com/Alvinnn-R)
- Bootcamp: Golang Advance - Lumoshive Academy

---

## License

This project is for educational purposes as part of Lumoshive Academy Bootcamp.

---

## Acknowledgments

- Lumoshive Academy - Golang Bootcamp
- Instructor & Mentors
- Fellow bootcamp participants
