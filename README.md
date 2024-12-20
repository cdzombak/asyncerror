# `github.com/cdzombak/asyncerror`

`asyncerror` provides an easy way to manage errors that occur in a heavily asynchronous Go program.

It is particularly useful if some number of errors are expected or permissible for your application. For example, if your program reads a sensor at high frequency, perhaps a 0.1% read failure rate is tolerable, but a higher error rate indicates a hardware problem.

## Installation

```
go get github.com/cdzombak/asyncerror
```

## Usage

### 1. Create an `asyncerror.Escalator`

Create an `asyncerror.Escalator`. Your application should listen to its `EscalationChannel()`, which will receive an error once an error policy determines something has gone wrong.

Your application typically should have a single `Escalator` instance.

_TK: example code_

### 2. Create and Register an `asyncerror.Policy`

Create and register one or more `asyncerror.Policy` instances with the `Escalator`. Policies determine when an error should be escalated.

Each area of responsibility in your application may have its own policy. For example, code that reads from a sensor may create a different policy than a routine that writes user data disk. 

`asyncerror` includes two built-in policies out of the box:
- `ImmediateEscalationPolicy` escalates every error received immediately.
- `ThresholdEscalationPolicy` escalates an error only if a certain number of errors have occurred within a certain time window.

Register each policy with the `Escalator` using the `Escalator`'s `RegisterPolicy` method. The channel returned by `RegisterPolicy` should be kept by the caller. 

_TK: example code_

### 3. Send errors to a policy

When an error occurs, send it to the appropriate policy's channel. (This channel was returned by `RegisterPolicy`.)

_TK: example code_

## License

MIT; see [LICENSE](LICENSE) for details.

## Author

Chris Dzombak ([dzombak.com](https://www.dzombak.com), [github.com/cdzombak](https://github.com/cdzombak))
