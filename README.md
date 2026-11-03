# Front-end service

This Go service renders the portfolio, case studies, request explorer, and privacy page with `html/template`. It gets localized content and CV-related data from the gateway over HTTP. Browser JavaScript runs the proof-of-work (PoW) solver, CV form, and request explorer.

## Routes and CV flow

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/` | Render About me. |
| `GET` | `/experience` | Render experience, career path, cases, and technologies. |
| `GET` | `/experience/monolith-migration`, `/experience/orders-and-payments`, `/experience/reservation-waitlist` | Render individual case studies. |
| `GET` | `/how-this-site-works` | Render the interactive request explorer. |
| `GET` | `/contact` | Render contact links. |
| `GET` | `/privacy-policy` | Render the privacy page. |
| `GET` | `/api/pow` | Fetch a signed PoW challenge through the gateway. |
| `POST` | `/api/captcha` | Submit a solved PoW and fetch an image CAPTCHA. |
| `POST` | `/api/captcha-verify` | Verify the entered CAPTCHA answer. |
| `POST` | `/api/cv-token` | Request a CV download token. |
| `GET` | `/api/download/cv` | Stream the PDF for a token and language. |

Polish is the default content language; `?lang=en` selects English. Unknown page paths render an HTML 404 page using the same visual layout as CV errors. Unknown `/api` paths also render this page when the request prefers `text/html`, as browser navigation does; API requests receive a JSON 404 error slug. This includes misspelled CV download paths. The error page offers a PL/EN switch.

The browser solves PoW, displays the CAPTCHA image, verifies the answer, and then requests a CV token with the chosen language and password. The download request passes through the front end, gateway, and content service. The download handler accepts only `GET`, rejects missing token or language, and sends `Cache-Control: no-store` and `Referrer-Policy: no-referrer`.

The browser solver currently requires four leading zero hexadecimal characters in the SHA-256 hash of `seed + nonce`. Keep this value aligned with the CAPTCHA service's configured difficulty. The form also contains a hidden honeypot field; normal requests leave it empty.

## Structure and configuration

Handlers accept HTTP requests, processes coordinate page or API flows, and small service interfaces call the gateway. The renderer loads layouts and partials from `internal/web/template/`; styles and images are in `internal/web/static/`. Dependencies and routes are assembled in `internal/app` without a global configuration variable.

The templates render schema-v2 localized content and error pages. The request explorer uses labels and descriptions from the content document; its request topology and animations live in `internal/web/static/js/explorer.js`. Images are in `internal/web/static/img/`. The CV modal refreshes an expired CAPTCHA challenge and starts the PDF download through `internal/web/static/js/cv.js`. The responsive navigation and email copy interaction use `internal/web/static/js/site.js`.

Inter and Fira Code are served from `internal/web/static/fonts/` as local WOFF2 files. Their OFL license notices are stored alongside them. The client accepts downstream errors only when a recognized JSON slug matches its HTTP status; a missing or inconsistent error payload becomes `502 error_bad_gateway`. The `error_message` slug remains accepted with its HTTP status for compatibility with gateway responses.

The module requires Go 1.25 (see `go.mod`). `APP_ENV=production` selects `config/production/config.yml`; otherwise the local file is used. `GATEWAY_URL` can override the gateway address in YAML. Both configurations set a 20-second HTTP server write timeout. The gateway client uses per-operation deadlines for ordinary requests; PDF streaming follows the incoming request context instead of a fixed client-wide body timeout.

## Run locally

1. Start the gateway and its dependencies. The front end fetches page content through the gateway, so it cannot render portfolio pages from this repository alone.
2. Set `GATEWAY_URL` to the gateway's HTTP base URL, for example `http://127.0.0.1:8082` when the gateway is running on the host. `config/local/config.yml` leaves this value empty.
3. Run `go run .` from this repository root, then open `http://127.0.0.1:8080/`. The local configuration uses port 8080 and loads templates and static assets relative to the repository root.

The `Dockerfile` packages the templates and static assets with the binary. In a container, set `GATEWAY_URL` to an address reachable from that container. The site text is owned by the content service; changes to page copy belong in its paired Polish and English JSON files.

## Verification

```sh
go test -race ./...
```

These tests cover handlers, processes, rendering, and gateway client behavior. They do not start the full service stack or exercise the browser's JavaScript.
