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
                "description": "Этот эндпоинт для получение текущего баланса пользователя",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение текущего баланса пользователя",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/user/balance/withdraw": {
            "post": {
                "description": "Этот эндпоинт на списание средств",
                "produces": [
                    "application/json"
                ],
                "summary": "Запрос на списание средств",
                "responses": {
                    "200": {
                        "description": "OK",
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
                "produces": [
                    "application/json"
                ],
                "summary": "Аутентификация пользователя",
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
                    "404": {
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
                "description": "Этот эндпоинт для получения списка загруженных номеров заказов",
                "produces": [
                    "application/json"
                ],
                "summary": "Получение списка загруженных номеров заказов",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Этот эндпоинт загружает номера заказа",
                "produces": [
                    "application/json"
                ],
                "summary": "Загрузка номера заказа",
                "responses": {
                    "200": {
                        "description": "OK",
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
                "produces": [
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
                            "$ref": "#/definitions/handlers.User"
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
                    "404": {
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
        "handlers.User": {
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
