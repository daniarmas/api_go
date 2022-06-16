## Demo

You can try Mindsdb ML SQL server here [(demo)](https://cloud.mindsdb.com).

## Mail server for development

To run the mail server for development:

```
docker run -p 25:25 -p 80:80 -p 143:143 -p 110:110 -p 443:443 -p 465:465 -p 587:587 -p 993:993 -p 995:995 -p 4190:4190 -e TZ=Europe/Prague -e DISABLE_CLAMAV="true" -e DISABLE_RSPAMD="true" -v posteio:/data --name "posteio" -h "email.com" --restart always -d docker.uclv.cu/analogic/poste.io:2.3.7
```

To run the Minio server:

```
docker run -p 9011:9000 -p 9010:9001 --restart always -e MINIO_ROOT_USER="root" -e MINIO_ROOT_PASSWORD="root1234" -v miniodata:/data -d docker.uclv.cu/minio/minio:RELEASE.2022-01-08T03-11-54Z server /data --console-address ":9001"
```


minio:
    image: docker.uclv.cu/minio/minio:RELEASE.2022-01-08T03-11-54Z
    container_name: minio
    volumes:
      - miniodata:/data
    restart: always
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: root
      MINIO_ROOT_PASSWORD: root1234
    ports:
      - 9011:9000
      - 9010:9001
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:9000/minio/health/live"]
      interval: 30s
      timeout: 20s
      retries: 3