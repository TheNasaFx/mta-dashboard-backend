package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "dashboard-backend/repository"
    "dashboard-backend/database"
)

func GetTaxAuditPapersHandler(c *gin.Context) {
    db := database.DB
    data, err := repository.GetTaxAuditPapers(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, data)
} 