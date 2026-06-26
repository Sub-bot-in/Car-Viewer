# HiveCars

HiveCars is a Go web application for browsing, filtering, viewing, and comparing cars from a local API.

The project uses Go's `net/http` package, HTML templates, static CSS, and data fetched from an external local API running on `localhost:3000`.

## Features

* Browse a catalog of cars
* View car details on a separate page
* Filter cars by:

  * manufacturer
  * category
  * year
* Pagination with 20 cars per page
* Preserve selected filters while navigating pages
* Compare up to 4 cars
* Add cars to comparison from:

  * catalog page
  * car details page
* Remove cars from comparison
* Clear the whole comparison list
* Show related cars from the same manufacturer
* Responsive layout with custom CSS

## Tech Stack

* Go
* HTML templates
* CSS
* Local JSON API
* Cookies for storing selected comparison cars

## Project Structure

```text
june_cars/
├── catalog/
│   └── catalog.go
├── internal/
│   ├── api/
│   │   └── fetchJSON.go
│   ├── handlers/
│   │   ├── homeHandler.go
│   │   ├── specs.go
│   │   └── compare.go
│   └── pagination/
│       └── pagination.go
├── static/
│   ├── css/
│   │   └── style.css
│   ├── fonts/
│   └── images/
├── templates/
│   ├── index.html
│   ├── car.html
│   └── compare.html
├── go.mod
└── main.go
```

## API

The application expects the cars API to be running at:

```text
http://localhost:3000/api

```

Used API endpoints:

```text
GET /api/models
GET /api/manufacturers
GET /api/categories
GET /api/images/{image}
```

## Running the Project

First, start the local API server on port `3000`.

To start the API, navigate to the api folder and run:

```bash
npm start
```

Then run the Go web server from the root of the project:

```bash
go run .
```

The application will start at:

```text
http://localhost:8080
```

## Main Pages

### Catalog Page

```text
/
```

The catalog page shows all cars with filters and pagination.

Users can:

* filter cars
* open car details
* add cars to comparison
* remove cars from comparison
* clear comparison

### Car Details Page

```text
/specifications/{id}
```

The car details page shows:

* car name
* year
* image
* basic information
* technical specifications
* manufacturer information
* related cars from the same manufacturer
* comparison actions

### Compare Page

```text
/compare
```

The compare page shows selected cars in a table.

Users can compare:

* manufacturer
* country
* category
* engine
* horsepower
* transmission
* drivetrain

A maximum of 4 cars can be compared at the same time.

## Comparison Logic

The selected cars are stored in a browser cookie.

Cookie name:

```text
compare_cars
```

When a user adds a car to comparison, its ID is saved in the cookie.

If the user already has 4 cars selected, the interface shows a message:

```text
You can compare up to 4 cars.
```

This message links to the comparison page.

## Filtering and Pagination

Filtering is handled before pagination.

The application first filters the full car list by:

* manufacturer
* category
* year

Then it applies pagination and shows 20 cars per page.

Pagination links preserve the selected filters.

Example:

```text
/?page=2&manufacturer=BMW&category=SUV&year=2020
```

## Related Cars

On the car details page, the application shows random related cars from the same manufacturer.

The current car is excluded from the related list.

## Notes

* The application uses server-side rendering with Go templates.
* Static files are served from the `/static/` path.
* The comparison feature does not require a database.
* The project is designed as a learning project for practicing Go web development, routing, templates, API fetching, filtering, pagination, and cookies.

## Possible Future Improvements

* Add search by car name
* Add sorting by year or horsepower
* Show a message when no cars match selected filters
* Improve mobile layout
* Add better error handling for API failures
* Add loading states with JavaScript
* Store comparison data in local storage instead of cookies
