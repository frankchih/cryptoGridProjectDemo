package orderApp

import (
	"gorm.io/gorm"
)

type OrderSymbolService struct {
	db *gorm.DB
}

func NewOrderSymbolService(db *gorm.DB) *OrderSymbolService {
	return &OrderSymbolService{db: db}
}

func (orderSymbolService *OrderSymbolService) GetOrderSymbol(orderSymbolId int) (*OrderSymbol, error) {
	var orderSymbol *OrderSymbol
	err := orderSymbolService.db.First(&orderSymbol, orderSymbolId).Error
	if err != nil {
		return nil, err
	}
	return orderSymbol, nil
}
func (orderSymbolService *OrderSymbolService) GetOrderSymbolByOrderId(orderId int64) (*OrderSymbol, error) {
	var orderSymbol *OrderSymbol
	err := orderSymbolService.db.First(&orderSymbol, orderId).Error
	if err != nil {
		return nil, err
	}
	return orderSymbol, nil
}
func (orderSymbolService *OrderSymbolService) GetOrderSymbolByOrderIdList(orderIdList []int64) ([]*OrderSymbol, error) {
	var orderSymbols []*OrderSymbol
	err := orderSymbolService.db.Where("order_id IN ? ", orderIdList).Find(&orderSymbols).Error
	if err != nil {
		return nil, err
	}
	return orderSymbols, nil
}
func (orderSymbolService *OrderSymbolService) GetOrderSymbols() (*[]*OrderSymbol, error) {
	var orderSymbols []*OrderSymbol
	err := orderSymbolService.db.Find(&orderSymbols).Error
	if err != nil {
		return nil, err
	}
	return &orderSymbols, nil
}
func (orderSymbolService *OrderSymbolService) CreateOrderSymbol(orderSymbol *OrderSymbol) error {
	return orderSymbolService.db.Create(orderSymbol).Error
}

func (orderSymbolService *OrderSymbolService) UpdateOrderSymbol(orderSymbol *OrderSymbol, updateOrderSymbol *OrderSymbol) error {
	return orderSymbolService.db.Model(&orderSymbol).Updates(updateOrderSymbol).Error
}
func (orderSymbolService *OrderSymbolService) DeleteOrderSymbol(orderSymbolId int) error {
	return orderSymbolService.db.Delete(&OrderSymbol{}, orderSymbolId).Error
}
