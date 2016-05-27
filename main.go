package main

import (
	"encoding/json"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/mewben/config-echo"
	"github.com/mewben/db-go-env"
	"projects/onix/controllers"
	"projects/onix/theme"
)

// Initialize Port and DB Connection config
func init() {
	type Config struct {
		SERVERPORT string
		DB         db.Database
	}

	configFile, err := os.Open("env.json")
	if err != nil {
		panic(err)
	}

	var devConfig Config
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&devConfig); err != nil {
		panic(err)
	}

	// setup postgres db connection
	db.Setup(devConfig.DB)

	// setup port
	// This sets the global Port string
	// If you set an environment variable DATABASE_URL,
	// it sets Mode = "prod" and uses the env Port instead
	config.Setup(devConfig.SERVERPORT)
}

func main() {
	app := echo.New()

	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(middleware.Gzip())
	//app.Use(middleware.Static("public"))

	app.Static("/admin", "public/admin")

	// Public
	// Setup Theme
	theme.Setup(app)

	api := app.Group("/api")

	if config.Mode == "dev" {
		// Enable Debug
		app.SetDebug(true)
		cors_enabled := middleware.CORSWithConfig(middleware.CORSConfig{
			AllowHeaders: []string{
				echo.HeaderOrigin,
			},
		})
		// Enable CORS /
		app.Use(cors_enabled)
		// Enable CORS /api
		api.Use(cors_enabled)
	}

	users := controllers.UsersController{}
	app.Post("/auth/login", users.Login)

	// get api routes
	/* api.Use(middleware.JWTAuthWithConfig(middleware.JWTAuthConfig{
		SigningMethod: middleware.AlgorithmHS256,
		SigningKey:    []byte("evdzpwadminsing"),
		Extractor:     middleware.JWTFromHeader,
	})) */
	ApiRoutes(api)

	// ======= SITES =====
	site := controllers.SiteController{}
	app.Get("/", site.Home)

	app.Run(fasthttp.New(config.Port))
}

// ========== API ROUTES ======
func ApiRoutes(api *echo.Group) {

	// ======= ADMIN API ======
	admin := controllers.AdminController{}
	api.Get("/utctime", admin.GetUTCTime)

	// CRUD /api/tags
	/* tags := controllers.TagsController{}
	api.Get("/tags", tags.Get)
	api.Post("/tags", tags.Save)
	api.Put("/tags/:id", tags.Update)
	api.Delete("/tags/:id", tags.Destroy)
	*/
}

/*
package main

import (
	"log"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/mewben/config-echo"
	"github.com/mewben/db-go-env"
	"github.com/rs/cors"
	"github.com/thoas/stats"

	"projects/onix/controllers"
	"projects/onix/theme"
)

func init() {
	// get env
	conn, port := utils.GetEnv()

	// setup DB
	db.Setup(conn)

	// setup config
	config.Setup(port)
}

func main() {
	e := echo.New()

	e.Use(mw.Recover())
	e.Use(mw.Gzip())

	key, err := utils.GetConfigString("admin_signingkey")
	if err != nil {
		panic(err)
	}

	api := e.Group("/api")

	if config.Mode == "dev" {
		e.SetDebug(true)
		e.Use(mw.Logger())

		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
			AllowCredentials: true,
		})
		e.Use(cors.Default().Handler)
		api.Use(c.Handler)
	}

	// Public
	// Setup Theme
	theme.Setup(e)

	e.ServeFile("/admin*", "public/admin/index.html")
	e.ServeFile("/admin/*", "public/admin/index.html")
	e.Static("/admin/build", "public/admin/build")

	// Public Routes
	s := stats.New()
	e.Get("/stats", func(c *echo.Context) error {
		return c.JSON(200, s.Data())
	})
	e.Get("/mew", func(c *echo.Context) error {
		theme.Change("mew", e)
		return c.String(200, "Mew")
	})

	e.Get("/eevee", func(c *echo.Context) error {
		theme.Change("eevee", e)
		return c.String(200, "Eevee")
	})

	site := &controllers.SiteController{}
	e.Get("/", site.Home)

	// ======== API
	content := controllers.ContentController{}
	api.Post("/content", content.Save)
	api.Put("/content/:id", content.Save)

	log.Println("Listening at " + config.Port)
	e.Run(config.Port)
}
*/
