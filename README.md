## Demo

You can try Mindsdb ML SQL server here [(demo)](https://cloud.mindsdb.com).

## Mail server for development

To run the mail server for development:

```
docker run -p 25:25 -p 80:80 -p 143:143 -p 110:110 -p 443:443 -p 465:465 -p 587:587 -p 993:993 -p 995:995 -p 4190:4190 -e TZ=Europe/Prague -e DISABLE_CLAMAV="true" -e DISABLE_RSPAMD="true" -v posteio:/data --name "posteio" -h "email.com" -d docker.uclv.cu/analogic/poste.io:2.3.7
```


