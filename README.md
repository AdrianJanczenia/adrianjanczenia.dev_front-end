# Adrian Janczenia - Front-End Service

The **Front-End Service** is the presentation layer of the portfolio ecosystem. It is a high-performance web application built with Go, utilizing Server-Side Rendering (SSR) to deliver a fast, SEO-friendly, and accessible user experience.

## Service Role

This service manages the user-facing part of the application. Its primary responsibilities include:

- **Dynamic Server-Side Rendering**: Generating HTML on the fly using Go's native 'html/template' engine for optimal performance and security.
- **Localization (i18n)**: Handling multi-language support by injecting content retrieved from the microservices into the templates.
- **Gateway Integration**: Orchestrating communication with the Gateway Service to fetch page data and handle CV access requests.
- **Interactive Captcha Management**: Handling the Proof of Work (PoW) process and Captcha image verification directly within the UI.

## Architecture and Design

The service follows the system-wide architectural standards, ensuring consistency and maintainability.

### Layered Pattern: Handler -> Process -> Task
1. Handler: Manages HTTP requests, parses form data, and selects the appropriate templates for rendering.
2. Process: Orchestrates high-level UI logic, such as coordinating content fetching with error handling fallback.
3. Task / Service: Executes atomic operations like calling the Gateway API or managing internal template data mapping.

### Advanced Rendering Engine
The service features a decoupled Rendering Logic. The renderer uses a dynamic template mapping system, allowing for flexible layout management and professional error page rendering (e.g., localized 404/500 pages) without hardcoded file paths in the handlers.

### Security and UX Features
- **PoW and Captcha Integration**: Integrated anti-bot protection. If a Captcha session expires, the form automatically refreshes the challenge without closing the modal popup.
- **Enhanced CV Modal**: Optimized field layout (Password, Captcha Input, Captcha Image) for better readability and user flow.
- **Reliable Spinners**: CSS logic ensuring spinner animations restart every time a modal is opened or the global loader is activated.
- **Mobile & Safari Optimization**: Specialized logic for detecting Safari and mobile devices to handle specific browser behaviors like Safe Area Insets and smooth scrolling.

## Technical Specification

- Go: 1.23+
- SSR: Native Go templates for zero-dependency rendering.
- Modern CSS: Custom styling using CSS Variables and responsive design principles.
- Unit Tests: Full coverage for PoW fetching, Captcha fetching, and verification processes.

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
By using Go's SSR instead of a heavy client-side framework, the application achieves near-instant First Contentful Paint (FCP) and eliminates the need for complex client-side state management.

---
Adrian Janczenia