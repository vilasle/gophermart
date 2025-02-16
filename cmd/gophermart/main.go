package main

func main() {

	// init router
	/*	r := chi.NewRouter()


			// chi routing (to exclude some jwt middleware influence from register and login endpoints)
			r.Route("/api/user", func(r chi.Router) {
				// set logger for chi-route Group
				r.Use(l.LogMW) // TODO: implement logger middleware
				r.Use(gzipMW.GzMW) // TODO: remove it from here????

				r.Post("/register", c.UserRegister())
				r.Post("/login", c.UserLogin())

			}

			r.Route("/api/user", func(r chi.Router) {
				// set logger for chi-route Group
				r.Use(l.LogMW) // TODO: implement logger middleware
				r.Use(gzipMW.GzMW)
				r.Use(jwt) // TODO: implement jwt middleware

				r.Get("/orders", c.ListOrdersRelatedWithUser())
				r.Post("/orders", c.RelateOrderWithUser())
				r.Get("/balance", c.BalanceStateByUser())
				r.Post("/withdraw", c.Withdraw())
				r.Get("/withdrawals", c.ListOfWithdrawals())

			}

			r.Route("/api", func(r chi.Router) {
				// set logger for chi-route Group
				r.Use(l.LogMW) // TODO: implement logger middleware
				r.Use(gzipMW.GzMW) // TODO: remove it or not????????????

				r.Get("/orders/{number}", c.SOMEHANDLER())


			}
		}

	*/
}
