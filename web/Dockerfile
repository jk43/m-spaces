FROM nginx:stable-alpine as molyspace-webapp

copy dist/spa /usr/share/nginx/html
copy docker/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80