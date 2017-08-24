FROM alpine:3.6
ARG API_VERSION="0.2.3"
ARG EXEC_FILE_NAME="swagchat-realtime_alpine_amd64"
RUN apk --update add tzdata curl \
  && curl -LJO https://github.com/fairway-corp/swagchat-realtime/releases/download/v${API_VERSION}/${EXEC_FILE_NAME} \
  && chmod 700 ${EXEC_FILE_NAME} \
  && mv ${EXEC_FILE_NAME} /bin/swagchat-realtime
EXPOSE 9000
CMD ["swagchat-realtime"]
