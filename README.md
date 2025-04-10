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
| 1    | [part1](https://github.com/smartrics/golang-tutorial/tree/part1) | Syntax & Language Fundamentals |
| 2    | [part2](https://github.com/smartrics/golang-tutorial/tree/part2) | Structs, Methods, Interfaces |
| 3    | [part3](https://github.com/smartrics/golang-tutorial/tree/part3) | Error Handling, Testing, Tooling |
| 4    | [part4](https://github.com/smartrics/golang-tutorial/tree/part4) | Application Structure, Project Organisation, and CI/CD Readiness |
| 5    | [part5](https://github.com/smartrics/golang-tutorial/tree/part5) | Functional Composition, OO Patterns, and Advanced Struct Techniques |
| 6    | [part6](https://github.com/smartrics/golang-tutorial/tree/part6) | Concurrency, Channels and Project Architecture |
| 7    | [part7](https://github.com/smartrics/golang-tutorial/tree/part7) | Interface Composition, APIs, and I/O Integration |
| 8    | [part8](https://github.com/smartrics/golang-tutorial/tree/part8) | Common Golang gotchas |

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

## Part 1

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

## Part 2

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

## Part 3

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
