# TravelSphere

TravelSphere is a Beego-based full-stack travel discovery and trip planning application. It allows users to explore countries, view destination details, discover nearby attractions, and manage a personal travel wishlist.

This project was built as a Beego full-stack intern assessment and demonstrates MVC structure, server-side rendering, JSON APIs, AJAX partial updates, service-layer business logic, external API integration, middleware/filter usage, and unit testing.

---

## Features

* Server-rendered home page, country explorer, destination detail page, wishlist page, and dashboard.
* REST Countries API integration for country data.
* OpenTripMap API integration for nearby attractions.
* Search and region filtering for countries.
* AJAX-powered partial updates without full page reloads.
* Auth-protected wishlist and dashboard pages.
* In-memory wishlist storage with JSON API CRUD operations.
* Dashboard summary showing total saved, planned, and visited destinations.
* Beego layout template with header and footer partials.
* Reusable service layer and utility layer.
* Unit tests for services layer and utilities layer.

---

## Tech Stack

* Go
* Beego Framework
* Beego Templates
* Vanilla JavaScript
* CSS
* REST Countries API
* OpenTripMap API

---

## Project Structure

```txt
.
├── conf
│   ├── app.conf
│   └── app.example.conf
├── controllers
│   ├── api
│   │   ├── countries.go
│   │   ├── dashboard.go
│   │   └── wishlist.go
│   ├── auth.go
│   ├── base.go
│   ├── country.go
│   ├── dashboard.go
│   ├── home.go
│   └── wishlist.go
├── go.mod
├── go.sum
├── main.go
├── models
│   ├── attraction.go
│   ├── country.go
│   ├── response.go
│   └── wishlist.go
├── README.md
├── routers
│   └── router.go
├── services
│   ├── attraction_service.go
│   ├── attraction_service_test.go
│   ├── country_service.go
│   ├── country_service_test.go
│   ├── dashboard_service.go
│   ├── dashboard_service_test.go
│   ├── registry.go
│   ├── wishlist_service.go
│   └── wishlist_service_test.go
├── static
│   ├── css
│   │   └── app.css
│   └── js
│       ├── countries.js
│       ├── dashboard.js
│       ├── home.js
│       └── wishlist.js
├── utils
│   ├── formatter.go
│   ├── formatter_test.go
│   ├── opentripmap_client.go
│   ├── opentripmap_client_test.go
│   ├── restcountries_client.go
│   ├── restcountries_test.go
│   ├── validation.go
│   └── validation_test.go
└── views
    ├── countries.tpl
    ├── dashboard.tpl
    ├── destination.tpl
    ├── errors
    │   └── 404.tpl
    ├── home.tpl
    ├── layout.tpl
    ├── login.tpl
    ├── partials
    │   ├── footer.tpl
    │   └── header.tpl
    └── wishlist.tpl

```

---

## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/rangonroyutsab/travel-sphere-rangon.git
cd travel-sphere-rangon
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Configure environment values

Copy the example configuration file:

```bash
cp conf/app.example.conf conf/app.conf
```

Update `conf/app.conf` with your local settings.

Example:

```ini
appname = TravelSphere
httpport = 8080
runmode = dev

sessionon = true
copyrequestbody = true

DEMO_USERNAME = user
DEMO_PASSWORD = admin

REST_COUNTRIES_BASE_URL = https://api.restcountries.com/countries/v5
REST_COUNTRIES_API_KEY = your_rest_countries_api_key_here
OPENTRIPMAP_BASE_URL = https://api.opentripmap.com/0.1/en/places
OPENTRIPMAP_API_KEY = your_opentripmap_api_key_here
```

The REST Countries v5 API key is required for country data. The application sends it as a bearer token and requests paginated country records from the configured v5 endpoint.

The OpenTripMap API key is required for nearby attractions. If the key is missing or invalid, the application will still run, but the nearby attractions section may show an empty or fallback state.

---

## Running the Application

Run the project with:

```bash
bee run
```

Then open:

```txt
http://localhost:8080
```

---

## Demo Login

The application uses a simple demo login system for protected pages.

Default credentials:

```txt
Username: user
Password: admin
```

The username and password can be configured from `conf/app.conf`.

Protected routes:

```txt
/wishlist
/dashboard
/api/wishlist
/api/dashboard/summary
```

---

## Routes

### Server-Side Rendered Routes

| Method | Route              | Description             |
| ------ | ------------------ | ----------------------- |
| GET    | `/`                | Home page               |
| GET    | `/login`           | Login page              |
| GET    | `/logout`          | Logout action           |
| GET    | `/countries`       | Country explorer page   |
| GET    | `/countries/:slug` | Destination detail page |
| GET    | `/wishlist`        | User wishlist page      |
| GET    | `/dashboard`       | Travel dashboard page   |

### JSON API Routes

| Method | Route                    | Description                                 |
| ------ | ------------------------ | ------------------------------------------- |
| GET    | `/api/countries`         | Returns country list for AJAX search/filter |
| GET    | `/api/countries/:slug`   | Returns one country by slug                 |
| GET    | `/api/wishlist`          | Returns authenticated user's wishlist       |
| POST   | `/api/wishlist`          | Creates a wishlist item                     |
| PUT    | `/api/wishlist/:id`      | Updates note/status of a wishlist item      |
| DELETE | `/api/wishlist/:id`      | Deletes a wishlist item                     |
| GET    | `/api/dashboard/summary` | Returns dashboard counters                  |

---

## Slug Format

Country detail pages use lowercase slug values generated from the country name.

Examples:

```txt
/countries/bangladesh
/countries/japan
/countries/united-states
/countries/new-zealand
```

Slugs are generated by lowercasing the country name, removing selected punctuation, and replacing spaces with hyphens.

---

## Wishlist Storage Approach

This project uses an in-memory wishlist store.

No database, ORM, SQLite, MySQL, PostgreSQL, or external storage is used.

Wishlist data is managed through the service layer and exposed through JSON API endpoints.

Each wishlist item contains:

| Field          | Description                           |
| -------------- | ------------------------------------- |
| `id`           | Unique wishlist item ID               |
| `username`     | Owner of the wishlist item            |
| `country_name` | Destination country name              |
| `note`         | Optional editable note                |
| `status`       | Travel status: `Planned` or `Visited` |
| `created_at`   | Creation timestamp                    |

Because storage is in-memory, wishlist data resets when the server restarts.

---

## AJAX Behavior

The project uses vanilla JavaScript and the Fetch API for partial page updates.

AJAX actions do not reload the full page.

| Page               | Container             | Behavior                                    |
| ------------------ | --------------------- | ------------------------------------------- |
| `/`                | `#search-suggestions` | Updates destination suggestions             |
| `/countries`       | `#country-results`    | Updates country cards after search/filter   |
| `/countries/:slug` | `#wishlist-feedback`  | Shows add-to-wishlist success/error message |
| `/wishlist`        | `#wishlist-rows`      | Updates wishlist rows after edit/delete     |
| `/dashboard`       | `#dashboard-stats`    | Updates dashboard counters                  |

---

## External API Integration

### REST Countries

Used for:

* Country name
* Official name
* Capital
* Region
* Subregion
* Population
* Flag
* Currency
* Languages
* Latitude and longitude

Default base URL:

```txt
https://api.restcountries.com/countries/v5
```

Required config:

```ini
REST_COUNTRIES_API_KEY = your_rest_countries_api_key_here
```

The client sends `Authorization: Bearer <REST_COUNTRIES_API_KEY>` and reads the v5 paginated `data.objects` response.

### OpenTripMap

Used for nearby attractions, museums, landmarks, historical sites, architecture, and natural places.

Default base URL:

```txt
https://api.opentripmap.com/0.1/en/places
```

Required config:

```ini
OPENTRIPMAP_API_KEY = your_opentripmap_api_key_here
```

---

## Architecture Notes

### Controllers

Controllers handle HTTP requests, template rendering, request parsing, redirects, and JSON responses.

Examples:

* `controllers/country.go`
* `controllers/wishlist.go`
* `controllers/api/wishlist.go`

### Services

Services contain business logic and coordinate between controllers, models, and utilities.

Examples:

* `CountryService`
* `AttractionService`
* `WishlistService`
* `DashboardService`

### Models

Models define application data structures.

Examples:

* `Country`
* `Attraction`
* `WishlistItem`
* `ErrorResponse`
* `SuccessResponse`

### Utilities

Utilities contain reusable helpers and API clients.

Examples:

* REST Countries client
* OpenTripMap client
* Validation helpers
* Formatting helpers
* Error response helpers

### Middleware / Filters

The router registers reusable filters for:

* Request logging
* Authentication protection for wishlist and dashboard routes

The logging filter records request method, URL path, and request duration.

---

## Testing

Run all tests with:

```bash
go test ./...
```

---


## Error Handling

The application handles common error cases such as:

* Invalid wishlist payloads.
* Missing country names.
* Invalid wishlist status values.
* Unknown country slugs.
* Missing or invalid API keys.
* External API errors.
* Invalid external API responses.

API routes return JSON error responses with a message and status code.

SSR routes show user-friendly template-based error states where applicable.

---

## Configuration Example

`conf/app.example.conf` should include:

```ini
appname = TravelSphere
httpport = 8080
runmode = dev

sessionon = true
copyrequestbody = true

DEMO_USERNAME = user
DEMO_PASSWORD = admin

REST_COUNTRIES_BASE_URL = https://api.restcountries.com/countries/v5
REST_COUNTRIES_API_KEY = your_rest_countries_api_key_here
OPENTRIPMAP_BASE_URL = https://api.opentripmap.com/0.1/en/places
OPENTRIPMAP_API_KEY = your_opentripmap_api_key_here
```

> Do not commit real API keys to the repository.

---

## Known Limitations

* Wishlist data is stored in memory and resets when the server restarts.
* REST Countries and OpenTripMap require valid API keys.
* Weather API integration is not included unless added as a bonus feature.
* This project is designed for assessment purposes and does not use persistent user accounts or a database.

---

## Assessment Checklist

* Beego MVC architecture implemented.
* SSR routes and JSON API routes separated.
* REST Countries API integrated.
* OpenTripMap API integrated.
* Wishlist CRUD implemented through JSON API.
* No database used.
* AJAX partial updates implemented.
* `Prepare()` used for shared template data.
* Authentication and logging filters implemented.
* Layout, header partial, and footer partial used.
* Unit tests included.
* README includes setup, configuration, run, and test instructions.
