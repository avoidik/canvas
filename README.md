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
go get -u github.com/jmoiron/sqlx github.com/jackc/pgx/v4
go get -u github.com/maragudk/migrate
```

## Steps

Install prerequisites

```bash
#
# docker
#
curl -fsSL https://get.docker.com | bash
#
# docker-compose
#
sudo mkdir -p /usr/local/lib/docker/cli-plugins
sudo curl -fsSL https://github.com/docker/compose/releases/download/v2.6.0/docker-compose-linux-x86_64 -o /usr/local/lib/docker/cli-plugins/docker-compose
sudo chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
#
# golang
#
curl -fsSLO https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz
rm -f go1.18.3.linux-amd64.tar.gz
cat <<'EOF' | tee -a ~/.profile > /dev/null
if [ -d "$HOME/go/bin" ] ; then
    PATH="$HOME/go/bin:$PATH"
fi
if [ -d "/usr/local/go/bin" ] ; then
    PATH="$PATH:/usr/local/go/bin"
fi
EOF
source ~/.profile
#
# jq
#
sudo curl -fsSL https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 -o /usr/local/bin/jq
sudo chmod +x /usr/local/bin/jq
#
# awscliv2
#
curl -fsSL https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip -o awscliv2.zip
unzip -q awscliv2.zip
sudo ./aws/install
rm -rf awscliv2.zip aws/
#
# lightsailctl
#
sudo curl -fsSL https://s3.us-west-2.amazonaws.com/lightsailctl/latest/linux-amd64/lightsailctl -o /usr/local/bin/lightsailctl
sudo chmod +x /usr/local/bin/lightsailctl
```

Deploy

```bash
aws lightsail get-container-services
aws lightsail create-container-service --service-name canvas --power micro --scale 1
aws lightsail push-container-image --service-name canvas --label app --image canvas
aws lightsail create-container-service-deployment --service-name canvas \
	--containers '{"app":{"image":":canvas.app.3","environment":{"HOST":"","PORT":"8080","LOG_ENV":"production"},"ports":{"8080":"HTTP"}}}' \
	--public-endpoint '{"containerName":"app","containerPort":8080,"healthCheck":{"path":"/health"}}'
```

Database

```bash
docker compose up -d
docker run --rm --add-host=host.docker.internal:host-gateway -e PGPASSWORD=123 postgres:12 psql -h host.docker.internal -U canvas -A -c "select 'excellent' as partytime"
```
