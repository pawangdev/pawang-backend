package seeder

import (
	"fmt"
	"pawang-backend/entity"

	"gorm.io/gorm"
)

var dataCategories = []entity.Category{
	{
		Name: "Belanja",
		Icon: "/api/v1/storage/assets/categories/Belanja.svg",
		Type: "outcome",
	},
	{
		Name: "Bensin",
		Icon: "/api/v1/storage/assets/categories/Bensin.svg",
		Type: "outcome",
	},
	{
		Name: "Bisnis",
		Icon: "/api/v1/storage/assets/categories/Bisnis.svg",
		Type: "outcome",
	},
	{
		Name: "Donasi",
		Icon: "/api/v1/storage/assets/categories/Donasi.svg",
		Type: "outcome",
	},
	{
		Name: "Edukasi",
		Icon: "/api/v1/storage/assets/categories/Edukasi.svg",
		Type: "outcome",
	},
	{
		Name: "Gaji",
		Icon: "/api/v1/storage/assets/categories/Gaji.svg",
		Type: "income",
	},
	{
		Name: "Hiburan",
		Icon: "/api/v1/storage/assets/categories/Hiburan.svg",
		Type: "outcome",
	},
	{
		Name: "Kesehatan",
		Icon: "/api/v1/storage/assets/categories/Kesehatan.svg",
		Type: "outcome",
	},
	{
		Name: "Makanan & Minuman",
		Icon: "/api/v1/storage/assets/categories/Makanan & Minuman.svg",
		Type: "outcome",
	},
	{
		Name: "Pakaian",
		Icon: "/api/v1/storage/assets/categories/Pakaian.svg",
		Type: "outcome",
	},
	{
		Name: "Peliharaan",
		Icon: "/api/v1/storage/assets/categories/Peliharaan.svg",
		Type: "outcome",
	},
	{
		Name: "Pemasukkan Lainnya",
		Icon: "/api/v1/storage/assets/categories/Pemasukkan Lainnya.svg",
		Type: "income",
	},
	{
		Name: "Pengeluaran Lainnya",
		Icon: "/api/v1/storage/assets/categories/Pengeluaran Lainnya.svg",
		Type: "outcome",
	},
	{
		Name: "Perbaikan",
		Icon: "/api/v1/storage/assets/categories/Perbaikan.svg",
		Type: "outcome",
	},
	{
		Name: "Tagihan",
		Icon: "/api/v1/storage/assets/categories/Tagihan.svg",
		Type: "outcome",
	},
	{
		Name: "Transportasi",
		Icon: "/api/v1/storage/assets/categories/Transportasi.svg",
		Type: "outcome",
	},
}

func SeederCategory(db *gorm.DB) {
	var categories []entity.Category

	db.Find(&categories)

	if len(categories) == 0 {
		db.Create(&dataCategories)
		fmt.Println("Seed OK")
	} else {
		fmt.Println("Seed Not Running")
	}

}
