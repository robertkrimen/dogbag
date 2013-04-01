# dogbag
--
Command dogbag is way to bundle data with your Go executable.

(Take your data to Go)

To dogbag a file and use it:

       # Create template.tmpl_dogbag.go from the file template.tmpl:
       $ dogbag template.tmpl .
       ...
       // var template []byte
       template := template_tmpl()

To dogbag a directory and use it

       # Create assets_dogbag.go from the directory ./assets:
       $ dogbag ./assets .
       ...
       // dogbag is a Bag
       dogbag, err := assets()

### Install

     go get github.com/robertkrimen/dogbag/dogbag

### Usage

     dogbag <input> <output> [...]

         <input> Can be a file, directory, or - (stdin)
         If a file, then a filebag will be built
         If a directory, then a zipbag will be built
         If omitted, the default is stdin (and a filebag will be built)

         <output> Can be a file, directory, or - (stdout)
         If a directory, then an <input>.go file will be put in the directory
         If omitted, the default is stdout

         -function=""      The name of the function returning the dogbag
         -package=""       The package to put the dogbag .go file in
         -empty=false      Make an empty (zip) dogbag (for development)
         -fmt=true         Postprocess through gofmt
         -usage=false      More help: Bag interface description, etc.

--
**godocdown** http://github.com/robertkrimen/godocdown
