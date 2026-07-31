# Homelab

Requirements: Docker Compose and Tailscale connected to your tailnet on the host.

Start Bolt:

```sh
docker compose up --build --detach
```

Bolt listens on `127.0.0.1:8080`. Caddy redirects HTTP requests on
`127.0.0.1:8081` to the HTTPS Tailscale Service hostname. Its SQLite database
is persisted in the `bolt-data` Docker volume.

Publish it to the tailnet:

```sh
scripts/configure-tailscale-service.sh
tailscale serve get-config --all
```

Configure `svc:blt` with `tcp:80` and `tcp:443` endpoints in Tailscale Services.
The default redirect hostname is `blt.tailb538a8.ts.net`; override it with
`BOLT_HOST` when starting Compose if needed.

`http://blt/<alias>` redirects to `https://blt.tailb538a8.ts.net/<alias>`.

Stop Bolt without removing its data:

```sh
docker compose down
```
