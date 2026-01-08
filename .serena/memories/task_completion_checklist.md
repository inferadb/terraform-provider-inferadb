# Task Completion Checklist

Before considering a task complete, ensure:

## 1. Code Quality
- [ ] Run `gofmt -s -w -e .` - Code is properly formatted
- [ ] Run `golangci-lint run` - No linting errors
- [ ] Code follows existing patterns in the codebase

## 2. Testing
- [ ] Run `go test -v -cover -timeout=120s -parallel=4 ./...` - Unit tests pass
- [ ] Add/update tests for new functionality
- [ ] For resource changes, consider acceptance tests

## 3. Documentation
- [ ] Update schema `MarkdownDescription` for new attributes
- [ ] Run docs generation if schemas changed
- [ ] Update examples in `examples/` directory if needed

## 4. Final Verification
```bash
gofmt -s -w -e .
golangci-lint run
go test -v -cover -timeout=120s -parallel=4 ./...
go build -v ./...
```

## For New Resources/Data Sources
- [ ] Add to provider's `Resources()` or `DataSources()` method
- [ ] Create example file in `examples/resources/<name>/resource.tf`
- [ ] Regenerate docs: `go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs`
