FROM ubuntu:latest
FROM golang:1.20

WORKDIR /resume_website

# Copy application files
COPY main main
COPY loclx loclx
COPY start_server.sh start_server.sh
COPY app/geolite app/geolite
COPY app/public app/public
COPY app/blacklisted_ips.txt app/blacklisted_ips.txt
COPY app/blacklisted_sites.txt app/blacklisted_sites.txt
COPY app/whitelisted_countries.txt app/whitelisted_countries.txt

RUN chmod +x loclx

EXPOSE 8080

# Start the application and LocalXpose with a reserved domain
CMD ["sh", "-c", "./main -port 8080"]

