package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductId int
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	CategoryID   int
	Category     Category
	SerialNumber SerialNumber
	gorm.Model   // adiciona os campos ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)
}

func main() {
	dsn := "sa:sasa@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	// cretate category
	category := Category{Name: "Cozinha"}
	db.Create(&category)
	//
	//////create product
	//db.Create(&Product{
	//	Name:       "Panela",
	//	Price:      100.90,
	//	CategoryID: 1,
	//})
	//
	//// create serial number
	//db.Create(&SerialNumber{
	//	Number:    "123456",
	//	ProductId: 1,
	//})

	//var products []Product
	//db.Preload("Category").Preload("SerialNumber").Find(&products)
	//for _, product := range products {
	//	fmt.Println(product.Name, product.Category.Name)
	//}

	var categories []Category
	err = db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			println("-", product.Name, "Serial Number:", product.SerialNumber.Number)
		}
	}

}
