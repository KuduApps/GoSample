package main

import (
    "fmt"
     "log"
     "time"
    "net/http"
    "os" 
    "io" 
    "bufio"
    "github.com/go-martini/martini"
    "github.com/gin-gonic/gin"
)

func main() {
    logFile := "testlogfile"
    port := "3001"
    if os.Getenv("HTTP_PLATFORM_PORT") != "" {
        logFile = "D:\\home\\site\\wwwroot\\testlogfile"
        port = os.Getenv("HTTP_PLATFORM_PORT")
    }

    f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
           
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // fmt.Fprintf(w, "Hello form Go! Error: %v", err)
        fmt.Fprintf(w, `
        <html>
            <body>
                <h1>Hello from Go!</h1>
                <br />
                <a href='/g'>Try Gin</a>
                <br />
                <a href='/m'>Try Martini</a>
                <br />
                <pre>`)
     
        rf, _ := os.Open(logFile)
        defer rf.Close()
        scanner := bufio.NewScanner(rf)
        lineCount := 0
        for scanner.Scan() {
            lineStr := scanner.Text()
            fmt.Fprintf(w, lineStr)
            fmt.Fprintf(w, "<br />")
            lineCount++
        }
        
        fmt.Fprintf(w, "<br />")
        fmt.Fprintf(w, "Log Count: %v/1000", lineCount)
        fmt.Fprintf(w, "<br />")
        fmt.Fprintf(w, `
                </pre>
            </body>
        </html>`)
        
        if lineCount > 1000 {
            wf, _ := os.OpenFile(logFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0777)
            defer wf.Close()
            w := bufio.NewWriter(wf)
            w.WriteString("")
            w.Flush()
        }
    })

    if err != nil {
        http.ListenAndServe(":"+port, nil)
    } else {
         defer f.Close()
         log.SetOutput(f)
         log.Println("--->   UP @ " + port +"  <------")
    }

    m := martini.Classic()
    m.Get("/m/", func() string {
      return `
        <html>
            <body>
                <h1>Hello from Martini!</h1>
                <br />
                <a href='/'>Home</a>
                <br />
                <a href='//github.com/go-martini/martini' target='_blank'>Martini @ Github</a>
            </body>
        </html>`;
    })
    m.Map(log.New(f, "[MARTINI]", log.LstdFlags))
    http.Handle("/m/", m)

    g := gin.Default()
    g.Use(GinLogger(f))
    if os.Getenv("HTTP_PLATFORM_PORT") != "" {
      g.LoadHTMLFiles("D:\\home\\site\\wwwroot\\index-gin.html")
    } else {
      g.LoadHTMLFiles("index-gin.html")
    }
    g.GET("/g/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index-gin.html", gin.H{
            "title" : "Hello from Gin!",
        })
    })
    http.Handle("/g/", g)
    
    http.ListenAndServe(":"+port, nil)
}

func GinLogger(out io.Writer) gin.HandlerFunc {
	stdlogger := log.New(out, "[GIN]", log.LstdFlags)
	
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		
		stdlogger.Printf("%v |%3d| %12v | %s |%-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			c.Request.URL.Path,
			c.Errors.String(),
		)
	}
}
