#Foster

Finds potentially orphaned files in your project.

## Dependencies
go get github.com/cheggaaa/pb
or git submodule update

## Building
`make`

## TODO
- Unit Tests
- Better output
    - Use https://github.com/andrew-d/go-termutil to hide meta output for cleaner piping 
- More robust guessing for determining if a file is used or not
    - Handle relative file paths? 

## Known Issues
- Any top level file will be picked up
- Files are flagged as used via basename bun not path name, this may result in unused file being ignored
    but it's safer for now until relative file paths are implemented

##Warning
This may be useful, but it's a learning project. 
I may or may (probably) not have done things right.

##Double Warning
Make sure --base is set to the root of your project tree for the safest results. 


#License
MIT
