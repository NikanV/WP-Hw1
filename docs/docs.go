// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/auth/reqdh": {
            "get": {
                "description": "second step of registeration which we send auth info and public keys.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Request DH",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The message ID (even and greater than zero).",
                        "name": "message_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The nonce (20 characters long).",
                        "name": "nonce",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The server_nonce (20 characters long).",
                        "name": "server_nonce",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "public key from client",
                        "name": "a",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WP-Hw1_proto.DHResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/reqpq": {
            "get": {
                "description": "first step of registeration which we send user info.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Request PQ",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "The message ID (even and greater than zero).",
                        "name": "message_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The nonce (20 characters long).",
                        "name": "nonce",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WP-Hw1_proto.PQResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/biz/getusers": {
            "get": {
                "description": "after checking authentication , gets the information that you desire.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Biz server"
                ],
                "summary": "get users of database.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The nonce (20 characters long).",
                        "name": "nonce",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "The server_nonce (20 characters long).",
                        "name": "server_nonce",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "gets first 100 users if negetive",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "auth key",
                        "name": "auth_key",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "The message ID (even and greater than zero).",
                        "name": "message_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/WP-Hw1_proto.GetUsersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "WP-Hw1_proto.DHResponse": {
            "type": "object",
            "properties": {
                "b": {
                    "description": "public key generated by server",
                    "type": "integer"
                },
                "message_id": {
                    "description": "odd and greater than zero",
                    "type": "integer"
                },
                "nonce": {
                    "description": "the exact nonce from the clients request",
                    "type": "string"
                },
                "server_nonce": {
                    "description": "generated in the previous step",
                    "type": "string"
                }
            }
        },
        "WP-Hw1_proto.GetUsersResponse": {
            "type": "object",
            "properties": {
                "message_id": {
                    "description": "odd and greater than zero",
                    "type": "integer"
                },
                "users": {
                    "description": "array of users",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/WP-Hw1_proto.USER"
                    }
                }
            }
        },
        "WP-Hw1_proto.PQResponse": {
            "type": "object",
            "properties": {
                "g": {
                    "description": "primitive root modulo of p",
                    "type": "integer"
                },
                "message_id": {
                    "description": "odd and greater than zero",
                    "type": "integer"
                },
                "nonce": {
                    "description": "the exact nonce from the clients request",
                    "type": "string"
                },
                "p": {
                    "description": "prime number",
                    "type": "integer"
                },
                "server_nonce": {
                    "description": "max_width -\u003e 20",
                    "type": "string"
                }
            }
        },
        "WP-Hw1_proto.USER": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "family": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "sex": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "WebPrograming homework 1",
	Description:      "a service which you can register in and get access to the users database",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
