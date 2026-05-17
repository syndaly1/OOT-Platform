package usecase

import (
	"context"
	"streaming-platform/auth-service/internal/domain"
	"time"
)

type AccountUsecase struct {
	accountRepo domain.AccountRepository
}

func NewAccountUsecase(accountRepo domain.AccountRepository) *AccountUsecase {
	return &AccountUsecase{accountRepo: accountRepo}
}

func (uc *AccountUsecase) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return uc.accountRepo.GetByID(ctx, id)
}

func (uc *AccountUsecase) UpdateAccount(ctx context.Context, id, username string) error {
	account, err := uc.accountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	account.Username = username
	account.UpdatedAt = time.Now()
	return uc.accountRepo.Update(ctx, account)
}

func (uc *AccountUsecase) DeleteAccount(ctx context.Context, id string) error {
	return uc.accountRepo.Delete(ctx, id)
}
