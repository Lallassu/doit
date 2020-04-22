package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type Header struct {
	Token string `header: token`
}

var timeFormat = "2006-01-02T15:04:00.000Z"

func main() {

	hostPort := flag.String("hostport", "localhost:8443", "host:port to run webserver on.")
	database := flag.String("database", "doit.db", "SQLite database file.")
	tlsCert := flag.String("tlscert", "server.crt", "TLS certificate")
	skipTLS := flag.Bool("skiptls", false, "Skip using TLS, PWA will not be available")
	tlsKey := flag.String("tlskey", "server.key", "TLS key")
	emailHost := flag.String("mailhost", "localhost:25", "Host:port to SMTP server (empty = disable email notifications)")
	emailFrom := flag.String("mailfrom", "doit@example.com", "From email for email reminders")

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	db := &DB{}
	db.Init(*database)

	// Send reminders for alarms.
	go func() {
		if *emailHost == "" {
			return
		}
		for {
			lists := db.AllListsWithItems()
			for _, l := range lists {
				for _, item := range l.Items {
					if item.Time == "" || item.ReminderSent || item.Complete {
						continue
					}

					t, err := time.Parse(timeFormat, item.Time)
					if err != nil {
						fmt.Printf("Error parsing time for item: %d\n", item.ID)
						continue
					}

					if time.Since(t).Seconds() >= 0 {
						// Send email to owner and shares.
						acc := db.GetAccount(item.AccountID)
						emails := []string{acc.Email}
						for _, e := range l.Share {
							emails = append(emails, e.Email)
						}

						SendMail(*emailHost,
							*emailFrom,
							fmt.Sprintf("[DoIt]%s", item.Title),
							fmt.Sprintf("This is a reminder for: %s", item.Title),
							emails,
						)
						fmt.Printf("Send email alert for '%s' to %v\n", item.Title, emails)

						item.ReminderSent = true
						db.UpdateItemSent(&item, acc)
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	router.Use(static.Serve("/", static.LocalFile("./dist", false)))

	router.GET("/hasadm", func(c *gin.Context) {
		if ok := db.HasAdminAccount(); ok {
			c.JSON(http.StatusOK, true)
		} else {
			c.JSON(http.StatusOK, false)
		}
	})

	router.POST("/signup", func(c *gin.Context) {
		// Just to make it slower to spam for now.
		// TBD: Throttle requests
		time.Sleep(time.Duration(2) * time.Second)

		acc := struct {
			User  string
			Pass  string
			Email string
			Adm   bool
		}{}

		// Make sure we only allow creation of one admin.
		if db.HasAdminAccount() {
			acc.Adm = false
		}

		if err := c.BindJSON(&acc); err == nil {
			if ok := db.CreateUser(acc.User, acc.Pass, acc.Email, acc.Adm); ok {
				c.JSON(http.StatusOK, true)
			} else {
				c.JSON(http.StatusOK, false)
			}
		}
	})

	router.POST("/login", func(c *gin.Context) {
		// Just to make it harder to test passwords for now
		// TBD: Throttle requests
		time.Sleep(time.Duration(1) * time.Second)

		acc := struct {
			User string
			Pass string
		}{}

		if err := c.BindJSON(&acc); err == nil {
			if token, ok := db.ValidateLogin(acc.User, acc.Pass); ok {
				c.JSON(http.StatusOK, token)
			} else {
				c.JSON(http.StatusUnauthorized, "")
			}
		}
	})

	api := router.Group("/api", func(c *gin.Context) {
		h := Header{}

		if err := c.ShouldBindHeader(&h); err == nil {
			if acc, ok := db.ValidateToken(h.Token); ok {
				c.Set("account", acc)
				c.Set("token", h.Token)
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, err)
			}
		}
	})

	api.PUT("/validate/:id", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			if id, err := strconv.Atoi(c.Param("id")); err == nil {
				db.AdminValidateAccount(id, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.PUT("/removeaccount/:id", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			if id, err := strconv.Atoi(c.Param("id")); err == nil {
				db.AdminRemoveAccount(id, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.PUT("/removetoken/:id", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			if id, err := strconv.Atoi(c.Param("id")); err == nil {
				db.AdminRemoveToken(id, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.GET("/allaccounts", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			accounts := db.AdminAllAccounts(acc.(*Account))
			c.JSON(http.StatusOK, &accounts)
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.GET("/alltokens", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			tokens := db.AdminAllTokens(acc.(*Account))
			c.JSON(http.StatusOK, &tokens)
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.GET("/lists", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			lists := db.AllLists(acc.(*Account))
			c.JSON(http.StatusOK, &lists)
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.GET("/refresh/:id/:uts", func(c *gin.Context) {
		id := c.Param("id")
		uts := c.Param("uts")

		lastUts, err := strconv.Atoi(uts)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		listId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		lastUpdated := db.LastUpdate(listId)
		if int(lastUpdated.Unix()) > lastUts {
			c.JSON(http.StatusOK, "update")
		} else {
			c.JSON(http.StatusOK, "")
		}
	})

	api.GET("/items/:id", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			id := c.Param("id")
			listId, err := strconv.Atoi(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, "")
				return
			}
			items := db.ItemsForList(listId, acc.(*Account))
			c.JSON(http.StatusOK, items)
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/list", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			list := List{}
			if c.ShouldBind(&list) == nil {
				db.SaveList(&list, acc.(*Account))
				c.JSON(http.StatusOK, list)
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/sort", func(c *gin.Context) {
		items := []Item{}
		if err := c.BindJSON(&items); err == nil {
			// TBD: Check if we own items.
			db.UpdateSortOrder(items)
			c.JSON(http.StatusOK, nil)
		}
	})

	api.PUT("/rename", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			list := List{}
			if err := c.BindJSON(&list); err == nil {
				db.RenameList(list, acc.(*Account))
				c.JSON(http.StatusOK, nil)
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/sharelist", func(c *gin.Context) {
		inp := struct {
			User string
			List int
		}{}
		if acc, ok := c.Get("account"); ok {
			if c.BindJSON(&inp) == nil {
				user := db.ShareList(inp.User, inp.List, acc.(*Account))
				if user != nil {
					c.JSON(http.StatusOK, user)
				} else {
					// Just dont send anything if error or not exists.
					c.JSON(http.StatusOK, "")
				}
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/favorite", func(c *gin.Context) {
		inp := struct {
			Favorite bool
			List     int
		}{}
		if acc, ok := c.Get("account"); ok {
			if c.BindJSON(&inp) == nil {
				db.Favorite(inp.Favorite, inp.List, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/removesharelist", func(c *gin.Context) {
		inp := struct {
			User string
			List int
		}{}
		if acc, ok := c.Get("account"); ok {
			if c.BindJSON(&inp) == nil {
				user := db.RemoveShareList(inp.User, inp.List, acc.(*Account))
				if user != nil {
					c.JSON(http.StatusOK, user)
				} else {
					c.JSON(http.StatusBadRequest, "")
				}
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/deletecompleted", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			list := List{}
			if c.BindJSON(&list) == nil {
				db.DeleteCompleted(&list, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/deletelist", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			list := List{}
			if c.BindJSON(&list) == nil {
				db.DeleteList(&list, acc.(*Account))
				c.JSON(http.StatusOK, "")
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/item", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			item := Item{}
			if c.ShouldBind(&item) == nil {
				db.SaveItem(&item, acc.(*Account))
				c.JSON(http.StatusOK, item)
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.PUT("/item", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			item := Item{}
			if c.ShouldBind(&item) == nil {
				db.UpdateItem(&item, acc.(*Account))
			} else {
				c.JSON(http.StatusBadRequest, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	api.POST("/logout", func(c *gin.Context) {
		if acc, ok := c.Get("account"); ok {
			if tok, ok := c.Get("token"); ok {
				db.Logout(tok.(string), acc.(*Account))
				c.JSON(http.StatusOK, "")
			}
		} else {
			c.JSON(http.StatusUnauthorized, "")
		}
	})

	if *skipTLS {
		log.Fatal(http.ListenAndServe(*hostPort, router))
	} else {
		log.Fatal(http.ListenAndServeTLS(*hostPort, *tlsCert, *tlsKey, router))
	}
}

func SendMail(addr, from, subject, body string, to []string) error {
	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
