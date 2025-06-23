package handlers

import (
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrganization returns a list of organizations (GET /api/v1/organization)
func ListOrganization(c *gin.Context) {
	name := c.Query("name")
	code := c.Query("code")
	status := c.Query("status")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page_number", "1"))

	orgs, err := repository.GetOrgList(name, code, status, pageSize, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

// GetOrganizationByID returns a single organization by ID (GET /api/v1/organization/:id)
func GetOrganizationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	org, err := repository.FindOrgByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}
	c.JSON(http.StatusOK, org)
}
