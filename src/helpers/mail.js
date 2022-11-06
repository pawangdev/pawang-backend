const nodemailer = require('nodemailer');

const emailService = nodemailer.createTransport({
  service: 'gmail',
  auth: {
    user: process.env.CONFIG_AUTH_EMAIL,
    pass: process.env.CONFIG_AUTH_PASSWORD,
  },
});

module.exports = emailService;
