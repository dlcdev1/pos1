package main

import (
	"github.com/dlcdev1/pos1/9-apis/api/configs"
	_ "github.com/dlcdev1/pos1/9-apis/api/docs"
	entity2 "github.com/dlcdev1/pos1/9-apis/api/internal/entity"
	database2 "github.com/dlcdev1/pos1/9-apis/api/internal/infra/database"
	handlers2 "github.com/dlcdev1/pos1/9-apis/api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title GoExample API
// @version 1.0
// @description This is a sample server Petstore server.

// @termsOfService http://swagger.io/terms/
// @contact.name Dlcdev
// @Contact.url dlc@

// @license.name Dlcdev License
// @license.url http://dlcdev@

// @host localhost:8000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.AutoMigrate(&entity2.Product{}, &entity2.User{})
	productDB := database2.NewProduct(db)
	productHandler := handlers2.NewProductHandler(productDB)

	userDB := database2.NewUser(db)
	userHandler := handlers2.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // não deixa a aplicação cair
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JwtExpiresIn))
	//r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)

	})

	r.Post("/users", userHandler.Create)
	r.Get("/users/generate_token", userHandler.GetJWT)
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

// reqeust -> Middleware(usa os dados, faz alguma coisa)| outro middle
//-> Handler -> Response

// client -> request -> http.handle -> response -> client
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
