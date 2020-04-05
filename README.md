# VScale service server deployer

Build: `go build -o .`

Docker run:
```
docker run -d --name vscale --restart always \
                -e 'VSCALE_TOKEN=secret-token' \
                -p 80:3000 \
                visteras/vscale:latest
```

Creating 3 servers

Server password is unique, he sended on email in vscale account for every created server

```
curl --request POST \
  --url http://app-url:3000/create/3
```
If successfully
```
{
  "code": 200,
  "msg": [
    {
      "status": "queued",
      "deleted": "",
      "public_address": {
        "netmask": "",
        "gateway": "",
        "address": ""
      },
      "active": true,
      "location": "spb0",
      "locked": true,
      "hostname": "",
      "created": "05.04.2020 10:12:48",
      "keys": [],
      "private_address": {
        "netmask": "",
        "gateway": "",
        "address": ""
      },
      "made_from": "ubuntu_18.04_64_001_master",
      "name": "TmpSrv0",
      "ctid": 1479849,
      "rplan": "small"
    }
  ]
}
```
If error:
```
{
  "code": 500,
  "error": "created servers not successfully"
}
```

Delete all servers
```
curl --request DELETE \
  --url http://app-url:3000/delete
```
If successfully
```
{
  "code": 200,
  "msg": [
    {
      "status": "started",
      "deleted": "",
      "public_address": {
        "netmask": "",
        "gateway": "",
        "address": ""
      },
      "active": true,
      "location": "spb0",
      "locked": true,
      "hostname": "tmpsrv0",
      "created": "05.04.2020 10:12:48",
      "keys": [],
      "private_address": {
        "netmask": "",
        "gateway": "",
        "address": ""
      },
      "made_from": "ubuntu_18.04_64_001_master",
      "name": "TmpSrv0",
      "ctid": 1479849,
      "rplan": "small"
    }
  ]
}
```
If servers not found:
```
{
  "code": 200,
  "msg": null
}
```
If error:
```
{
  "code": 500,
  "error": "deleted servers not successfully"
}
```