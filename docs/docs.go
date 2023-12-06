// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/login": {
            "post": {
                "description": "Логин пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Вход пользователя",
                "parameters": [
                    {
                        "description": "добавить данные в тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.LoginPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/profile": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить инфо о пользователе",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Инфо о профиле",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.UserSwaggerJsonMap"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/roles": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить список всех ролей",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "roles"
                ],
                "summary": "Список ролей",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/roles.RoleSwaggerJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/secrets": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить список всех секретов",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Список секретов/ключей",
                "parameters": [
                    {
                        "type": "string",
                        "description": "список секретов по типу",
                        "name": "type",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "список секретов от пользователя",
                        "name": "user",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/secrets.SecretSwaggerJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создать секрет/ключ",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Создать секрет/ключ",
                "parameters": [
                    {
                        "description": "добавить данные в тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/secrets.SecretSwaggerReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Секрет сохранен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/secrets/:secret_id": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить ключ по id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Объект ключа",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID ключа",
                        "name": "secret_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/secrets.SecretSwaggerJsonObj"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Администратор может изменять все ключи, а пользователь только свои.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Редактировать ключ",
                "parameters": [
                    {
                        "description": "добавить данные в тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/secrets.SecretReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Секрет изменен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/secrets/share": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Поделиться секретом с пользователями",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "secrets"
                ],
                "summary": "Поделиться секретом с пользователями",
                "parameters": [
                    {
                        "description": "добавить данные в тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/secrets.UsersSecret"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Секрет сохранен",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Получить список всех пользователей",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Список пользователей",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.UserSwaggerJson"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создать пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Создать пользователя",
                "parameters": [
                    {
                        "description": "добавить данные в тело запроса",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.UserSwaggerRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь создан",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Err"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.LoginPayload": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "example": "kamilla.n"
                },
                "password": {
                    "type": "string",
                    "example": "password123"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDA3OTgzMzcsIm5hbWUiOiJiZWthcnlzIiwidXNlcl9pZCI6MX0.YoLwVoMnGvA7q44dFQJ4E4qBghr3zYDKbJNrhV1yrC0"
                }
            }
        },
        "resp.Err": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "internal server error"
                }
            }
        },
        "roles.Role": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string",
                    "example": "user"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                },
                "id": {
                    "type": "integer",
                    "example": 12
                },
                "name": {
                    "type": "string",
                    "example": "Пользователь"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                }
            }
        },
        "roles.RoleSwaggerData": {
            "type": "object",
            "properties": {
                "attributes": {
                    "$ref": "#/definitions/roles.Role"
                },
                "id": {
                    "type": "string",
                    "example": "11"
                },
                "type": {
                    "type": "string",
                    "example": "users"
                }
            }
        },
        "roles.RoleSwaggerJson": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/roles.RoleSwaggerData"
                    }
                }
            }
        },
        "secrets.Secret": {
            "type": "object",
            "required": [
                "author_id",
                "data",
                "key",
                "title"
            ],
            "properties": {
                "author_id": {
                    "type": "integer",
                    "example": 1
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                },
                "data": {
                    "type": "string",
                    "example": "localhost:5432"
                },
                "id": {
                    "type": "integer",
                    "example": 11
                },
                "key": {
                    "type": "string",
                    "example": "dbHost"
                },
                "title": {
                    "type": "string",
                    "example": "адрес базы данных"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "auth",
                        "ssh",
                        "env"
                    ],
                    "example": "env"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                }
            }
        },
        "secrets.SecretReq": {
            "type": "object",
            "required": [
                "author_id",
                "data",
                "key",
                "title"
            ],
            "properties": {
                "author_id": {
                    "type": "integer",
                    "example": 1
                },
                "data": {
                    "type": "string",
                    "example": "localhost:5432"
                },
                "key": {
                    "type": "string",
                    "example": "dbHost"
                },
                "title": {
                    "type": "string",
                    "example": "адрес базы данных"
                },
                "type": {
                    "type": "string",
                    "enum": [
                        "auth",
                        "ssh",
                        "env"
                    ],
                    "example": "env"
                }
            }
        },
        "secrets.SecretSwaggerData": {
            "type": "object",
            "properties": {
                "attributes": {
                    "$ref": "#/definitions/secrets.Secret"
                },
                "id": {
                    "type": "string",
                    "example": "11"
                },
                "type": {
                    "type": "string",
                    "example": "secrets"
                }
            }
        },
        "secrets.SecretSwaggerJson": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/secrets.SecretSwaggerData"
                    }
                }
            }
        },
        "secrets.SecretSwaggerJsonObj": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/secrets.SecretSwaggerData"
                }
            }
        },
        "secrets.SecretSwaggerReq": {
            "type": "object",
            "required": [
                "author_id",
                "data",
                "key",
                "title"
            ],
            "properties": {
                "author_id": {
                    "type": "integer",
                    "example": 1
                },
                "data": {
                    "type": "string",
                    "example": "localhost:5432"
                },
                "key": {
                    "type": "string",
                    "example": "dbHost"
                },
                "title": {
                    "type": "string",
                    "example": "адрес базы данных"
                },
                "type": {
                    "type": "string",
                    "example": "env"
                }
            }
        },
        "secrets.UsersSecret": {
            "type": "object",
            "properties": {
                "secret_id": {
                    "type": "integer",
                    "example": 41
                },
                "user_ids": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    },
                    "example": [
                        13
                    ]
                }
            }
        },
        "users.UserResp": {
            "type": "object",
            "required": [
                "email",
                "name"
            ],
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                },
                "email": {
                    "type": "string",
                    "example": "bekarys.t@evrika.com"
                },
                "id": {
                    "type": "integer",
                    "example": 32
                },
                "is_active": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "bekarys"
                },
                "role": {
                    "$ref": "#/definitions/roles.Role"
                },
                "role_id": {
                    "type": "integer",
                    "example": 2
                },
                "updated_at": {
                    "type": "string",
                    "example": "2023-11-20T11:15:37Z"
                }
            }
        },
        "users.UserSwaggerData": {
            "type": "object",
            "properties": {
                "attributes": {
                    "$ref": "#/definitions/users.UserResp"
                },
                "id": {
                    "type": "string",
                    "example": "11"
                },
                "type": {
                    "type": "string",
                    "example": "users"
                }
            }
        },
        "users.UserSwaggerJson": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/users.UserSwaggerData"
                    }
                }
            }
        },
        "users.UserSwaggerJsonMap": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/users.UserSwaggerData"
                }
            }
        },
        "users.UserSwaggerRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "role_id"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "yelena.a@evrika.com"
                },
                "is_active": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "yelena.a"
                },
                "password": {
                    "type": "string",
                    "maxLength": 20,
                    "minLength": 7,
                    "example": "password123"
                },
                "role_id": {
                    "type": "integer",
                    "enum": [
                        1,
                        2
                    ],
                    "example": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "10.10.1.59:44044",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Evrika Secrets API",
	Description:      "Platform for managing application secrets and keys.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
