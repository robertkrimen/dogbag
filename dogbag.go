// Package dogbag is way to bundle data with your Go executable.
//
// (Take your data to Go)
//
// This is a placeholder package, the actual command is:
// http://github.com/robertkrimen/dogbag/dogbag
//
//      $ go get http://github.com/robertkrimen/dogbag/dogbag
//
//      $ dogbag -help      # Standard command help
//      $ dogbag -usage     # More help, Bag interface description, etc.
//
// To dogbag a file:
//
//      # Create/overwrite template.tmpl.go
//      $ dogbag template.tmpl .
//       
// To dogbag a directory:
//
//      # Create/overwrite assets.go
//      $ dogbag ./assets .
//
package dogbag

// TODO Try and detect the current package
// FIXME The empty flag being used for a filebag?
