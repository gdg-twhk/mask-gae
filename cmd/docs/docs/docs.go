// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-03-23 17:53:18.80603 +0800 CST m=+0.084743138

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
        "contact": {
            "name": "API Support",
            "url": "https://github.com/tnstiger/mask-gdg/issues",
            "email": "cage.chung@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/feedback": {
            "post": {
                "description": "The endpoint for Retailbase to fetch specific pharmacy feedbacks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feedback"
                ],
                "summary": "specific pharmacy feedbacks",
                "parameters": [
                    {
                        "description": "Feedback",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.FeedBackRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/endpoints.FeedBackResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    }
                }
            }
        },
        "/api/feedback/options": {
            "get": {
                "description": "The endpoint for Retailbase to fetch feedback options",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feedback"
                ],
                "summary": "feedback options",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/endpoints.OptionsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    }
                }
            }
        },
        "/api/feedback/pharmacies/{pharmacy_id}": {
            "get": {
                "description": "The endpoint for Retailbase to fetch specific pharmacy feedbacks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feedback"
                ],
                "summary": "specific pharmacy feedbacks",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "date, yyyy_mmdd",
                        "name": "date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Pharmacy ID",
                        "name": "pharmacy_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.FeedbackItemPage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    }
                }
            }
        },
        "/api/feedback/users/{user_id}": {
            "get": {
                "description": "The endpoint for Retailbase to fetch specific user feedbacks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feedback"
                ],
                "summary": "specific user feedbacks",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "date, yyyy_mmdd",
                        "name": "date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.FeedbackItemPage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    }
                }
            }
        },
        "/api/pharmacies": {
            "post": {
                "description": "The endpoint for Mask to fetch pharmacies",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pharmacy"
                ],
                "parameters": [
                    {
                        "description": "Fetch Pharmacies",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/endpoints.QueryRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/endpoints.QueryResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorRes"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "endpoints.Bounds": {
            "type": "object",
            "properties": {
                "ne": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.LatLng"
                },
                "nw": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.LatLng"
                },
                "se": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.LatLng"
                },
                "sw": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.LatLng"
                }
            }
        },
        "endpoints.FeedBackRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "optionId": {
                    "type": "string"
                },
                "pharmacyId": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "endpoints.FeedBackResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "endpoints.LatLng": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "number"
                },
                "lng": {
                    "type": "number"
                }
            }
        },
        "endpoints.OptionsResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Option"
                    }
                }
            }
        },
        "endpoints.QueryRequest": {
            "type": "object",
            "properties": {
                "bounds": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.Bounds"
                },
                "center": {
                    "type": "object",
                    "$ref": "#/definitions/endpoints.LatLng"
                },
                "max": {
                    "type": "integer"
                }
            }
        },
        "endpoints.QueryResponse": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Pharmacy"
                    }
                }
            }
        },
        "errors.Errors": {
            "type": "object",
            "properties": {
                "domain": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "locationType": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "model.Feedback": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "optionId": {
                    "type": "string"
                },
                "pharmacyId": {
                    "type": "string"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "model.FeedbackItemPage": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Feedback"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "model.Option": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.Pharmacy": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "available": {
                    "type": "string"
                },
                "county": {
                    "type": "string"
                },
                "cunli": {
                    "type": "string"
                },
                "customNote": {
                    "type": "string"
                },
                "distance": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "maskAdult": {
                    "type": "integer"
                },
                "maskChild": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "note": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "serviceNote": {
                    "type": "string"
                },
                "servicePeriods": {
                    "type": "string"
                },
                "town": {
                    "type": "string"
                },
                "updated": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "responses.ErrorRes": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "$ref": "#/definitions/responses.ErrorResItem"
                }
            }
        },
        "responses.ErrorResItem": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/errors.Errors"
                    }
                },
                "message": {
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
	Version:     "0.2.0",
	Host:        "mask.goodideas-studio.com",
	BasePath:    "/",
	Schemes:     []string{"https"},
	Title:       "Mask API",
	Description: "This is a Mask server celler server.",
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
