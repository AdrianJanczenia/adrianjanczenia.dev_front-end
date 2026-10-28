# Front-end service

This Go service renders the portfolio and privacy pages with `html/template`. It gets localized content and CV-related data from the gateway over HTTP. A small amount of browser JavaScript runs the proof-of-work (PoW) solver and CV form; there is no client-side application framework.

## Routes and CV flow

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/` | Render the portfolio page. |
| `GET` | `/privacy-policy`, `/polityka-prywatnosci` | Render the privacy page. |
| `GET` | `/api/pow` | Fetch a signed PoW challenge through the gateway. |
| `POST` | `/api/captcha` | Submit a solved PoW and fetch an image CAPTCHA. |
| `POST` | `/api/captcha-verify` | Verify the entered CAPTCHA answer. |
| `POST` | `/api/cv-token` | Request a CV download token. |
| `GET` | `/api/download/cv` | Stream the PDF for a token and language. |

Unknown page paths render an HTML 404 page using the same visual layout as CV errors. Unknown `/api` paths also render this page when the request prefers `text/html`, as browser navigation does; API requests receive a JSON 404 error slug. This includes misspelled CV download paths. Polish is the default; `?lang=en` selects English, and the error page offers a PL/EN switch.

The browser solves PoW, displays the CAPTCHA image, verifies the answer, and then requests a CV token with the chosen language and password. The download request passes through the front end, gateway, and content service. The download handler accepts only `GET`, rejects missing token or language, and sends `Cache-Control: no-store` and `Referrer-Policy: no-referrer`.

The browser solver currently requires four leading zero hexadecimal characters in the SHA-256 hash of `seed + nonce`. Keep this value aligned with the CAPTCHA service's configured difficulty. The form also contains a hidden honeypot field; normal requests leave it empty.

## Structure and configuration

Handlers accept HTTP requests, processes coordinate page or API flows, and small service interfaces call the gateway. The renderer loads layouts and partials from `internal/web/template/`; styles and images are in `internal/web/static/`. Dependencies and routes are assembled in `internal/app` without a global configuration variable.

The templates render localized content and error pages. The CV modal refreshes an expired CAPTCHA challenge, and the browser opens the PDF link differently on mobile and desktop. These browser behaviors live in `internal/web/template/partials/scripts.html`.

Inter and Fira Code are served from `internal/web/static/fonts/` as local WOFF2 files. Their OFL license notices are stored alongside them. The browser does not need to contact Google Fonts. The client accepts downstream errors only when a recognized JSON slug matches its HTTP status; a missing or inconsistent error payload becomes `502 error_bad_gateway`. The old `error_message` slug remains accepted with its HTTP status during a service rollout.

The module requires Go 1.25 (see `go.mod`). `APP_ENV=production` selects `config/production/config.yml`; otherwise the local file is used. `GATEWAY_URL` can override the gateway address in YAML. Both configurations set a 20-second HTTP server write timeout. The gateway client uses per-operation deadlines for ordinary requests; PDF streaming follows the incoming request context instead of a fixed client-wide body timeout.

From the repository root, start the gateway, provide its address in configuration, and run `go run .`. The `Dockerfile` is available for container builds.

## Verification

```sh
go test -race ./...
```

These tests cover handlers, processes, rendering, and gateway client behavior. They do not start the full service stack or exercise the browser's JavaScript.
