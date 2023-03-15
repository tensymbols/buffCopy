# buffCopy
UTF-8 data duplicator. Stdin/stdout and file
### Parameters
* all parameters are optional 
* `-from` - path to source file. By default stdin is used as source
* `-to` - copy path. By default stdout is used
* `-offset` - offset (data to skip) in bytes
* `-limit` - limit of read bytes. By default all data from offset is copied
* `-block-size` - Block size in bytes, amount of read and written data per iteration
* `-conv` - Conversions:
   - `upper_case` - convert to uppercase
   - `lower_case` - convert to lowercase
     если указаны и `lower_case`, и `upper_case` - возвращаем ошибку;
   - `trim_spaces` - trim space symbols 
