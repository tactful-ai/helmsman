ARG GO_VERSION="1.15.2"
ARG ALPINE_VERSION="3.12"
ARG GLOBAL_KUBE_VERSION="v1.19.0"
ARG GLOBAL_HELM_VERSION="v3.3.4"
ARG GLOBAL_HELM_DIFF_VERSION="v3.1.3"

### Go Builder & Tester ###
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

ARG GLOBAL_KUBE_VERSION
ARG GLOBAL_HELM_VERSION
ARG GLOBAL_HELM_DIFF_VERSION
ENV KUBE_VERSION=$GLOBAL_KUBE_VERSION
ENV HELM_VERSION=$GLOBAL_HELM_VERSION
ENV HELM_DIFF_VERSION=$GLOBAL_HELM_DIFF_VERSION

RUN apk add --update --no-cache ca-certificates git openssh ruby curl tar gzip make bash

RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

RUN curl -Lk https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz | tar zxv -C /tmp
RUN mv /tmp/linux-amd64/helm /usr/local/bin/helm && rm -rf /tmp/linux-amd64
RUN chmod +x /usr/local/bin/helm

RUN helm plugin install https://github.com/hypnoglow/helm-s3.git
RUN helm plugin install https://github.com/nouney/helm-gcs
RUN helm plugin install https://github.com/databus23/helm-diff --version ${HELM_DIFF_VERSION}
RUN helm plugin install https://github.com/futuresimple/helm-secrets
RUN rm -r /tmp/helm-diff /tmp/helm-diff.tgz

RUN gem install hiera-eyaml --no-doc
RUN update-ca-certificates

WORKDIR /go/src/github.com/tactful-ai/robban

COPY . .
# RUN make test
RUN LastTag=$(git describe --abbrev=0 --tags) \
    && TAG=$LastTag-$(date +"%d%m%y") \
    && LT_SHA=$(git rev-parse ${LastTag}^{}) \
    && LC_SHA=$(git rev-parse HEAD) \
    && if [ ${LT_SHA} != ${LC_SHA} ]; then TAG=latest-$(date +"%d%m%y"); fi \
    && make build

### Final Image ###
FROM alpine:${ALPINE_VERSION} as base

RUN apk add --update --no-cache ca-certificates git openssh ruby curl bash gnupg
RUN gem install hiera-eyaml --no-doc
RUN update-ca-certificates

COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm
COPY --from=builder /usr/local/bin/sops /usr/local/bin/sops
COPY --from=builder /root/.cache/helm/plugins/ /root/.cache/helm/plugins/
COPY --from=builder /root/.local/share/helm/plugins/ /root/.local/share/helm/plugins/

WORKDIR /opt
COPY --from=builder /go/src/github.com/tactful-ai/robban/public/ /opt/public/
COPY --from=builder /go/src/github.com/tactful-ai/robban/robban /opt/robban

EXPOSE 8080
ENTRYPOINT [ "robban" ]