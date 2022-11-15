const { PrismaClient } = require("@prisma/client");
const joi = require("joi");
const eventsEmitter = require("../helpers/event");

const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    try {
      const wallets = await prisma.wallets.findMany({
        where: {
          user_id: req.user.id,
        },
        include: {
          transactions: {
            where: {
              user_id: req.user.id,
            },
            orderBy: {
              date: "desc",
            },
            include: {
              wallet: {},
              category: {},
              subcategory: {},
            },
          },
        },
      });

      wallets.forEach((wallet) => {
        let totalIncome = 0;
        let totalOutcome = 0;

        wallet.transactions.forEach((transaction) => {
          if (transaction.type === "income") {
            totalIncome += transaction.amount;
          } else {
            totalOutcome += transaction.amount;
          }
        });

        wallet["total_income"] = totalIncome;
        wallet["total_outcome"] = totalOutcome;
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: wallets,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  show: async (req, res) => {
    try {
      const { id } = req.params;

      const checkWallet = await prisma.wallets.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
        include: {
          transactions: {
            where: {
              user_id: req.user.id,
            },
            orderBy: {
              date: "desc",
            },
            include: {
              wallet: {},
              category: {},
              subcategory: {},
            },
          },
        },
      });

      if (!checkWallet) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      let totalIncome = 0;
      let totalOutcome = 0;

      checkWallet.transactions.forEach((transaction) => {
        if (transaction.type === "income") {
          totalIncome += transaction.amount;
        } else {
          totalOutcome += transaction.amount;
        }
      });

      checkWallet["total_income"] = totalIncome;
      checkWallet["total_outcome"] = totalOutcome;

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: checkWallet,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  create: async (req, res) => {
    const createWalletSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
        balance: joi.number().required().messages({
          "string.base": "Nominal hanya bisa dimasukkan angka",
          "string.empty": "Nominal tidak boleh dikosongi",
          "any.required": "Nominal wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { name, balance = 0 } = req.body;

      const { error, value } = createWalletSchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const newWallet = await prisma.wallets.create({
        data: {
          name,
          balance: Number(balance),
          user_id: req.user.id,
        },
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
          },
        });
      }

      res.status(201).json({
        status: true,
        message: "SUCCESS_CREATE_DATA",
        data: newWallet,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  update: async (req, res) => {
    const updateWalletSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
        balance: joi.number().required().messages({
          "string.base": "Nominal hanya bisa dimasukkan angka",
          "string.empty": "Nominal tidak boleh dikosongi",
          "any.required": "Nominal wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { id } = req.params;
      const { name, balance } = req.body;

      const { error, value } = updateWalletSchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const checkWallet = await prisma.wallets.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
      });

      if (!checkWallet) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      if (
        Number(balance) != Number(checkWallet.balance) &&
        Number(balance) != 0
      ) {
        if (checkWallet.balance < Number(balance)) {
          await prisma.transactions.create({
            data: {
              category_id: 13,
              type: "income",
              user_id: req.user.id,
              wallet_id: checkWallet.id,
              date: new Date().toISOString(),
              amount: Number(balance) - checkWallet.balance,
            },
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
            },
          });
        }
      }

      const updateWallet = await prisma.wallets.update({
        where: {
          id: Number(id),
        },
        data: {
          name,
          balance: Number(balance),
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_UPDATE_DATA",
        data: updateWallet,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  destroy: async (req, res) => {
    try {
      const { id } = req.params;

      const checkWallet = await prisma.wallets.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
        include: {
          transactions: {},
        },
      });

      if (!checkWallet) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      checkWallet.transactions.forEach((transaction) => {
        if (transaction.image_url) {
          eventsEmitter.emit("deleteFile", transaction.image_url);
        }
      });

      await prisma.wallets.delete({
        where: {
          id: Number(id),
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_DELETE_DATA",
        data: null,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
};
