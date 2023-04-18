# Usage

```
POST /sync {
  "host": "host.example.com"
  aliases: [
    "alias1.example.com"
    "alias2.example.com"
  ]
}
``` 

# Docker Compose

```yaml
version: "3"
services: 
  opnsense-proxy-api:
    image: ghcr.io/markfortedev/opnsenseproxyapi:master
    environment:
      - OPNSENSE_ADDRESS=address
      - API_KEY=key
      - API_SECRET=secret
      - DOMAIN_NAME=example.com
    ports:
      - "9657:9657"
```