// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
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
        "/games": {
            "get": {
                "description": "returns paginated games with extended properties",
                "produces": [
                    "application/json"
                ],
                "summary": "List games info",
                "operationId": "get-all-games-info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "last fetched Id",
                        "name": "lastId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/game.GameInfoResp"
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
                            "$ref": "#/definitions/game.CreateGame"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/game.GameResp"
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
        "/games/rate": {
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
                        "description": "game rating",
                        "name": "rating",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/game.CreateRating"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/game.Rating"
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
        "/games/{id}": {
            "get": {
                "description": "returns game with extended properties by ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get game info",
                "operationId": "get-game-info-by-id",
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
                            "$ref": "#/definitions/game.GameInfoResp"
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
                "operationId": "delete-game-by-id",
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
                        "description": ""
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
                "operationId": "update-game-by-id",
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
                            "$ref": "#/definitions/game.UpdateGame"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/game.GameResp"
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
        "/games/{id}/sales": {
            "get": {
                "description": "returns sales for specified game",
                "produces": [
                    "application/json"
                ],
                "summary": "List game sales",
                "operationId": "get-game-sales-by-game-id",
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
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/game.GameSaleResp"
                            }
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
            "post": {
                "description": "adds game on sale",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add game on sale",
                "operationId": "add-game-on-sale",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Game ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "game sale",
                        "name": "gamesale",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/game.CreateGameSale"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/game.GameSaleResp"
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
        "/sales": {
            "get": {
                "description": "Returns all sales",
                "produces": [
                    "application/json"
                ],
                "summary": "List all sales",
                "operationId": "get-sales",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/game.SaleResp"
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
                "description": "Creates new sale",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create sale",
                "operationId": "create-sale",
                "parameters": [
                    {
                        "description": "create sale",
                        "name": "sale",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/game.CreateSale"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/game.SaleResp"
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
        "/user/ratings": {
            "post": {
                "description": "returns user ratings for specified games",
                "produces": [
                    "application/json"
                ],
                "summary": "Get user ratings for specified games",
                "operationId": "get-user-ratings",
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
        "game.CreateGame": {
            "type": "object",
            "required": [
                "developer",
                "name",
                "publisher"
            ],
            "properties": {
                "developer": {
                    "type": "string"
                },
                "genre": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "publisher": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                }
            }
        },
        "game.CreateGameSale": {
            "type": "object",
            "properties": {
                "discountPercent": {
                    "type": "integer"
                },
                "saleId": {
                    "type": "integer"
                }
            }
        },
        "game.CreateRating": {
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
        "game.CreateSale": {
            "type": "object",
            "required": [
                "beginDate",
                "endDate"
            ],
            "properties": {
                "beginDate": {
                    "type": "string"
                },
                "endDate": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "game.GameInfoResp": {
            "type": "object",
            "properties": {
                "currentPrice": {
                    "type": "number"
                },
                "developer": {
                    "type": "string"
                },
                "genre": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "publisher": {
                    "type": "string"
                },
                "rating": {
                    "type": "number"
                },
                "releaseDate": {
                    "type": "string"
                }
            }
        },
        "game.GameResp": {
            "type": "object",
            "properties": {
                "developer": {
                    "type": "string"
                },
                "genre": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "publisher": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                }
            }
        },
        "game.GameSaleResp": {
            "type": "object",
            "properties": {
                "beginDate": {
                    "type": "string"
                },
                "discountPercent": {
                    "type": "integer"
                },
                "endDate": {
                    "type": "string"
                },
                "gameId": {
                    "type": "integer"
                },
                "sale": {
                    "type": "string"
                },
                "saleId": {
                    "type": "integer"
                }
            }
        },
        "game.Rating": {
            "type": "object",
            "properties": {
                "gameID": {
                    "type": "integer"
                },
                "rating": {
                    "type": "integer"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "game.SaleResp": {
            "type": "object",
            "properties": {
                "beginDate": {
                    "type": "string"
                },
                "endDate": {
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
        "game.UpdateGame": {
            "type": "object",
            "properties": {
                "developer": {
                    "type": "string"
                },
                "genre": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "publisher": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
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

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "localhost:8000",
	BasePath:    "/api",
	Schemes:     []string{"http"},
	Title:       "Game library API",
	Description: "API for game library service",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
