const express = require("express");
const authMiddleware = require("../middleware/auth");
const adminMiddleware = require("../middleware/admin");
const userController = require("../controllers/userController");
const walletController = require("../controllers/walletController");
const transactionController = require("../controllers/transactionController");
const categoryController = require("../controllers/categoryController");
const taskReminderController = require("../controllers/taskReminderController");
const appVersioningController = require("../controllers/appVersioningController");

const swaggerUi = require("swagger-ui-express");
const swaggerFile = require("../../swagger_output.json");

const router = express();

// Documentation
router.use("/doc", swaggerUi.serve, swaggerUi.setup(swaggerFile));

// Auth
router.post("/auth/register", userController.register);
router.post("/auth/login", userController.login);
router.post("/auth/login/google", userController.loginWithGoogle);
router.get("/auth/profile", authMiddleware, userController.profile);
router.post(
  "/auth/change-password",
  authMiddleware,
  userController.changePassword
);
router.post(
  "/auth/change-profile",
  authMiddleware,
  userController.changeProfile
);
router.post(
  "/auth/reset-password/request",
  userController.requestResetPasswordToken
);
router.post(
  "/auth/reset-password/verify",
  userController.verifyResetPasswordToken
);
router.post("/auth/reset-password", userController.resetPassword);

// Wallets
router.get("/wallets", authMiddleware, walletController.index);
router.get("/wallets/:id", authMiddleware, walletController.show);
router.post("/wallets/create", authMiddleware, walletController.create);
router.put("/wallets/update/:id", authMiddleware, walletController.update);
router.delete("/wallets/delete/:id", authMiddleware, walletController.destroy);

// Transactions
router.get("/transactions", authMiddleware, transactionController.index);
router.get(
  "/transactions/details",
  authMiddleware,
  transactionController.detail
);
router.get("/transactions/:id", authMiddleware, transactionController.show);
router.post(
  "/transactions/create",
  authMiddleware,
  transactionController.create
);
router.put(
  "/transactions/update/:id",
  authMiddleware,
  transactionController.update
);
router.delete(
  "/transactions/delete/:id",
  authMiddleware,
  transactionController.destroy
);

// Categories
router.get("/categories", authMiddleware, categoryController.index);
router.get("/categories/:id", authMiddleware, categoryController.show);
router.post(
  "/categories/create/:id",
  authMiddleware,
  categoryController.create
);
router.put(
  "/categories/update/:id/:subcategoryId",
  authMiddleware,
  categoryController.update
);
router.delete(
  "/categories/delete/:id/:subcategoryId",
  authMiddleware,
  categoryController.destroy
);

// Task Reminder
router.get("/task-reminders", authMiddleware, taskReminderController.index);
router.get("/task-reminders/:id", authMiddleware, taskReminderController.show);
router.post(
  "/task-reminders/create",
  authMiddleware,
  taskReminderController.create
);
router.put(
  "/task-reminders/update/:id",
  authMiddleware,
  taskReminderController.update
);
router.delete(
  "/task-reminders/delete/:id",
  authMiddleware,
  taskReminderController.destroy
);

// App Versioning
router.get("/app-version", appVersioningController.index);
router.put(
  "/app-version/update",
  authMiddleware,
  adminMiddleware,
  appVersioningController.update
);

module.exports = router;
