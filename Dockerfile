FROM node:lts-alpine3.16

# Create app directory
WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install
COPY . .

RUN npx prisma migrate deploy

EXPOSE 5000

CMD [ "npm", "run", "start" ]
