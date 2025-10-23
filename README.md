# Anime API

Anime API is a RESTful API built with Go and the Gin framework, designed to manage anime-related data, including user favorites, watch history, timecodes, and collections, while integrating with external anime data sources like Zoro and Anilibria.

## Features

- **Device Management**: Generate and manage unique device IDs for user tracking.
- **Anime Search**: Query anime data from Zoro and Anilibria, including search by ID, genres, latest releases, and recommended anime.
- **Timecode Tracking**: Store and retrieve timecodes for specific episodes per user.
- **Watch History**: Record and retrieve users' anime watch history.
- **Favourites**: Allow users to add, remove, and view favorite anime.
- **Collections**: Manage user-defined collections of anime.
- **Middleware**: Logging and device-based authentication for secure access.

### Public URL

The API is publicly hosted and available at:
https://anime-api-rsc7.onrender.com

All endpoints below should be prefixed with this base URL.
Absolutely! Since your API uses a **device-based authentication** system via the `DeviceMiddleware`, we can clearly document the **auth flow**, how to provide the device ID, and how it interacts with requests. Here's a polished version you can add to your API docs:

---

## Authentication

### Overview

The API uses **device-based authentication**. Each user/device is identified by a unique **device ID**, which is required for most endpoints except `/users/device`. This allows tracking user-specific data such as favorites, watch history, and timecodes without traditional user accounts.

### Obtaining a Device ID

1. **Endpoint:** `POST /users/device`
2. **Description:** Generates a new device ID for the user.
3. **Request Body:** None
4. **Response:**

```json
{
  "id": "<device-uuid>"
}
```

5. **Errors:**

   * `400 Bad Request`: If device creation fails.

### Using the Device ID

* Include the device ID in the `Authorization` header with the prefix `Device `.
* Example:

```
Authorization: Device 123e4567-e89b-12d3-a456-426614174000
```

### Header Requirements

| Header        | Value Format         | Required                                      |
| ------------- | -------------------- | --------------------------------------------- |
| Authorization | `Device <device-id>` | For all endpoints except `/users/device` |

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Device Management

- **POST /users/device**
  - Description: Generates a new device ID for the user.
  - Response: `200 OK` with JSON `{ "id": "<device-uuid>" }`
  - Errors:
    - `400 Bad Request`: If device creation fails.

### Anime Routes

Anime routes are split into two groups: `/anime/consumet` and `/anime/anilibria`.

#### Consumet Routes

- **GET /anime/consumet/**
  - Description: Search for anime by query.
  - Query Parameters: `query` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid query.
- **GET /anime/consumet/genres**
  - Description: Retrieve available genres.
  - Response: `200 OK` with JSON `{ "results": [genres] }`
  - Errors:
    - `400 Bad Request`: Failed to fetch genres.
- **GET /anime/consumet/latest**
  - Description: Get latest anime releases.
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Failed to fetch releases.
- **GET /anime/consumet/recommended**
  - Description: Get recommended anime.
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Failed to fetch recommendations.
- **GET /anime/consumet/genre/releases**
  - Description: Search anime by genre.
  - Query Parameters: `genre` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid genre.
- **GET /anime/consumet/:id**
  - Description: Get anime details by Consumet ID.
  - Path Parameters: `id` (required)
  - Response: `200 OK` with JSON `{ "result": anime }`
  - Errors:
    - `400 Bad Request`: Missing or invalid ID.
- **GET /anime/consumet/episode/:id**
  - Description: Get episode info by Consumet ID.
  - Path Parameters: `id` (required)
  - Query Parameters: `title` (optional), `ordinal` (optional)
  - Response: `200 OK` with JSON `{ "result": episode }`
  - Errors:
    - `400 Bad Request`: Missing or invalid ID or parameters.

#### Anilibria Routes

- **GET /anime/anilibria/**
  - Description: Search for anime by query.
  - Query Parameters: `query` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid query.
- **GET /anime/anilibria/genres**
  - Description: Retrieve available genres.
  - Response: `200 OK` with JSON `{ "results": [genres] }`
  - Errors:
    - `400 Bad Request`: Failed to fetch genres.
- **GET /anime/anilibria/latest**
  - Description: Get latest anime releases.
  - Query Parameters: `limit` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid limit.
- **GET /anime/anilibria/random**
  - Description: Get random anime releases.
  - Query Parameters: `limit` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid limit.
- **GET /anime/anilibria/recommended**
  - Description: Get recommended anime.
  - Query Parameters: `limit` (optional, default: 14)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Invalid limit.
- **GET /anime/anilibria/genre/releases**
  - Description: Search anime by genre.
  - Query Parameters: `genre` (required), `limit` (required)
  - Response: `200 OK` with JSON `{ "results": [anime] }`
  - Errors:
    - `400 Bad Request`: Missing or invalid genre or limit.
- **GET /anime/anilibria/:id**
  - Description: Get anime details by Anilibria ID.
  - Path Parameters: `id` (required)
  - Response: `200 OK` with JSON `{ "result": anime }`
  - Errors:
    - `400 Bad Request`: Missing or invalid ID.
- **GET /anime/anilibria/episode/:id**
  - Description: Get episode info by Anilibria ID.
  - Path Parameters: `id` (required)
  - Response: `200 OK` with JSON `{ "result": episode }`
  - Errors:
    - `400 Bad Request`: Missing or invalid ID.
- **GET /anime/search/:id**
  - Description: Search anime by ID (generic).
  - Path Parameters: `id` (required)
  - Response: `200 OK` with JSON `{ "result": anime }`
  - Errors:
    - `400 Bad Request`: Missing or invalid ID.

### Timecode Routes

Requires `DeviceMiddleware` for authentication.

- **GET /timecode**
  - Description: Get timecode for a specific episode.
  - Query Parameters: `episodeID` (required)
  - Response: `200 OK` with JSON timecode object or `404 Not Found` if not found.
  - Errors:
    - `400 Bad Request`: Missing deviceID or episodeID.
    - `500 Internal Server Error`: Failed to fetch timecode.
- **GET /timecode/anime**
  - Description: Get timecodes for a specific anime.
  - Query Parameters: `animeID` (required)
  - Response: `200 OK` with JSON array of timecodes or `404 Not Found` if none found.
  - Errors:
    - `400 Bad Request`: Missing deviceID or animeID.
    - `500 Internal Server Error`: Failed to fetch timecodes.
- **GET /timecode/all**
  - Description: Get all timecodes for the user.
  - Response: `200 OK` with JSON array of timecodes.
  - Errors:
    - `400 Bad Request`: Missing deviceID.
    - `500 Internal Server Error`: Failed to fetch timecodes.
- **POST /timecode**
  - Description: Add or update a timecode.
  - Body: JSON `Timecode` object
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid request body.
    - `500 Internal Server Error`: Failed to add/update timecode.

### History Routes

Requires `DeviceMiddleware` for authentication.

- **POST /history**
  - Description: Add a history entry.
  - Body: JSON `History` object
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid request body.
    - `500 Internal Server Error`: Failed to add history.
- **GET /history**
  - Description: Get paginated history.
  - Query Parameters: `page` (optional, default: 1), `limit` (optional, default: 10)
  - Response: `200 OK` with JSON array of history entries.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid pagination parameters.
    - `500 Internal Server Error`: Failed to fetch history.
- **GET /history/all**
  - Description: Get all history entries for the user.
  - Response: `200 OK` with JSON array of history entries.
  - Errors:
    - `400 Bad Request`: Missing deviceID.
    - `500 Internal Server Error`: Failed to fetch history.

### Favourite Routes

Requires `DeviceMiddleware` for authentication.

- **POST /favourite**
  - Description: Add a favorite anime.
  - Body: JSON `Favourite` object
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid request body.
    - `500 Internal Server Error`: Failed to add favorite.
- **DELETE /favourite**
  - Description: Remove a favorite anime.
  - Body: JSON `Favourite` object
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid request body.
    - `500 Internal Server Error`: Failed to remove favorite.
- **GET /favourite**
  - Description: Get paginated favorites.
  - Query Parameters: `page` (optional, default: 1), `limit` (optional, default: 10)
  - Response: `200 OK` with JSON array of favorites.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid pagination parameters.
    - `500 Internal Server Error`: Failed to fetch favorites.
- **GET /favourite/all**
  - Description: Get all favorite anime for the user.
  - Response: `200 OK` with JSON array of favorites.
  - Errors:
    - `400 Bad Request`: Missing deviceID.
    - `500 Internal Server Error`: Failed to fetch favorites.

### Collection Routes

Requires `DeviceMiddleware` for authentication.

- **POST /collection**
  - Description: Add a collection.
  - Body: JSON `Collection` object
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID or invalid request body.
    - `500 Internal Server Error`: Failed to add collection.
- **DELETE /collection**
  - Description: Remove a collection.
  - Body: JSON `{ "anime_id": string, "type": string }`
  - Response: `204 No Content` on success.
  - Errors:
    - `400 Bad Request`: Missing deviceID, anime_id, or type.
    - `500 Internal Server Error`: Failed to remove collection.
- **GET /collection**
  - Description: Get paginated collections.
  - Query Parameters: `type` (required), `page` (optional, default: 1), `limit` (optional, default: 10)
  - Response: `200 OK` with JSON array of collections.
  - Errors:
    - `400 Bad Request`: Missing deviceID, type, or invalid pagination parameters.
    - `500 Internal Server Error`: Failed to fetch collections.
- **GET /collection/all**
  - Description: Get all collections for the user.
  - Response: `200 OK` with JSON array of collections.
  - Errors:
    - `400 Bad Request`: Missing deviceID.
    - `500 Internal Server Error`: Failed to fetch collections.

## Error Handling

- **400 Bad Request**: Returned for missing or invalid parameters.
- **404 Not Found**: Returned when resources (e.g., timecodes) are not found.
- **500 Internal Server Error**: Returned for server-side errors.
- Error responses include a JSON object with an `error` field describing the issue.
