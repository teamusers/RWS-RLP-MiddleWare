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
        "/auth": {
            "post": {
                "description": "Validates AppID header and HMAC signature, then returns a JWT access token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Generate authentication token",
                "parameters": [
                    {
                        "type": "string",
                        "default": "app1234",
                        "description": "Client system AppID",
                        "name": "AppID",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Authentication request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/responses.AuthResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/gr/verify": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Checks if a GR member ID is already registered.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "gr"
                ],
                "summary": "Verify GR member existence",
                "parameters": [
                    {
                        "description": "GR registration check payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterGr"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "email registered",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "400": {
                        "description": "invalid JSON",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/grcms/profile/{reg_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieves a temporarily cached GR CMS profile by registration ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "grcms"
                ],
                "summary": "Get cached GR CMS profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Registration ID",
                        "name": "reg_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "successful",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.GrMember"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "missing reg_id",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "404": {
                        "description": "not found or expired",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/grcms/verify": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Checks if a GR CMS member email is in the system and caches their profile for follow‑up.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "grcms"
                ],
                "summary": "Verify and cache GR CMS member",
                "parameters": [
                    {
                        "description": "GR CMS register payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterGrCms"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "email not registered",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "400": {
                        "description": "invalid JSON",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "409": {
                        "description": "email already registered",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/member/burn-pin": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Updates the burn PIN for a given email address.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Update user burn PIN",
                "parameters": [
                    {
                        "description": "Email + new burn PIN",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.UpdateBurnPin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "update successful",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "400": {
                        "description": "invalid JSON or missing fields",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "update unsuccessful",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/member/{external_id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieves the profile (including phone numbers) for a given member by external_id.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Get member profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Member external ID",
                        "name": "external_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "profile found",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.User"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "missing or invalid external_id",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "member not found",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Updates a member’s profile fields (non‐zero values in the JSON body).",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "member"
                ],
                "summary": "Update member profile",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Member external ID",
                        "name": "external_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Profile fields to update",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "update successful",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.User"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "invalid input",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "member not found",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Validates user email, generates an OTP, emails it, and returns the OTP details plus a login session token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Start login flow via email",
                "parameters": [
                    {
                        "description": "Login request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/responses.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "invalid JSON",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "404": {
                        "description": "email not found",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/user/register": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Registers a new user record in the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create new user",
                "parameters": [
                    {
                        "description": "User create payload",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "user created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.User"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "invalid JSON",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        },
        "/user/register/verify": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Checks if an email is already registered; if not, sends an OTP for signup.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Verify email for registration",
                "parameters": [
                    {
                        "description": "Registration request payload",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.Register"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "email not registered, OTP sent",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/responses.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/model.Otp"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "invalid JSON",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "409": {
                        "description": "email already registered",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/responses.APIResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.GrMember": {
            "type": "object",
            "properties": {
                "dob": {
                    "description": "Dob is the date of birth in YYYY-MM-DD format.",
                    "type": "string",
                    "example": "1985-04-12"
                },
                "email": {
                    "description": "Email is the member's email address.",
                    "type": "string",
                    "example": "jane.doe@example.com"
                },
                "f_name": {
                    "description": "FirstName is the given name of the GR member.",
                    "type": "string",
                    "example": "Jane"
                },
                "gr_id": {
                    "description": "GrId is the unique GR member identifier.",
                    "type": "string",
                    "example": "GR12345"
                },
                "l_name": {
                    "description": "LastName is the family name of the GR member.",
                    "type": "string",
                    "example": "Doe"
                },
                "mobile": {
                    "description": "Mobile is the contact phone number.",
                    "type": "string",
                    "example": "98765432"
                }
            }
        },
        "model.Otp": {
            "type": "object",
            "properties": {
                "otp": {
                    "description": "Otp is the one‑time password sent to the user.\nexample: \"123456\"",
                    "type": "string",
                    "example": "123456"
                },
                "otp_expiry": {
                    "description": "OtpExpiry is the Unix timestamp (seconds since epoch) when the OTP expires.\nexample: 1744176000",
                    "type": "integer",
                    "example": 1744176000
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "burn_pin": {
                    "type": "integer",
                    "example": 1234
                },
                "country": {
                    "type": "string",
                    "example": "SGP"
                },
                "created_at": {
                    "type": "string",
                    "example": "2025-04-19T10:00:00Z"
                },
                "dob": {
                    "type": "string",
                    "example": "2007-08-05"
                },
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "external_id": {
                    "type": "string",
                    "example": "abc123"
                },
                "external_id_type": {
                    "type": "string",
                    "example": "EMAIL"
                },
                "first_name": {
                    "type": "string",
                    "example": "Brendan"
                },
                "id": {
                    "type": "integer",
                    "example": 42
                },
                "last_name": {
                    "type": "string",
                    "example": "Test"
                },
                "opted_in": {
                    "type": "boolean",
                    "example": true
                },
                "phone_numbers": {
                    "description": "PhoneNumbers holds zero or more phone numbers associated with this user.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.UserPhoneNumber"
                    }
                },
                "updated_at": {
                    "type": "string",
                    "example": "2025-04-19T11:00:00Z"
                }
            }
        },
        "model.UserPhoneNumber": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2025-04-19T10:05:00Z"
                },
                "id": {
                    "type": "integer",
                    "example": 101
                },
                "phone_number": {
                    "type": "string",
                    "example": "+6598765432"
                },
                "phone_type": {
                    "type": "string",
                    "example": "mobile"
                },
                "preference_flags": {
                    "type": "string",
                    "example": "primary"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2025-04-19T10:05:00Z"
                },
                "user_id": {
                    "type": "integer",
                    "example": 42
                }
            }
        },
        "requests.AuthRequest": {
            "type": "object",
            "required": [
                "nonce",
                "signature",
                "timestamp"
            ],
            "properties": {
                "nonce": {
                    "description": "A unique random string for each request to prevent replay attacks.",
                    "type": "string",
                    "example": "API"
                },
                "signature": {
                    "description": "HMAC-SHA256 signature of \"appID|timestamp|nonce\" hex-encoded.\nComputed by concatenating the appID, timestamp, and nonce to form a base string,\nthen applying HMAC-SHA256 with the secret key and hex-encoding the resulting digest.",
                    "type": "string",
                    "example": "1558850cb1b48e826197c48d6a14c5f3bf4b644bcb0065ceb0b07978296116bc"
                },
                "timestamp": {
                    "description": "Unix timestamp (seconds since epoch) when the request was generated.",
                    "type": "string",
                    "example": "1744075148"
                }
            }
        },
        "requests.Login": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "description": "Email address of the user attempting to log in",
                    "type": "string",
                    "example": "user@example.com"
                }
            }
        },
        "requests.Register": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "description": "Email address to check for existing registration.",
                    "type": "string",
                    "example": "user@example.com"
                }
            }
        },
        "requests.RegisterGr": {
            "type": "object",
            "required": [
                "gr_id",
                "gr_pin"
            ],
            "properties": {
                "gr_id": {
                    "description": "GR system identifier for the member.",
                    "type": "string",
                    "example": "GR12345"
                },
                "gr_pin": {
                    "description": "PIN code associated with the GR member.",
                    "type": "string",
                    "example": "9876"
                }
            }
        },
        "requests.RegisterGrCms": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "dob": {
                    "description": "Dob is the date of birth in YYYY-MM-DD format.",
                    "type": "string",
                    "example": "1985-04-12"
                },
                "email": {
                    "description": "Email is the member's email address.",
                    "type": "string",
                    "example": "jane.doe@example.com"
                },
                "f_name": {
                    "description": "FirstName is the given name of the GR member.",
                    "type": "string",
                    "example": "Jane"
                },
                "gr_id": {
                    "description": "GrId is the unique GR member identifier.",
                    "type": "string",
                    "example": "GR12345"
                },
                "l_name": {
                    "description": "LastName is the family name of the GR member.",
                    "type": "string",
                    "example": "Doe"
                },
                "mobile": {
                    "description": "Mobile is the contact phone number.",
                    "type": "string",
                    "example": "98765432"
                },
                "url": {
                    "description": "URL to send the registration confirmation link to.",
                    "type": "string",
                    "example": "https://example.com/confirm?reg_id=abc123"
                }
            }
        },
        "requests.UpdateBurnPin": {
            "type": "object",
            "required": [
                "burn_pin",
                "email"
            ],
            "properties": {
                "burn_pin": {
                    "description": "BurnPin is the new numeric PIN to set.",
                    "type": "integer",
                    "example": 4321
                },
                "email": {
                    "description": "Email of the user whose burn PIN is being updated.",
                    "type": "string",
                    "example": "user@example.com"
                }
            }
        },
        "requests.User": {
            "type": "object",
            "properties": {
                "burn_pin": {
                    "description": "BurnPin is the numeric PIN used for burn operations.",
                    "type": "integer",
                    "example": 1234
                },
                "email": {
                    "description": "Email is the user’s email address.",
                    "type": "string",
                    "example": "user@example.com"
                },
                "external_id": {
                    "description": "ExternalID is the client system’s unique identifier for this user.",
                    "type": "string",
                    "example": "abc123"
                },
                "external_id_type": {
                    "description": "ExternalTYPE describes the type or source of the external ID.",
                    "type": "string",
                    "example": "EMAIL"
                },
                "gr_id": {
                    "description": "GR_ID is the group or partner system identifier for the user.",
                    "type": "string",
                    "example": "GR12345"
                },
                "rlp_id": {
                    "description": "RLP_ID is the RLP system identifier for the user.",
                    "type": "string",
                    "example": "RLP67890"
                },
                "rws_membership_id": {
                    "description": "RWS_Membership_ID is the RWS membership ID assigned to this user.",
                    "type": "string",
                    "example": "RWS54321"
                },
                "rws_membership_number": {
                    "description": "RWS_Membership_Number is the numeric membership number in the RWS system.",
                    "type": "integer",
                    "example": 987654
                },
                "session_expiry": {
                    "description": "SessionExpiry is the Unix timestamp (seconds since epoch) when the session token expires.",
                    "type": "integer",
                    "example": 1712345678
                },
                "session_token": {
                    "description": "SessionToken is the login session token issued to the user.",
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1..."
                }
            }
        },
        "responses.APIResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Data holds the response payload. Its type depends on the endpoint:\ne.g. AuthResponse for /auth, LoginResponse for /user/login, etc."
                },
                "message": {
                    "description": "Message provides a human‑readable status or result description.\nExample: \"user created\", \"email found\"",
                    "type": "string",
                    "example": ""
                }
            }
        },
        "responses.AuthResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "description": "AccessToken is the JWT issued to the client for subsequent requests.\nExample: \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"",
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                }
            }
        },
        "responses.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "description": "Error provides the error detail.\nExample: \"invalid json request body\"",
                    "type": "string",
                    "example": "invalid json request body"
                }
            }
        },
        "responses.LoginResponse": {
            "type": "object",
            "properties": {
                "login_session_token": {
                    "description": "LoginSessionToken is the JWT issued after successful authentication.\nexample: \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...\"",
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                },
                "login_session_token_expiry": {
                    "description": "LoginSessionTokenExpiry is the Unix timestamp (seconds since epoch) when the token expires.\nexample: 1744176000",
                    "type": "integer",
                    "example": 1744176000
                },
                "otp": {
                    "description": "Otp is the one‑time password sent to the user.\nexample: \"123456\"",
                    "type": "string",
                    "example": "123456"
                },
                "otp_expiry": {
                    "description": "OtpExpiry is the Unix timestamp (seconds since epoch) when the OTP expires.\nexample: 1744176000",
                    "type": "integer",
                    "example": 1744176000
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "Type \"Bearer \u003cyour-jwt\u003e\" to authorize",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:18080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "LBE API",
	Description:      "Endpoints for authentication, login and register",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
