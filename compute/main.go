package main

import (
	"context"
	"fmt"
	"os"

	"github.com/darkweak/go-esi/esi"
	"github.com/fastly/compute-sdk-go/fsthttp"
)

// The entry point for your application.
//
// Use this function to define your main request handling logic. It could be
// used to route based on the request properties (such as method or path), send
// the request to a backend, make completely new requests, and/or generate
// synthetic responses.

func main() {
	// Log service version
	fmt.Println("FASTLY_SERVICE_VERSION:", os.Getenv("FASTLY_SERVICE_VERSION"))

	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		// Filter requests that have unexpected methods.
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "This method is not allowed\n")
			return
		}

		if r.URL.Path == "/" {
			res := esi.Parse(ctx, []byte(body), r)

			w.WriteHeader(fsthttp.StatusOK)
			fmt.Fprint(w, string(res))

			return
		}

		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprint(w, "not found\n")
	})
}

const body = `
<html>

<head>
  <title>
    <esi:vars>Hello from $(HTTP_HOST)</esi:vars>
  </title>
  <esi:remove>
    <esi:include src="https://angusd.com" />
  </esi:remove>
</head>

<body>
  <!--esi
        <esi:include src="domain.com:9080/not-interpreted"/>
        -->
  <esi:include src="https://www.angusd.com" />
  <esi:choose>
    <esi:when test="$(HTTP_COOKIE{group})=='Advanced'">
      <span>
        So Advanced
      </span>
    </esi:when>
    <esi:when test="$(QUERY_STRING{cat})=='yes'">
      <img src="https://cdn2.thecatapi.com/images/TdxQ2VvJK.jpg" />
    </esi:when>
    <esi:otherwise>
      <div>
        That's the wrong choice
      </div>
    </esi:otherwise>
  </esi:choose>
</body>

</html>
`
