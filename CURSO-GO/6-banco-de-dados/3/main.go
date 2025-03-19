package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type Product struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Category   Category
	gorm.Model // adiciona os campos ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)
}

func main() {
	dsn := "sa:sasa@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{})

	// cretate category
	//category := Category{Name: "Eletrônicos"}
	//db.Create(&category)
	//
	////create product
	db.Create(&Product{
		Name:       "mouse",
		Price:      1889.90,
		CategoryID: 1,
	})

	var products []Product
	db.Preload("Category").Find(&products)
	for _, product := range products {
		fmt.Println(product.Name, product.Category.Name)
	}

}
