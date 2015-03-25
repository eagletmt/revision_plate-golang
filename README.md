# revision_plate-golang
[![Build Status](https://travis-ci.org/eagletmt/revision_plate-golang.svg?branch=master)](https://travis-ci.org/eagletmt/revision_plate-golang)

Serve application's REVISION.

Golang version of [revision_plate](https://github.com/sorah/revision_plate) .

## Usage

```go
import (
	"net/http"

	"github.com/eagletmt/revision_plate-golang"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/site/sha", revision_plate.New("REVISION"))
	return mux
}
```
