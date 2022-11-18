const { PrismaClient } = require("@prisma/client");
const joi = require("joi");
const eventsEmitter = require("../helpers/event");
const { multer, upload } = require("../helpers/uploadFile");

const store = upload.single("receipt_image");

const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    try {
      const { wallet } = req.query;

      if (wallet != null) {
        const checkWallet = await prisma.wallets.findFirst({
          where: {
            AND: {
              id: Number(wallet),
              user_id: req.user.id,
            },
          },
        });
        if (!checkWallet) {
          return res.status(404).json({
            message: "wallet not found!",
          });
        }

        const transactions = await prisma.transactions.findMany({
          where: {
            AND: {
              wallet_id: Number(wallet),
              user_id: req.user.id,
            },
          },
          orderBy: {
            date: "desc",
          },
          include: {
            wallet: {},
            category: {},
            subcategory: {},
          },
        });

        res.status(200).json({
          status: true,
          message: "SUCCESS_GET_DATA",
          data: transactions,
        });
      }

      let transactions = await prisma.transactions.findMany({
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
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: transactions,
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
      const transaction = await prisma.transactions.findUnique({
        where: {
          id: Number(id),
        },
        include: {
          wallet: {},
          category: {},
          subcategory: {},
        },
      });

      if (!transaction) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: transaction,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  detail: async (req, res) => {
    try {
      let totalIncome = 0;
      let totalOutcome = 0;
      let totalBalanceWallet = 0;

      const { wallet } = req.query;

      if (wallet != null) {
        const checkWallet = await prisma.wallets.findFirst({
          where: {
            AND: {
              id: Number(wallet),
              user_id: req.user.id,
            },
          },
        });
        if (!checkWallet) {
          throw { status: 404, message: "WALLET_NOT_FOUND", data: null };
        }

        const transactions = await prisma.transactions.findMany({
          where: {
            AND: {
              wallet_id: Number(wallet),
              user_id: req.user.id,
            },
          },
          orderBy: {
            date: "desc",
          },
          include: {
            wallet: {},
            category: {},
            subcategory: {},
          },
        });

        transactions.forEach((transaction) => {
          if (transaction.type === "income") {
            totalIncome += transaction.amount;
          } else {
            totalOutcome += transaction.amount;
          }
        });

        totalBalanceWallet = totalIncome - totalOutcome;

        return res.status(200).json({
          status: true,
          message: "SUCCESS_GET_DATA",
          data: {
            total_income: totalIncome,
            total_outcome: totalOutcome,
            total_balance_wallet: totalBalanceWallet,
          },
        });
      }

      const transactions = await prisma.transactions.findMany({
        where: {
          user_id: req.user.id,
        },
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
          user_id: req.user.id,
        },
      });

      wallets.forEach((item) => {
        totalBalanceWallet += item.balance;
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: {
          total_income: totalIncome,
          total_outcome: totalOutcome,
          total_balance: totalBalanceWallet,
        },
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
    const createTransactionSchema = joi
      .object({
        amount: joi.number().required().messages({
          "string.base": "Nominal hanya bisa dimasukkan angka",
          "string.empty": "Nominal tidak boleh dikosongi",
          "any.required": "Nominal wajib diisi",
        }),
        category_id: joi.number().required().messages({
          "string.base": "Kategori hanya bisa dimasukkan angka",
          "string.empty": "Kategori tidak boleh dikosongi",
          "any.required": "Kategori wajib diisi",
        }),
        wallet_id: joi.number().required().messages({
          "string.base": "Wallet hanya bisa dimasukkan angka",
          "string.empty": "Wallet tidak boleh dikosongi",
          "any.required": "Wallet wajib diisi",
        }),
        subcategory_id: joi.number().allow(null).messages({
          "string.base": "Sub Kategori hanya bisa dimasukkan angka",
          "string.empty": "Sub Kategori tidak boleh dikosongi",
          "any.required": "Sub Kategori wajib diisi",
        }),
        date: joi.date().required().messages({
          "string.base": "Tanggal hanya bisa dimasukkan tanggal",
          "string.empty": "Tanggal tidak boleh dikosongi",
          "any.required": "Tanggal wajib diisi",
        }),
      })
      .unknown(true);

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
          const {
            amount,
            category_id,
            wallet_id,
            subcategory_id,
            description,
            date,
          } = req.body;

          let newBalance = 0;

          const { error, value } = createTransactionSchema.validate(req.body, {
            abortEarly: false,
          });
          if (error) {
            let message = error.details[0].message;
            throw { status: 400, message: "INVALID_INPUT", data: message };
          }

          const checkWallet = await prisma.wallets.findFirst({
            where: {
              AND: {
                id: Number(wallet_id),
                user_id: req.user.id,
              },
            },
          });
          if (!checkWallet) {
            throw { status: 404, message: "WALLET_NOT_FOUND", data: null };
          }

          const checkCategory = await prisma.categories.findUnique({
            where: { id: Number(category_id) },
          });
          if (!checkCategory) {
            throw { status: 404, message: "CATEGORY_NOT_FOUND", data: null };
          }

          if (subcategory_id) {
            const checkSubcategory = await prisma.sub_categories.findFirst({
              where: {
                AND: {
                  id: Number(subcategory_id),
                  category_id: Number(category_id),
                  user_id: req.user.id,
                },
              },
            });
            if (!checkSubcategory) {
              throw {
                status: 404,
                message: "SUB_CATEGORY_NOT_FOUND",
                data: null,
              };
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
              user_id: req.user.id,
            },
          });

          if (checkCategory.type == "income") {
            newBalance = checkWallet.balance + parseFloat(amount);
          } else if (checkCategory.type == "outcome") {
            newBalance = checkWallet.balance - parseFloat(amount);
          } else {
            throw { status: 400, message: "INVALID_INPUT", data: null };
          }

          await prisma.wallets.update({
            where: {
              id: checkWallet.id,
            },
            data: {
              balance: newBalance,
            },
          });

          res.status(201).json({
            status: true,
            message: "SUCCESS_CREATE_DATA",
            data: newTransaction,
          });
        } catch (error) {
          return res.status(error.status || 500).json({
            status: false,
            message: error.message || "INTERNAL_SERVER_ERROR",
          });
        }
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
    const updateTransactionSchema = joi
      .object({
        amount: joi.number().required().messages({
          "string.base": "Nominal hanya bisa dimasukkan angka",
          "string.empty": "Nominal tidak boleh dikosongi",
          "any.required": "Nominal wajib diisi",
        }),
        category_id: joi.number().required().messages({
          "string.base": "Kategori hanya bisa dimasukkan angka",
          "string.empty": "Kategori tidak boleh dikosongi",
          "any.required": "Kategori wajib diisi",
        }),
        wallet_id: joi.number().required().messages({
          "string.base": "Wallet hanya bisa dimasukkan angka",
          "string.empty": "Wallet tidak boleh dikosongi",
          "any.required": "Wallet wajib diisi",
        }),
        subcategory_id: joi.number().allow(null).messages({
          "string.base": "Sub Kategori hanya bisa dimasukkan angka",
          "string.empty": "Sub Kategori tidak boleh dikosongi",
          "any.required": "Sub Kategori wajib diisi",
        }),
        date: joi.date().required().messages({
          "string.base": "Tanggal hanya bisa dimasukkan tanggal",
          "string.empty": "Tanggal tidak boleh dikosongi",
          "any.required": "Tanggal wajib diisi",
        }),
      })
      .unknown(true);

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
          const { id } = req.params;
          const {
            amount,
            category_id,
            wallet_id,
            subcategory_id,
            description,
            date,
          } = req.body;

          let newBalance = 0;

          const { error, value } = updateTransactionSchema.validate(req.body, {
            abortEarly: false,
          });
          if (error) {
            let message = error.details[0].message;
            throw { status: 400, message: "INVALID_INPUT", data: message };
          }

          const checkTransaction = await prisma.transactions.findFirst({
            where: {
              AND: {
                id: Number(id),
                user_id: req.user.id,
              },
            },
          });

          if (!checkTransaction) {
            throw { status: 404, message: "DATA_NOT_FOUND", data: null };
          }

          const checkWallet = await prisma.wallets.findFirst({
            where: {
              AND: {
                id: Number(wallet_id),
                user_id: req.user.id,
              },
            },
          });
          if (!checkWallet) {
            throw { status: 404, message: "WALLET_NOT_FOUND", data: null };
          }

          const checkCategory = await prisma.categories.findUnique({
            where: { id: Number(category_id) },
          });
          if (!checkCategory) {
            throw { status: 404, message: "CATEGORY_NOT_FOUND", data: null };
          }

          if (subcategory_id) {
            const checkSubcategory = await prisma.sub_categories.findFirst({
              where: {
                AND: {
                  id: Number(subcategory_id),
                  category_id: Number(category_id),
                  user_id: req.user.id,
                },
              },
            });
            if (!checkSubcategory) {
              throw {
                status: 404,
                message: "SUB_CATEGORY_NOT_FOUND",
                data: null,
              };
            }
          }

          if (req.file) {
            if (checkTransaction.image_url) {
              eventsEmitter.emit("deleteFile", checkTransaction.image_url);
            }
          }

          let updateTransaction = await prisma.transactions.update({
            where: {
              id: Number(id),
            },
            data: {
              amount: parseFloat(amount),
              category_id: checkCategory.id,
              wallet_id: checkWallet.id,
              subcategory_id: subcategory_id ? parseInt(subcategory_id) : null,
              type: checkCategory.type,
              description,
              date: new Date(date),
              user_id: req.user.id,
            },
          });

          if (req.file) {
            updateTransaction = await prisma.transactions.update({
              where: {
                id: updateTransaction.id,
              },
              data: {
                image_url: req.file ? req.file.path : null,
              },
            });
          }

          if (checkTransaction.type == "income") {
            newBalance =
              checkWallet.balance -
              checkTransaction.amount +
              parseFloat(amount);
          } else if (checkTransaction.type == "outcome") {
            newBalance =
              checkWallet.balance +
              checkTransaction.amount -
              parseFloat(amount);
          } else {
            throw { status: 400, message: "INVALID_INPUT", data: null };
          }

          await prisma.wallets.update({
            where: {
              id: checkWallet.id,
            },
            data: {
              balance: newBalance,
            },
          });

          res.status(200).json({
            status: true,
            message: "SUCCESS_UPDATE_DATA",
            data: updateTransaction,
          });
        } catch (error) {
          return res.status(error.status || 500).json({
            status: false,
            message: error.message || "INTERNAL_SERVER_ERROR",
          });
        }
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

      let newBalance = 0;

      const checkTransaction = await prisma.transactions.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
      });

      if (!checkTransaction) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      const checkWallet = await prisma.wallets.findFirst({
        where: {
          AND: {
            id: checkTransaction.wallet_id,
            user_id: req.user.id,
          },
        },
      });
      if (!checkWallet) {
        throw { status: 404, message: "WALLET_NOT_FOUND", data: null };
      }

      if (checkTransaction.type == "income") {
        newBalance = checkWallet.balance - checkTransaction.amount;
      } else if (checkTransaction.type == "outcome") {
        newBalance = checkWallet.balance + checkTransaction.amount;
      } else {
        throw { status: 400, message: "INVALID_INPUT", data: null };
      }

      if (checkTransaction.image_url) {
        eventsEmitter.emit("deleteFile", checkTransaction.image_url);
      }

      await prisma.wallets.update({
        where: {
          id: checkWallet.id,
        },
        data: {
          balance: newBalance,
        },
      });

      await prisma.transactions.delete({
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
