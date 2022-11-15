const { PrismaClient } = require("@prisma/client");
const joi = require("joi");

const prisma = new PrismaClient();

module.exports = {
  index: async (req, res) => {
    let categories = await prisma.categories.findMany({
      include: {
        subcategories: {
          where: {
            user_id: req.user.id,
          },
        },
      },
    });

    if (req.query.type) {
      categories = await prisma.categories.findMany({
        include: {
          subcategories: {
            where: {
              user_id: req.user.id,
            },
          },
        },
        where: {
          type: req.query.type,
        },
      });
    }

    res.status(200).json({
      status: true,
      message: "SUCCESS_GET_DATA",
      data: categories,
    });
    try {
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
      });
    }
  },
  show: async (req, res) => {
    try {
      const { id } = req.params;

      const category = await prisma.categories.findUnique({
        where: { id: Number(id) },
        include: {
          subcategories: {
            where: {
              user_id: req.user.id,
            },
          },
        },
      });

      if (!category) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      res.status(200).json({
        status: true,
        message: "SUCCESS_GET_DATA",
        data: category,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
      });
    }
  },
  create: async (req, res) => {
    const createSubCategorySchema = joi.object({
      name: joi.string().required().messages({
        "string.base": "Nama hanya bisa dimasukkan text",
        "string.empty": "Nama tidak boleh dikosongi",
        "any.required": "Nama wajib diisi",
      }),
    });

    try {
      const { id } = req.params;

      const { name } = req.body;

      const { error, value } = createSubCategorySchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const category = await prisma.categories.findUnique({
        where: { id: Number(id) },
      });
      if (!category) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      const newSubCategory = await prisma.sub_categories.create({
        data: {
          name,
          type: category.type,
          category_id: category.id,
          user_id: req.user.id,
        },
      });

      res.status(201).json({
        status: true,
        message: "SUCCESS_CREATE_DATA",
        data: newSubCategory,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
      });
    }
  },
  update: async (req, res) => {
    const updateSubCategorySchema = joi.object({
      name: joi.string().required().messages({
        "string.base": "Nama hanya bisa dimasukkan text",
        "string.empty": "Nama tidak boleh dikosongi",
        "any.required": "Nama wajib diisi",
      }),
    });

    try {
      const { id, subcategoryId } = req.params;
      const { name } = req.body;

      const { error, value } = updateSubCategorySchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        throw { status: 400, message: "INVALID_INPUT", data: message };
      }

      const category = await prisma.categories.findFirst({
        where: {
          id: Number(id),
        },
      });

      if (!category) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      const subCategory = await prisma.sub_categories.findFirst({
        where: {
          AND: {
            id: Number(subcategoryId),
            category_id: Number(id),
            user_id: req.user.id,
          },
        },
      });

      if (!subCategory) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      const updateSubCategory = await prisma.sub_categories.update({
        where: {
          id: Number(subcategoryId),
        },
        data: {
          name,
          type: category.type,
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_UPDATE_DATA",
        data: updateSubCategory,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
      });
    }
  },
  destroy: async (req, res) => {
    try {
      const { id, subcategoryId } = req.params;

      const category = await prisma.categories.findUnique({
        where: {
          id: Number(id),
        },
      });

      if (!category) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      const subCategory = await prisma.sub_categories.findFirst({
        where: {
          AND: {
            id: Number(subcategoryId),
            category_id: Number(id),
            user_id: req.user.id,
          },
        },
      });

      if (!subCategory) {
        throw { status: 404, message: "DATA_NOT_FOUND" };
      }

      await prisma.sub_categories.delete({
        where: {
          id: Number(subcategoryId),
        },
      });

      res.status(200).json({
        status: true,
        message: "SUCCESS_DELETE_DATA",
        data: null,
      });
    } catch (error) {
      return res.status(error.status || 500).json({
        status: false,
        message: error.message || "INTERNAL_SERVER_ERROR",
      });
    }
  },
};
