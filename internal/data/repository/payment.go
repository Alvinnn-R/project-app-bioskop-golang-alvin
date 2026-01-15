package repository

import (
	"context"
	"project-app-bioskop/internal/data/entity"
)

type PaymentRepoInterface interface {
	GetAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, id int) (entity.PaymentMethod, error)
	CreatePayment(ctx context.Context, payment entity.Payment) (int, error)
	GetPaymentByBookingID(ctx context.Context, bookingID int) (entity.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error
}

type PaymentRepo struct {
	DB DBPool
}

func NewPaymentRepo(db DBPool) PaymentRepoInterface {
	return &PaymentRepo{DB: db}
}

// GetAllPaymentMethods retrieves all available payment methods
func (r *PaymentRepo) GetAllPaymentMethods(ctx context.Context) ([]entity.PaymentMethod, error) {
	query := `SELECT id, name FROM payment_methods ORDER BY id`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []entity.PaymentMethod
	for rows.Next() {
		var m entity.PaymentMethod
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return methods, nil
}

// GetPaymentMethodByID retrieves payment method by ID
func (r *PaymentRepo) GetPaymentMethodByID(ctx context.Context, id int) (entity.PaymentMethod, error) {
	query := `SELECT id, name FROM payment_methods WHERE id = $1`
	var m entity.PaymentMethod
	err := r.DB.QueryRow(ctx, query, id).Scan(&m.ID, &m.Name)
	if err != nil {
		return m, err
	}
	return m, nil
}

// CreatePayment creates a new payment record
func (r *PaymentRepo) CreatePayment(ctx context.Context, payment entity.Payment) (int, error) {
	query := `INSERT INTO payments (booking_id, payment_method_id, status, payment_details, paid_at) 
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	var id int
	err := r.DB.QueryRow(ctx, query,
		payment.BookingID, payment.PaymentMethodID, payment.Status, payment.PaymentDetails,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetPaymentByBookingID retrieves payment by booking ID
func (r *PaymentRepo) GetPaymentByBookingID(ctx context.Context, bookingID int) (entity.Payment, error) {
	query := `SELECT id, booking_id, payment_method_id, status, payment_details, paid_at 
			  FROM payments WHERE booking_id = $1`
	var p entity.Payment
	err := r.DB.QueryRow(ctx, query, bookingID).Scan(
		&p.ID, &p.BookingID, &p.PaymentMethodID, &p.Status, &p.PaymentDetails, &p.PaidAt,
	)
	if err != nil {
		return p, err
	}
	return p, nil
}

// UpdatePaymentStatus updates payment status
func (r *PaymentRepo) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	query := `UPDATE payments SET status = $1, paid_at = NOW() WHERE id = $2`
	_, err := r.DB.Exec(ctx, query, status, paymentID)
	return err
}
