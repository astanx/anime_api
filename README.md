

# Anime API v1 — Routes

## **Users**

| Method | Route           | Query | Description         |
| ------ | --------------- | ----- | ------------------- |
| GET    | `/users/device` | —     | Add a new device ID |

---

## **Anime**

### **Consumet**

| Method | Route                            | Query   | Description                             |
| ------ | -------------------------------- | ------- | --------------------------------------- |
| GET    | `/anime/consumet/`               | `query` | Search anime by query                   |
| GET    | `/anime/consumet/genres`         | —       | Get all genres                          |
| GET    | `/anime/consumet/latest`         | —       | Get latest releases                     |
| GET    | `/anime/consumet/recommended`    | —       | Get recommended releases                |
| GET    | `/anime/consumet/genre/releases` | `genre` | Get releases by genre                   |
| GET    | `/anime/consumet/:id`            | —       | Get anime info by Consumet ID           |
| GET    | `/anime/consumet/episode/:id`    | —       | Get episode info by Consumet episode ID |

### **Anilibria**

| Method | Route                             | Query          | Description                              |
| ------ | --------------------------------- | -------------- | ---------------------------------------- |
| GET    | `/anime/anilibria/`               | `query`        | Search anime by query                    |
| GET    | `/anime/anilibria/genres`         | —              | Get all genres                           |
| GET    | `/anime/anilibria/latest`         | `limit`        | Get latest releases                      |
| GET    | `/anime/anilibria/random`         | `limit`        | Get random releases                      |
| GET    | `/anime/anilibria/recommended`    | `limit`        | Get recommended releases                 |
| GET    | `/anime/anilibria/genre/releases` | `genre, limit` | Get releases by genre ID                 |
| GET    | `/anime/anilibria/:id`            | —              | Get anime info by Anilibria ID           |
| GET    | `/anime/anilibria/episode/:id`    | —              | Get episode info by Anilibria episode ID |

### **General**

| Method | Route               | Query | Description                     |
| ------ | ------------------- | ----- | ------------------------------- |
| GET    | `/anime/search/:id` | —     | Get anime by internal search ID |

---

## **Timecode**

| Method | Route           | Query / Body           | Description                           |
| ------ | --------------- | ---------------------- | ------------------------------------- |
| GET    | `/timecode/all` | `deviceID`             | Get all timecodes for a device        |
| GET    | `/timecode`     | `deviceID, episodeID`  | Get a timecode for a specific episode |
| POST   | `/timecode`     | `deviceID` + JSON body | Add or update a timecode for a device |

---

## **History**

| Method | Route          | Query / Body            | Description                          |
| ------ | -------------- | ----------------------- | ------------------------------------ |
| POST   | `/history`     | `deviceID` + JSON body  | Add a history entry for a device     |
| GET    | `/history/all` | `deviceID`              | Get all history entries for a device |
| GET    | `/history`     | `deviceID, page, limit` | Get paginated history entries        |

---

## **Favourite**

| Method | Route            | Query / Body            | Description                             |
| ------ | ---------------- | ----------------------- | --------------------------------------- |
| POST   | `/favourite`     | `deviceID` + JSON body  | Add an anime to favourites for a device |
| DELETE | `/favourite`     | JSON body               | Remove an anime from favourites         |
| GET    | `/favourite/all` | `deviceID`              | Get all favourites for a device         |
| GET    | `/favourite`     | `deviceID, page, limit` | Get paginated favourites                |

---

## **Collection**

| Method | Route             | Query / Body            | Description                       |
| ------ | ----------------- | ----------------------- | --------------------------------- |
| POST   | `/collection`     | `deviceID` + JSON body  | Add an anime to a collection      |
| DELETE | `/collection`     | JSON body               | Remove an anime from a collection |
| GET    | `/collection/all` | `deviceID`              | Get all collections for a device  |
| GET    | `/collection`     | `deviceID, page, limit` | Get paginated collections         |

---

