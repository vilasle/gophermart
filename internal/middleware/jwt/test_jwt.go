package jwt

/*
import (
	"fmt"
	"net/http"
	"strings"
)

package jwt

import (
"fmt"
"github.com/Painkiller675/url_shortener_6750/internal/repository"
"net/http"
)

type storHelp struct {
	storage repository.URLStorage
}

func (s storHelp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("implemented!")
}

func JWTMW(storH storHelp) func(h http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("JWT MW Is available!")
			err := storH.storage.Ping(r.Context())
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println(err)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}


func JWTTESTMW(some ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("JWT MW Is available!")
			next.ServeHTTP(w,r)
		})


		}
	}




func JWTM_W_W_W(h http.Handler) http.Handler {
	gzipFunc := func(res http.ResponseWriter, req *http.Request, ss storHelp) {
		fmt.Println("JWT MW 2 Is available!")
		err := ss.storage.Ping(req.Context())
		if err != nil {
			fmt.Println(err)
		}
		h.ServeHTTP(res, req)

	}
	return http.HandlerFunc(gzipFunc)
}

func JWTM_W_W(h http.Handler) http.Handler {
	gzipFunc := func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("JWT MW 2 Is available!")
		err := repository.URLStorage.Ping(req.Context())
		h.ServeHTTP(res, req)

	}
	return http.HandlerFunc(gzipFunc)
}
*/
