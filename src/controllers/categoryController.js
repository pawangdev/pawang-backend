const { PrismaClient } = require('@prisma/client');
const joi = require('joi');

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
          type: req.query.type
        }
      });
    }

    res.status(200).json({
      message: 'success retreived data!',
      data: categories
    });
    try {
    } catch (error) {
      res.status(500).send(error.message);
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

      res.status(200).json({
        message: 'success retreived data!',
        data: category
      });
    } catch (error) {
      res.status(500).send(error.message);
    }
  },
  create: async (req, res) => {
    const createSubCategorySchema = joi.object({
      name: joi.string().required().messages({
        'string.base': 'Nama hanya bisa dimasukkan text',
        'string.empty': 'Nama tidak boleh dikosongi',
        'any.required': 'Nama wajib diisi',
      }),
    });

    try {
      const { id } = req.params;

      const {
        name,
      } = req.body;

      const { error, value } = createSubCategorySchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        res.status(422).json({ message: "Format Tidak Valid", data: message });

        return;
      }

      const category = await prisma.categories.findUnique({ where: { id: Number(id) } });
      if (!category) {
        res.status(404).json({
          message: 'failed retreived data!',
          data: null
        });
      }

      const newSubCategory = await prisma.sub_categories.create({
        data: {
          name, type: category.type, category_id: category.id, user_id: req.user.id
        }
      });

      res.status(201).json({
        message: 'success create data!',
        data: newSubCategory
      })
    } catch (error) {
      res.status(500).send(error.message);
    }
  },
  update: async (req, res) => {
    const updateSubCategorySchema = joi.object({
      name: joi.string().required().messages({
        'string.base': 'Nama hanya bisa dimasukkan text',
        'string.empty': 'Nama tidak boleh dikosongi',
        'any.required': 'Nama wajib diisi',
      }),
    });

    try {
      const { id, subcategoryId } = req.params;
      const {
        name,
      } = req.body;

      const { error, value } = updateSubCategorySchema.validate(req.body, {
        abortEarly: false,
      });
      if (error) {
        let message = error.details[0].message;
        res.status(422).json({ message: "Format Tidak Valid", data: message });

        return;
      }

      const category = await prisma.categories.findUnique({
        where: {
          id: Number(id)
        }
      });

      if (!category) {
        res.status(404).json({
          message: 'failed retreived data!',
          data: null
        });

        return;
      }

      const subCategory = await prisma.sub_categories.findUnique({
        where: {
          id: Number(subcategoryId)
        }
      });

      if (!subCategory) {
        res.status(404).json({
          message: 'failed retreived data!',
          data: null
        });

        return;
      } else {
        if (subCategory.user_id != req.user.id) {
          res.status(404).json({
            message: 'failed retreived data!',
            data: null
          });

          return;
        }
      }

      const updateSubCategory = await prisma.sub_categories.update({
        where: {
          id: Number(subcategoryId)
        },
        data: {
          name, type: category.type
        }
      })

      res.status(200).json({
        message: 'success update data!',
        data: updateSubCategory
      })

    } catch (error) {
      res.status(500).send(error.message);
    }
  },
  destroy: async (req, res) => {
    try {
      const { id, subcategoryId } = req.params;

      const category = await prisma.categories.findUnique({
        where: {
          id: Number(id)
        }
      });

      if (!category) {
        res.status(404).json({
          message: 'failed retreived data!',
          data: null
        });

        return;
      }

      const subCategory = await prisma.sub_categories.findUnique({
        where: {
          id: Number(subcategoryId)
        }
      });

      if (!subCategory) {
        res.status(404).json({
          message: 'failed retreived data!',
          data: null
        });

        return;
      } else {
        if (subCategory.user_id != req.user.id) {
          res.status(404).json({
            message: 'failed retreived data!',
            data: null
          });

          return;
        }
      }

      await prisma.sub_categories.delete({
        where: {
          id: Number(subcategoryId)
        },
      });

      res.status(200).json({
        message: 'success delete data!',
        data: null
      })
    } catch (error) {
      res.status(500).send(error.message);
    }
  },
}
