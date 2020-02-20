// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-02-26 11:25:00.81335 +0800 CST m=+0.056945557

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
        },
        "/api/pharmacies/health_check": {
            "get": {
                "description": "The endpoint for Mask health check\n07-23 will return the lastest pharmacy update time, other wise will return ok",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "pharmacy"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/endpoints.HealthCheckResponse"
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
        "endpoints.HealthCheckResponse": {
            "type": "object",
            "properties": {
                "updated": {
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
                "message": {
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
	Version:     "v1",
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