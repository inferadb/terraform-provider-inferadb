# Code Style & Conventions

## File Headers
All files must include the copyright and license header:
```go
// Copyright 2025 InferaDB
// SPDX-License-Identifier: Apache-2.0
```

## Package Organization
```
internal/
├── provider/       # Terraform resources and data sources
│   ├── provider.go           # Main provider implementation
│   ├── resource_*.go         # Resource implementations
│   ├── datasource_*.go       # Data source implementations
│   └── *_test.go             # Test files alongside code
└── client/         # HTTP client for InferaDB API
    ├── client.go             # Base HTTP client
    ├── models.go             # API request/response types
    └── <entity>.go           # Per-entity API methods
```

## Terraform Resource Pattern
Each resource follows this structure:
1. Interface compliance variables: `var _ resource.Resource = &TypeResource{}`
2. Resource struct with `client *client.Client` field
3. Model struct with `tfsdk` tags
4. Constructor: `NewTypeResource() resource.Resource`
5. Methods: `Metadata`, `Schema`, `Configure`, `Create`, `Read`, `Update`, `Delete`
6. Optional: `ImportState` for import support

## Naming Conventions
- Resource files: `resource_<entity>.go`
- Data source files: `datasource_<entity>.go`
- Test files: `<entity>_test.go`
- Resource types: `<Entity>Resource`
- Model types: `<Entity>ResourceModel`

## Error Handling
- Use `resp.Diagnostics.AddError()` for Terraform errors
- Check for 404 errors using `apiErr.IsNotFound()` in Read/Delete
- On 404 in Read: call `resp.State.RemoveResource(ctx)`

## Type Mappings
- Use `types.String` from terraform-plugin-framework for all attributes
- Convert with `.ValueString()` to get Go string
- Convert with `types.StringValue()` to set from string
- Use `types.StringNull()` for optional null values

## Comments
- Brief inline comments for interface compliance
- Method comments describing the purpose
- Markdown in schema descriptions for docs generation
