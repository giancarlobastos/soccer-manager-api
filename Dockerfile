FROM ubuntu:latest
ADD soccer-manager-api_unix .
CMD ["./soccer-manager-api_unix"]
EXPOSE 8080