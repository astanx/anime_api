# Anime API v1 — Routes

## **Users**

| Method | Route         | Query | Description         |
| ------ | ------------- | ----- | ------------------- |
| GET    | /users/device | —     | Add a new device ID |

---

## **Anime**

### **Consumet**

| Method | Route                          | Query | Description                             |
| ------ | ------------------------------ | ----- | --------------------------------------- |
| GET    | /anime/consumet/               | query | Search anime by query                   |
| GET    | /anime/consumet/genres         | —     | Get all genres                          |
| GET    | /anime/consumet/latest         | limit | Get latest releases                     |
| GET    | /anime/consumet/recommended    | limit | Get recommended releases                |
| GET    | /anime/consumet/genre/releases | genre | Get releases by genre                   |
| GET    | /anime/consumet/:id            | —     | Get anime info by Consumet ID           |
| GET    | /anime/consumet/episode/:id    | —     | Get episode info by Consumet episode ID |

### **Anilibria**

| Method | Route                           | Query          | Description                              |
| ------ | ------------------------------- | -------------- | ---------------------------------------- |
| GET    | /anime/anilibria/               | query          | Search anime by query                    |
| GET    | /anime/anilibria/genres         | —              | Get all genres                           |
| GET    | /anime/anilibria/latest         | limit          | Get latest releases                      |
| GET    | /anime/anilibria/random         | limit          | Get random releases                      |
| GET    | /anime/anilibria/recommended    | limit          | Get recommended releases                 |
| GET    | /anime/anilibria/genre/releases | genreId, limit | Get releases by genre ID                 |
| GET    | /anime/anilibria/:id            | —              | Get anime info by Anilibria ID           |
| GET    | /anime/anilibria/episode/:id    | —              | Get episode info by Anilibria episode ID |

### **General**

| Method | Route             | Query | Description                     |
| ------ | ----------------- | ----- | ------------------------------- |
| GET    | /anime/search/:id | —     | Get anime by internal search ID |

---
