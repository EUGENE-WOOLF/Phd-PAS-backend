FROM golang:1.18-bullseye

# Set the working directory inside the container
WORKDIR $GOPATH/src/github.com/spo-iitk/ras-backend

RUN apt-get update && \
    apt-get install -y vim nginx git && \
    rm -rf /var/lib/apt/lists/*

RUN git config --global user.name "SPO Web Team"
RUN git config --global user.email "pas@iitk.ac.in"

COPY . .

RUN cp secret.yml.template secret.yml
RUN cp secret.GCPcredentials.json.template secret.GCPcredentials.json

# Configure nginx
RUN rm /etc/nginx/sites-enabled/default
RUN ln -s $GOPATH/src/github.com/spo-iitk/ras-backend/container/nginx.conf /etc/nginx/sites-enabled/default

# Expose port 80
EXPOSE 80

# Start the server
CMD ["./scripts/production.sh"]
