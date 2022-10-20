const { PrismaClient } = require('@prisma/client')
const schedule = require('node-schedule');
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
        try {
            const {
                name,
                type,
                date,
                is_active,
            } = req.body;

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
        try {
            const { id } = req.params;
            const {
                name,
                type,
                date,
                is_active,
            } = req.body;

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
                    await sendNotification({ title: item.name, subtitle: moment().format('LL'), playerId: item.user.onesignal_id });

                    await prisma.task_reminders.update({
                        where: {
                            id: item.id
                        },
                        data: {
                            is_active: false,
                        }
                    });
                }
            })
        });
    }
}