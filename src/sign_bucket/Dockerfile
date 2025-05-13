FROM node:18

WORKDIR /app

COPY ./package.json .
COPY ./package-lock.json .
RUN npm install


COPY ./dist/sign.js .

CMD ["node", "/app/sign.js"]