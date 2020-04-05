# VScale service server deployer

Build: `go build -o .`

Docker run:
```
docker run -d --name vscale --restart always \
                -e 'VSCALE_TOKEN=secret-token' \
                -p 80:3000 \
                visteras/vscale:latest```