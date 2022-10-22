const { PrismaClient } = require('@prisma/client');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const crypto = require('crypto');
const joi = require('joi');
const emailService = require('../helpers/mail');
const { sendNotification, addPlayer } = require('../helpers/notification');
const prisma = new PrismaClient();

module.exports = {
    register: async (req, res) => {
        const registerSchema = joi.object({
            name: joi.string().required().messages({
                'string.base': 'Nama hanya bisa dimasukkan text',
                'string.empty': 'Nama tidak boleh dikosongi',
                'any.required': 'Nama wajib diisi',
            }),
            email: joi.string().email().required().messages({
                'string.base': 'Email hanya bisa dimasukkan email',
                'string.email': 'Email hanya bisa dimasukkan email',
                'string.empty': 'Email tidak boleh dikosongi',
                'any.required': 'Email wajib diisi',
            }),
            password: joi.string().min(8).required().messages({
                'string.empty': 'Password tidak boleh dikosongi',
                'string.min': 'Password tidak boleh kurang dari 8 karakter',
                'any.required': 'Password wajib diisi',
            }),
            gender: joi.string().valid('male', 'female').messages({
                'string.base': 'Gender hanya bisa dimasukkan male / female',
                'string.empty': 'Gender tidak boleh dikosongi',
                'any.valid': 'Gender hanya bisa dimasukkan male / female',
                'any.required': 'Gender wajib diisi',
            }),
        }).unknown(true);

        try {
            const { name, email, password, gender, phone, onesignal_id } = req.body;

            const { error, value } = registerSchema.validate(req.body, {
                abortEarly: false,
            });

            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            // Check Email Duplicated
            const checkEmail = await prisma.users.findUnique({
                where: {
                    email
                }
            });

            if (checkEmail) {
                res.status(422).json({
                    message: 'Email telah terdaftar',
                    data: null
                });

                return;
            }


            const hashedPassword = await bcrypt.hash(password, 10);
            const newUser = await prisma.users.create({
                data: {
                    name, email, password: hashedPassword, gender, phone, onesignal_id, wallets: {
                        create: {
                            name: "Dompet",
                            balance: 0,
                        }
                    }
                },
            });

            if (newUser) {
                // Onesignal Update
                await prisma.users.update({
                    where: {
                        id: newUser.id
                    },
                    data: {
                        onesignal_id: onesignal_id
                    }
                });

                // Add Player OneSignal
                await addPlayer({ onesignal_id });


                delete newUser.password;
                delete newUser.google_id;
                delete newUser.onesignal_id;

                const payload = { id: newUser.id };
                const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY)
                res.status(201).json({
                    message: "success created user!",
                    data: {
                        user: newUser,
                        access_token: accessToken
                    }
                })

                return;
            }
        } catch (error) {
            res.status(500).send(error)
        }
    },
    login: async (req, res) => {
        const loginSchema = joi.object({
            email: joi.string().email().required().messages({
                'string.base': 'Email hanya bisa dimasukkan email',
                'string.empty': 'Email tidak boleh dikosongi',
                'any.required': 'Email wajib diisi',
            }),
            password: joi.string().min(8).required().messages({
                'string.empty': 'Password tidak boleh dikosongi',
                'string.min': 'Password tidak boleh kurang dari 8 karakter',
                'any.required': 'Password wajib diisi',
            }),
        }).unknown(true);

        try {
            const { email, password, onesignal_id } = req.body;

            const { error, value } = loginSchema.validate(req.body, {
                abortEarly: false,
            });

            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const user = await prisma.users.findUnique({
                where: {
                    email
                }
            })

            if (user) {
                try {
                    if (await bcrypt.compare(password, user.password)) {
                        // Onesignal Update
                        await prisma.users.update({
                            where: {
                                id: user.id
                            },
                            data: {
                                onesignal_id: onesignal_id
                            }
                        });

                        // Add Player OneSignal
                        await addPlayer({ onesignal_id });

                        delete user.password;
                        delete user.google_id;
                        delete user.onesignal_id;

                        const payload = { id: user.id };
                        const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY);
                        res.status(200).json({
                            message: 'success login user!',
                            data: {
                                user,
                                access_token: accessToken
                            }
                        });

                        return;
                    } else {
                        res.status(401).json({
                            message: "Password invalid!",
                            data: null,
                        });

                        return;
                    }
                } catch (error) {
                    res.status(500).send(error.message);
                }
            } else {
                res.status(404).json({
                    message: 'Email not found',
                    data: null,
                });

                return;
            }
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    profile: async (req, res) => {
        try {
            const user_id = req.user.id;

            const user = await prisma.users.findUnique({
                where: {
                    id: user_id
                }
            });

            delete user.password;
            delete user.google_id;
            delete user.onesignal_id;

            res.status(200).json({
                message: 'Berhasil mengambil profile user',
                data: {
                    user,
                }
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    changeProfile: async (req, res) => {
        const changeProfileSchema = joi.object({
            name: joi.string().required().messages({
                'string.base': 'Nama hanya bisa dimasukkan text',
                'string.empty': 'Nama tidak boleh dikosongi',
                'any.required': 'Nama wajib diisi',
            }),
            gender: joi.string().valid('male', 'female').messages({
                'string.base': 'Gender hanya bisa dimasukkan male / female',
                'string.empty': 'Gender tidak boleh dikosongi',
                'any.valid': 'Gender hanya bisa dimasukkan male / female',
                'any.required': 'Gender wajib diisi',
            }),
        }).unknown(true);

        try {
            const user_id = req.user.id;
            const { name, gender, phone } = req.body;

            const { error, value } = changeProfileSchema.validate(req.body, {
                abortEarly: false,
            });

            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const updateUser = await prisma.users.update({
                where: {
                    id: user_id
                },
                data: {
                    name, gender, phone
                }
            });

            delete updateUser.password;
            delete updateUser.google_id;
            delete updateUser.onesignal_id;

            res.status(200).json({
                message: 'berhasil memperbarui profile',
                data: {
                    updateUser,
                }
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    changePassword: async (req, res) => {
        const changePasswordSchema = joi.object({
            old_password: joi.string().min(8).required().messages({
                'string.empty': 'Password lama tidak boleh dikosongi',
                'string.min': 'Password lama tidak boleh kurang dari 8 karakter',
                'any.required': 'Password lama wajib diisi',
            }),
            new_password: joi.string().min(8).required().messages({
                'string.empty': 'Password baru tidak boleh dikosongi',
                'string.min': 'Password baru tidak boleh kurang dari 8 karakter',
                'any.required': 'Password baru wajib diisi',
            }),
            new_password_confirm: joi.ref('new_password').messages({
                'string.empty': 'Password baru konfirmasi tidak boleh dikosongi',
                'string.min': 'Password baru konfirmasi tidak boleh kurang dari 8 karakter',
                'any.required': 'Password baru konfirmasi wajib diisi',
                'any.ref': 'Password baru konfirmasi tidak sesuai dengan password baru',
            }),
        });

        try {
            const user_id = req.user.id;
            const { old_password, new_password, new_password_confirm } = req.body;

            const { error, value } = changePasswordSchema.validate(req.body, {
                abortEarly: false,
            });

            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const checkUser = await prisma.users.findUnique({ where: { id: user_id } });

            if (checkUser) {
                try {
                    if (await bcrypt.compare(old_password, checkUser.password)) {
                        if (new_password !== new_password_confirm) {
                            res.status(422).json({
                                message: 'password baru dan password konfirmasi tidak sesuai',
                                data: null
                            });

                            return;
                        }

                        const newPasswordHashed = await bcrypt.hash(new_password, 10);

                        const updateUser = await prisma.users.update({
                            where: {
                                id: user_id
                            },
                            data: {
                                password: newPasswordHashed
                            }
                        });

                        delete updateUser.password;

                        res.status(200).json({
                            message: 'berhasil memperbarui password',
                            data: {
                                updateUser,
                            }
                        });

                        return;
                    } else {
                        res.status(401).json({
                            message: 'password lama tidak cocok',
                            data: null
                        });

                        return;
                    }
                } catch (error) {
                    res.status(500).send(error.message);
                }
            }
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    requestResetPasswordToken: async (req, res) => {
        const requestResetPasswordTokenSchema = joi.object({
            email: joi.string().email().required(),
        });

        try {
            const { email } = req.body;

            const { error, value } = requestResetPasswordTokenSchema.validate(req.body, {
                abortEarly: false,
            });

            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const checkEmail = await prisma.users.findUnique({
                where: {
                    email
                }
            });

            if (!checkEmail) {
                res.status(404).json({
                    message: 'Email tidak terdaftar',
                    data: null
                });

                return;
            }

            const code = await crypto.randomBytes(3).toString('hex');
            const resetToken = parseInt(code.toString('hex'), 16).toString().substring(0, 6);
            const hashToken = await bcrypt.hashSync(resetToken, 8);
            const dateNow = new Date();
            dateNow.setMinutes(dateNow.getMinutes() + 10);
            const expiredToken = new Date(dateNow).toISOString();

            await prisma.user_reset_passwords.upsert({
                create: {
                    email, token: hashToken, expired_at: expiredToken
                },
                update: {
                    email, token: hashToken, expired_at: expiredToken
                },
                where: {
                    email
                }
            });

            await emailService.sendMail({
                from: 'Pawang <teampawang.dev@gmail.com>',
                to: email,
                subject: "Kode Lupa Kata Sandi",
                text: `Gunakan kode ini untuk mengatur ulang kata sandi akun Anda: ${resetToken}. Kode hanya berlaku 10 menit.`,
                html: `<p>Gunakan kode ini untuk mengatur ulang kata sandi akun Anda: <strong>${resetToken}</strong>. Kode hanya berlaku 10 menit.</p>`,
            });

            res.status(200).json({
                message: 'berhasil mengirim token untuk reset password, silahkan cek email anda',
                data: null
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    verifyResetPasswordToken: async (req, res) => {
        const verifyResetPasswordTokenSchema = joi.object({
            email: joi.string().email().required().messages({
                'string.base': 'Email hanya bisa dimasukkan email',
                'string.empty': 'Email tidak boleh dikosongi',
                'any.required': 'Email wajib diisi',
            }),
            token: joi.string().min(6).required().messages({
                'string.empty': 'Token tidak boleh dikosongi',
                'string.min': 'Token tidak boleh kurang dari 6 karakter',
                'any.required': 'Token wajib diisi',
            }),
        });

        try {
            const { token, email } = req.body;

            const { error, value } = verifyResetPasswordTokenSchema.validate(req.body, {
                abortEarly: false,
            });
            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const checkEmail = await prisma.user_reset_passwords.findFirst({
                where: {
                    email
                }
            });

            if (!checkEmail) {
                res.status(422).json({
                    message: 'email tidak valid',
                    data: null
                });

                return;
            }

            const compareToken = await bcrypt.compare(token, checkEmail.token);

            if (!compareToken) {
                res.status(422).json({
                    message: 'token tidak valid',
                    data: null
                });

                return;
            }

            const now = new Date();
            const tokenExpired = new Date(checkEmail.expired_at);

            if (now < tokenExpired) {
                res.status(200).json({
                    message: 'token valid, silahkan buat password baru',
                    data: null
                });

                return;
            } else {
                res.status(410).json({
                    message: 'token tidak valid, silahkan request token kembali',
                    data: null
                });

                return;
            }

        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    resetPassword: async (req, res) => {
        const resetPasswordSchema = joi.object({
            email: joi.string().email().required().messages({
                'string.base': 'Email hanya bisa dimasukkan email',
                'string.empty': 'Email tidak boleh dikosongi',
                'any.required': 'Email wajib diisi',
            }),
            token: joi.string().min(6).required().messages({
                'string.empty': 'Token tidak boleh dikosongi',
                'string.min': 'Token tidak boleh kurang dari 6 karakter',
                'any.required': 'Token wajib diisi',
            }),
            password: joi.string().min(8).required().messages({
                'string.empty': 'Password baru tidak boleh dikosongi',
                'string.min': 'Password baru tidak boleh kurang dari 8 karakter',
                'any.required': 'Password baru wajib diisi',
            }),
            password_confirm: joi.ref('password').messages({
                'string.empty': 'Password baru konfirmasi tidak boleh dikosongi',
                'string.min': 'Password baru konfirmasi tidak boleh kurang dari 8 karakter',
                'any.required': 'Password baru konfirmasi wajib diisi',
                'any.ref': 'Password baru konfirmasi tidak sesuai dengan password baru',
            }),
        });

        try {
            const { email, token, password, password_confirm } = req.body;

            const { error, value } = resetPasswordSchema.validate(req.body, {
                abortEarly: false,
            });
            if (error) {
                let message = error.details[0].message;
                res.status(422).json({ message: "Format Tidak Valid", data: message });

                return;
            }

            const checkToken = await prisma.user_reset_passwords.findFirst({
                where: {
                    email
                }
            });

            if (!checkToken) {
                res.status(422).json({
                    message: 'email tidak valid',
                    data: null
                });

                return;
            }

            const compareToken = await bcrypt.compare(token, checkToken.token);

            if (!compareToken) {
                res.status(422).json({
                    message: 'token tidak valid',
                    data: null
                });

                return;
            }

            if (password != password_confirm) {
                res.status(422).json({
                    message: 'password dan password konfirmasi tidak sesuai',
                    data: null
                });

                return;
            }

            const passwordHash = await bcrypt.hash(password, 10);

            await prisma.users.update({
                where: {
                    email: checkToken.email
                },
                data: {
                    password: passwordHash
                }
            });

            await prisma.user_reset_passwords.delete({
                where: {
                    id: checkToken.id
                }
            });

            res.status(200).json({
                message: 'reset password berhasil, silahkan login kembali',
                data: null
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
    loginWithGoogle: async (req, res) => {
        try {
            const { google_id, email, name, image_profile, onesignal_id } = req.body;

            const passwordHash = await bcrypt.hash(`${google_id}_${Date.now().toLocaleString()}`, 10)

            const checkUser = await prisma.users.upsert({
                where: {
                    email,
                },
                create: {
                    google_id,
                    email,
                    name,
                    password: passwordHash,
                    image_profile,
                    onesignal_id,
                    wallets: {
                        create: {
                            name: "Dompet",
                            balance: 0,
                        }
                    }
                },
                update: {
                    image_profile, google_id, onesignal_id
                }
            });

            // Onesignal Update
            await prisma.users.update({
                where: {
                    id: checkUser.id
                },
                data: {
                    onesignal_id: onesignal_id
                }
            });

            // Add Player OneSignal
            await addPlayer({ onesignal_id });

            const payload = { id: checkUser.id };
            const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY);

            delete checkUser.password;
            delete checkUser.google_id;
            delete checkUser.onesignal_id;

            res.status(200).json({
                message: 'success login user!',
                data: {
                    user: checkUser,
                    access_token: accessToken
                }
            });
        } catch (error) {
            res.status(500).send(error.message);
        }
    },
};