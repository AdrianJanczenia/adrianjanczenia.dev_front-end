# Adrian Janczenia - Front-End Service

The **Front-End Service** is the presentation layer of the portfolio ecosystem. It is a high-performance web application built with Go, utilizing Server-Side Rendering (SSR) to deliver a fast, SEO-friendly, and accessible user experience.

## Service Role

This service manages the user-facing part of the application. Its primary responsibilities include:

- **Dynamic Server-Side Rendering**: Generating HTML on the fly using Go's native 'html/template' engine for optimal performance and security.
- **Localization (i18n)**: Handling multi-language support by injecting content retrieved from the microservices into the templates.
- **Gateway Integration**: Orchestrating communication with the Gateway Service to fetch page data and handle CV access requests.
- **Secure Download Orchestration**: Managing the multi-step process of password validation, token retrieval, and binary file streaming for CV downloads.

## Architecture and Design

The service follows the system-wide architectural standards, ensuring consistency and maintainability.

### Layered Pattern: Handler -> Process -> Task
1. Handler: Manages HTTP requests, parses form data, and selects the appropriate templates for rendering.
2. Process: Orchestrates high-level UI logic, such as coordinating content fetching with error handling fallback.
3. Task / Service: Executes atomic operations like calling the Gateway API or managing internal template data mapping.

### Advanced Rendering Engine
The service features a decoupled Rendering Logic. Unlike standard implementations, the renderer uses a dynamic template mapping system. This allows for flexible layout management and professional error page rendering (e.g., localized 404/500 pages) without hardcoded file paths in the handlers.

### Security and UX Features
- **Anti-Spam Protection**: Implements "Honeypot" mechanisms in forms to prevent automated bot submissions without requiring intrusive CAPTCHAs.
- **Server-Side Translated Errors**: Error messages are mapped to user-friendly, localized strings before reaching the browser, ensuring a professional UX.
- **Static Asset Optimization**: Efficiently serves optimized CSS and assets required for a modern look and feel.

## Technical Specification

- Go: 1.23+ (utilizing the latest html/template security features).
- Server-Side Rendering: Native Go templates for zero-dependency rendering.
- Tailwind CSS: For modern, responsive styling.
- Docker: Optimized multi-stage builds on Alpine Linux, specifically configured to include necessary web assets (HTML/CSS) in the final image.

## Environment Configuration

| Variable | Description |
|----------|-------------|
| APP_ENV | Runtime environment (local/production) |
| GATEWAY_URL | The internal or external URL of the Gateway Service |

## Development and Deployment

### Build Optimized Docker Image
docker build -t frontend-service .

### Execute Unit Tests
go test -v ./...

## Performance Note
By using Go's SSR instead of a heavy client-side framework, the application achieves near-instant First Contentful Paint (FCP) and eliminates the need for complex client-side state management, making it an ideal "business card" project.

---
Adrian Janczenia