FROM node:12.14.1

ARG ng_version=latest
RUN npm install -g @angular/cli@$ng_version --unsafe-perms && \
    ng version

ENTRYPOINT ["ng"]
