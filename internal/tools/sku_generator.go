package tools

import (
	"MyShoo/internal/domain/entities"
	"fmt"
)

// sku generator
func GenerateNameAndSKUCode(dimensionalVariant *entities.DimensionalVariant, size *string) (string, string) {
	var skuCode, name string
	var brandName string = dimensionalVariant.FkColourVariant.FkModel.FkBrand.Name
	var modelName string = dimensionalVariant.FkColourVariant.FkModel.Name
	var colour string = dimensionalVariant.FkColourVariant.Colour
	var dimensionalVariantName string = entities.DimensionalVariations[dimensionalVariant.DVIndex].Name

	var brandcode string = dimensionalVariant.FkColourVariant.FkModel.FkBrand.Name[0:2]
	var modelAlphabet string = string(dimensionalVariant.FkColourVariant.FkModel.Name[0])
	var colourAlphabet string = string(dimensionalVariant.FkColourVariant.Colour[0])
	var colourvariantid uint = dimensionalVariant.ColourVariantID
	var dVCode = entities.DimensionalVariations[dimensionalVariant.DVIndex].Code

	//skuCode
	skuCode = fmt.Sprint(
		brandcode, "_",
		modelAlphabet, "_",
		colourAlphabet,
		colourvariantid, "_",
		dVCode, "_",
		*size)
	name = fmt.Sprint(
		brandName, " ",
		modelName, " ",
		colour, " ",
		dimensionalVariantName, " ",
		"size:", *size, "(US)")

	return name, skuCode
}
