package seeder

import (
	"fmt"
	"pawang-backend/entity"

	"gorm.io/gorm"
)

var dataCategories = []entity.Category{
	{
		Name: "Belanja",
		Icon: "/api/storage/assets/categories/Belanja.svg",
		Type: "outcome",
	},
	{
		Name: "Bensin",
		Icon: "/api/storage/assets/categories/Bensin.svg",
		Type: "outcome",
	},
	{
		Name: "Bisnis",
		Icon: "/api/storage/assets/categories/Bisnis.svg",
		Type: "outcome",
	},
	{
		Name: "Donasi",
		Icon: "/api/storage/assets/categories/Donasi.svg",
		Type: "outcome",
	},
	{
		Name: "Edukasi",
		Icon: "/api/storage/assets/categories/Edukasi.svg",
		Type: "outcome",
	},
	{
		Name: "Gaji",
		Icon: "/api/storage/assets/categories/Gaji.svg",
		Type: "income",
	},
	{
		Name: "Hiburan",
		Icon: "/api/storage/assets/categories/Hiburan.svg",
		Type: "outcome",
	},
	{
		Name: "Kesehatan",
		Icon: "/api/storage/assets/categories/Kesehatan.svg",
		Type: "outcome",
	},
	{
		Name: "Makanan & Minuman",
		Icon: "/api/storage/assets/categories/Makanan & Minuman.svg",
		Type: "outcome",
	},
	{
		Name: "Pakaian",
		Icon: "/api/storage/assets/categories/Pakaian.svg",
		Type: "outcome",
	},
	{
		Name: "Peliharaan",
		Icon: "/api/storage/assets/categories/Peliharaan.svg",
		Type: "outcome",
	},
	{
		Name: "Pemasukkan Lainnya",
		Icon: "/api/storage/assets/categories/Pemasukkan Lainnya.svg",
		Type: "income",
	},
	{
		Name: "Pengeluaran Lainnya",
		Icon: "/api/storage/assets/categories/Pengeluaran Lainnya.svg",
		Type: "outcome",
	},
	{
		Name: "Perbaikan",
		Icon: "/api/storage/assets/categories/Perbaikan.svg",
		Type: "outcome",
	},
	{
		Name: "Tagihan",
		Icon: "/api/storage/assets/categories/Tagihan.svg",
		Type: "outcome",
	},
	{
		Name: "Transportasi",
		Icon: "/api/storage/assets/categories/Transportasi.svg",
		Type: "outcome",
	},
}

func SeederCategory(db *gorm.DB) {
	var categories []entity.Category

	db.Find(&categories)

	if len(categories) == 0 {
		db.Create(&dataCategories)
		fmt.Println("Seed OK")
	}

}
