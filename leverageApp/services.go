package leverageApp

import (
	"gorm.io/gorm"
)

type LeverageSymbolService struct {
	db *gorm.DB
}

func NewLeverageSymbolService(db *gorm.DB) *LeverageSymbolService {
	return &LeverageSymbolService{db: db}
}

func (leverageSymbolService *LeverageSymbolService) GetLeverageSymbol(leverageSymbolId int) (*LeverageSymbol, error) {
	var leverageSymbol *LeverageSymbol
	err := leverageSymbolService.db.First(&leverageSymbol, leverageSymbolId).Error
	if err != nil {
		return nil, err
	}
	return leverageSymbol, nil
}
func (leverageSymbolService *LeverageSymbolService) GetLeverageSymbolBySymbol(symbol string) (*LeverageSymbol, error) {
	var leverageSymbol *LeverageSymbol
	err := leverageSymbolService.db.Where("symbol = ?", symbol).First(&leverageSymbol).Error
	if err != nil {
		return nil, err
	}
	return leverageSymbol, nil
}

func (leverageSymbolService *LeverageSymbolService) GetLeverageSymbols() (*[]*LeverageSymbol, error) {
	var leverageSymbols []*LeverageSymbol
	err := leverageSymbolService.db.Find(&leverageSymbols).Error
	if err != nil {
		return nil, err
	}
	return &leverageSymbols, nil
}
func (leverageSymbolService *LeverageSymbolService) CreateLeverageSymbol(leverageSymbol *LeverageSymbol) error {
	return leverageSymbolService.db.Create(leverageSymbol).Error
}

func (leverageSymbolService *LeverageSymbolService) UpdateLeverageSymbol(leverageSymbol *LeverageSymbol, updateLeverageSymbol *LeverageSymbol) error {
	return leverageSymbolService.db.Model(&leverageSymbol).Updates(updateLeverageSymbol).Error
}
func (leverageSymbolService *LeverageSymbolService) DeleteLeverageSymbol(leverageSymbolId int) error {
	return leverageSymbolService.db.Delete(&LeverageSymbol{}, leverageSymbolId).Error
}
