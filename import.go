package main

// Import represents a declared import
type Import struct {
	Name     string // the package's name
	Alias    string // the local alias for the package name
	Path     string // import path for the package
	Required bool   // is the import required in the charlatan output?
}

// ImportSet contains all the import declarations encountered
type ImportSet struct {
	imports []*Import
}

// Add inserts the given value into the set if it doesn't already exist
func (r *ImportSet) Add(value *Import) {
	if r.imports == nil {
		r.imports = []*Import{value}
	} else if !r.Contains(value) {
		r.imports = append(r.imports, value)
	}
}

// Contains returns true if the given value is in the set
func (r *ImportSet) Contains(value *Import) bool {
	for _, i := range r.imports {
		if i.Name == value.Name && i.Path == value.Path {
			return true
		}
	}

	return false
}

// GetRequired returns all imports referenced by the target interface
func (r *ImportSet) GetRequired() []*Import {
	result := make([]*Import, 0, len(r.imports))
	for _, imp := range r.imports {
		if imp.Required {
			result = append(result, imp)
		}
	}
	return result
}

// RequireByName marks an import symbol as required
func (r *ImportSet) RequireByName(s string) {
	for i, imp := range r.imports {
		if imp.Name == s || imp.Alias == s {
			r.imports[i].Required = true
		}
	}
}
