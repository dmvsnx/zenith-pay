package usecase

import (
	"errors"
	"time"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
)

type ShiftUsecase interface {
	OpenShift(cashierID string, req dtos.OpenShiftRequest) (*dtos.ShiftResponse, error)
	CloseShift(cashierID string, req dtos.CloseShiftRequest) (*dtos.ShiftResponse, error)
	GetActiveShift(cashierID string) (*dtos.ShiftResponse, error)
}

type shiftUsecase struct {
	shiftRepo       repository.ShiftRepository
	transactionRepo repository.TransactionRepository
}

func NewShiftUsecase(sr repository.ShiftRepository, tr repository.TransactionRepository) ShiftUsecase {
	return &shiftUsecase{
		shiftRepo:       sr,
		transactionRepo: tr,
	}
}

func (u *shiftUsecase) OpenShift(cashierID string, req dtos.OpenShiftRequest) (*dtos.ShiftResponse, error) {
	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errors.New("invalid cashier id")
	}

	existing, _ := u.shiftRepo.FindActiveShiftByCashier(cashierUUID.String())
	if existing != nil {
		return nil, errors.New("cashier already has an active shift")
	}

	shift := &model.Shift{
		CashierID: cashierUUID,
		Status: model.ShiftOpen,
		OpeningBalance: req.OpeningBalance,
		OpenedAt: time.Now(),
	}

	if err := u.shiftRepo.Create(shift); err != nil {
		return nil, err
	}

	res := &dtos.ShiftResponse{
		ID: shift.ID.String(),
		CashierID: shift.CashierID.String(),
		Status: string(shift.Status),
		OpeningBalance: shift.OpeningBalance,
		OpenedAt: shift.OpenedAt,
	}

	return res, nil
}

func (u *shiftUsecase) CloseShift(cashierID string, req dtos.CloseShiftRequest) (*dtos.ShiftResponse, error) {
	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errors.New("invalid cashier id")
	}
	shiftUUID, err := uuid.Parse(req.ShiftID)
	if err != nil {
		return nil, errors.New("invalid shift id")
	}

	shift, err := u.shiftRepo.FindByID(shiftUUID.String())
	if err != nil {
		return nil, errors.New("shift not found")
	}

	if shift.CashierID != cashierUUID {
		return nil, errors.New("not your shift")
	}

	if shift.Status != model.ShiftOpen {
		return nil, errors.New("shift already closed")
	}

	cashIncome, debitIncome, qrisIncome, err := u.transactionRepo.SumByShiftIDGrouped(shiftUUID.String())
	if err != nil {
		return nil, errors.New("failed to calculate expected closing balance")
	}

	expectedClosingBalance := shift.OpeningBalance + cashIncome
	variance := req.ClosingBalance - expectedClosingBalance

	now := time.Now()
	shift.Status = model.ShiftClose
	shift.ClosingBalance = &req.ClosingBalance
	shift.ClosedAt = now

	if err := u.shiftRepo.CloseShift(shift); err != nil {
		return nil, err
	}

	res := &dtos.ShiftResponse{
		ID:                     shift.ID.String(),
		CashierID:              shift.CashierID.String(),
		Status:                 string(shift.Status),
		OpeningBalance:         shift.OpeningBalance,
		ClosingBalance:         shift.ClosingBalance,
		ExpectedClosingBalance: &expectedClosingBalance,
		Variance:               &variance,
		CashIncome:             &cashIncome,
		DebitIncome:            &debitIncome,
		QrisIncome:             &qrisIncome,
		OpenedAt:               shift.OpenedAt,
		ClosedAt:               &shift.ClosedAt,
	}

	return res, nil
}

func (u *shiftUsecase) GetActiveShift(cashierID string) (*dtos.ShiftResponse, error) {
	cashierUUID, err := uuid.Parse(cashierID)
	if err != nil {
		return nil, errors.New("invalid cashier id")
	}

	shift, err := u.shiftRepo.FindActiveShiftByCashier(cashierUUID.String())
	if err != nil {
		return nil, err
	}

	if shift == nil {
		return nil, errors.New("no active shift")
	}

	res := &dtos.ShiftResponse{
		ID: shift.ID.String(),
		CashierID: shift.CashierID.String(),
		Status: string(shift.Status),
		OpeningBalance: shift.OpeningBalance,
		ClosingBalance: shift.ClosingBalance,
		OpenedAt: shift.OpenedAt,
	}

	return res, nil
}
