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

	//create
	//db.Create(&Product{
	//	Name:  "Notebook",
	//	Price: 1001.00,
	//})

	//create batch

	//products := []Product{
	//	{ID: uuid.New().String(), Name: "Notebook", Price: 1004.00},
	//	{ID: uuid.New().String(), Name: "Smartphone", Price: 1051.00},
	//	{ID: uuid.New().String(), Name: "Tablet", Price: 10221.00},
	//}
	//db.Create(&products)

	//var product Product
	//db.First(&product, 1)
	//fmt.Println(product)
	//db.First(&product, "name = ?", "Notebook")
	//fmt.Println(product)

	//select all
	//var products []Product
	//db.Find(&products)
	//for _, product := range products {
	//	fmt.Println(product)
	//}

	//var products []Product
	//db.Limit(2).Offset(2).Find(&products)
	//for _, product := range products {
	//	fmt.Println(product)
	//}

	//var products []Product
	//db.Where("price > ?", 1000).Find(&products)
	//for _, product := range products {
	//	fmt.Println(product)
	//}
	//var products []Product
	//db.Where("name LIKE ?", "%K%").Find(&products)
	//for _, product := range products {
	//	fmt.Println(product)
	//}

	var p Product
	db.First(&p, 1)
	p.Name = "New mouse"
	fmt.Println(p.Name)
	//db.Save(&p)
	db.Delete(&p)
	//
	//var p2 Product
	//db.First(&p2, "name = ?", "Notebook 2")
	//fmt.Println(p2.Name)
	//
	//db.Delete(&p2)

}
