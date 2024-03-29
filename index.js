const express = require("express");
const dotenv = require("dotenv");
const cors = require("cors");
const compression = require("compression");
const morgan = require("morgan");

dotenv.config();

const { scheduler } = require("./src/controllers/taskReminderController");

const routes = require("./src/routes/index");

const app = express();

app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cors({}));
app.use(morgan("tiny"));

// Compress all HTTP responses
app.use(compression());

// Routing
app.get("/", (req, res) => {
    res.send("Team Pawang API 2023");
});

// Static File
app.use("/public", express.static("./public"));

// All Routes
app.use("/api", routes);

// Handling Error Page Not Found
app.use((req, res) => {
    res.status(404).json({ status: false, message: "404_NOT_FOUND", data: null });
});

// Notification Scheduler
scheduler();

app.listen(process.env.APP_PORT || 5000, () => {
    console.log("🚀 Server Started on port", process.env.APP_PORT || 5000);
});
