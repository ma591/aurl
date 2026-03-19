# arc

A command line tool for turning any API into a CLI command, supporting OpenAPI 3.0, OpenAPI 3.1, Swagger 2.0, and GraphQL.

arc auto-detects authentication, validates requests against the spec, and generates documentation from introspection.

## Install

**Homebrew:**

```bash
brew install shawnpana/tap/arc
```

**Go:**

```bash
go install github.com/shawnpana/arc@latest
```

**From source:**

```bash
git clone https://github.com/ShawnPana/arc.git
cd arc
make install
```

## Quick Start

Register an API, then use it:

```bash
# Register a REST API
arc add petstore https://petstore3.swagger.io/api/v3/openapi.json

# Register a GraphQL API
arc add --graphql linear https://api.linear.app/graphql

# See what's available
arc petstore --help
arc linear --help

# Make requests
arc petstore GET /pet/1
arc linear '{ viewer { name email } }'
```

## Usage

### Register an API

```bash
# From a URL
arc add petstore https://petstore3.swagger.io/api/v3/openapi.json

# From a local file
arc add myapi ./openapi.json

# With a base URL override (if the spec doesn't include one)
arc add myapi ./openapi.json --base-url https://api.example.com

# With auth headers
arc add myapi https://api.example.com/openapi.json --header "Authorization: Bearer token"

# GraphQL endpoint
arc add --graphql linear https://api.linear.app/graphql
```

When registering, arc parses the spec's `securitySchemes` and prompts you for credentials:

```
Registered "petstore" (Swagger Petstore v1.0.27)
Base URL: https://petstore3.swagger.io/api/v3

Auth detected:
  [1] api_key (API Key in header "api_key")
  [2] petstore_auth (OAuth2 (provide token manually))

Enter value for api_key (or press Enter to skip): sk-xxxxx
Enter value for petstore_auth (or press Enter to skip):
```

### Explore an API

```bash
# See all endpoints grouped by tag
arc petstore --help
```

```
Swagger Petstore - OpenAPI 3.0 v1.0.27
Base URL: https://petstore3.swagger.io/api/v3

A sample Pet Store Server based on the OpenAPI 3.0 specification.
Docs: https://swagger.io

[pet] Everything about your Pets
  POST   /pet                           # Add a new pet to the store.
    body: '{"id":0,"name":"...","photoUrls":["..."],"status":"available"}'
    200: Pet  400: Invalid input
  GET    /pet/findByStatus              # Finds Pets by status.
    params: status* (available|pending|sold) - Status values to filter by
    200: Pet[]  400: Invalid status value
  GET    /pet/{petId}                   # Find pet by ID.
    params: petId* - ID of pet to return
    200: Pet  400: Invalid ID supplied  404: Pet not found

(*) = required
```

```bash
# Detailed docs for a single endpoint
arc petstore describe GET /pet/{petId}
```

```
GET /pet/{petId}
Find pet by ID.

Returns a single pet.

Parameters:
  petId* (path, integer)
    ID of pet to return

Responses:
  200 - successful operation
    Schema: Pet
  400 - Invalid ID supplied
  404 - Pet not found
```

```bash
# Open external documentation in browser
arc petstore docs
```

### Make Requests

**REST:**

```bash
# GET
arc petstore GET /pet/1

# GET with query params
arc petstore GET '/pet/findByStatus?status=available'

# POST with JSON body
arc petstore POST /pet '{"name":"doggie","photoUrls":["http://example.com"]}'

# PUT
arc petstore PUT /pet '{"id":1,"name":"updated","photoUrls":["http://example.com"]}'

# DELETE
arc petstore DELETE /pet/1
```

**GraphQL:**

```bash
# Query
arc linear '{ viewer { name email } }'

# Query with variables
arc linear '{ issue(id: $id) { title state { name } } }' '{"id":"ABC-123"}'

# Mutation
arc linear 'mutation { issueCreate(input: {title: "Bug", teamId: "xxx"}) { issue { id } } }'
```

### Validation

arc validates your requests against the spec before sending:

**Enum violations** are caught immediately:

```
$ arc petstore GET '/pet/findByStatus?status=invalid'
Error: status: invalid value "invalid"
  allowed: available, pending, sold
```

**Missing required fields** trigger a warning:

```
$ arc petstore POST /pet '{"tag":"dog"}'
Warning: missing required fields: name, photoUrls
  expected: '{"name":"...","photoUrls":["..."]}'
Send anyway? [y/N]:
```

**4xx errors** suggest the expected body from the spec:

```
$ arc petstore POST /pet
HTTP 415
Expected body for POST /pet:
  '{"name":"...","photoUrls":["..."],"status":"available"}'
```

### Manage APIs

```bash
# List all registered APIs
arc list

NAME      TYPE     TITLE                      VERSION  ENDPOINT
petstore  api      Swagger Petstore           1.0.27   https://petstore3.swagger.io/api/v3
linear    graphql                                       https://api.linear.app/graphql

# Reconfigure auth
arc auth petstore --header "Authorization: Bearer new-token"

# Remove an API
arc remove petstore
```

## Auth

arc supports all standard OpenAPI security schemes:

| Scheme | What arc does |
|--------|--------------|
| `apiKey` | Detects header/query param name from spec, prompts for value |
| `http` + `bearer` | Prompts for token, sets `Authorization: Bearer <token>` |
| `http` + `basic` | Prompts for username + password, encodes as Basic auth |
| `oauth2` / `openIdConnect` | Prompts for a token manually |

Auth is stored in `~/.config/arc/auth/` with `0600` permissions.

For GraphQL APIs or APIs without `securitySchemes`, use `--header`:

```bash
arc add myapi https://api.example.com/openapi.json --header "X-Api-Key: secret"
arc auth myapi --header "Authorization: Bearer new-token"
```

## How It Works

arc stores API specs and auth config locally:

```
~/.config/arc/
├── apis/          # OpenAPI spec files
├── graphql/       # GraphQL introspection results
└── auth/          # Auth headers per API (0600 permissions)
```

When you run `arc [name]`, it:

1. Loads the spec from `~/.config/arc/apis/` or `~/.config/arc/graphql/`
2. Parses it to understand endpoints, parameters, types, and auth
3. For `--help`: generates documentation from the spec
4. For requests: validates against the spec, attaches auth headers, executes, and pretty-prints the response

Specs are parsed lazily — arc only reads the spec for the command you invoke. Adding 50 APIs doesn't slow down startup.

## Shell Completions

```bash
# Zsh
arc completion zsh > "${fpath[1]}/_arc"

# Bash
arc completion bash > /etc/bash_completion.d/arc

# Fish
arc completion fish > ~/.config/fish/completions/arc.fish
```

## License

MIT
