FROM nginx
RUN mkdir -p /app
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY ./frontend/admin/dist /app/admin
COPY ./frontend/appvideo/dist /app/app
COPY ./appvideo.exe /app/av.exe
COPY ./start.sh /start.sh
COPY ./config/conf.example.yml /app/config.yaml
RUN chmod +x /start.sh

EXPOSE 80

CMD ["/start.sh"]