const { PrismaClient } = require("@prisma/client");
const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    try {
      const versioning = await prisma.app_versioning.findFirst();
      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: versioning,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
        data: null,
      });
    }
  },
  update: async (req, res) => {
    try {
      const { version, force_update } = req.body;

      if (!version || !force_update) {
        throw { status: 400, message: "INVALID_INPUT", data: null };
      }

      const versioning = await prisma.app_versioning.findFirst();

      const updateVersioning = await prisma.app_versioning.update({
        where: {
          id: versioning.id,
        },
        data: {
          version,
          force_update,
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_UPDATE_DATA",
        data: updateVersioning,
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
