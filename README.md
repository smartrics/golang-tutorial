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

|     | Gotcha| Description |
| --- |  ---  |    ---      |
| ❌  | Pointer vs Value Receiver   | Value receivers won't mutate original struct             |
| ❌  | String-based error matching | Use `errors.Is()` instead                                |
| ❌  | Range copies values         | Use `&slice[i]` or index access                          |
| ❌  | Defer order & scope         | Defer runs LIFO after return                             |
| ❌  | Ignoring errors             | Always check returned errors                             |
| ❌  | Interface `nil != nil`      | Use `errors.Is()` or `errors.As()`                       |
| ❌  | Lowercase fields            | Not exported outside package                             |
| ✅  | Zero values work            | e.g. `Account{}` is safe unless you enforce constructors |

