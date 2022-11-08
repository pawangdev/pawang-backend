const { PrismaClient } = require('@prisma/client')
const schedule = require('node-schedule');
const joi = require('joi');
const { sendNotification } = require('../helpers/notification');
const moment = require('moment-timezone');

const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    try {
      const reminders = await prisma.task_reminders.findMany({
        where: {
          user_id: req.user.id
        }
      });

      res.status(200).json({
        message: 'success retreived data!',
        data: reminders,
      });
    } catch (error) {
      res.status(500).send(error.message)
    }
  },
  show: async (req, res) => {
    try {
      const { id } = req.params;

      const reminder = await prisma.task_reminders.findUnique({
        where: {
          id
        }
      });

      if (!reminder) {
        res.status(404).json({
          message: 'failed get data!',
          data: null,
        });

        return;
      }

      res.status(200).json({
        message: 'success retreived data!',
        data: reminder,
      });
    } catch (error) {
      res.status(500).send(error.message)
    }
  },
  create: async (req, res) => {
    const taskReminderSchema = joi.object({
      name: joi.string().required().messages({
        'string.base': 'Nama hanya bisa dimasukkan text',
        'string.empty': 'Nama tidak boleh dikosongi',
        'any.required': 'Nama wajib diisi',
      }),
      date: joi.date().required().messages({
        'string.base': 'Tanggal hanya bisa dimasukkan tanggal',
        'string.empty': 'Tanggal tidak boleh dikosongi',
        'any.required': 'Tanggal wajib diisi',
      }),
    }).unknown(true);

    try {
      const {
        name,
        type,
        date,
        is_active,
      } = req.body;

      const { error, value } = taskReminderSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        res.status(422).json({ message: "Format Tidak Valid", data: message });

        return;
      }

      const typeValid = ['once', 'daily', 'weekly', 'monthly', 'yearly'];
      if (!typeValid.includes(type)) {
        res.status(422).json({ message: "Format Tidak Valid", data: 'Type Tidak Valid' });

        return;
      }

      const newReminder = await prisma.task_reminders.create({
        data: {
          name,
          type,
          date: new Date(date).toISOString(),
          is_active: Boolean(is_active),
          user_id: req.user.id
        }
      });

      res.status(201).json({
        message: 'success create data!',
        data: newReminder,
      });
    } catch (error) {
      res.status(500).send(error.message)
    }
  },
  update: async (req, res) => {
    const taskReminderSchema = joi.object({
      name: joi.string().required().messages({
        'string.base': 'Nama hanya bisa dimasukkan text',
        'string.empty': 'Nama tidak boleh dikosongi',
        'any.required': 'Nama wajib diisi',
      }),
      date: joi.date().required().messages({
        'string.base': 'Tanggal hanya bisa dimasukkan tanggal',
        'string.empty': 'Tanggal tidak boleh dikosongi',
        'any.required': 'Tanggal wajib diisi',
      }),
    }).unknown(true);

    try {
      const { id } = req.params;
      const {
        name,
        type,
        date,
        is_active,
      } = req.body;

      const { error, value } = taskReminderSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        res.status(422).json({ message: "Format Tidak Valid", data: message });

        return;
      }

      const checkReminder = await prisma.task_reminders.findUnique({
        where: {
          id: Number(id)
        }
      });
      if (!checkReminder) {
        res.status(404).json({
          message: 'failed get data!',
          data: null,
        });

        return;
      }

      const typeValid = ['once', 'daily', 'weekly', 'monthly', 'yearly'];
      if (!typeValid.includes(type)) {
        res.status(422).json({ message: "Format Tidak Valid", data: 'Type Tidak Valid' });

        return;
      }

      const updateReminder = await prisma.task_reminders.update({
        where: {
          id: Number(id)
        },
        data: {
          name,
          type,
          date: new Date(date).toISOString(),
          is_active: Boolean(is_active),
          user_id: req.user.id
        }
      });

      res.status(200).json({
        message: 'success update data!',
        data: updateReminder,
      });
    } catch (error) {
      res.status(500).send(error.message)
    }
  },
  destroy: async (req, res) => {
    try {
      const { id } = req.params;

      const checkReminder = await prisma.task_reminders.findUnique({
        where: {
          id: Number(id)
        }
      });
      if (!checkReminder) {
        res.status(404).json({
          message: 'failed get data!',
          data: null,
        });

        return;
      }

      await prisma.task_reminders.delete({
        where: {
          id: Number(id)
        }
      });

      res.status(200).json({
        message: 'success delete data!',
        data: null,
      });
    } catch (error) {
      res.status(500).send(error.message)
    }
  },
  scheduler: async () => {
    schedule.scheduleJob('*/5 * * * * *', async () => {
      const reminders = await prisma.task_reminders.findMany({
        include: {
          user: {}
        },
        where: {
          is_active: true
        }
      });

      reminders.forEach(async (item) => {
        if (moment().format('HH:mm') === moment(item.date).format("HH:mm")) {
          await sendNotification({ title: "Pengingat", subtitle: `Jangan Lupa ${item.name}, ${moment().format('LL')}`, playerId: item.user.onesignal_id });
          const newDate = new Date(item.date);

          if (item.type === 'daily') {
            newDate.setDate(newDate.getDate() + 1);
          } else if (item.type === 'weekly') {
            newDate.setDate(newDate.getDate() + 7);
          } else if (item.type === 'monthly') {
            newDate.setMonth(newDate.getMonth() + 1);
          } else if (item.type === 'yearly') {
            newDate.setFullYear(newDate.getFullYear() + 1);
          } else if (item.type === 'once') {
            await prisma.task_reminders.delete({
              where: {
                id: item.id
              },
              data: {
                is_active: false
              }
            });
          }

          if (item.type !== 'once') {
            await prisma.task_reminders.update({
              where: {
                id: item.id
              },
              data: {
                is_active: true,
                date: newDate,
              }
            });
          }
        }
      })
    });
  }
}
