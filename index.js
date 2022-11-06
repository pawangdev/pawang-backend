const express = require("express")
const dotenv = require("dotenv")
const cors = require("cors")
const compression = require("compression");
dotenv.config();

const { scheduler } = require('./src/controllers/taskReminderController');
const emailService = require("./src/helpers/mail");


const routes = require("./src/routes/index");

const app = express();

app.use(express.json())
app.use(express.urlencoded({ extended: true }))
app.use(cors({}));

// Compress all HTTP responsess
app.use(compression());


// Routing
app.get('/', (req, res) => {
  res.send('Team Pawang API 2022')
});

// Static File
app.use('/public', express.static('./public'))

// All Routes
app.use('/api', routes)

// Handling Error
app.use((error, req, res, next) => {
  const status = error.errorStatus || 500;
  const message = error.message;
  const data = error.data;

  res.status(status).json({ message: message, data: data });
});

// Notification Scheduler
scheduler();

app.listen(process.env.APP_PORT || 5000, () => {
  console.log("🚀 Server Started on port", process.env.APP_PORT || 5000);
});
