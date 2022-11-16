const { PrismaClient } = require("@prisma/client");
const schedule = require("node-schedule");
const joi = require("joi");
const { sendNotification } = require("../helpers/notification");
const moment = require("moment-timezone");

const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    try {
      const reminders = await prisma.task_reminders.findMany({
        where: {
          user_id: req.user.id,
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: reminders,
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

      const reminder = await prisma.task_reminders.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
      });

      if (!reminder) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: reminder,
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
    const taskReminderSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
        date: joi.date().required().messages({
          "string.base": "Tanggal hanya bisa dimasukkan tanggal",
          "string.empty": "Tanggal tidak boleh dikosongi",
          "any.required": "Tanggal wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { name, type, date, is_active } = req.body;

      const { error, value } = taskReminderSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const typeValid = ["once", "daily", "weekly", "monthly", "yearly"];
      if (!typeValid.includes(type)) {
        throw { status: 400, message: "INVALID_INPUT", data: null };
      }

      const newReminder = await prisma.task_reminders.create({
        data: {
          name,
          type,
          date: new Date(date).toISOString(),
          is_active: Boolean(is_active),
          user_id: req.user.id,
        },
      });

      res.status(201).json({
        status: true,
        message: "SUCCESS_CREATE_DATA",
        data: newReminder,
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
    const taskReminderSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
        date: joi.date().required().messages({
          "string.base": "Tanggal hanya bisa dimasukkan tanggal",
          "string.empty": "Tanggal tidak boleh dikosongi",
          "any.required": "Tanggal wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { id } = req.params;
      const { name, type, date, is_active } = req.body;

      const { error, value } = taskReminderSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const typeValid = ["once", "daily", "weekly", "monthly", "yearly"];
      if (!typeValid.includes(type)) {
        throw { status: 400, message: "INVALID_INPUT", data: null };
      }

      const checkReminder = await prisma.task_reminders.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
      });
      if (!checkReminder) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      const updateReminder = await prisma.task_reminders.update({
        where: {
          id: Number(id),
        },
        data: {
          name,
          type,
          date: new Date(date).toISOString(),
          is_active: Boolean(is_active),
          user_id: req.user.id,
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_UPDATE_DATA",
        data: updateReminder,
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

      const checkReminder = await prisma.task_reminders.findFirst({
        where: {
          AND: {
            id: Number(id),
            user_id: req.user.id,
          },
        },
      });
      if (!checkReminder) {
        throw { status: 404, message: "DATA_NOT_FOUND", data: null };
      }

      await prisma.task_reminders.delete({
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
  scheduler: async () => {
    schedule.scheduleJob("*/3 * * * * *", async () => {
      const reminders = await prisma.task_reminders.findMany({
        include: {
          user: {
            include: {
              users_onesignals: {},
            },
          },
        },
        where: {
          is_active: true,
        },
      });

      reminders.forEach(async (item) => {
        if (moment().format("LLL") === moment(item.date).format("LLL")) {
          item.user.users_onesignals.forEach(async (user) => {
            await sendNotification({
              title: "Pengingat",
              subtitle: `Jangan Lupa ${item.name}, ${moment().format("LL")}`,
              playerId: user.onesignal_id,
            });

            console.log("send notification success to", user.onesignal_id);
          });

          const newDate = new Date(item.date);

          if (item.type == "daily") {
            newDate.setDate(newDate.getDate() + 1);
          } else if (item.type == "weekly") {
            newDate.setDate(newDate.getDate() + 7);
          } else if (item.type == "monthly") {
            newDate.setMonth(newDate.getMonth() + 1);
          } else if (item.type == "yearly") {
            newDate.setFullYear(newDate.getFullYear() + 1);
          }

          if (item.type == "once") {
            await prisma.task_reminders.update({
              where: {
                id: item.id,
              },
              data: {
                is_active: false,
              },
            });
          } else {
            await prisma.task_reminders.update({
              where: {
                id: item.id,
              },
              data: {
                is_active: true,
                date: newDate,
              },
            });
          }
        }
      });
    });
  },
};
