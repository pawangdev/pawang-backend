const { PrismaClient } = require('@prisma/client');
const joi = require('joi');
const fileHelper = require('../helpers/file')
const { multer, upload } = require('../helpers/uploadFile');

const store = upload.single('receipt_image');

const prisma = new PrismaClient();

module.exports = {
    index: async (req, res) => {
        try {
            let transactions = await prisma.transactions.findMany({
                where: {
                    user_id: req.user.id
                },
                orderBy: {
                    date: 'desc'
                },
                include: {
                    wallet: {},
                    category: {},
                    subcategory: {},
                }
            });

            res.status(200).json({
                message: 'success retreived data!',
                data: transactions,
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    show: async (req, res) => {
        try {
            const { id } = req.params;
            const transaction = await prisma.transactions.findUnique({
                where: {
                    id: Number(id)
                },
                include: {
                    wallet: {},
                    category: {},
                    subcategory: {},
                }
            });

            if (!transaction) {
                res.status(404).json({
                    message: 'failed retreived data!',
                    data: null
                });

                return;
            }

            res.status(200).json({
                message: 'success retreived data!',
                data: transaction,
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    detail: async (req, res) => {
        try {
            let totalIncome = 0;
            let totalOutcome = 0;
            let totalBalanceWallet = 0;
            const transactions = await prisma.transactions.findMany({
                where: {
                    user_id: req.user.id
                }
            });

            transactions.forEach((item) => {
                if (item.type == "income") {
                    totalIncome += item.amount;
                } else {
                    totalOutcome += item.amount;
                }
            });

            const wallets = await prisma.wallets.findMany({
                where: {
                    user_id: req.user.id
                }
            });

            wallets.forEach((item) => {
                totalBalanceWallet += item.balance
            });

            res.status(200).json({
                message: 'success retreived data!',
                data: {
                    total_income: totalIncome,
                    total_outcome: totalOutcome,
                    total_balance: totalBalanceWallet,
                },
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    create: async (req, res) => {
        const createTransactionSchema = joi.object({
            amount: joi.number().required().messages({
                'string.base': 'Nominal hanya bisa dimasukkan angka',
                'string.empty': 'Nominal tidak boleh dikosongi',
                'any.required': 'Nominal wajib diisi',
            }),
            category_id: joi.number().required().messages({
                'string.base': 'Kategori hanya bisa dimasukkan angka',
                'string.empty': 'Kategori tidak boleh dikosongi',
                'any.required': 'Kategori wajib diisi',
            }),
            wallet_id: joi.number().required().messages({
                'string.base': 'Wallet hanya bisa dimasukkan angka',
                'string.empty': 'Wallet tidak boleh dikosongi',
                'any.required': 'Wallet wajib diisi',
            }),
            subcategory_id: joi.number().allow(null).messages({
                'string.base': 'Sub Kategori hanya bisa dimasukkan angka',
                'string.empty': 'Sub Kategori tidak boleh dikosongi',
                'any.required': 'Sub Kategori wajib diisi',
            }),
            date: joi.date().required().messages({
                'string.base': 'Tanggal hanya bisa dimasukkan tanggal',
                'string.empty': 'Tanggal tidak boleh dikosongi',
                'any.required': 'Tanggal wajib diisi',
            }),
        }).unknown(true);

        try {
            store(req, res, async (err) => {
                if (err instanceof multer.MulterError) {
                    return res.status(400).send({
                        error: "maximum file size is 2MB",
                    });
                } else if (req.fileValidationError) {
                    return res.status(400).send({
                        error: req.fileValidationError,
                    });
                } else if (err) {
                    return res.status(400).send({
                        error: err,
                    });
                }

                try {
                    const { amount, category_id, wallet_id, subcategory_id, description, date } = req.body;

                    let newBalance = 0;

                    const { error, value } = createTransactionSchema.validate(req.body, {
                        abortEarly: false,
                    });
                    if (error) {
                        let message = error.details[0].message;
                        res.status(422).json({ message: "Format Tidak Valid", data: message });

                        return;
                    }

                    const checkWallet = await prisma.wallets.findUnique({ where: { id: Number(wallet_id) } });
                    if (!checkWallet) {
                        res.status(404).json({
                            message: 'failed retreived wallet!',
                            data: null
                        });

                        return;
                    } else {
                        if (checkWallet.user_id != req.user.id) {
                            res.status(404).json({
                                message: 'failed retreived wallet!',
                                data: null
                            });

                            return;
                        }
                    }

                    const checkCategory = await prisma.categories.findUnique({ where: { id: Number(category_id) } });
                    if (!checkCategory) {
                        res.status(404).json({
                            message: 'failed retreived category!',
                            data: null
                        });

                        return;
                    }

                    if (subcategory_id) {
                        const checkSubcategory = await prisma.sub_categories.findUnique({ where: { id: Number(subcategory_id) } });
                        if (!checkSubcategory) {
                            res.status(404).json({
                                message: 'failed retreived subcategory!',
                                data: null
                            });

                            return;
                        }
                    }

                    const newTransaction = await prisma.transactions.create({
                        data: {
                            amount: parseFloat(amount),
                            category_id: checkCategory.id,
                            wallet_id: checkWallet.id,
                            subcategory_id: subcategory_id ? Number(subcategory_id) : null,
                            type: checkCategory.type,
                            description,
                            date: new Date(date),
                            image_url: req.file ? req.file.path : null,
                            user_id: req.user.id
                        }
                    });


                    if (checkCategory.type == "income") {
                        newBalance = checkWallet.balance + parseFloat(amount);
                    } else if (checkCategory.type == "outcome") {
                        newBalance = checkWallet.balance - parseFloat(amount);
                    } else {
                        res.status(422).json({
                            message: 'type not valid',
                            data: null
                        });

                        return;
                    }

                    await prisma.wallets.update({
                        where: {
                            id: checkWallet.id
                        },
                        data: {
                            balance: newBalance,
                        }
                    });

                    res.status(201).json({
                        message: 'success create data!',
                        data: newTransaction
                    });
                } catch (error) {
                    res.status(500).send(error.message);
                }
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    update: async (req, res, next) => {
        const updateTransactionSchema = joi.object({
            amount: joi.number().required().messages({
                'string.base': 'Nominal hanya bisa dimasukkan angka',
                'string.empty': 'Nominal tidak boleh dikosongi',
                'any.required': 'Nominal wajib diisi',
            }),
            category_id: joi.number().required().messages({
                'string.base': 'Kategori hanya bisa dimasukkan angka',
                'string.empty': 'Kategori tidak boleh dikosongi',
                'any.required': 'Kategori wajib diisi',
            }),
            wallet_id: joi.number().required().messages({
                'string.base': 'Wallet hanya bisa dimasukkan angka',
                'string.empty': 'Wallet tidak boleh dikosongi',
                'any.required': 'Wallet wajib diisi',
            }),
            subcategory_id: joi.number().allow(null).messages({
                'string.base': 'Sub Kategori hanya bisa dimasukkan angka',
                'string.empty': 'Sub Kategori tidak boleh dikosongi',
                'any.required': 'Sub Kategori wajib diisi',
            }),
            date: joi.date().required().messages({
                'string.base': 'Tanggal hanya bisa dimasukkan tanggal',
                'string.empty': 'Tanggal tidak boleh dikosongi',
                'any.required': 'Tanggal wajib diisi',
            }),
        }).unknown(true);

        try {
            store(req, res, async (err) => {
                if (err instanceof multer.MulterError) {
                    return res.status(400).send({
                        error: "maximum file size is 2MB",
                    });
                } else if (req.fileValidationError) {
                    return res.status(400).send({
                        error: req.fileValidationError,
                    });
                } else if (err) {
                    return res.status(400).send({
                        error: err,
                    });
                }

                console.log(req.file);
                try {
                    const { id } = req.params;
                    const { amount, category_id, wallet_id, subcategory_id, description, date } = req.body;

                    let newBalance = 0;

                    const { error, value } = updateTransactionSchema.validate(req.body, {
                        abortEarly: false,
                    });
                    if (error) {
                        let message = error.details[0].message;
                        res.status(422).json({ message: "Format Tidak Valid", data: message });

                        return;
                    }

                    const checkTransaction = await prisma.transactions.findUnique({
                        where: {
                            id: Number(id)
                        }
                    });

                    if (!checkTransaction) {
                        res.status(404).json({
                            message: 'failed retreived transaction!',
                            data: null
                        });

                        return;
                    } else {
                        if (checkTransaction.user_id != req.user.id) {
                            res.status(404).json({
                                message: 'failed retreived transaction!',
                                data: null
                            });

                            return;
                        }
                    }

                    const checkWallet = await prisma.wallets.findUnique({ where: { id: Number(wallet_id) } });
                    if (!checkWallet) {
                        res.status(404).json({
                            message: 'failed retreived wallet!',
                            data: null
                        });

                        return;
                    } else {
                        if (checkWallet.user_id != req.user.id) {
                            res.status(404).json({
                                message: 'failed retreived wallet!',
                                data: null
                            });

                            return;
                        }
                    }

                    const checkCategory = await prisma.categories.findUnique({ where: { id: Number(category_id) } });
                    if (!checkCategory) {
                        res.status(404).json({
                            message: 'failed retreived category!',
                            data: null
                        });

                        return;
                    }

                    if (subcategory_id) {
                        const checkSubcategory = await prisma.sub_categories.findUnique({ where: { id: Number(subcategory_id) } });
                        if (!checkSubcategory) {
                            res.status(404).json({
                                message: 'failed retreived subcategory!',
                                data: null
                            });

                            return;
                        }
                    }

                    if (req.file) {
                        fileHelper.deleteFile(checkTransaction.image_url);
                    }

                    let updateTransaction = await prisma.transactions.update({
                        where: {
                            id: Number(id)
                        },
                        data: {
                            amount: parseFloat(amount),
                            category_id: checkCategory.id,
                            wallet_id: checkWallet.id,
                            subcategory_id: subcategory_id ? parseInt(subcategory_id) : null,
                            type: checkCategory.type,
                            description,
                            date: new Date(date),
                            user_id: req.user.id
                        }
                    });

                    if (req.file) {
                        updateTransaction = await prisma.transactions.update({
                            where: {
                                id: updateTransaction.id
                            },
                            data: {
                                image_url: req.file ? req.file.path : null,
                            }
                        });
                    }

                    if (checkTransaction.type == "income") {
                        newBalance = (checkWallet.balance - checkTransaction.amount) + parseFloat(amount);
                    } else if (checkTransaction.type == "outcome") {
                        newBalance = (checkWallet.balance + checkTransaction.amount) - parseFloat(amount);
                    } else {
                        res.status(422).json({
                            message: 'input type not valid!',
                            data: null
                        });

                        return;
                    }

                    await prisma.wallets.update({
                        where: {
                            id: checkWallet.id
                        },
                        data: {
                            balance: newBalance,
                        }
                    });

                    res.status(200).json({
                        message: 'success update data',
                        data: updateTransaction
                    });
                } catch (error) {
                    res.status(500).send(error.message);
                }
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    destroy: async (req, res) => {
        try {
            const { id } = req.params;

            let newBalance = 0;

            const checkTransaction = await prisma.transactions.findUnique({
                where: {
                    id: Number(id)
                }
            });

            if (!checkTransaction) {
                res.status(404).json({
                    message: 'failed retreived transaction!',
                    data: null
                });

                return;
            } else {
                if (checkTransaction.user_id != req.user.id) {
                    res.status(404).json({
                        message: 'failed retreived transaction!',
                        data: null
                    });

                    return;
                }
            }

            const checkWallet = await prisma.wallets.findUnique({ where: { id: checkTransaction.wallet_id } });
            if (!checkWallet) {
                res.status(404).json({
                    message: 'failed retreived wallet!',
                    data: null
                });

                return;
            } else {
                if (checkWallet.user_id != req.user.id) {
                    res.status(404).json({
                        message: 'failed retreived wallet!',
                        data: null
                    });

                    return;
                }
            }

            if (checkTransaction.type == "income") {
                newBalance = (checkWallet.balance - checkTransaction.amount);
            } else if (checkTransaction.type == "outcome") {
                newBalance = (checkWallet.balance + checkTransaction.amount);
            } else {
                res.status(422).json({
                    message: 'type not valid',
                    data: null
                });

                return;
            }

            await prisma.wallets.update({
                where: {
                    id: checkWallet.id
                },
                data: {
                    balance: newBalance,
                }
            });

            await prisma.transactions.delete({
                where: {
                    id: Number(id)
                }
            });

            res.status(200).json({
                message: 'success delete data',
                data: null
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
}