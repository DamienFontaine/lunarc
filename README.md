# Lunarc

[![Build Status](https://travis-ci.org/DamienFontaine/lunarc.svg?branch=master)](https://travis-ci.org/DamienFontaine/lunarc)
[![Coverage Status](https://coveralls.io/repos/github/DamienFontaine/lunarc/badge.svg?branch=master)](https://coveralls.io/github/DamienFontaine/lunarc?branch=master)
[![GoDoc](https://godoc.org/github.com/DamienFontaine/lunarc?status.png)](https://godoc.org/github.com/DamienFontaine/lunarc)

## Download and install

``` sh
$ go get github.com/DamienFontaine/lunarc
$ cd $GOPATH/src/github.com/DamienFontaine/lunarc
$ godep restore
```

## Example

### Serve static files

```
.
+-- config.yml
+-- main.go
+-- Public
|   +-- index.html
```
config.yml
```yml
production:
  server:
    port: 8888
```
main.go
``` go
package main

import (
	"log"
	"net/http"

	"github.com/DamienFontaine/lunarc/web"
)

func main() {
	s, err := web.NewServer("config.yml", "production")
	if err != nil {
		log.Printf("Error: %v", err)
	}
	m := s.Handler.(*web.LoggingServeMux)
	m.Handle("/", http.FileServer(http.Dir("public/")))

	go s.Start()

	select {
	case <-s.Done:
		return
	case <-s.Error:
		log.Println("Error: server terminate")
		return
	}
}
```
public/index.html
```html
<!DOCTYPE html>
<html>
  <head>
<title>Example</title>
  </head>
  <body>
    <h1>Hello Wolrd!</h1>
  </body>
</html>
```
Start the server
```sh
$ go run main.go
```

## License
GNU Affero General Public License version 3: <http://www.gnu.org/licenses/agpl-3.0.txt>
