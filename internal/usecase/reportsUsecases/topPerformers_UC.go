package reportsusecases

import (
	e "MyShoo/internal/domain/customErrors"
	response "MyShoo/internal/models/responseModels"
)

// GetTopModels implements usecase.ReportsUC.
func (uc *ReportsUseCase) GetTopModels(limit int) (*[]response.TopModels, *e.Error) {
	return uc.reportsRepo.GetTopModels(limit)
}

// GetTopProducts implements usecase.ReportsUC.
func (uc *ReportsUseCase) GetTopProducts(limit int) (*[]response.TopProducts, *e.Error) {
	return uc.reportsRepo.GetTopProducts(limit)
}

// GetTopBrands implements usecase.ReportsUC.
func (uc *ReportsUseCase) GetTopBrands(limit int) (*[]response.TopBrands, *e.Error) {
	return uc.reportsRepo.GetTopBrands(limit)
}

// GetTopSellers implements usecase.ReportsUC.
func (uc *ReportsUseCase) GetTopSellers(limit int) (*[]response.TopSellers, *e.Error) {
	return uc.reportsRepo.GetTopSellers(limit)
}
