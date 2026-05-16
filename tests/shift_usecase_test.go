package tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/tests/mocks"
)

var (
	validCashierID = "00000000-0000-0000-0000-000000000001"
	validShiftID   = "00000000-0000-0000-0000-000000000002"
)

func TestShiftUsecase_OpenShift_Success(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return nil, nil
		},
		CreateFn: func(shift *model.Shift) error {
			shift.ID = uuid.MustParse(validShiftID)
			return nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	res, err := uc.OpenShift(validCashierID, dtos.OpenShiftRequest{OpeningBalance: 100000})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != "open" {
		t.Fatalf("expected open, got %s", res.Status)
	}
}

func TestShiftUsecase_OpenShift_InvalidCashierID(t *testing.T) {
	repo := &mocks.ShiftRepo{}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.OpenShift("not-a-uuid", dtos.OpenShiftRequest{OpeningBalance: 100000})

	if err == nil || err.Error() != "invalid cashier id" {
		t.Fatalf("expected invalid cashier id, got %v", err)
	}
}

func TestShiftUsecase_OpenShift_AlreadyActive(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return &model.Shift{Status: model.ShiftOpen}, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.OpenShift(validCashierID, dtos.OpenShiftRequest{OpeningBalance: 100000})

	if err == nil || err.Error() != "cashier already has an active shift" {
		t.Fatalf("expected cashier already has an active shift, got %v", err)
	}
}

func TestShiftUsecase_OpenShift_FailCreate(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return nil, nil
		},
		CreateFn: func(shift *model.Shift) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.OpenShift(validCashierID, dtos.OpenShiftRequest{OpeningBalance: 100000})

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestShiftUsecase_CloseShift_Success(t *testing.T) {
	openingBalance := int64(100000)
	cashTotal := int64(50000)
	closingBalance := int64(150000)
	expectedBal := openingBalance + cashTotal
	variance := closingBalance - expectedBal

	repo := &mocks.ShiftRepo{
		FindByIDFn: func(id string) (*model.Shift, error) {
			return &model.Shift{
				ID:             uuid.MustParse(validShiftID),
				CashierID:      uuid.MustParse(validCashierID),
				Status:         model.ShiftOpen,
				OpeningBalance: openingBalance,
			}, nil
		},
		CloseShiftFn: func(shift *model.Shift) error {
			return nil
		},
	}
	transactionRepo := &mocks.TransactionRepo{
		SumCashByShiftIDFn: func(shiftID string) (int64, error) {
			return cashTotal, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, transactionRepo)

	res, err := uc.CloseShift(validCashierID, dtos.CloseShiftRequest{
		ShiftID:        validShiftID,
		ClosingBalance: closingBalance,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *res.ClosingBalance != closingBalance {
		t.Fatalf("expected %d, got %d", closingBalance, *res.ClosingBalance)
	}
	if *res.ExpectedClosingBalance != expectedBal {
		t.Fatalf("expected expected closing balance %d, got %d", expectedBal, *res.ExpectedClosingBalance)
	}
	if *res.Variance != variance {
		t.Fatalf("expected variance %d, got %d", variance, *res.Variance)
	}
}

func TestShiftUsecase_CloseShift_NotFound(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindByIDFn: func(id string) (*model.Shift, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.CloseShift(validCashierID, dtos.CloseShiftRequest{
		ShiftID:        validShiftID,
		ClosingBalance: 200000,
	})

	if err == nil || err.Error() != "shift not found" {
		t.Fatalf("expected shift not found, got %v", err)
	}
}

func TestShiftUsecase_CloseShift_NotYourShift(t *testing.T) {
	otherCashier := "00000000-0000-0000-0000-000000000099"
	repo := &mocks.ShiftRepo{
		FindByIDFn: func(id string) (*model.Shift, error) {
			return &model.Shift{
				ID:       uuid.MustParse(validShiftID),
				CashierID: uuid.MustParse(otherCashier),
				Status:   model.ShiftOpen,
			}, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.CloseShift(validCashierID, dtos.CloseShiftRequest{
		ShiftID:        validShiftID,
		ClosingBalance: 200000,
	})

	if err == nil || err.Error() != "not your shift" {
		t.Fatalf("expected not your shift, got %v", err)
	}
}

func TestShiftUsecase_CloseShift_AlreadyClosed(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindByIDFn: func(id string) (*model.Shift, error) {
			return &model.Shift{
				ID:       uuid.MustParse(validShiftID),
				CashierID: uuid.MustParse(validCashierID),
				Status:   model.ShiftClose,
			}, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.CloseShift(validCashierID, dtos.CloseShiftRequest{
		ShiftID:        validShiftID,
		ClosingBalance: 200000,
	})

	if err == nil || err.Error() != "shift already closed" {
		t.Fatalf("expected shift already closed, got %v", err)
	}
}

func TestShiftUsecase_GetActiveShift_Success(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return &model.Shift{
				ID:       uuid.MustParse(validShiftID),
				CashierID: uuid.MustParse(cashierID),
				Status:   model.ShiftOpen,
			}, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	res, err := uc.GetActiveShift(validCashierID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != "open" {
		t.Fatalf("expected open, got %s", res.Status)
	}
}

func TestShiftUsecase_GetActiveShift_NoActiveShift(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return nil, nil
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.GetActiveShift(validCashierID)

	if err == nil || err.Error() != "no active shift" {
		t.Fatalf("expected no active shift, got %v", err)
	}
}

func TestShiftUsecase_GetActiveShift_DBError(t *testing.T) {
	repo := &mocks.ShiftRepo{
		FindActiveShiftByCashierFn: func(cashierID string) (*model.Shift, error) {
			return nil, errors.New("db error")
		},
	}
	uc := usecase.NewShiftUsecase(repo, &mocks.TransactionRepo{})

	_, err := uc.GetActiveShift(validCashierID)

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}
