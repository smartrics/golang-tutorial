# GOLANG Tutorial

This is a step-by-step tutorial designed with the help of ChatGPT.  
It walks through the design and implementation of a simple banking system in Go, with each part building incrementally on top of the previous one.

It is designed to guide mid-to-senior engineers through building a real-world backend system from scratch — a banking platform — using idiomatic Go, clean architecture principles, and test-driven development.

The purpose is to help experienced Go developers refresh their knowledge through practical, progressive exercises — with a focus on writing production-quality code, structuring real-world projects, and revisiting concurrency, interfaces, and clean design.

## 🎯 What You'll Build

  * Core banking logic (accounts, transfers)
  * Clean interfaces and structs
  * Transfer pipelines with decorators and mixins
  * An async job processor (fan-in/fan-out, concurrency patterns)
  * A testable, composable transfer engine
  * A full HTTP API (routes, JSON I/O, handlers)
  * End-to-end integration tests
  * Monitoring/debug endpoints for observability

Each **part is implemented in a dedicated branch** named `partX` (e.g. `part1`, `part2`, ... `part8`).

---

## 🔁 Tutorial Structure

| Part | Branch   | Focus                                         |
|------|----------|-----------------------------------------------|
| 1    | [part1](https://github.com/smartrics/golang-tutorial/tree/part1) | Structs, interfaces, basic testing |
| 2    | [part2](https://github.com/smartrics/golang-tutorial/tree/part2) | Table-driven tests, test coverage |
| 3    | [part3](https://github.com/smartrics/golang-tutorial/tree/part3) | Mocking, interfaces, error handling |
| 4    | [part4](https://github.com/smartrics/golang-tutorial/tree/part4) | Functional pipelines, decorators |
| 5    | [part5](https://github.com/smartrics/golang-tutorial/tree/part5) | Transfer engine + integration tests |
| 6    | [part6](https://github.com/smartrics/golang-tutorial/tree/part6) | Async processor with concurrency |
| 7    | [part7](https://github.com/smartrics/golang-tutorial/tree/part7) | Transfer engine + processor coordination |
| 8    | [part8](https://github.com/smartrics/golang-tutorial/tree/part8) | HTTP server, monitoring, E2E testing |

> 💡 You can check out any part using `git checkout partX`.

---

## 🚀 Getting Started

Clone the repo and initialise your Go module:

```bash
git clone https://github.com/smartrics/golang-tutorial.git
cd golang-tutorial
git checkout part1   # or any part you want to start from

go mod tidy
```

*Note: Use your own module path if you're adapting the tutorial.*

## 🧱 Build Philosophy
  * Test-first (TDD) development
  * Realistic banking domain use case
  * Clean architecture + Hexagonal layering
  * Async processing, observability, and end-to-end testing

## Part 1: Syntax & Language Fundamentals

### 🎯 Goal:

 * Refresh Go syntax and control flow
 * Implement a practical Bank Account Simulator
 * Write unit tests using idiomatic Go

### ✅ Requirements

*Objective*: Build a basic in-memory bank account system in Go, focusing on fundamental language constructs and idiomatic style.

#### 💡 Functional Requirements

 * Create accounts with an ID and an initial balance.
 * Deposit funds into an account.
 * Withdraw funds, with an error returned if funds are insufficient.
 * View balance and represent an account as a formatted string.

#### 🧠 Non-Functional Requirements

 * Use Go structs, methods, and error handling idioms.
 * Follow Go's naming conventions and receiver patterns.
 * Make Account immutable.
 * Encapsulate account behaviour behind methods (not just raw field access).
 * Add unit tests to validate logic using the built-in testing package.

### Golang gotchas

|  #  | Gotcha| Description |
| --- |  ---  |    ---      |
| 1  | Pointer vs Value Receiver   | Value receivers won't mutate original struct             |
| 2  | String-based error matching | Use `errors.Is()` instead                                |
| 3  | Range copies values         | Use `&slice[i]` or index access                          |
| 4  | Defer order & scope         | Defer runs LIFO after return                             |
| 5  | Ignoring errors             | Always check returned errors                             |
| 6  | Interface `nil != nil` | Use `errors.Is()` or `errors.As()` |
| 7  | Lowercase fields            | Not exported outside package                             |
| 8  | Zero values work            | e.g. `Account{}` is safe unless you enforce constructors |

## Part 2: Structs, Methods, Interfaces

### 🎯 Objective

To learn how Go handles:

 * Structs and encapsulation
 * Methods and receiver types
 * Interfaces and implicit satisfaction
 * Composition via embedding (instead of inheritance)

You’ll build a small `BankService` abstraction that operates on different types of accounts ( `SavingsAccount` , `CheckingAccount` ) and uses interfaces to decouple behaviour. This introduces polymorphism in Go.

### ✅ What You'll Build

You'll extend your previous immutable Account into:
 * Multiple account types:
   * `SavingsAccount` : supports interest
   * `CheckingAccount` : may support overdraft
 * A service interface:
   * `BankService` defines operations on accounts (e.g., `Transfer()` , `GetStatement()` )
 * You’ll also add:
   * Method overloading patterns via interface
   * Unit tests to verify type behaviour
   * Embedding to reuse base logic without inheritance

### Requirements (Functional & Design)

#### 📦 Types and Structs

 1. Define a common base account struct for shared logic:
    + Internal fields: id, balance
    + Expose via `Balance()` and `ID()`
 2. Define `SavingsAccount` and `CheckingAccount` :
    + Embed the base account struct
    + Add type-specific logic:
      - `SavingsAccount`: support for an interest rate between 0 and 1 (inclusive) and add `ApplyInterest()`
      - `CheckingAccount`: support an overdraft limit
 3. Design these types as immutable: return new instances on state change

#### 🧠 Behaviour & Interfaces

 * Define a `BankAccount` interface
 * Implement the interface in both `SavingsAccount` and `CheckingAccount` (implicitly)
 * Define a `BankService` interface with:
   * `Transfer(from, to BankAccount, amount float64, reference string) (BankAccount, BankAccount, error)`
   * `GetStatement(acc BankAccount)`

#### 🧪 Testing

Use table-driven tests to test:
 * Withdraw and deposit with specific rules per account type
 * Transfers via BankService
 * Verify code coverage
   * Run tests with `go test -v -coverprofile="coverage.out"`
   * Observe coverage with `go tool cover -html coverage.out`

### ⚙️ Non-Functional Requirements

 * Code should follow Go idioms (zero-value safe, unexported fields where appropriate)
 * Use interface embedding sparingly but demonstrate it
 * No pointers unless required (e.g., in test helpers)
 * Encapsulation: hide fields, expose behaviour via methods
 * Add `String()` implementations for debugging

### Gotchas

|#| Gotcha| Example / Impact| Fix / Idiom|
|---|---|---|---|
|1| Implicit interface satisfaction| No implements keyword |Use `var _ Interface = Type{}` |
|2| Pointer vs value receiver for interfaces| `*T` needed to satisfy interface| Use pointer receiver for mutability|
|3| Interface nil ≠ typed nil| `var x *T = nil → interface != nil` |Use `errors.Is` , handle nil explicitly|
|4| Embedding ≠ override| Method sets do not override |Understand embedding is composition|
|5| Interface pollution |Massive interfaces| Keep interfaces small, single-purpose|
|6| Copying structs| Shared slice/map state| Use deep copy or avoid direct copy|
|7| Method set mismatch| Value receiver ≠ pointer receiver| Be consistent, prefer pointer receivers|
|8| Unsafe zero-values| Accessing unset fields| Make types safe by default|

## Part 3: Error Handling, Testing, Tooling

### 🎯 Goal

 * Understand and use Go’s idiomatic error handling (error, `fmt.Errorf`, `errors.Is`/`As`)
 * Write table-driven, benchmark, and mock-driven tests
 * Use Go's built-in tooling: `go test`, `go vet`, `golangci-lint`, `go fmt`
 * Improve code quality with static analysis
 * Use custom errors to support richer business logic

### 📋 Requirements

#### ✅ Functional Requirements

1. Extend your banking code to return and handle rich errors:
   * Define custom error types (ErrInsufficientFunds, etc.)
   * Use errors.Is() and errors.As() to match and extract
   * Use fmt.Errorf(...%w...) for wrapping
2. Add new unit tests for error flows:
   * Insufficient funds
   * Invalid input (negative amount)
   * Self-transfer
3. Convert unit tests to table-driven tests for readability and coverage 
3. Add a benchmark test for Transfer() performance
4. Introduce basic mocking:
   * Use a fake/mock BankAccount for testing BankService
   * use `moq` as mocking framework

#### ⚙️ Tooling Requirements

1. Format, vet, and lint your code:
   * Use `go fmt`, `go vet`
   * Use `golangci-lint` (optional)
2. Add Makefile to automate testing, linting, and formatting
3. Optionally: Add GitHub Actions to enforce tests/quality

### ⚠️ Gotchas & Tips
|Gotcha|Why It Matters|
| ---  | --- |
|`errors.New()` vs `fmt.Errorf(...%w...)`  |Use `fmt.Errorf` to wrap underlying causes                    |
|Comparing `err.Error() strings`           |Fragile — prefer `errors.Is()` or `errors.As()`               |
|Interface method returns `nil` typed value|Still non-nil interface! Use explicit `nil`                   |
|Forgetting to run `go vet`                |It catches subtle bugs — always run it with tests             |
|Linting ignored                           |`golangci-lint` catches bad practices beyond `go vet`         |
|Benchmarks require naming convention      |Must start with `BenchmarkXxx` to run with `go test -bench`   |
|Tests without assertions                  |Always compare expected values or use libraries like `testify`|

## Part 4: Application Structure, Project Organisation, and CI/CD Readiness

### 🎯 Goal

 * Organise your project using idiomatic Go layout (cmd/, pkg/, internal/, etc.)
 * Split domain, service, and transport logic
 * Configure a simple CI workflow (GitHub Actions or local runner)
 * Integrate formatting, vetting, linting, and test automation
 * Write integration-style tests that simulate end-to-end usage

### 📋 Requirements

#### ✅ Functional Requirements

Project layout should support:
 * Clean separation between domain logic and orchestration
 * A `cmd/` folder with a real `main.go` entry point
 * A `pkg/` or `internal/` directory for reusable business code
 * Easy import for mocks and unit tests

Make application testable:
 * Interfaces for external dependencies
 * All domain logic isolated in pure functions or services
 * Setup logic (e.g., `NewBankService`) extracted into `factory.go` or `main.go`

Testing:
 * All existing table-driven + benchmark tests still work
 * Add at least one integration test: simulate a full transfer + statement fetch

CI/CD Readiness:
 * Linting with `golangci-lint`
 * Testing with `go test ./...`
 * Format check with `go fmt` or `gofmt -l`
 * Optional: GitHub Actions config (`.github/workflows/ci.yml`)

#### 🧰 Tooling
|Tool|Purpose|
| ---                         | ---                     |
|`go test ./...`              | Run all tests           |
|`go test -bench=. -benchmem` | Run benchmarks          |
|`go vet ./...`               | Static analysis         |
|`gofmt`, `go fmt ./...`      | Code formatting         |
|`golangci-lint`              | Aggregated linter       |
|`make` / `build.ps1`         | Automate commands       |
|GitHub Actions               | Run checks on push/PR   |
|`moq`                        | Generate interface mocks|

### ⚠️ Gotchas
|# |Gotcha|Tip|
|---|---|---|
|1 | Mixing domain and I/O logic | Split handlers/controllers from services and models|
|2 | No dependency boundaries    | Use interfaces for services and ports/adapters for I/O|
|3 | Everything in `main.go`     | Move wiring logic to `cmd/` or internal/ packages|
|4 | CI flaky or missing         | Use local make test before relying on GitHub CI|
|5 | Test files in wrong folders | Put integration tests in `/test/`, not `/pkg/`|
|6 | Skipping `go vet`           | Always run vet in CI — catches actual bugs|
|7 | gofmt vs `go fmt` confusion | `go fmt` is a wrapper for `gofmt`; use `go fmt ./...` for simplicity|

## Part 5: Functional Composition, OO Patterns, and Advanced Struct Techniques

### 🎯 Goal

 * Apply object-oriented composition idiomatically in Go (struct embedding, interfaces)
 * Use functional programming techniques like function literals, decorators, and pipelines
 * Understand and apply mixins via composition and embedding
 * Build and test reusable behaviours (e.g., logging, metrics, validation wrappers)
 * Maintain immutability and testability in composable patterns

### 📋 Requirements

#### ✅ Object-Oriented Composition (Go idioms)

 * Use interface composition to split behaviors (Withdrawer, Depositor, Statementable)
 * Use struct embedding to reuse logic between different account types
 * Implement a mixin-style helper (e.g., to track audit logs, transaction counts)
 * Enable method override via embedding shadowing

#### ✅ Functional Composition

 * Write a transfer pipeline using decorator functions:
   * e.g., `withLogging`, `withValidation`, `withAuditing`
 * Demonstrate functional chaining or middleware-style layering:
   * `transfer := withLogging(withValidation(realTransfer))`
 * Write pure functions for operations like `ApplyInterest`, `FeeDeduction`

#### ✅ Advanced Techniques
 * Demonstrate composition vs inheritance clearly
 * Use interfaces + embedding to achieve reusable but isolated logic
 * Avoid reflect and generics unless absolutely needed
 * Structure composable business logic like policy evaluation chains or validators

### ⚠️ Gotchas
| # | Gotcha|Tip|
|---|---|---|
|1|Expecting inheritance|Go uses composition (struct embedding, interfaces) instead|
|2|Embedding != polymorphism|Embedding provides reuse, but not virtual dispatch|
|3|Function types aren't interfaces|Decorators need to match the signature exactly|
|4|Over reliance on global state|Keep functional wrappers pure where possible|
|5|Testing embedded behavior|Write tests for the outer type, not just the inner struct|
|6|Overly rigid interface hierarchies|Compose interfaces from minimal responsibilities|
|7|Confusing struct method sets|Remember: value vs pointer receivers matter in composition|


### 🧰 Tools and Concepts
|Concept|Why It Matters|
|Struct embedding|Simulates mixins / method reuse|
|Functional wrapping|Enables cross-cutting concerns (log, auth, audit)|
|Small interfaces|Increases testability and reusability|
|Closures|Maintain internal state (e.g., counters, contexts)|
|Middleware chaining|Enables business logic orchestration|
|Immutability + composition|Safer for concurrency and testing|

### 🧪 What You’ll End Up With
 * Multiple account types sharing common logic, but with different behaviour
 * Decorators for logging, validation, and error wrapping
 * A composable Transfer pipeline (like functional middleware)
 * A clearly isolated, testable, modular service design

## Part 6: Concurrency, Channels and Project Architecture

### 🎯 Goal

Understand and correctly implement all core concurrency patterns in Go:
 * Goroutines
 * Channels (buffered/unbuffered)
 * `select` statements
 * Worker pools
 * Cancellation using `context`
 * Fan-out / fan-in patterns

Apply these patterns to banking-relevant scenarios, such as:
 * Parallel transaction processing
 * Concurrent balance aggregation
 * Timed or cancellable transfer operations

Establish a resilient and maintainable application architecture, with:
 * Separation of orchestration vs business logic
 * Lifecycle-safe goroutine management
 * Support for observability and graceful shutdown

### 📋 Requirements
 1. Concurrency Fundamentals
    * Use goroutines to execute transfer operations concurrently
    * Use channels to coordinate:
      * Event dispatch
      * Result aggregation
      * Backpressure and throttling
    * Use `context.Context` for timeout and cancellation propagation

 2. Concurrency Patterns
    * Implement and test:
      * Fan-out: Splitting a stream of tasks to N workers
      * Fan-in: Aggregating results from multiple sources
      * Bounded worker pool: Processing a channel of transfer requests with limited workers
      * Timeout pattern: Enforcing time limits per transfer using context.WithTimeout

 3. Application Architecture
    * Extract transfer processing into a dedicated goroutine-managed service
    * Establish lifecycle hooks for starting/stopping workers
    * Ensure thread-safety in shared state (transaction history, audit log, etc.)
    * Make unit tests deterministic using channels and contexts

### ⚠️ Common Gotchas
|Problem|Cause|Mitigation
|---|---|---|
|Goroutine leaks|Forgetting to exit on cancel/timeout|Always check `<-ctx.Done()`|
|Data races|Shared mutable state|Use immutability or mutexes when sharing|
|Channel deadlocks|Incorrect send/receive balance|Use buffered channels or fan-out with care|
|Non-deterministic tests|Async timing issues|Use `sync.WaitGroup`, channels, or mocks to control flow|
|Overusing goroutines|Thinking each task needs one|Prefer pooled execution for I/O-bound work|

### 🧰 Tools & Concepts
|Concept|Usage|
|---|---|
|`go func()`|Launch lightweight concurrent unit|
|`chan T`|Synchronised communication between routines|
|`context.Context`|Deadline, cancellation, and propagation|
|`sync.WaitGroup`|Deterministic test orchestration|
|`select`|Multiplex multiple channel operations|
|`time.After`, `time.Ticker`|Timers and periodic scheduling|

### 📁 Directory Suggestions for This Phase
You may introduce:

```bash
internal/bank/
├── async/
│   ├── processor.go       # concurrent transfer processor
│   └── ...
```
## Part 7: Interface Composition, APIs, and I/O Integration

### 🎯 Goal

 * Defining clear, decoupled interfaces
 * Designing a public API for your banking processor
 * Handling I/O and interaction boundaries like HTTP, CLI, or gRPC
 * Testing integrations using mock services and dependency injection

### ✅ Requirements

1. Design an API-friendly Interface. Think of an external caller that wants to:
  * Submit a transfer
  * Get account statements
  * Possibly subscribe to events or results
  * Design the TransferEngine interface that wraps your processor like:
      ```go
      type TransferEngine interface {
        SubmitTransfer(fromID, toID string, amount float64, ref string) error
        GetStatement(accountID string) ([]bank.Transaction, error)
      }
      ```
    to:
      * Encapsulates your internal types
      * Easy to expose via HTTP or CLI
      * Hides processor and account structs behind strings

2. Account Registry
    * Map from accountID string → BankAccount object
    * Handle unknown accounts
    * Return proper errors if not found

3. Hook Up I/O Boundary (e.g. HTTP Handler or CLI)
    * Accept a JSON payload or CLI command
    * Decode it
    * Call `TransferEngine.SubmitTransfer(...)`
    * HTTP API:
      * Register account: 
      ```
        curl -X POST http://localhost:8080/accounts -d "{\"id\":\"cliuser\", \"balance\":999}" -H "Content-Type: application/json"
      ```
      * Get Statement: 
        ```
        curl http://localhost:8080/statement/cliuser`
        ```
      * Make Transfer: 
        ```
        curl -X POST http://localhost:8080/transfer \
          -H "Content-Type: application/json" \
          -d '{
            "from_id": "cliuser1",
            "to_id": "cliuser2",
            "amount": 500,
            "reference": "cli transfer"
        }'
        ```

4. Support Observability
    * Log when transfers are submitted
    * Optionally, expose results or errors via a callback or audit

5. Write Integration Tests
    * Submit a transfer through the API
    * Assert the account balances and statements

### ⚠️ Gotchas to Watch For
| # |Area|Gotcha|
| --- | --- | --- |
| 1 |Account lookups|Don't panic on unknown ID|
| 2 |API layer|Avoid leaking internal types like BankAccount|
| 3 |Tests|Don't forget to test error paths (missing account, bad amount)|
| 4 |Concurrency|Account registry access must be thread-safe if shared|

### Clean approach to architecture

At this stage, the project follows the Hexagonal Architecture to ensure clear separation of concerns, testability and maintainability.

#### Principles

  * Core logic is independent of frameworks, databases, or transport layers
  * Dependencies point inward: outer layers depend on interfaces defined in inner layers
  * Adapters implement interfaces — not the other way around
  * Easy to test and replace parts in isolation

#### Layer Overview

|Layer|Description|Example Packages
| --- | --- | ---
|Core|Defines domain interfaces (ports), orchestrators, and business rules|`internal/engine`, `internal/ports`
|Adapters|Implement core interfaces using concrete logic (banking, async processing)|`internal/bank`, `internal/async`
|Drivers|Expose the system externally (e.g. CLI, HTTP)|cmd/api, cmd/cli 

#### Example Flow: Transfer Execution

1. `TransferEngine.SubmitTransfer(...)` is called (API-friendly interface)
2. It resolves from and to accounts via a `Registry` (core port)
3. It delegates the transfer to an `async.Processor` (adapter)
4. The `Processor` executes a decorated `TransferFunc` (core logic)
5. Results can be observed via a callback or queried via `GetStatement(...)`
