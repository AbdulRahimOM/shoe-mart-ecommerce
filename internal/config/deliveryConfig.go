package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	myMath "github.com/AbdulRahimOM/shoe-mart-ecommerce/pkg/math"
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

func LoadDeliveryConfig() error {
	relativePath := "config/shippingCharges.json"
	filePath := filepath.Join(ExecutableDir, relativePath)
	if err := loadDeliveryConfig(filePath); err != nil {
		fmt.Println("Error loading delivery config (via executing from binary mode). err= ", err)
		fmt.Println("Trying to restart config(in dev mode., i.e. via 'go run cmd/main.go' command))....")
		if err2 := loadDeliveryConfig(relativePath); err2 != nil {
			fmt.Println("Error loading config(in dev mode). err= ", err2)
			return fmt.Errorf("from binary run mode: %v, from normal go run mode: %v", err, err2)
		}
		return nil
	}
	return nil
}

func RestartDeliveryConfig() error {
	return LoadDeliveryConfig()
}

func loadDeliveryConfig(filePath string) error {
	var err error
	defer func() {
		if err != nil {
			failedDiagram()
		}
	}()

	preDiagram()

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file. err= ", err)
		return err
	}

	if err := json.Unmarshal(data, &configRead); err != nil {
		fmt.Println("Error unmarshalling data. err= ", err)
		return err
	}

	for _, v := range configRead.FreeDeliveryPincodeRanges {
		var start, end uint
		_, err := fmt.Sscanf(v, "%d-%d", &start, &end)
		if err != nil {
			fmt.Println("Error scanning data. err= ", err, "v= ", v)
			return err
		}
		DeliveryConfig.FreeDeliveryPincodeRanges = append(DeliveryConfig.FreeDeliveryPincodeRanges, struct{ Start, End uint }{start, end})
	}
	for _, v := range configRead.IntermediatePincodeRanges {
		var start, end uint
		_, err := fmt.Sscanf(v, "%d-%d", &start, &end)
		if err != nil {
			fmt.Println("Error scanning data. err= ", err, "v= ", v)
			return err
		}
		DeliveryConfig.IntermediatePincodeRanges = append(DeliveryConfig.IntermediatePincodeRanges, struct{ Start, End uint }{start, end})
	}
	DeliveryConfig.IntermediateDeliveryCharge = myMath.RoundFloat32(configRead.IntermediateDeliveryCharge, 2)
	DeliveryConfig.DistantDeliveryCharge = myMath.RoundFloat32(configRead.DistantDeliveryCharge, 2)
	DeliveryConfig.MaxOrderAmountForCOD = myMath.RoundFloat32(configRead.MaxOrderAmountForCOD, 2)
	DeliveryConfig.OrderAmountForFreeDelivery = myMath.RoundFloat32(configRead.OrderAmountForFreeDelivery, 2)
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

	successDiagram()
	return nil
}

func preDiagram() { //who doesn't love art?
	fmt.Println()
	fmt.Println("        ██████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗ ")
	fmt.Println("       ██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝ ")
	fmt.Println("       ██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗")
	fmt.Println("       ██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║")
	fmt.Println("       ╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝")
	fmt.Println("        ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝ ")
	fmt.Println()
	fmt.Println("                ╦═╗╔═╗╔═╗╔╦╗╔═╗╦═╗╔╦╗╦╔╗╔╔═╗")
	fmt.Println("                ╠╦╝║╣ ╚═╗ ║ ╠═╣╠╦╝ ║ ║║║║║ ╦")
	fmt.Println("                ╩╚═╚═╝╚═╝ ╩ ╩ ╩╩╚═ ╩ ╩╝╚╝╚═╝")
}

func successDiagram() {
	fmt.Println()
	fmt.Println("░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░")
	fmt.Println("░░      ░░░  ░░░░  ░░░      ░░░░      ░░░        ░░░      ░░░░      ░░")
	fmt.Println("▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒  ▒▒  ▒▒▒▒  ▒▒  ▒▒▒▒  ▒▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒")
	fmt.Println("▓▓      ▓▓▓  ▓▓▓▓  ▓▓  ▓▓▓▓▓▓▓▓  ▓▓▓▓▓▓▓▓      ▓▓▓▓▓      ▓▓▓▓      ▓▓")
	fmt.Println("███████  ██  ████  ██  ████  ██  ████  ██  ██████████████  ████████  █")
	fmt.Println("██      ████      ████      ████      ███        ███      ████      ██")
	fmt.Println("██████████████████████████████████████████████████████████████████████")
	fmt.Print("\n\n")
}

func failedDiagram() {
	fmt.Println("            _____ _    ___ _     _____ ____  ")
	fmt.Println("           |  ___/ \\  |_ _| |   | ____|  _ \\ ")
	fmt.Println("           | |_ / _ \\  | || |   |  _| | | | |")
	fmt.Println("           |  _/ ___ \\ | || |___| |___| |_| |")
	fmt.Println("           |_|/_/   \\_\\___|_____|_____|____/ ")
	fmt.Println("===================================================================")

}
