const { PrismaClient } = require('@prisma/client');
const e = require('express');
const prisma = new PrismaClient();

const categoryData = [
    {
        name: 'Belanja',
        icon: '/public/assets/categories/belanja.svg',
        type: 'outcome',
    },
    {
        name: 'Makan & Minum',
        icon: '/public/assets/categories/makan-minum.svg',
        type: 'outcome',
    },
    {
        name: 'Transportasi',
        icon: '/public/assets/categories/transportasi.svg',
        type: 'outcome',
    },
    {
        name: 'Pakaian',
        icon: '/public/assets/categories/pakaian.svg',
        type: 'outcome',
    },
    {
        name: 'Pakaian',
        icon: '/public/assets/categories/pakaian.svg',
        type: 'outcome',
    },
    {
        name: 'Pendidikan',
        icon: '/public/assets/categories/pendidikan.svg',
        type: 'outcome',
    },
    {
        name: 'Hiburan',
        icon: '/public/assets/categories/hiburan.svg',
        type: 'outcome',
    },
    {
        name: 'Elektronik',
        icon: '/public/assets/categories/elektronik.svg',
        type: 'outcome',
    },
    {
        name: 'Kesehatan',
        icon: '/public/assets/categories/kesehatan.svg',
        type: 'outcome',
    },
    {
        name: 'Asuransi',
        icon: '/public/assets/categories/asuransi.svg',
        type: 'outcome',
    },
    {
        name: 'Pengeluaran Lainnya',
        icon: '/public/assets/categories/pengeluaran-lainnya.svg',
        type: 'outcome',
    },
    {
        name: 'Gaji',
        icon: '/public/assets/categories/gaji.svg',
        type: 'income',
    },
    {
        name: 'Orang Tua',
        icon: '/public/assets/categories/orang-tua.svg',
        type: 'income',
    },
    {
        name: 'Pemasukkan Lainnya',
        icon: '/public/assets/categories/pemasukan-lainnya.svg',
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