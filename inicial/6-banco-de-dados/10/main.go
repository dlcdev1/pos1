package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2many:products_categories;"`
}

type Product struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Price      float64
	Categories []Category `gorm:"many2many:products_categories;"`
	gorm.Model            // adiciona os campos ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)
}

func main() {
	dsn := "sa:sasa@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{}, &Category{})

	tx := db.Begin() //iniciando uma transaction
	var c Category
	err = tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, 1).Error
	if err != nil {
		panic(err)
	}
	c.Name = "Eletrônicos"
	tx.Debug().Save(&c)
	tx.Commit()

}

// select * from products where id = 1 for update;

// Otimista
//name, email, versao
//w   , w@w  , 2.0.0

// Pessimista
// Locka a table, a linha
