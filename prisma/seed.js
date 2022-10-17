const { PrismaClient } = require('@prisma/client');
const prisma = new PrismaClient();

const categoryData = [
    {
        name: 'Belanja',
        icon: '/public/assets/categories/Belanja.svg',
        type: 'outcome',
    },
    {
        name: 'Bensin',
        icon: '/public/assets/categories/Bensin.svg',
        type: 'outcome',
    },
    {
        name: 'Donasi',
        icon: '/public/assets/categories/Donasi.svg',
        type: 'outcome',
    },
    {
        name: 'Edukasi',
        icon: '/public/assets/categories/Edukasi.svg',
        type: 'outcome',
    },
    {
        name: 'Gaji',
        icon: '/public/assets/categories/Gaji.svg',
        type: 'income',
    },
    {
        name: 'Orang Tua',
        icon: '/public/assets/categories/Bisnis.svg',
        type: 'income',
    },
    {
        name: 'Hiburan',
        icon: '/public/assets/categories/Hiburan.svg',
        type: 'outcome',
    },
    {
        name: 'Kesehatan',
        icon: '/public/assets/categories/Kesehatan.svg',
        type: 'outcome',
    },
    {
        name: 'Makanan & Minuman',
        icon: '/public/assets/categories/Makanan & Minuman.svg',
        type: 'outcome',
    },
    {
        name: 'Pakaian',
        icon: '/public/assets/categories/Pakaian.svg',
        type: 'outcome',
    },
    {
        name: 'Peliharaan',
        icon: '/public/assets/categories/Peliharaan.svg',
        type: 'outcome',
    },
    {
        name: 'Perbaikan',
        icon: '/public/assets/categories/Perbaikan.svg',
        type: 'outcome',
    },
    {
        name: 'Tagihan',
        icon: '/public/assets/categories/Tagihan.svg',
        type: 'outcome',
    },
    {
        name: 'Transportasi',
        icon: '/public/assets/categories/Transportasi.svg',
        type: 'outcome',
    },
    {
        name: 'Pemasukkan Lainnya',
        icon: '/public/assets/categories/Pemasukkan Lainnya.svg',
        type: 'income',
    },
    {
        name: 'Pengeluaran Lainnya',
        icon: '/public/assets/categories/Pengeluaran Lainnya.svg',
        type: 'outcome',
    },
];

const main = async () => {
    console.log(`🔥 Start seeding ...`);
    for (const i of categoryData) {
        const category = await prisma.categories.create({
            data: i
        })
        console.log(`Success create category ${category.id}`);
    }
    console.log(`🚀 Seeding finished ... `);
}

main()
    .then(async () => {
        await prisma.$disconnect();
    })
    .catch(async (e) => {
        console.error(e)
        await prisma.$disconnect();
        process.exit(1)
    })