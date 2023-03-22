// Code generated by swaggo/swag. DO NOT EDIT
package http

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Dat Boi",
            "url": "https://datboi420.t.me"
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
        "/accounts": {
            "get": {
                "description": "Returns account states and its parsed data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "account data",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "only given addresses",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "only latest account states",
                        "name": "latest",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "filter by interfaces",
                        "name": "interface",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "filter FT wallets or NFT items by owner address",
                        "name": "owner_address",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "filter FT wallets or NFT items by minter address",
                        "name": "minter_address",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ASC",
                            "DESC"
                        ],
                        "type": "string",
                        "default": "DESC",
                        "description": "order by last_tx_lt",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start from this last_tx_lt",
                        "name": "after",
                        "in": "query"
                    },
                    {
                        "maximum": 10000,
                        "type": "integer",
                        "default": 3,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.AccountStateFilterResults"
                        }
                    }
                }
            }
        },
        "/accounts/aggregated": {
            "get": {
                "description": "Aggregates FT or NFT data filtered by minter address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "account"
                ],
                "summary": "aggregated account data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "NFT collection or FT master address",
                        "name": "minter_address",
                        "in": "query"
                    },
                    {
                        "maximum": 1000000,
                        "type": "integer",
                        "default": 25,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.AccountStateAggregation"
                        }
                    }
                }
            }
        },
        "/blocks": {
            "get": {
                "description": "Returns filtered blocks",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "block"
                ],
                "summary": "block info",
                "parameters": [
                    {
                        "type": "integer",
                        "default": -1,
                        "description": "workchain",
                        "name": "workchain",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "shard",
                        "name": "shard",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "seq_no",
                        "name": "seq_no",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "default": false,
                        "description": "include transactions",
                        "name": "with_transactions",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ASC",
                            "DESC"
                        ],
                        "type": "string",
                        "default": "DESC",
                        "description": "order by seq_no",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start from this seq_no",
                        "name": "after",
                        "in": "query"
                    },
                    {
                        "maximum": 100,
                        "type": "integer",
                        "default": 3,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.BlockFilterResults"
                        }
                    }
                }
            }
        },
        "/contract/interfaces": {
            "get": {
                "description": "Returns known contract interfaces",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "contract interfaces",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.GetInterfacesRes"
                        }
                    }
                }
            }
        },
        "/contract/operations": {
            "get": {
                "description": "Returns known contract message payloads schema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "contract"
                ],
                "summary": "contract operations",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/http.GetOperationsRes"
                        }
                    }
                }
            }
        },
        "/messages": {
            "get": {
                "description": "Returns filtered messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transaction"
                ],
                "summary": "transaction messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "msg hash",
                        "name": "hash",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "source address",
                        "name": "src_address",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "destination address",
                        "name": "dst_address",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "source contract interface",
                        "name": "src_contract",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "destination contract interface",
                        "name": "dst_contract",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "filter by contract operation names",
                        "name": "operation_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "filter FT or NFT operations by minter address",
                        "name": "minter_address",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ASC",
                            "DESC"
                        ],
                        "type": "string",
                        "default": "DESC",
                        "description": "order by created_lt",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start from this created_lt",
                        "name": "after",
                        "in": "query"
                    },
                    {
                        "maximum": 10000,
                        "type": "integer",
                        "default": 3,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.MessageFilterResults"
                        }
                    }
                }
            }
        },
        "/transactions": {
            "get": {
                "description": "Returns transactions, states and messages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transaction"
                ],
                "summary": "transactions data",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "only given addresses",
                        "name": "address",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by tx hash",
                        "name": "hash",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "search by incoming message hash",
                        "name": "in_msg_hash",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "filter by workchain",
                        "name": "workchain",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "ASC",
                            "DESC"
                        ],
                        "type": "string",
                        "default": "DESC",
                        "description": "order by created_lt",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start from this created_lt",
                        "name": "after",
                        "in": "query"
                    },
                    {
                        "maximum": 10000,
                        "type": "integer",
                        "default": 3,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/core.TransactionFilterResults"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "bunbig.Int": {
            "type": "object"
        },
        "core.AccountData": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "admin_addr": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "balance": {
                    "type": "string"
                },
                "content_description": {
                    "type": "string"
                },
                "content_image": {
                    "type": "string"
                },
                "content_image_data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "content_name": {
                    "type": "string"
                },
                "content_uri": {
                    "type": "string"
                },
                "editor_address": {
                    "description": "CollectionAddress *addr.Address ` + "`" + `ch:\"type:String\" bun:\"type:bytea\" json:\"collection_address,omitempty\"` + "`" + `",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "error": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "initialized": {
                    "type": "boolean"
                },
                "item_index": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "last_tx_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "last_tx_lt": {
                    "type": "integer"
                },
                "mintable": {
                    "type": "boolean"
                },
                "minter_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "next_item_index": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "owner_address": {
                    "description": "common fields for FT and NFT",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "royalty_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "royalty_base": {
                    "type": "integer"
                },
                "royalty_factor": {
                    "type": "integer"
                },
                "total_supply": {
                    "type": "string"
                },
                "types": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "core.AccountState": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "balance": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "code": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "code_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "data_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "get_method_hashes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "is_active": {
                    "type": "boolean"
                },
                "last_tx_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "last_tx_lt": {
                    "type": "integer"
                },
                "state_data": {
                    "$ref": "#/definitions/core.AccountData"
                },
                "state_hash": {
                    "description": "only if account is frozen",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "status": {
                    "description": "TODO: ch enum",
                    "type": "string"
                }
            }
        },
        "core.AccountStateAggregation": {
            "type": "object",
            "properties": {
                "items": {
                    "description": "NFT minter",
                    "type": "integer"
                },
                "owned_balance": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.OwnedBalance"
                    }
                },
                "owned_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.OwnedItems"
                    }
                },
                "owners_count": {
                    "type": "integer"
                },
                "total_supply": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "unique_owners": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.UniqueOwners"
                    }
                },
                "wallets": {
                    "description": "FT minter",
                    "type": "integer"
                }
            }
        },
        "core.AccountStateFilterResults": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.AccountState"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "core.Block": {
            "type": "object",
            "properties": {
                "file_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "master": {
                    "$ref": "#/definitions/core.BlockID"
                },
                "root_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "seq_no": {
                    "type": "integer"
                },
                "shard": {
                    "type": "integer"
                },
                "shards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Block"
                    }
                },
                "transactions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Transaction"
                    }
                },
                "workchain": {
                    "type": "integer"
                }
            }
        },
        "core.BlockFilterResults": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Block"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "core.BlockID": {
            "type": "object",
            "properties": {
                "seq_no": {
                    "type": "integer"
                },
                "shard": {
                    "type": "integer"
                },
                "workchain": {
                    "type": "integer"
                }
            }
        },
        "core.ContractInterface": {
            "type": "object",
            "properties": {
                "addresses": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        }
                    }
                },
                "code": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "code_hash": {
                    "description": "TODO: match by code hash",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "get_method_hashes": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "get_methods": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "core.ContractOperation": {
            "type": "object",
            "properties": {
                "contract_name": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "operation_id": {
                    "type": "integer"
                },
                "outgoing": {
                    "description": "if operation is going from contract",
                    "type": "boolean"
                },
                "schema": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.Message": {
            "type": "object",
            "properties": {
                "amount": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "body": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "body_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "bounce": {
                    "type": "boolean"
                },
                "bounced": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string"
                },
                "created_lt": {
                    "type": "integer"
                },
                "dst_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "fwd_fee": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "ihr_disabled": {
                    "type": "boolean"
                },
                "ihr_fee": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "operation_id": {
                    "type": "integer"
                },
                "payload": {
                    "$ref": "#/definitions/core.MessagePayload"
                },
                "source": {
                    "description": "TODO: join it",
                    "allOf": [
                        {
                            "$ref": "#/definitions/core.Transaction"
                        }
                    ]
                },
                "source_tx_hash": {
                    "description": "SourceTx initiates outgoing message.\nFor external incoming messages SourceTx == nil.",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "source_tx_lt": {
                    "type": "integer"
                },
                "src_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "state_init_code": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "state_init_data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "transfer_comment": {
                    "type": "string"
                },
                "type": {
                    "description": "TODO: ch enum",
                    "type": "string"
                }
            }
        },
        "core.MessageFilterResults": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Message"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "core.MessagePayload": {
            "type": "object",
            "properties": {
                "amount": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "body_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "created_at": {
                    "type": "string"
                },
                "created_lt": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "dst_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "dst_contract": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "minter_address": {
                    "description": "can be used to show all jetton or nft item transfers",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "operation_id": {
                    "type": "integer"
                },
                "operation_name": {
                    "type": "string"
                },
                "src_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "src_contract": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "core.OwnedBalance": {
            "type": "object",
            "properties": {
                "balance": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "owner_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "wallet_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.OwnedItems": {
            "type": "object",
            "properties": {
                "items_count": {
                    "type": "integer"
                },
                "owner_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "core.Transaction": {
            "type": "object",
            "properties": {
                "account": {
                    "$ref": "#/definitions/core.AccountState"
                },
                "address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "block_seq_no": {
                    "type": "integer"
                },
                "block_shard": {
                    "type": "integer"
                },
                "block_workchain": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "created_lt": {
                    "type": "integer"
                },
                "description": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "end_status": {
                    "type": "string"
                },
                "hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "in_amount": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "in_msg": {
                    "$ref": "#/definitions/core.Message"
                },
                "in_msg_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "orig_status": {
                    "type": "string"
                },
                "out_amount": {
                    "$ref": "#/definitions/bunbig.Int"
                },
                "out_msg": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Message"
                    }
                },
                "out_msg_count": {
                    "type": "integer"
                },
                "prev_tx_hash": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "prev_tx_lt": {
                    "type": "integer"
                },
                "state_update": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "total_fees": {
                    "$ref": "#/definitions/bunbig.Int"
                }
            }
        },
        "core.TransactionFilterResults": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.Transaction"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "core.UniqueOwners": {
            "type": "object",
            "properties": {
                "item_address": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "owners_count": {
                    "type": "integer"
                }
            }
        },
        "http.GetInterfacesRes": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.ContractInterface"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "http.GetOperationsRes": {
            "type": "object",
            "properties": {
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/core.ContractOperation"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "anton.tools",
	BasePath:         "/api/v0",
	Schemes:          []string{"https"},
	Title:            "tonidx",
	Description:      "Project fetches data from TON blockchain.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
