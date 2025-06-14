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
        "/api/client/{index}/enabled": {
            "put": {
                "description": "Updates the enabled state for a specific client",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Update client enabled state",
                "operationId": "updateClientEnabled",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Client index",
                        "name": "index",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New enabled state",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateClientEnabledRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid client index",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Client not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/client/{index}/group": {
            "put": {
                "description": "Updates the group for a specific client",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Update client group",
                "operationId": "updateClientGroup",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Client index",
                        "name": "index",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New group name",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateClientGroupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid client index",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Client not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/clients": {
            "get": {
                "description": "Returns a list of all clients with their details",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "summary": "Get all clients",
                "operationId": "getClients",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.ClientEntry"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/scenarios": {
            "get": {
                "description": "Returns a list of all scenarios",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Scenario"
                ],
                "summary": "Get all scenarios",
                "operationId": "getScenarios",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.ScenarioEntries"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Failure",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/scenarios/{name}/config": {
            "get": {
                "description": "Returns the default configuration for a specific scenario",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Scenario"
                ],
                "summary": "Get scenario configuration",
                "operationId": "getScenarioConfig",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Scenario name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "YAML configuration",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Scenario not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer": {
            "post": {
                "description": "Creates a new spammer with the given configuration",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Create a new spammer",
                "operationId": "createSpammer",
                "parameters": [
                    {
                        "description": "Spammer configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateSpammerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Spammer ID",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "integer"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}": {
            "get": {
                "description": "Returns detailed information about a specific spammer",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Get spammer details",
                "operationId": "getSpammerDetails",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/api.SpammerDetails"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates an existing spammer's configuration",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Update a spammer",
                "operationId": "updateSpammer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateSpammerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a spammer and stops it if running",
                "tags": [
                    "Spammer"
                ],
                "summary": "Delete a spammer",
                "operationId": "deleteSpammer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}/logs": {
            "get": {
                "description": "Returns the most recent logs for a specific spammer",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Get spammer logs",
                "operationId": "getSpammerLogs",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.LogEntry"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}/logs/stream": {
            "get": {
                "description": "Streams logs for a specific spammer using Server-Sent Events",
                "produces": [
                    "text/event-stream"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Stream spammer logs",
                "operationId": "streamSpammerLogs",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Timestamp to start from (RFC3339Nano)",
                        "name": "since",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SSE stream of log entries",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Streaming unsupported",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}/pause": {
            "post": {
                "description": "Pauses a running spammer",
                "tags": [
                    "Spammer"
                ],
                "summary": "Pause a spammer",
                "operationId": "pauseSpammer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}/reclaim": {
            "post": {
                "description": "Reclaims funds from a spammer's wallet pool back to the root wallet",
                "tags": [
                    "Spammer"
                ],
                "summary": "Reclaim funds from a spammer",
                "operationId": "reclaimFunds",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammer/{id}/start": {
            "post": {
                "description": "Starts a specific spammer",
                "tags": [
                    "Spammer"
                ],
                "summary": "Start a spammer",
                "operationId": "startSpammer",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Spammer ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "400": {
                        "description": "Invalid spammer ID",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "404": {
                        "description": "Spammer not found",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammers": {
            "get": {
                "description": "Returns a list of all configured spammers",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Get all spammers",
                "operationId": "getSpammerList",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.SpammerListEntry"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/spammers/export": {
            "post": {
                "description": "Exports specified spammers or all spammers to YAML format",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Export spammers to YAML",
                "operationId": "exportSpammers",
                "parameters": [
                    {
                        "description": "Spammer IDs to export (optional)",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/api.ExportSpammersRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "YAML configuration",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        },
        "/api/spammers/import": {
            "post": {
                "description": "Imports spammers from YAML data or URL with validation and deduplication",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Spammer"
                ],
                "summary": "Import spammers from YAML data or URL",
                "operationId": "importSpammers",
                "parameters": [
                    {
                        "description": "Import configuration",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.ImportSpammersRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/daemon.ImportResult"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ClientEntry": {
            "type": "object",
            "properties": {
                "block_height": {
                    "type": "integer"
                },
                "enabled": {
                    "type": "boolean"
                },
                "group": {
                    "type": "string"
                },
                "index": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "ready": {
                    "type": "boolean"
                },
                "rpc_host": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "api.CreateSpammerRequest": {
            "type": "object",
            "properties": {
                "config": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "scenario": {
                    "type": "string"
                },
                "startImmediately": {
                    "type": "boolean"
                }
            }
        },
        "api.ExportSpammersRequest": {
            "type": "object",
            "properties": {
                "spammer_ids": {
                    "description": "If empty, exports all spammers",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "api.ImportSpammersRequest": {
            "type": "object",
            "properties": {
                "input": {
                    "description": "Can be YAML data or a URL",
                    "type": "string"
                }
            }
        },
        "api.LogEntry": {
            "type": "object",
            "properties": {
                "fields": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "level": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                }
            }
        },
        "api.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "error": {
                    "type": "string"
                }
            }
        },
        "api.ScenarioEntries": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "api.SpammerDetails": {
            "type": "object",
            "properties": {
                "config": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "scenario": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "api.SpammerListEntry": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "RFC3339Nano formatted timestamp",
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "scenario": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "api.UpdateClientEnabledRequest": {
            "type": "object",
            "properties": {
                "enabled": {
                    "type": "boolean"
                }
            }
        },
        "api.UpdateClientGroupRequest": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                }
            }
        },
        "api.UpdateSpammerRequest": {
            "type": "object",
            "properties": {
                "config": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "daemon.ImportResult": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "imported": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/daemon.ImportedSpammerInfo"
                    }
                },
                "imported_count": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "validation": {
                    "$ref": "#/definitions/daemon.ImportValidationResult"
                },
                "warnings": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "daemon.ImportValidationResult": {
            "type": "object",
            "properties": {
                "duplicates": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "invalid_scenarios": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "spammers": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/daemon.SpammerValidationInfo"
                    }
                },
                "total_spammers": {
                    "type": "integer"
                },
                "valid_spammers": {
                    "type": "integer"
                }
            }
        },
        "daemon.ImportedSpammerInfo": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "scenario": {
                    "type": "string"
                }
            }
        },
        "daemon.SpammerValidationInfo": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "issues": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "scenario": {
                    "type": "string"
                },
                "valid": {
                    "type": "boolean"
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
