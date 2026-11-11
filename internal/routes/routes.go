package routes

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"se69-lab04-starter/internal/config"
	"se69-lab04-starter/internal/controllers"
	"se69-lab04-starter/internal/middleware"
	"se69-lab04-starter/internal/utils"
)

func Setup(db *gorm.DB, cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery()) // log every request, turn panics into 500
	r.SetTrustedProxies(nil)                   // do not trust X-Forwarded-For from clients
	r.HandleMethodNotAllowed = true            // wrong method on a known path = 405

	auth := &controllers.AuthController{DB: db, Config: cfg}
	equipment := &controllers.EquipmentController{DB: db}
	loans := &controllers.LoanController{DB: db}
	reports := &controllers.ReportController{DB: db}
	web := &controllers.WebController{DB: db}

	// who can call each route (no role = any logged-in user)
	anyUser := middleware.RequireAuth(cfg.JWTSecret)
	staff := middleware.RequireAuth(cfg.JWTSecret, "staff")
	student := middleware.RequireAuth(cfg.JWTSecret, "student")
	teacher := middleware.RequireAuth(cfg.JWTSecret, "teacher")

	// --- public ---
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.POST("/sessions", auth.Login)

	// --- equipment ---
	r.GET("/equipment", anyUser, equipment.List)
	r.GET("/equipment/:id", anyUser, equipment.Get)
	r.POST("/equipment", staff, equipment.Create)
	r.PUT("/equipment/:id", staff, equipment.Update)
	r.GET("/deleteEquipment/:id", staff, equipment.Delete)

	// --- loans ---
	r.POST("/loans", student, loans.Create)
	r.GET("/loans", anyUser, loans.List)
	r.GET("/loans/:id", anyUser, loans.Get)

	// --- reports ---
	r.GET("/reports/usage", teacher, reports.Usage)

	// --- web app (session cookie instead of Bearer token) ---
	r.POST("/web/sessions", web.Login)
	r.DELETE("/web/sessions", web.Logout)
	r.GET("/web/loans", middleware.RequireSession(db), web.MyLoans)
	r.POST("/web/loans", middleware.RequireSession(db), web.Borrow)

	// unknown path or method: answer in the same error format
	r.NoRoute(func(c *gin.Context) { utils.WriteError(c, http.StatusNotFound, "No route for this path") })
	r.NoMethod(func(c *gin.Context) {
		utils.WriteError(c, http.StatusMethodNotAllowed, "Method not allowed for this path")
	})

	routes := r.Routes()
	sort.Slice(routes, func(i, j int) bool { return routes[i].Path < routes[j].Path })
	for _, route := range routes {
		log.Printf("route: %-6s %s", route.Method, route.Path)
	}
	return r
}
