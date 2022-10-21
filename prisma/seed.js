const { PrismaClient } = require('@prisma/client');
const e = require('express');
const prisma = new PrismaClient();

const categoryData = [
    {
        name: 'Belanja',
        icon: '/public/assets/categories/belanja.png',
        type: 'outcome',
    },
    {
        name: 'Makan & Minum',
        icon: '/public/assets/categories/makan-minum.png',
        type: 'outcome',
    },
    {
        name: 'Transportasi',
        icon: '/public/assets/categories/transportasi.png',
        type: 'outcome',
    },
    {
        name: 'Pakaian',
        icon: '/public/assets/categories/pakaian.png',
        type: 'outcome',
    },
    {
        name: 'Pakaian',
        icon: '/public/assets/categories/pakaian.png',
        type: 'outcome',
    },
    {
        name: 'Pendidikan',
        icon: '/public/assets/categories/pendidikan.png',
        type: 'outcome',
    },
    {
        name: 'Hiburan',
        icon: '/public/assets/categories/hiburan.png',
        type: 'outcome',
    },
    {
        name: 'Elektronik',
        icon: '/public/assets/categories/elektronik.png',
        type: 'outcome',
    },
    {
        name: 'Kesehatan',
        icon: '/public/assets/categories/kesehatan.png',
        type: 'outcome',
    },
    {
        name: 'Asuransi',
        icon: '/public/assets/categories/asuransi.png',
        type: 'outcome',
    },
    {
        name: 'Pengeluaran Lainnya',
        icon: '/public/assets/categories/pengeluaran-lainnya.png',
        type: 'outcome',
    },
    {
        name: 'Gaji',
        icon: '/public/assets/categories/gaji.png',
        type: 'income',
    },
    {
        name: 'Orang Tua',
        icon: '/public/assets/categories/orang-tua.png',
        type: 'income',
    },
    {
        name: 'Pemasukkan Lainnya',
        icon: '/public/assets/categories/pemasukan-lainnya.png',
        type: 'income',
    },
];

const main = async () => {
    console.log(`🔥 Start seeding ...`);
    for (const i of categoryData) {
        const category = await prisma.categories.upsert({
            where: {
                name: i.name
            },
            create: {

                name: i.name,
                icon: i.icon,
                type: i.type,
            },
            update: {
                name: i.name,
                icon: i.icon,
                type: i.type,
            }
        })
        console.log(`Success create category ${category.name}`);
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