package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/controller"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/domain/repository"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/domain/service"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/store/connection"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/store/postgres"
	"github.com/wittyjudge/todo-api/internal/app/todo-api/usecase"
)

var (
	postgresConn                             = connection.Connect()
	todoRepo       repository.TodoRepository = postgres.NewPostgresTodo(postgresConn)
	todoUsecase    usecase.TodoUsecase       = usecase.NewTodoUsecase(todoRepo)
	todoController controller.TodoController = controller.NewTodoController(todoUsecase)

	userRepo       repository.UserRepository = postgres.NewPostgresUser(postgresConn)
	userService    service.UserService
	userUsecase    usecase.UserUsecase       = usecase.NewUserUsecase(userRepo, userService)
	userController controller.UserController = controller.NewUserController(userUsecase)
)

func main() {
	route := mux.NewRouter()

	api := route.PathPrefix("/api").Subrouter()
	api.HandleFunc("/test", func(rw http.ResponseWriter, req *http.Request) { rw.Write([]byte("TEST")) }).Methods("GET")
	api.HandleFunc("/todos", todoController.AllTodos()).Methods("GET")
	api.HandleFunc("/todos", todoController.CreateTodo()).Methods("POST")
	api.HandleFunc("/todos/{id}", todoController.DeleteTodo()).Methods("DELETE")

	api.HandleFunc("/login", userController.Login()).Methods("POST")
	api.HandleFunc("/signup", userController.Signup()).Methods("POST")

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		// handlers.AllowCredentials(),
	)

	server := &http.Server{
		Addr:         "127.0.0.1:8085",
		Handler:      cors(route),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
