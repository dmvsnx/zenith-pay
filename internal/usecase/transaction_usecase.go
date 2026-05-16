package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
	"gorm.io/gorm"
)

type TransactionUsecase interface {
	CreateTransaction(userID, shiftID string, req *dtos.TransactionRequest) (*dtos.TransactionResponse, error)
	GetTransactionByID(id string) (*dtos.TransactionResponse, error)
	GetAllTransaction(page, limit int) ([]*dtos.TransactionResponse, int64, error)
}

type transactionUsecase struct {
	db *gorm.DB
	transactionRepo repository.TransactionRepository
	itemRepo repository.TransactionItemRepository
	productRepo repository.ProductRepository
}

func NewTransactionUsecase(db *gorm.DB, tr repository.TransactionRepository, ir repository.TransactionItemRepository, pr repository.ProductRepository) TransactionUsecase {
	return &transactionUsecase{
		db: db,
		transactionRepo: tr,
		itemRepo: ir,
		productRepo: pr,
	}
}

func (u *transactionUsecase) CreateTransaction(userID, shiftID string, req *dtos.TransactionRequest) (*dtos.TransactionResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("transaction items cannot be empty")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	paymentMethod := model.PaymentMethod(req.PaymentMethod)
	if !paymentMethod.IsValid() {
		return nil, errors.New("invalid payment method")
	}

	var (
		totalAmount     int64
		transaction     model.Transaction
		transactionItems []model.TransactionItems
		responseItems   []dtos.TransactionItemResponse
	)

	if err := u.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Items {
			product, err := u.productRepo.FindByIDForUpdate(tx, item.ProductID)
			if err != nil {
				return fmt.Errorf("product not found: %s", item.ProductID)
			}

			if product.Stock < item.Quantity {
				return fmt.Errorf("insufficient stock for product %s", product.Name)
			}

			subTotal := product.Price * int64(item.Quantity)
			totalAmount += subTotal

			product.Stock -= item.Quantity
			if err := u.productRepo.UpdateTx(tx, product); err != nil {
				return err
			}

			transactionItems = append(transactionItems, model.TransactionItems{
				ProductID:    product.ID,
				ProductName:  product.Name,
				ProductPrice: product.Price,
				Quantity:     item.Quantity,
				Subtotal:     subTotal,
			})

			responseItems = append(responseItems, dtos.TransactionItemResponse{
				ProductID:    product.ID.String(),
				ProductName:  product.Name,
				ProductPrice: product.Price,
				Quantity:     item.Quantity,
				SubTotal:     subTotal,
			})
		}

		paymentAmount := req.PaymentAmount
		var changeAmount int64

		if paymentMethod == model.Cash {
			if paymentAmount < totalAmount {
				return errors.New("payment amount is less than total amount")
			}
			changeAmount = paymentAmount - totalAmount
		} else {
			paymentAmount = totalAmount
			changeAmount = 0
		}

		shiftUUID, err := uuid.Parse(shiftID)
		if err != nil {
			return errors.New("invalid shift ID")
		}

		transaction = model.Transaction{
			UserID:          userUUID,
			ShiftID:         shiftUUID,
			TransactionDate: time.Now(),
			PaymentMethod:   paymentMethod,
			TotalAmount:     totalAmount,
			PaymentAmount:   paymentAmount,
			ChangeAmount:    changeAmount,
		}

		if err := u.transactionRepo.Create(tx, &transaction); err != nil {
			return err
		}

		for i := range transactionItems {
			transactionItems[i].TransactionID = transaction.ID
		}

		if err := u.itemRepo.CreateMany(tx, transactionItems); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	res := &dtos.TransactionResponse{
		ID:              transaction.ID.String(),
		UserID:          transaction.UserID.String(),
		TransactionDate: transaction.TransactionDate,
		PaymentMethod:   string(transaction.PaymentMethod),
		TotalAmount:     transaction.TotalAmount,
		PaymentAmount:   transaction.PaymentAmount,
		ChangeAmount:    transaction.ChangeAmount,
		Items:           responseItems,
	}

	return res, nil
}

func (u *transactionUsecase) GetTransactionByID(id string) (*dtos.TransactionResponse, error) {
	transaction, err := u.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	resItems := []dtos.TransactionItemResponse{}
	for _, item := range transaction.TransactionItems {
		resItems = append(resItems, dtos.TransactionItemResponse{
			ProductID: item.ProductID.String(),
			ProductName: item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity: item.Quantity,
			SubTotal: item.Subtotal,
		})
	}

	res := &dtos.TransactionResponse{
		ID: transaction.ID.String(),
		UserID: transaction.UserID.String(),
		TransactionDate: transaction.TransactionDate,
		PaymentMethod: string(transaction.PaymentMethod),
		TotalAmount: transaction.TotalAmount,
		PaymentAmount: transaction.PaymentAmount,
		ChangeAmount: transaction.ChangeAmount,
		Items: resItems,
	}

	return res, nil
}

func (u *transactionUsecase) GetAllTransaction(page, limit int) ([]*dtos.TransactionResponse, int64, error) {
	transactions, total, err := u.transactionRepo.FindAllPaginated((page-1)*limit, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dtos.TransactionResponse
	for _, t := range transactions {
		resItems := []dtos.TransactionItemResponse{}
		for _, item := range t.TransactionItems {
			resItems = append(resItems, dtos.TransactionItemResponse{
				ProductID:    item.ProductID.String(),
				ProductName:  item.ProductName,
				ProductPrice: item.ProductPrice,
				Quantity:     item.Quantity,
				SubTotal:     item.Subtotal,
			})
		}

		responses = append(responses, &dtos.TransactionResponse{
			ID:              t.ID.String(),
			UserID:          t.UserID.String(),
			TransactionDate: t.TransactionDate,
			PaymentMethod:   string(t.PaymentMethod),
			PaymentAmount:   t.PaymentAmount,
			TotalAmount:     t.TotalAmount,
			ChangeAmount:    t.ChangeAmount,
			Items:           resItems,
		})
	}

	return responses, total, nil
}
