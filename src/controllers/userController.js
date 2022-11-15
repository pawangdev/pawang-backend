const { PrismaClient } = require("@prisma/client");
const bcrypt = require("bcryptjs");
const jwt = require("jsonwebtoken");
const crypto = require("crypto");
const joi = require("joi");
const eventsEmitter = require("../helpers/event");
const { addPlayer } = require("../helpers/notification");

const prisma = new PrismaClient();

module.exports = {
  register: async (req, res) => {
    const registerSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
        email: joi.string().email().required().messages({
          "string.base": "Email hanya bisa dimasukkan email",
          "string.email": "Email hanya bisa dimasukkan email",
          "string.empty": "Email tidak boleh dikosongi",
          "any.required": "Email wajib diisi",
        }),
        password: joi.string().min(8).required().messages({
          "string.empty": "Password tidak boleh dikosongi",
          "string.min": "Password tidak boleh kurang dari 8 karakter",
          "any.required": "Password wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { name, email, password, gender, phone, onesignal_id } = req.body;

      const { error, value } = registerSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      // Check Email Duplicated
      const checkEmail = await prisma.users.findUnique({
        where: {
          email,
        },
      });

      if (checkEmail) {
        throw { status: 422, message: "EMAIL_HAS_REGISTERED", data: null };
      }

      const hashedPassword = await bcrypt.hash(password, 10);
      const newUser = await prisma.users.create({
        data: {
          name,
          email,
          password: hashedPassword,
          gender,
          phone,
          wallets: {
            create: {
              name: "Dompet",
              balance: 0,
            },
          },
          users_onesignals: {
            create: {
              onesignal_id,
            },
          },
        },
      });

      // Add Player OneSignal
      await addPlayer({ email, onesignal_id });

      delete newUser.password;
      delete newUser.google_id;
      delete newUser.onesignal_id;

      const payload = { id: newUser.id };
      const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY);
      res.status(201).json({
        status: true,
        message: "SUCCESS_CREATE_USER",
        data: {
          user: newUser,
          access_token: accessToken,
        },
      });

      return;
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  login: async (req, res) => {
    const loginSchema = joi
      .object({
        email: joi.string().email().required().messages({
          "string.base": "Email hanya bisa dimasukkan email",
          "string.empty": "Email tidak boleh dikosongi",
          "any.required": "Email wajib diisi",
        }),
        password: joi.string().min(8).required().messages({
          "string.empty": "Password tidak boleh dikosongi",
          "string.min": "Password tidak boleh kurang dari 8 karakter",
          "any.required": "Password wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const { email, password, onesignal_id } = req.body;

      const { error, value } = loginSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const user = await prisma.users.findUnique({
        where: {
          email,
        },
      });

      if (user) {
        try {
          if (await bcrypt.compare(password, user.password)) {
            // Onesignal Update
            await prisma.users.update({
              where: {
                id: user.id,
              },
              data: {
                users_onesignals: {
                  create: {
                    onesignal_id,
                  },
                },
              },
            });

            // Add Player OneSignal
            await addPlayer({ email, onesignal_id });

            delete user.password;
            delete user.google_id;
            delete user.onesignal_id;

            const payload = { id: user.id };
            const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY);
            res.status(200).json({
              status: true,
              message: "SUCCESS_LOGIN",
              data: {
                user,
                access_token: accessToken,
              },
            });

            return;
          } else {
            throw { status: 401, message: "WRONG_PASSWORD", data: null };
          }
        } catch (error) {
          return res.status(error.status || 500).json({
            status: false,
            message: error.message || "INTERNAL_SERVER_ERROR",
          });
        }
      } else {
        throw { status: 404, message: "EMAIL_NOT_FOUND", data: null };
      }
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  profile: async (req, res) => {
    try {
      const user_id = req.user.id;

      const user = await prisma.users.findUnique({
        where: {
          id: user_id,
        },
      });

      delete user.password;
      delete user.google_id;
      delete user.onesignal_id;

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_PROFILE",
        data: user,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  changeProfile: async (req, res) => {
    const changeProfileSchema = joi
      .object({
        name: joi.string().required().messages({
          "string.base": "Nama hanya bisa dimasukkan text",
          "string.empty": "Nama tidak boleh dikosongi",
          "any.required": "Nama wajib diisi",
        }),
      })
      .unknown(true);

    try {
      const user_id = req.user.id;
      const { name, gender, phone } = req.body;

      const { error, value } = changeProfileSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const updateUser = await prisma.users.update({
        where: {
          id: user_id,
        },
        data: {
          name,
          gender,
          phone,
        },
      });

      delete updateUser.password;
      delete updateUser.google_id;
      delete updateUser.onesignal_id;

      res.status(200).json({
        status: true,
        message: "SUCCESS_UPDATE_DATA",
        data: {
          user: updateUser,
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
  changePassword: async (req, res) => {
    const changePasswordSchema = joi.object({
      old_password: joi.string().min(8).required().messages({
        "string.empty": "Password lama tidak boleh dikosongi",
        "string.min": "Password lama tidak boleh kurang dari 8 karakter",
        "any.required": "Password lama wajib diisi",
      }),
      new_password: joi.string().min(8).required().messages({
        "string.empty": "Password baru tidak boleh dikosongi",
        "string.min": "Password baru tidak boleh kurang dari 8 karakter",
        "any.required": "Password baru wajib diisi",
      }),
      new_password_confirm: joi.ref("new_password"),
    });

    try {
      const user_id = req.user.id;
      const { old_password, new_password, new_password_confirm } = req.body;

      const { error, value } = changePasswordSchema.validate(req.body, {
        abortEarly: false,
      });

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const checkUser = await prisma.users.findUnique({
        where: { id: user_id },
      });

      if (checkUser) {
        try {
          if (await bcrypt.compare(old_password, checkUser.password)) {
            if (new_password !== new_password_confirm) {
              let message = error.details[0].message;
              throw {
                status: 400,
                message: "PASSWORD_CONFIRM_NOT_THE_SAME_NEW_PASSWORD",
                data: null,
              };
            }

            const newPasswordHashed = await bcrypt.hash(new_password, 10);

            const updateUser = await prisma.users.update({
              where: {
                id: user_id,
              },
              data: {
                password: newPasswordHashed,
              },
            });

            delete updateUser.password;

            res.status(200).json({
              status: true,
              message: "SUCCESS_UPDATE_DATA",
              data: {
                updateUser,
              },
            });

            return;
          } else {
            throw {
              status: 401,
              message: "WRONG_PASSWORD",
              data: null,
            };
          }
        } catch (error) {
          return res.status(error.status || 500).json({
            status: false,
            message: error.message || "INTERNAL_SERVER_ERROR",
          });
        }
      }
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  requestResetPasswordToken: async (req, res) => {
    const requestResetPasswordTokenSchema = joi.object({
      email: joi.string().email().required(),
    });

    try {
      const { email } = req.body;

      const { error, value } = requestResetPasswordTokenSchema.validate(
        req.body,
        {
          abortEarly: false,
        }
      );

      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const checkEmail = await prisma.users.findUnique({
        where: {
          email,
        },
      });

      if (!checkEmail) {
        throw { status: 404, message: "EMAIL_NOT_FOUND", data: null };
      }

      const code = crypto.randomBytes(3).toString("hex");
      const resetToken = parseInt(code.toString("hex"), 16)
        .toString()
        .substring(0, 6);
      const hashToken = await bcrypt.hashSync(resetToken, 8);
      const dateNow = new Date();
      dateNow.setMinutes(dateNow.getMinutes() + 10);
      const expiredToken = new Date(dateNow).toISOString();

      await prisma.user_reset_passwords.upsert({
        create: {
          email,
          token: hashToken,
          expired_at: expiredToken,
        },
        update: {
          email,
          token: hashToken,
          expired_at: expiredToken,
        },
        where: {
          email,
        },
      });

      const data = {
        from: "Pawang <teampawang.dev@gmail.com>",
        to: email,
        subject: "Kode Lupa Kata Sandi",
        text: `Gunakan kode ini untuk mengatur ulang kata sandi akun Anda: ${resetToken}. Kode hanya berlaku 10 menit.`,
        html: `<p>Gunakan kode ini untuk mengatur ulang kata sandi akun Anda: <strong>${resetToken}</strong>. Kode hanya berlaku 10 menit.</p>`,
      };

      eventsEmitter.emit("sendEmail", data);

      res.status(200).json({
        status: true,
        message: "SUCCESS_SEND_TOKEN",
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
  verifyResetPasswordToken: async (req, res) => {
    const verifyResetPasswordTokenSchema = joi.object({
      email: joi.string().email().required().messages({
        "string.base": "Email hanya bisa dimasukkan email",
        "string.empty": "Email tidak boleh dikosongi",
        "any.required": "Email wajib diisi",
      }),
      token: joi.string().min(6).required().messages({
        "string.empty": "Token tidak boleh dikosongi",
        "string.min": "Token tidak boleh kurang dari 6 karakter",
        "any.required": "Token wajib diisi",
      }),
    });

    try {
      const { token, email } = req.body;

      const { error, value } = verifyResetPasswordTokenSchema.validate(
        req.body,
        {
          abortEarly: false,
        }
      );
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const checkEmail = await prisma.user_reset_passwords.findFirst({
        where: {
          email,
        },
      });

      if (!checkEmail) {
        throw { status: 404, message: "EMAIL_NOT_FOUND", data: null };
      }

      const compareToken = await bcrypt.compare(token, checkEmail.token);

      if (!compareToken) {
        throw { status: 404, message: "TOKEN_NOT_MISMATCH", data: null };
      }

      const now = new Date();
      const tokenExpired = new Date(checkEmail.expired_at);

      if (now < tokenExpired) {
        res.status(200).json({
          status: true,
          message: "SUCCESS_VERIFY_TOKEN",
          data: null,
        });

        return;
      } else {
        throw { status: 410, message: "TOKEN_EXPIRED", data: null };
      }
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  resetPassword: async (req, res) => {
    const resetPasswordSchema = joi.object({
      email: joi.string().email().required().messages({
        "string.base": "Email hanya bisa dimasukkan email",
        "string.empty": "Email tidak boleh dikosongi",
        "any.required": "Email wajib diisi",
      }),
      token: joi.string().min(6).required().messages({
        "string.empty": "Token tidak boleh dikosongi",
        "string.min": "Token tidak boleh kurang dari 6 karakter",
        "any.required": "Token wajib diisi",
      }),
      password: joi.string().min(8).required().messages({
        "string.empty": "Password baru tidak boleh dikosongi",
        "string.min": "Password baru tidak boleh kurang dari 8 karakter",
        "any.required": "Password baru wajib diisi",
      }),
      password_confirm: joi.ref("password"),
    });

    try {
      const { email, token, password, password_confirm } = req.body;

      const { error, value } = resetPasswordSchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const checkToken = await prisma.user_reset_passwords.findFirst({
        where: {
          email,
        },
      });

      if (!checkToken) {
        throw { status: 404, message: "EMAIL_NOT_FOUND", data: null };
      }

      const compareToken = await bcrypt.compare(token, checkToken.token);

      if (!compareToken) {
        throw { status: 404, message: "TOKEN_NOT_MISMATCH", data: null };
      }

      if (password != password_confirm) {
        throw {
          status: 400,
          message: "PASSWORD_CONFIRM_NOT_THE_SAME_NEW_PASSWORD",
          data: null,
        };
      }

      const passwordHash = await bcrypt.hash(password, 10);

      await prisma.users.update({
        where: {
          email: checkToken.email,
        },
        data: {
          password: passwordHash,
        },
      });

      await prisma.user_reset_passwords.delete({
        where: {
          id: checkToken.id,
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_RESET_PASSWORD",
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
  loginWithGoogle: async (req, res) => {
    try {
      const { google_id, email, name, image_profile, onesignal_id } = req.body;

      const passwordHash = await bcrypt.hash(
        `${google_id}_${Date.now().toLocaleString()}`,
        10
      );

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
          wallets: {
            create: {
              name: "Dompet",
              balance: 0,
            },
          },
          users_onesignals: {
            create: {
              onesignal_id,
            },
          },
        },
        update: {
          image_profile,
          google_id,
          users_onesignals: {
            create: {
              onesignal_id,
            },
          },
        },
      });

      // Add Player OneSignal
      await addPlayer({ email, onesignal_id });

      const payload = { id: checkUser.id };
      const accessToken = jwt.sign(payload, process.env.TOKEN_SECRET_KEY);

      delete checkUser.password;
      delete checkUser.google_id;
      delete checkUser.onesignal_id;

      res.status(200).json({
        status: true,
        message: "SUCCESS_LOGIN",
        data: {
          user: checkUser,
          access_token: accessToken,
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
  logout: async (req, res) => {
    try {
      const { onesignal_id } = req.body;

      const checkOneSignal = await prisma.users_onesignals.findMany({
        where: {
          AND: {
            onesignal_id,
            user_id: req.user.id,
          },
        },
      });

      if (checkOneSignal.length > 0) {
        await prisma.users_onesignals.deleteMany({
          where: {
            onesignal_id,
          },
        });
      }

      res.status(200).json({
        status: true,
        message: "SUCCESS_LOGOUT",
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
