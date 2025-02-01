FROM node:alpine

WORKDIR /app
COPY package.json .
# Install all dependencies since we need typescript and ts-node in production
RUN npm install
COPY . .

CMD [ "npm", "run", "start:prod" ]
