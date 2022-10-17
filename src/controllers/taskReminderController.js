const { PrismaClient } = require('@prisma/client')
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
            const { name, type, days, is_active } = req.body;

            const newReminder = await prisma.task_reminders.create({
                data: {
                    name, type, days, is_active, user_id: req.user.id
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
            const { name, type, days, is_active } = req.body;

            const checkReminder = await prisma.task_reminders.findUnique({
                where: {
                    id
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
                    id
                },
                data: {
                    name, type, days, is_active, user_id: req.user.id
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
                    id
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
                    id
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
}