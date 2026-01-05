# Task Completion Checklist

Before considering a task complete, ensure:

## 1. Code Quality
- [ ] Run `make fmt` - Code is properly formatted
- [ ] Run `make lint` - No linting errors
- [ ] Code follows existing patterns in the codebase

## 2. Testing
- [ ] Run `make test` - Unit tests pass
- [ ] Add/update tests for new functionality
- [ ] For resource changes, consider acceptance tests

## 3. Documentation
- [ ] Update schema `MarkdownDescription` for new attributes
- [ ] Run `make docs` if schemas changed
- [ ] Update examples in `examples/` directory if needed

## 4. Final Verification
- [ ] Run `make check` (combines fmt, lint, test)
- [ ] Build succeeds: `make build`

## Quick Command
```bash
make check && make build
```

## For New Resources/Data Sources
- [ ] Add to provider's `Resources()` or `DataSources()` method
- [ ] Create example file in `examples/resources/<name>/resource.tf`
- [ ] Regenerate docs with `make docs`
