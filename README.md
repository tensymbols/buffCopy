# Buffer copy
UTF-8 data duplicator. Stdin/stdout and files
### Parameters
* All parameters are optional 
* `-from` - Path to source file. By default stdin is used as source
* `-to` - Copy path. By default stdout is used
* `-offset` - Offset (data to skip) in bytes
* `-limit` - Limit of read bytes. By default all data from offset is copied
* `-block-size` - Block size in bytes, amount of read and written data per iteration
* `-conv` - Conversions:
   - `upper_case` - Convert to uppercase
   - `lower_case` - Convert to lowercase
   - `trim_spaces` - Trim space symbols 
