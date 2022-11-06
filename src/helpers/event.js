const events = require("events");
const { deleteFile } = require("./file");
const emailService = require("./mail");

// Events Emitter for sending emails
const eventsEmitter = new events.EventEmitter();

eventsEmitter.on("sendEmail", async (data) => {
  await emailService.sendMail(data);
});

eventsEmitter.on("deleteFile", async (data) => {
  await deleteFile(data);
})

module.exports = eventsEmitter;
