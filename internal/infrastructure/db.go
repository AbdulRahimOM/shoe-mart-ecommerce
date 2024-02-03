package infra

import (
	"MyShoo/internal/domain/entities"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ConnectionError error = nil

func ConnectToDB() error {
	dsn := os.Getenv("DB_URL")
	DB, ConnectionError = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if ConnectionError != nil {
		fmt.Println("Couldn't connect to DB")
		return ConnectionError
	}

	//account related tables
	if err := DB.AutoMigrate(&entities.User{}); err != nil {
		fmt.Println("Failed to migrate table user. Error:", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.Admin{}); err != nil {
		fmt.Println("Failed to migrate table admin", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.Seller{}); err != nil {
		fmt.Println("Failed to migrate table seller", err)
		return err
	}

	//product related tables
	if err := DB.AutoMigrate(&entities.Categories{}); err != nil {
		fmt.Println("Failed to migrate table categories", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.Brands{}); err != nil {
		fmt.Println("Failed to migrate table brands", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.Models{}); err != nil {
		fmt.Println("Failed to migrate table models", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.Product{}); err != nil {
		fmt.Println("Failed to migrate table product", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.ColourVariant{}); err != nil {
		fmt.Println("Failed to migrate table colour_variants", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.DimensionalVariant{}); err != nil {
		fmt.Println("Failed to migrate table dimensional_variants", err)
		return err
	}

	//cart related tables
	if err := DB.AutoMigrate(&entities.Cart{}); err != nil {
		fmt.Println("Failed to migrate table cart", err)
		return err
	}

	//user address table
	if err := DB.AutoMigrate(&entities.UserAddress{}); err != nil {
		fmt.Println("Failed to migrate table user_address", err)
		return err
	}

	//order related tables
	if err := DB.AutoMigrate(&entities.Order{}); err != nil {
		fmt.Println("Failed to migrate table order", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.OrderItem{}); err != nil {
		fmt.Println("Failed to migrate table order_item", err)
		return err
	}

	//wishlist related tables
	if err := DB.AutoMigrate(&entities.WishList{}); err != nil {
		fmt.Println("Failed to migrate table wishlist", err)
		return err
	}
	if err := DB.AutoMigrate(&entities.WishListItems{}); err != nil {
		fmt.Println("Failed to migrate table wishListItems", err)
		return err
	}

	//coupon
	if err := DB.AutoMigrate(&entities.Coupon{}); err != nil {
		fmt.Println("Failed to migrate table coupon", err)
		return err
	}
	return nil
}
