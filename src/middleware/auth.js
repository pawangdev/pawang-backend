const jwt = require('jsonwebtoken');
const { PrismaClient } = require('@prisma/client');
const prisma = new PrismaClient();

module.exports = (req, res, next) => {
    const authHeader = req.headers['authorization'];
    const token = authHeader && authHeader.split(" ")[1];

    if (token == null) return res.status(401).send({ message: "Unauthorized" });

    jwt.verify(token, process.env.TOKEN_SECRET_KEY, async (err, user) => {
        try {
            if (err) return res.status(401).send({ message: "Unauthorized" });
            const auth = await prisma.users.findUnique({
                where: {
                    id: user.id
                }
            });

            if (!auth) return res.status(401).send({ message: "Unauthorized" });
            req.user = auth;
            next();
        } catch (error) {
            throw new Error(error.message);
        }
    });
}