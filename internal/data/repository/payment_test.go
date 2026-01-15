package repository

import (
	"context"
	"errors"
	"project-app-bioskop/internal/data/entity"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRepo_GetAllPaymentMethods(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepo(mock)

	t.Run("Success - Get All Payment Methods", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "name",
		}).
			AddRow(1, "Credit Card").
			AddRow(2, "Cash").
			AddRow(3, "QRIS")

		mock.ExpectQuery("SELECT (.+) FROM payment_methods").
			WillReturnRows(rows)

		methods, err := repo.GetAllPaymentMethods(context.Background())
		assert.NoError(t, err)
		assert.Len(t, methods, 3)
		assert.Equal(t, "Credit Card", methods[0].Name)
		assert.Equal(t, "Cash", methods[1].Name)
		assert.Equal(t, "QRIS", methods[2].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty Result", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "name",
		})

		mock.ExpectQuery("SELECT (.+) FROM payment_methods").
			WillReturnRows(rows)

		methods, err := repo.GetAllPaymentMethods(context.Background())
		assert.NoError(t, err)
		assert.Len(t, methods, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM payment_methods").
			WillReturnError(errors.New("database error"))

		methods, err := repo.GetAllPaymentMethods(context.Background())
		assert.Error(t, err)
		assert.Nil(t, methods)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentRepo_GetPaymentMethodByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepo(mock)

	t.Run("Success - Payment Method Found", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "name",
		}).AddRow(1, "Credit Card")

		mock.ExpectQuery("SELECT (.+) FROM payment_methods WHERE id").
			WithArgs(1).
			WillReturnRows(rows)

		method, err := repo.GetPaymentMethodByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, method.ID)
		assert.Equal(t, "Credit Card", method.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM payment_methods WHERE id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		method, err := repo.GetPaymentMethodByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, 0, method.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentRepo_CreatePayment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepo(mock)

	t.Run("Success", func(t *testing.T) {
		payment := entity.Payment{
			BookingID:       1,
			PaymentMethodID: 1,
			Status:          "pending",
			PaymentDetails:  "Payment for booking #1",
		}

		rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

		mock.ExpectQuery("INSERT INTO payments").
			WithArgs(payment.BookingID, payment.PaymentMethodID, payment.Status, payment.PaymentDetails).
			WillReturnRows(rows)

		id, err := repo.CreatePayment(context.Background(), payment)
		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		payment := entity.Payment{
			BookingID:       1,
			PaymentMethodID: 1,
			Status:          "pending",
			PaymentDetails:  "Payment for booking #1",
		}

		mock.ExpectQuery("INSERT INTO payments").
			WithArgs(payment.BookingID, payment.PaymentMethodID, payment.Status, payment.PaymentDetails).
			WillReturnError(errors.New("database error"))

		id, err := repo.CreatePayment(context.Background(), payment)
		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentRepo_GetPaymentByBookingID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepo(mock)

	t.Run("Success - Payment Found", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{
			"id", "booking_id", "payment_method_id", "status", "payment_details", "paid_at",
		}).AddRow(1, 1, 1, "success", "Payment completed", nil)

		mock.ExpectQuery("SELECT (.+) FROM payments WHERE booking_id").
			WithArgs(1).
			WillReturnRows(rows)

		payment, err := repo.GetPaymentByBookingID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, payment.ID)
		assert.Equal(t, 1, payment.BookingID)
		assert.Equal(t, "success", payment.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM payments WHERE booking_id").
			WithArgs(999).
			WillReturnError(pgx.ErrNoRows)

		payment, err := repo.GetPaymentByBookingID(context.Background(), 999)
		assert.Error(t, err)
		assert.Equal(t, 0, payment.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPaymentRepo_UpdatePaymentStatus(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewPaymentRepo(mock)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE payments SET status").
			WithArgs("success", 1).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := repo.UpdatePaymentStatus(context.Background(), 1, "success")
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error - Database Error", func(t *testing.T) {
		mock.ExpectExec("UPDATE payments SET status").
			WithArgs("success", 1).
			WillReturnError(errors.New("database error"))

		err := repo.UpdatePaymentStatus(context.Background(), 1, "success")
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
