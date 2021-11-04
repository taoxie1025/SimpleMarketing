# build stage
FROM node:lts-alpine as build-stage
WORKDIR /
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM openresty/openresty:1.11.2.3-xenial as production-stage
RUN opm get SkyLothar/lua-resty-jwt
COPY --from=build-stage /dist /usr/share/nginx/html
COPY nginx.conf /usr/local/openresty/nginx/conf/nginx.conf
COPY mime.types /usr/local/openresty/nginx/conf/conf/mime.types
COPY jwt-auth.lua  /usr/local/openresty/lualib/resty/jwt-auth.lua

EXPOSE 80
ENTRYPOINT ["nginx", "-g", "daemon off;"]