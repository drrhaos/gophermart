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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/user/balance": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Этот эндпоинт для получение текущего баланса пользователя",
                "consumes": [
                    "application/json"
                ],
                "summary": "Получение текущего баланса пользователя",
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/balance/withdraw": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Этот эндпоинт на списание средств",
                "consumes": [
                    "application/json"
                ],
                "summary": "Запрос на списание средств",
                "parameters": [
                    {
                        "description": "JSON тело запроса",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.BalanceWithdrawn"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "неверный номер заказа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "description": "Этот эндпоинт производит аутентификацию пользователя",
                "consumes": [
                    "application/json"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "JSON тело запроса",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "пользователь успешно аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "неверная пара логин/пароль",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/orders": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Этот эндпоинт для получения списка загруженных номеров заказов",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение списка загруженных номеров заказов",
                "responses": {
                    "200": {
                        "description": "успешная обработка запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "204": {
                        "description": "нет данных для ответа.",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не авторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Этот эндпоинт загружает номера заказа",
                "consumes": [
                    "text/plain"
                ],
                "summary": "Загрузка номера заказа",
                "parameters": [
                    {
                        "description": "номер заказа",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "номер заказа уже был загружен этим пользователем",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "202": {
                        "description": "новый номер заказа принят в обработку",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "пользователь не аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "номер заказа уже был загружен другим пользователем",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "422": {
                        "description": "неверный формат номера заказа",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "description": "Этот эндпоинт производит регистрацию пользователя",
                "consumes": [
                    "application/json"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "JSON тело запроса",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "пользователь успешно аутентифицирован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "неверный формат запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "409": {
                        "description": "логин уже занят",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/withdrawals": {
            "get": {
                "description": "Этот эндпоинт для получение информации о выводе средств",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение информации о выводе средств",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.BalanceWithdrawn": {
            "type": "object",
            "properties": {
                "order": {
                    "type": "string"
                },
                "sum": {
                    "type": "number"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "login": {
                    "description": "логин",
                    "type": "string"
                },
                "password": {
                    "description": "параметр, принимающий значение gauge или counter",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
