FROM nginx:latest

COPY ./api_gateway.conf /etc/nginx/conf.d/api_gateway.conf
COPY ./ssl/ssl.ext /etc/nginx/ssl/ssl.ext
COPY ./ssl/ssl.crt /etc/nginx/ssl/ssl.crt
COPY ./ssl/ssl.key /etc/nginx/ssl/ssl.key

RUN rm /etc/nginx/conf.d/default.conf
