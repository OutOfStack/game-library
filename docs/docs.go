// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/companies/top": {
            "get": {
                "description": "returns top companies based on amount of games having it",
                "produces": [
                    "application/json"
                ],
                "summary": "Get top companies",
                "operationId": "get-top-companies",
                "parameters": [
                    {
                        "enum": [
                            "pub",
                            "dev"
                        ],
                        "type": "string",
                        "description": "company type (dev or pub)",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Company"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid or missing company type",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    }
                }
            }
        },
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
                    },
                    {
                        "type": "integer",
                        "description": "genre filter",
                        "name": "genre",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "developer id filter",
                        "name": "developer",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "publisher id filter",
                        "name": "publisher",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GamesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.CreateGameRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.IDResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GameResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.UpdateGameRequest"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.CreateRatingRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.RatingResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                                "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Genre"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/genres/top": {
            "get": {
                "description": "returns top genres based on amount of games having it",
                "produces": [
                    "application/json"
                ],
                "summary": "Get top genres",
                "operationId": "get-top-genres",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Genre"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                                "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Platform"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GetUserRatingsRequest"
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
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Company": {
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.CreateGameRequest": {
            "type": "object",
            "required": [
                "developer",
                "name"
            ],
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.CreateRatingRequest": {
            "type": "object",
            "properties": {
                "rating": {
                    "description": "0 - remove rating",
                    "type": "integer",
                    "maximum": 5,
                    "minimum": 0
                }
            }
        },
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GameResponse": {
            "type": "object",
            "properties": {
                "developers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Company"
                    }
                },
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Genre"
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
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Platform"
                    }
                },
                "publishers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Company"
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GamesResponse": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "games": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GameResponse"
                    }
                }
            }
        },
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Genre": {
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.GetUserRatingsRequest": {
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.IDResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.Platform": {
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.RatingResponse": {
            "type": "object",
            "properties": {
                "gameId": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                }
            }
        },
        "github_com_OutOfStack_game-library_internal_app_game-library-api_api_model.UpdateGameRequest": {
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
        "github_com_OutOfStack_game-library_internal_app_game-library-api_web.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/github_com_OutOfStack_game-library_internal_app_game-library-api_web.FieldError"
                    }
                }
            }
        },
        "github_com_OutOfStack_game-library_internal_app_game-library-api_web.FieldError": {
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
