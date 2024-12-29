// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "produces": [
        "application/json"
    ],
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Chaika",
            "email": "chaika.contact@yandex.ru"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/packets": {
            "post": {
                "description": "Add a new packet of products to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "packets"
                ],
                "summary": "Add packet",
                "parameters": [
                    {
                        "description": "Packet details",
                        "name": "packet",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.AddPacketRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.AddPacketResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/packets/search": {
            "get": {
                "description": "Search for packets",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "packets"
                ],
                "summary": "Search packet",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search query",
                        "name": "query",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.SearchPacketResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/products": {
            "get": {
                "description": "Get all products from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get all products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.GetAllProductsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update product details in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Update product",
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.UpdateProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.UpdateProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new product to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Add product",
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.AddProductRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.AddProductResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/products/{id}": {
            "get": {
                "description": "Get product details by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get product by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.GetProductByIDResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a product from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Delete product",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Product ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.DeleteProductResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ChaikaGoods_internal_handler_schemas.AddPacketRequest": {
            "description": "Запрос на добавление пакета",
            "type": "object",
            "properties": {
                "packet": {
                    "description": "Сведения о новом пакете",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ChaikaGoods_internal_models.Package"
                        }
                    ]
                },
                "packet_content": {
                    "description": "Содержимое пакета",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ChaikaGoods_internal_models.PackageContent"
                    }
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.AddPacketResponse": {
            "description": "Ответ на запрос на добавление пакета",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "packet_id": {
                    "description": "ID созданного пакета",
                    "type": "integer"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.AddProductRequest": {
            "description": "Запрос на добавление продукта",
            "type": "object",
            "properties": {
                "product": {
                    "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ProductSchema"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.AddProductResponse": {
            "description": "Ответ на запрос на добавление продукта",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.DeleteProductResponse": {
            "description": "Ответ на запрос на удаление продукта",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.ErrorResponse": {
            "description": "Represents a standard error response for the API",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Error code",
                    "type": "integer"
                },
                "message": {
                    "description": "Error message",
                    "type": "string"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.GetAllProductsResponse": {
            "description": "Ответ на запрос на получение всех продуктов",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ProductSchema"
                    }
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.GetProductByIDResponse": {
            "description": "Ответ на запрос на получение продукта по его ID",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "product": {
                    "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ProductSchema"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.ProductSchema": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "imageurl": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.SearchPacketResponse": {
            "description": "Ответ на запрос на поиск пакетов",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                },
                "packets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ChaikaGoods_internal_models.Package"
                    }
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.UpdateProductRequest": {
            "description": "Запрос на обновление продукта",
            "type": "object",
            "properties": {
                "product": {
                    "$ref": "#/definitions/ChaikaGoods_internal_handler_schemas.ProductSchema"
                }
            }
        },
        "ChaikaGoods_internal_handler_schemas.UpdateProductResponse": {
            "description": "Ответ на запрос на обновление продукта",
            "type": "object",
            "properties": {
                "err": {
                    "type": "string"
                }
            }
        },
        "ChaikaGoods_internal_models.Package": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "package_name": {
                    "type": "string"
                }
            }
        },
        "ChaikaGoods_internal_models.PackageContent": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "package_id": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.2",
	Host:             "127.0.0.1:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "ChaikaGoods API",
	Description:      "This is a simple API to manage goods for the Chaika app.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
