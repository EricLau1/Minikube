FROM alpine:3.7

RUN mkdir Uploads

RUN apk update && apk add openssh-client python curl bash

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-228.0.0-linux-x86_64.tar.gz

RUN mkdir -p /usr/local/gcloud \
  && tar -C /usr/local/gcloud -xvf /google-cloud-sdk-228.0.0-linux-x86_64.tar.gz \
  && /usr/local/gcloud/google-cloud-sdk/install.sh

RUN rm /google-cloud-sdk-228.0.0-linux-x86_64.tar.gz

RUN curl -sSL https://sdk.cloud.google.com | bash && \
    echo 'source /root/google-cloud-sdk/path.bash.inc' >>~/.bashrc
