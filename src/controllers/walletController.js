const { PrismaClient } = require('@prisma/client');
const joi = require('joi');

const prisma = new PrismaClient();

module.exports = {
    index: async (req, res) => {
        try {
            const wallets = await prisma.wallets.findMany({
                where: {
                    user_id: req.user.id
                },
                include: {
                    transactions: {
                        where: {
                            user_id: req.user.id
                        },
                    }
                },
            });

            res.status(200).json({
                message: 'success retreived data!',
                data: wallets,
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    show: async (req, res) => {
        try {
            const { id } = req.params;

            const checkWallet = await prisma.wallets.findUnique({
                where: {
                    id: Number(id),
                },
                include: {
                    transactions: {
                        where: {
                            user_id: req.user.id
                        },
                    }
                }
            });

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

            const wallet = await prisma.wallets.findUnique({
                where: {
                    id: Number(id),
                }
            });
            res.status(200).json({
                message: 'success retreived data!',
                data: wallet
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    create: async (req, res) => {
        const createWalletSchema = joi.object({
            name: joi.string().required().messages({
                'string.base': 'Nama hanya bisa dimasukkan text',
                'string.empty': 'Nama tidak boleh dikosongi',
                'any.required': 'Nama wajib diisi',
            }),
            balance: joi.number().required().messages({
                'string.base': 'Nominal hanya bisa dimasukkan angka',
                'string.empty': 'Nominal tidak boleh dikosongi',
                'any.required': 'Nominal wajib diisi',
            }),
        }).unknown(true);

        try {
            const { name, balance = 0 } = req.body;

            const { error, value } = createWalletSchema.validate(req.body, {
                abortEarly: false,
            });
            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const newWallet = await prisma.wallets.create({
                data: {
                    name, balance: Number(balance), user_id: req.user.id
                }
            });

            if (Number(balance) != 0) {
                await prisma.transactions.create({
                    data: {
                        category_id: 13,
                        type: "income",
                        user_id: req.user.id,
                        wallet_id: newWallet.id,
                        date: new Date().toISOString(),
                        amount: Number(balance),
                    }
                });
            }

            res.status(201).json({
                message: 'success create wallet!',
                data: newWallet,
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    update: async (req, res) => {
        const updateWalletSchema = joi.object({
            name: joi.string().required().messages({
                'string.base': 'Nama hanya bisa dimasukkan text',
                'string.empty': 'Nama tidak boleh dikosongi',
                'any.required': 'Nama wajib diisi',
            }),
            balance: joi.number().required().messages({
                'string.base': 'Nominal hanya bisa dimasukkan angka',
                'string.empty': 'Nominal tidak boleh dikosongi',
                'any.required': 'Nominal wajib diisi',
            }),
        }).unknown(true);

        try {
            const { id } = req.params;
            const { name, balance } = req.body;

            const { error, value } = updateWalletSchema.validate(req.body, {
                abortEarly: false,
            });
            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const checkWallet = await prisma.wallets.findUnique({
                where: {
                    id: Number(id),
                }
            });

            if (!checkWallet) {
                res.status(404).json({
                    message: 'failed retreived data!',
                    data: null
                });

                return;
            } else {
                if (checkWallet.user_id != req.user.id) {
                    res.status(404).json({
                        message: 'failed retreived data!',
                        data: null
                    });

                    return;
                }
            }

            if (Number(balance) != Number(checkWallet.balance) && Number(balance) != 0) {
                if (checkWallet.balance < Number(balance)) {
                    await prisma.transactions.create({
                        data: {
                            category_id: 13,
                            type: "income",
                            user_id: req.user.id,
                            wallet_id: checkWallet.id,
                            date: new Date().toISOString(),
                            amount: Number(balance) - checkWallet.balance,
                        }
                    });
                } else {
                    await prisma.transactions.create({
                        data: {
                            category_id: 10,
                            type: "outcome",
                            user_id: req.user.id,
                            wallet_id: checkWallet.id,
                            date: new Date().toISOString(),
                            amount: checkWallet.balance - Number(balance),
                        }
                    });
                }
            }



            const updateWallet = await prisma.wallets.update({
                where: {
                    id: Number(id),
                },
                data: {
                    name, balance: Number(balance),
                }
            });

            res.status(200).json({
                message: 'success update wallet!',
                data: updateWallet,
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    destroy: async (req, res) => {
        try {
            const { id } = req.params;

            const checkWallet = await prisma.wallets.findUnique({
                where: {
                    id: Number(id),
                }
            });

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

            await prisma.wallets.delete({
                where: {
                    id: Number(id),
                }
            });

            res.status(200).json({
                message: 'success delete wallet!',
                data: null
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
}