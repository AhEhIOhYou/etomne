FROM node:18 AS build

RUN mkdir -p /app
WORKDIR /app

COPY . .
RUN npm i
RUN mkdir -p dist
RUN npm run build

# ---- nginx ----

FROM nginx:alpine3.17

RUN rm /etc/nginx/conf.d/default.conf
COPY ./nginx/frontend.conf /etc/nginx/conf.d/

COPY --from=build /app/dist/ /usr/share/nginx/html/

EXPOSE 80