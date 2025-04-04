# GOLANG tutorial

This is a tutorial entirely designed by ChatGPT. Each part builds on top of the previous. 
Each branch contains the code for the current part.

Understand the goal and requirements, then build the code. 

## Setup

```bash
mkdir golang-tutorial && cd golang-tutorial
go mod init github.com/smartrics/golang-tutorial
```

*Note*: Use your own module name for the project.

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
    - `SavingsAccount`: add ApplyInterest(rate float64)
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
