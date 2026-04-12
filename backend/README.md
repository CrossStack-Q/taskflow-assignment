# API Documentation

> Base URL: `/api/v1`

> Authorization : BearerToken `<token>`

---

## Project API

### 1. Get All Projects

**Endpoint:** `GET /projects`

**Request Body:** None

**Query Params:**

| Parameter | Type   | Values                          |
|-----------|--------|---------------------------------|
| `page`    | number |                                 |
| `limit`   | number |                                 |
| `sort`    | string | `created_at` `updated_at` `name`|
| `order`   | string | `asc` `desc`                    |
| `search`  | string |                                 |

**Response Body:**
```json
{
  "data": [
    {
      "id": "uuid",
      "userId": "string",
      "name": "string",
      "color": "string",
      "description": "string | null",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ],
  "total": 0,
  "page": 1,
  "limit": 10,
  "totalPages": 0
}
```

---

### 2. Create Project

**Endpoint:** `POST /projects`

**Request Body:**
```json
{
  "name": "string",
  "color": "string",
  "description": "string | null"
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "userId": "string",
  "name": "string",
  "color": "string",
  "description": "string | null",
  "createdAt": "string",
  "updatedAt": "string"
}
```

---

### 3. Update Project

**Endpoint:** `PATCH /projects/:id`

**Request Body:**
```json
{
  "name": "string",
  "color": "string",
  "description": "string | null"
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "userId": "string",
  "name": "string",
  "color": "string",
  "description": "string | null",
  "createdAt": "string",
  "updatedAt": "string"
}
```

---

### 4. Delete Project

**Endpoint:** `DELETE /projects/:id`

**Request Body:** None

**Response Body:**
```json
{
  "status": 204,
  "message": "No Content"
}
```

---

## Todo API

### 1. Get Todo by ID

**Endpoint:** `GET /tasks/:id`

**Request Body:** None

**Response Body:**
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string | null",
  "status": "draft | active | completed | archived",
  "priority": "low | medium | high",
  "projectId": "uuid | null",
  "parentTodoId": "uuid | null",
  "sortOrder": 0,
  "dueDate": "string | null",
  "completedAt": "string | null",
  "createdAt": "string",
  "updatedAt": "string",
  "userId": "string",
  "metadata": {},
  "category": {},
  "children": [],
  "comments": []
}
```

---

### 2. Create Todo

**Endpoint:** `POST /tasks`

**Request Body:**
```json
{
  "title": "string",
  "description": "string | null",
  "projectId": "uuid | null",
  "parentTodoId": "uuid | null",
  "priority": "low | medium | high",
  "dueDate": "string | null",
  "metadata": {}
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string | null",
  "status": "draft | active | completed | archived",
  "priority": "low | medium | high",
  "projectId": "uuid | null",
  "parentTodoId": "uuid | null",
  "sortOrder": 0,
  "dueDate": "string | null",
  "completedAt": "string | null",
  "createdAt": "string",
  "updatedAt": "string",
  "userId": "string",
  "metadata": {}
}
```

---

### 3. Get All Todos

**Endpoint:** `GET /tasks`

**Request Body:** None

**Query Params:**

| Parameter     | Type      | Values                                                      |
|---------------|-----------|-------------------------------------------------------------|
| `page`        | number    |                                                             |
| `limit`       | number    |                                                             |
| `sort`        | string    | `created_at` `updated_at` `title` `priority` `due_date` `status` |
| `order`       | string    | `asc` `desc`                                                |
| `search`      | string    |                                                             |
| `status`      | string    | `draft` `active` `completed` `archived`                     |
| `priority`    | string    | `low` `medium` `high`                                       |
| `projectId`  | uuid      |                                                             |
| `parentTodoId`| uuid      |                                                             |
| `dueFrom`     | datetime  |                                                             |
| `dueTo`       | datetime  |                                                             |
| `overdue`     | boolean   |                                                             |
| `completed`   | boolean   |                                                             |

**Response Body:**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "description": "string | null",
      "status": "draft | active | completed | archived",
      "priority": "low | medium | high",
      "projectId": "uuid | null",
      "parentTodoId": "uuid | null",
      "sortOrder": 0,
      "dueDate": "string | null",
      "completedAt": "string | null",
      "createdAt": "string",
      "updatedAt": "string",
      "userId": "string",
      "metadata": {},
      "category": {},
      "attachments": [],
      "children": [],
      "comments": []
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 100,
  "totalPages": 10
}
```

---

### 4. Update Task

**Endpoint:** `PATCH /tasks/:id`

**Request Body:**
```json
{
  "title": "string",
  "description": "string | null",
  "projectId": "uuid | null",
  "parentTodoId": "uuid | null",
  "priority": "low | medium | high",
  "status": "draft | active | completed | archived",
  "dueDate": "string | null",
  "metadata": {
    "color": "string",
    "difficulty": 0,
    "reminder": "string",
    "tags": ["string"]
  }
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "title": "string",
  "description": "string | null",
  "status": "draft | active | completed | archived",
  "priority": "low | medium | high",
  "projectId": "uuid | null",
  "parentTodoId": "uuid | null",
  "sortOrder": 0,
  "dueDate": "string | null",
  "completedAt": "string | null",
  "createdAt": "string",
  "updatedAt": "string",
  "userId": "string",
  "metadata": {
    "color": "string",
    "difficulty": 0,
    "reminder": "string",
    "tags": ["string"]
  }
}
```

---

### 5. Delete Task

**Endpoint:** `DELETE /tasks/:id`

**Request Body:** None

**Response Body:**
```json
{
  "status": 204,
  "message": "No Content"
}
```

---

### 6. Get Todo Statistics

**Endpoint:** `GET /tasks/statistics`

**Request Body:** None

**Response Body:**
```json
{
  "active": 0,
  "archived": 0,
  "completed": 0,
  "draft": 0,
  "overdue": 0,
  "total": 0
}
```

---

## Comment API

### 1. Add Comment to Task

**Endpoint:** `POST /todos/:taskId/comments`

**Request Body:**
```json
{
  "content": "string"
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "taskId": "uuid",
  "userId": "string",
  "content": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

---

### 2. Get Comments for Task

**Endpoint:** `GET /todos/:taskId/comments`

**Request Body:** None

**Response Body:**
```json
[
  {
    "id": "uuid",
    "todoId": "uuid",
    "userId": "string",
    "content": "string",
    "createdAt": "string",
    "updatedAt": "string"
  }
]
```

---

### 3. Update Comment

**Endpoint:** `PATCH /todos/:taskId/comments/:commentId`

**Request Body:**
```json
{
  "content": "string"
}
```

**Response Body:**
```json
{
  "id": "uuid",
  "taskId": "uuid",
  "userId": "string",
  "content": "string",
  "createdAt": "string",
  "updatedAt": "string"
}
```

---

### 4. Delete Comment

**Endpoint:** `DELETE /todos/:taskId/comments/:commentId`

**Request Body:** None

**Response Body:**
```json
{
  "status": 204,
  "message": "No Content"
}
```

---
