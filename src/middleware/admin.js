const { PrismaClient } = require("@prisma/client");
const prisma = new PrismaClient();

module.exports = async (req, res, next) => {
  try {
    const user_id = req.user.id;
    const user = await prisma.users.findFirst({
      where: {
        id: user_id,
      },
    });

    if (user.role !== "admin") {
      throw { status: 401, message: "Unauthorized", data: null };
    }

    next();
  } catch (error) {
    return res.status(error.status || 500).json({
      status: false,
      message: error.message || "INTERNAL_SERVER_ERROR",
      data: null,
    });
  }
};
