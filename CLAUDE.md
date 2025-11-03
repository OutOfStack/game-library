# Project Guidelines

## Understanding the Codebase
- Read README.md to understand what this service is about

## Code Documentation
- Write comments for all exported functions and structs
- Start comments that are inside funcs with lowercase
- Do NOT use periods (`.`) at the end of comments
- Write proper openapi-style specs for handlers.

# Code practices:
- DO NOT use env variables unless it is specified, for dynamic configuring refer to previously read README.md to find out how configuration works

## Testing Requirements
- Write tests for all exported functions
- Place tests in separate files using the `*package*_test` naming convention
- Test files should be in the same directory as the code being tested
- If test already uses mock or requires updates or new mock, use `make generate` or add new line into `generate` command in `Makefile` file
- DO NOT write comments in tests unless they explain something that is not self-evident
- Use @internal/pkg/td/random.go to generate random data
- Use t.Context() instead of context.Background() for context in tests
- For local API tests refers to @docs/swagger.json. 
Some endpoints require authorization, refers game-library-auth API in https://github.com/OutOfStack/game-library-auth/blob/main/docs/swagger.json. 
Use `aiuser:aiuser__` creds for `user` with user role, and `aipublisher:aipublisher` creds for user with `publisher` role.

## Build and Quality Checks
- Run validation commands before completing work:
  - `make build` - compile the project
  - `make test` - run all tests
  - `make lint` - check code quality
  - `make generate | grep -E "(error:|warning:|failed)` - generate swagger files if there were updates in definitions
- Fix any issues found by these commands

## Documentation
- If there are significant updates regarding what written in `README.md`, add it there

## Git Workflow Restrictions
- DO NOT run `git stage` or `git commit`
- DO NOT delete files—notify me if files become redundant

## File Ignoring Guidelines
- DO NOT review or analyze files in the following directories:
  - `docs/**` - documentation files
  - `**/mocks/**` - generated mock files
  - `vendor/**` - external dependencies
- DO NOT review or analyze files matching these patterns:
  - `**/*_mock.go` - mock files
  - `**/*.gen.go` - generated files
  - `*.pem`, `*.key`, `app.env` - data-sensitive files
