package main

import (
    "fmt"
    "net/http"
    "os"

    "cocopen-backend/config"
    "cocopen-backend/routes"
    "cocopen-backend/middleware"
    "cocopen-backend/utils"
)

func main() {

    config.LoadEnv()

    utils.InitJWTSecret()

    db := config.ConnectDB()
    defer db.Close()

    config.StartCleanupJob(db)

    mux := routes.Setup(db)

    handler := middleware.Recovery(middleware.Cors(mux))

    port := os.Getenv("PORT")
    if port == "" {
        port = os.Getenv("SERVER_PORT")
    }
    if port == "" {
        port = "5012"
    }

    fmt.Printf("Server berjalan di port %s\n", port)
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        panic(fmt.Sprintf("Error menjalankan server: %v", err))
    }
}
