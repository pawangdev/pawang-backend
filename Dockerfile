FROM node:lts-alpine3.16

# Create app directory
WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install
COPY . .
EXPOSE 5000

RUN npx prisma db push
CMD [ "npm", "run", "start" ]
