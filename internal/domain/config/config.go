package config

import (
	myMath "MyShoo/pkg/math"
	"encoding/json"
	"fmt"
	"os"
)

var configRead struct {
	FreeDeliveryPincodeRanges  []string `json:"freeDeliveryPincodeRanges"`
	IntermediatePincodeRanges  []string `json:"intermediatePincodeRanges"`
	IntermediateDeliveryCharge float32  `json:"intermediateDeliveryCharge"`
	DistantDeliveryCharge      float32  `json:"distantDeliveryCharge"`
	MaxOrderAmountForCOD       float32  `json:"maxOrderAmountForCOD"`
	OrderAmountForFreeDelivery float32  `json:"orderAmountForFreeDelivery"`
	CashOnDeliveryAvailable    bool     `json:"cashOnDeliveryAvailable"`
}

var DeliveryConfig struct {
	FreeDeliveryPincodeRanges  []struct{ Start, End uint }
	IntermediatePincodeRanges  []struct{ Start, End uint }
	IntermediateDeliveryCharge float32
	DistantDeliveryCharge      float32
	MaxOrderAmountForCOD       float32
	OrderAmountForFreeDelivery float32
	CashOnDeliveryAvailable    bool
}

func RestartConfig() error {
	return LoadConfig()
}
func LoadConfig() error {
	if err := loadDeliveryConfig(); err != nil {
		return err
	}
	return nil
}
func loadDeliveryConfig() error {
	filePath := "config/shippingCharges.json"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &configRead); err != nil {
		return err
	}

	for _, v := range configRead.FreeDeliveryPincodeRanges {
		var start, end uint
		_, err := fmt.Sscanf(v, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		DeliveryConfig.FreeDeliveryPincodeRanges = append(DeliveryConfig.FreeDeliveryPincodeRanges, struct{ Start, End uint }{start, end})
	}
	for _, v := range configRead.IntermediatePincodeRanges {
		var start, end uint
		_, err := fmt.Sscanf(v, "%d-%d", &start, &end)
		if err != nil {
			return err
		}
		DeliveryConfig.IntermediatePincodeRanges = append(DeliveryConfig.IntermediatePincodeRanges, struct{ Start, End uint }{start, end})
	}
	DeliveryConfig.IntermediateDeliveryCharge = myMath.RoundFloat32(configRead.IntermediateDeliveryCharge,2)
	DeliveryConfig.DistantDeliveryCharge = myMath.RoundFloat32(configRead.DistantDeliveryCharge,2)
	DeliveryConfig.MaxOrderAmountForCOD = myMath.RoundFloat32(configRead.MaxOrderAmountForCOD,2)
	DeliveryConfig.OrderAmountForFreeDelivery = myMath.RoundFloat32(configRead.OrderAmountForFreeDelivery,2)
	DeliveryConfig.CashOnDeliveryAvailable = configRead.CashOnDeliveryAvailable

	//validate
	for _, v1 := range DeliveryConfig.FreeDeliveryPincodeRanges {
		for _, v2 := range DeliveryConfig.IntermediatePincodeRanges {
			if v1.Start <= v2.Start && v1.End >= v2.Start {
				return fmt.Errorf("IntermediatePincodeRanges should not overlap values of FreeDeliveryPincodeRanges")
			} else if v1.Start <= v2.End && v1.End >= v2.End {
				return fmt.Errorf("IntermediatePincodeRanges should not overlap values of FreeDeliveryPincodeRanges")
			}
		}
	}
	if DeliveryConfig.IntermediateDeliveryCharge < 0 ||
		DeliveryConfig.DistantDeliveryCharge < 0 ||
		DeliveryConfig.MaxOrderAmountForCOD < 0 ||
		DeliveryConfig.OrderAmountForFreeDelivery < 0 {
		return fmt.Errorf("IntermediateDeliveryCharge, DistantDeliveryCharge, MaxOrderAmountForCOD should not be negative")
	}

	return nil
}
