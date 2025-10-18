# anime_api

---

# Anime API v1 — Routes

## **Users**

| Method | Route         | Query | Description         |
| ------ | ------------- | ----- | ------------------- |
| GET    | /users/device | —     | Add a new device ID |

---

## **Anime**

### **Consumet**

| Method | Route                          | Query | Description           |
| ------ | ------------------------------ | ----- | --------------------- |
| GET    | /anime/consumet/               | query | Search anime by query |
| GET    | /anime/consumet/genres         | —     | Get all genres        |
| GET    | /anime/consumet/latest         | limit | Get latest releases   |
| GET    | /anime/consumet/genre/releases | genre | Get releases by genre |

### **Anilibria**

| Method | Route                           | Query          | Description              |
| ------ | ------------------------------- | -------------- | ------------------------ |
| GET    | /anime/anilibria/               | query          | Search anime by query    |
| GET    | /anime/anilibria/genres         | —              | Get all genres           |
| GET    | /anime/anilibria/latest         | limit          | Get latest releases      |
| GET    | /anime/anilibria/random         | limit          | Get random releases      |
| GET    | /anime/anilibria/genre/releases | genreId, limit | Get releases by genre ID |

---
