const jwt = require("jsonwebtoken");
const { PrismaClient } = require("@prisma/client");
const prisma = new PrismaClient();

module.exports = (req, res, next) => {
  try {
    const authHeader = req.headers["authorization"];
    const token = authHeader && authHeader.split(" ")[1];

    if (token == null)
      throw { status: 401, message: "Unauthorized", data: null };

    jwt.verify(token, process.env.TOKEN_SECRET_KEY, async (err, user) => {
      try {
        if (err) throw { status: 401, message: "Unauthorized", data: null };
        const auth = await prisma.users.findUnique({
          where: {
            id: user.id,
          },
        });

        if (!auth) throw { status: 401, message: "Unauthorized", data: null };
        req.user = auth;
        next();
      } catch (error) {
        return res.status(error.status || 500).json({
          status: false,
          message: error.message || "INTERNAL_SERVER_ERROR",
          data: null,
        });
      }
    });
  } catch (error) {
    return res.status(error.status || 500).json({
      status: false,
      message: error.message || "INTERNAL_SERVER_ERROR",
      data: null,
    });
  }
};
