// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/games": {
            "get": {
                "description": "returns paginated games",
                "produces": [
                    "application/json"
                ],
                "summary": "Get games",
                "operationId": "get-games",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "default",
                            "name",
                            "releaseDate"
                        ],
                        "type": "string",
                        "description": "order by",
                        "name": "orderBy",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name filter",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handler.GameResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "creates new game",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create game",
                "operationId": "create-game",
                "parameters": [
                    {
                        "description": "create game",
                        "name": "game",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateGameRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/handler.IDResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/games/count": {
            "get": {
                "description": "returns games count",
                "produces": [
                    "application/json"
                ],
                "summary": "Get games count",
                "operationId": "get-games-count",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name filter",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handler.CountResponse"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/games/{id}": {
            "get": {
                "description": "returns game by ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get game",
                "operationId": "get-game-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Game ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GameResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "deletes game by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delete game",
                "operationId": "delete-game",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Game ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "updates game by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update game",
                "operationId": "update-game",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Game ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "update game",
                        "name": "game",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.UpdateGameRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/games/{id}/rate": {
            "post": {
                "description": "rates game",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Rate game",
                "operationId": "rate-game",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "game ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "game rating",
                        "name": "rating",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.CreateRatingRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.RatingResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/genres": {
            "get": {
                "description": "returns all genres",
                "produces": [
                    "application/json"
                ],
                "summary": "Get genres",
                "operationId": "get-genres",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handler.Genre"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/platforms": {
            "get": {
                "description": "returns all platforms",
                "produces": [
                    "application/json"
                ],
                "summary": "Get platforms",
                "operationId": "get-platforms",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handler.Platform"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user/ratings": {
            "post": {
                "description": "returns user ratings for specified games",
                "produces": [
                    "application/json"
                ],
                "summary": "Get user ratings for specified games",
                "operationId": "get-user-ratings",
                "parameters": [
                    {
                        "description": "games ids",
                        "name": "gameIds",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.GetUserRatingsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "integer"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/web.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.Company": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.CountResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                }
            }
        },
        "handler.CreateGameRequest": {
            "type": "object",
            "required": [
                "developer",
                "name"
            ],
            "properties": {
                "developer": {
                    "type": "string"
                },
                "genre": {
                    "description": "Deprecated: use Genres instead",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "genresIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "logoUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "platformsIDs": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "releaseDate": {
                    "type": "string"
                },
                "screenshots": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "slug": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "websites": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.CreateRatingRequest": {
            "type": "object",
            "properties": {
                "rating": {
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 1
                }
            }
        },
        "handler.GameResponse": {
            "type": "object",
            "properties": {
                "developer": {
                    "description": "Deprecated: use Developers instead",
                    "type": "string"
                },
                "developers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Company"
                    }
                },
                "genre": {
                    "description": "Deprecated: use Genres instead",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Genre"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "logoUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "platforms": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Platform"
                    }
                },
                "publisher": {
                    "description": "Deprecated: use Publishers instead",
                    "type": "string"
                },
                "publishers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/handler.Company"
                    }
                },
                "rating": {
                    "type": "number"
                },
                "releaseDate": {
                    "type": "string"
                },
                "screenshots": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "slug": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "websites": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.Genre": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.GetUserRatingsRequest": {
            "type": "object",
            "properties": {
                "gameIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "handler.IDResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "handler.Platform": {
            "type": "object",
            "properties": {
                "abbreviation": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handler.RatingResponse": {
            "type": "object",
            "properties": {
                "gameId": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "handler.UpdateGameRequest": {
            "type": "object",
            "properties": {
                "developer": {
                    "type": "string"
                },
                "genresIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "logoUrl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "platforms": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "releaseDate": {
                    "type": "string"
                },
                "screenshots": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "slug": {
                    "type": "string"
                },
                "summary": {
                    "type": "string"
                },
                "websites": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "web.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/web.FieldError"
                    }
                }
            }
        },
        "web.FieldError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "field": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.2",
	Host:             "localhost:8000",
	BasePath:         "/api",
	Schemes:          []string{"http"},
	Title:            "Game library API",
	Description:      "API for game library service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
