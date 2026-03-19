# arc

Turn any API spec into a CLI command. Register APIs by name, explore their endpoints, and make requests.

## Commands

```bash
arc add [name] [openapi.json URL or path]       # register an API
arc add [name] [spec] --base-url [URL]           # override base URL
arc add [name] [spec] --header "Key: Value"      # manual auth header
arc add --graphql [name] [endpoint]              # register a GraphQL API
arc list                                          # list registered APIs
arc remove [name]                                 # unregister
arc auth [name]                                   # reconfigure auth
```

## Using a Registered API

```bash
arc [name] --help                                # see all endpoints, params, enums, response codes
arc [name] describe METHOD /path                 # detailed docs for one endpoint
arc [name] describe [field]                      # detailed docs for a GraphQL field/type
arc [name] docs                                  # open external docs in browser
arc [name] METHOD /path                          # make a REST request
arc [name] METHOD /path '{"key":"value"}'        # make a REST request with JSON body
arc [name] '{ graphql query }'                   # run a GraphQL query
arc [name] '{ query }' '{"var":"val"}'           # GraphQL with variables
```

## Discovery Workflow

When working with an unfamiliar registered API:

1. `arc [name] --help` — scan all endpoints grouped by tag
2. `arc [name] describe METHOD /path` — read full docs for an endpoint before calling it
3. Make the request — enum params are validated before sending, missing required body fields trigger a warning

## Notes

- Auth is auto-detected from the spec's `securitySchemes` during `arc add`
- Query params go inline in the path: `arc [name] GET '/path?key=value'`
- Quote paths with `?` to prevent shell globbing
- On 4xx errors, the CLI suggests the expected request body from the spec
