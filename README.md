# loca
simple localization utility for text based files (HTML, CSS, Javascript, Bash, etc.)

## Reasoning
I wanted a simple tool to statically localize files, that
could be installed anywhere, and run in a Makefile.

## Files
`loca` depends on the existence of 2 files: a `locali` file
containing the localizations, and a source file to be localized.

### locali
`locali` is a jagged csv file, where the first row contains
the names of all the languages you are localizing for.  The rest
of the file contains the name of each string, and the string in
each language, for example:

    eng,spa,port
    h,ham,jamón,presunto
    e,eggs,huevos,ovos
    ...

### Source Files
Source files contain all of the boilerplate, with a special character specifying where to interpolate the strings.  By default, this character is `` ` `` (backtick).  For a simple text document example:

    I like to put `h on my `e's
    but not to put `e on `h, that's gross
    the only thing worse is a `h by itself

## Invocation
if a localfile is not specified, it defaults to locali

    loca [options]... <srcfile> [localfile]
    
    Options
    	-c specify custom SIC (String Interpolation Character) (default "`" (backtick))
    	-l specify language to generate (default "eng")
    	-o output filename (default to lang+basename(src))
    			if filename is -, output to standard output

    Examples
    
    	loca eggs.txt.loc
    		use file locali, language "eng", and backtick to generate eng.eggsandham.txt
    	loca -c + -o eng.txt -l spa eggs.txt.loc localize.csv
    		use file localize.csv, language "spa", and + to generate eng.txt

Using the first invocation, with the files above, the results would be:

    I like to put ham on my eggs's
    but not to put eggs on ham, that's gross
    the only thing worse is a ham by itself