# canvas
A blank canvas for your web app, in Go.

This repository is used in the course [Build Cloud Apps in Go](https://www.golang.dk/courses/build-cloud-apps-in-go).

## Dependencies

```bash
go get -u github.com/go-chi/chi/v5
go get -u github.com/matryer/is
go get -u go.uber.org/zap
go get -u golang.org/x/sync
go get -u github.com/maragudk/gomponents github.com/maragudk/gomponents-heroicons
```

## Steps

```bash
# awscliv2
curl -fsSL https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip
unzip -q awscliv2.zip
sudo ./aws/install
rm -f awscliv2.zip
# lightsailctl
curl -fsSL https://s3.us-west-2.amazonaws.com/lightsailctl/latest/linux-amd64/lightsailctl -o /usr/local/bin/lightsailctl
sudo chmod +x /usr/local/bin/lightsailctl
```

```bash
aws lightsail get-container-services
aws lightsail create-container-service --service-name canvas --power micro --scale 1
aws lightsail push-container-image --service-name canvas --label app --image canvas
aws lightsail create-container-service-deployment --service-name canvas \
	--containers '{"app":{"image":":canvas.app.3","environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}' \
	--public-endpoint '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}'
```
