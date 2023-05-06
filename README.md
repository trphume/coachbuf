# Coachbuf

Coachbuf is a Go library for data serialization by using bitpacking to minimize size. \
Understand that this is a personal learning project with limitations not meant for real world use case.

## Installation

Run the following command to add Coachbuf package to your project.

```bash
$ go get -u github.com/trphume/coachbuf
```

## Features

* Support 3 data types
    * Int32
    * Float32
    * Struct (including nested struct)
* Minimal data footprint
    * Bitpack Int32 (support two optional struct tag; min and max to specify range of available values)
    * Only metadata used is for ordering number (bitpacked ordering number as well)
* Simple to use
  * Encode and decode function just like JSON serialization package
  * Utilizes struct tags

## Usage

```go
package main

import "github.com/trphume/coachbuf"

type Example struct {
  Int32   int32   `coachbuf:"1,min=0,max=100"`
  Float32 float32 `coachbuf:"2"`
}

func main() {
  inputStruct := Example{Int32: 50, Float32: 100.567}

  data, err := coachbuf.Encode(inputStruct)
  if err != nil { /* Handle error */ }

  result := Example{}
  err = coachbuf.Decode(data, &result)
  if err != nil { /* Handle error */ }
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)