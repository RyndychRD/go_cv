package submodule

import (
	"rsc.io/quote"
)

func Hello(name string) string {
    var message string
    message=quote.Go()+" ..."+name
    return message
}