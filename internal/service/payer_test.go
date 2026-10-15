package service_test

import (
	"github.com/xegors/gocoffee/internal/domain"
	"github.com/xegors/gocoffee/internal/service"
)

var _ service.Payer = (*domain.Wallet)(nil)
