package reportsusecases

import response "MyShoo/internal/models/responseModels"

// GetTopModels implements usecaseInterface.ReportsUC.
func (uc *ReportsUseCase) GetTopModels(limit int) (*[]response.TopModels, error) {
	resp, err := uc.reportsRepo.GetTopModels(limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTopProducts implements usecaseInterface.ReportsUC.
func (uc *ReportsUseCase) GetTopProducts(limit int) (*[]response.TopProducts, error) {
	resp, err := uc.reportsRepo.GetTopProducts(limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTopBrands implements usecaseInterface.ReportsUC.
func (uc *ReportsUseCase) GetTopBrands(limit int) (*[]response.TopBrands, error) {
	resp, err := uc.reportsRepo.GetTopBrands(limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetTopSellers implements usecaseInterface.ReportsUC.
func (uc *ReportsUseCase) GetTopSellers(limit int) (*[]response.TopSellers, error) {
	resp, err := uc.reportsRepo.GetTopSellers(limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
