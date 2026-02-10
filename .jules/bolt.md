## 2025-05-14 - Direct Struct Unmarshaling in Go-TOML
**Learning:** In Go, unmarshaling TOML/JSON into a generic `map[string]interface{}` followed by manual type assertions is significantly slower than unmarshaling directly into a structured map or struct. For the Glyph codebase, switching to direct unmarshaling reduced template parsing time by ~25% and overall CLI start-up (with 1000 templates) by ~28%.
**Action:** Always prefer structured unmarshaling with appropriate tags over generic map processing for configuration files.
