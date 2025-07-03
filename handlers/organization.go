package handlers

import (
	"dashboard-backend/auth"
	"dashboard-backend/repository"
	"net/http"
	"strconv"

	"dashboard-backend/database"
	"dashboard-backend/database/model"

	"github.com/gin-gonic/gin"
)

// ListOrganization returns a list of organizations (GET /api/v1/organization)
// Now returns: ID, NAME, OFFICE_CODE, REGNO, KHO_CODE, BUILD_FLOOR, ADDRESS, LNG, LAT
func ListOrganization(c *gin.Context) {
	name := c.Query("name")
	code := c.Query("code")
	status := c.Query("status")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("page_number", "1"))

	orgs, err := repository.GetOrgList(name, code, status, pageSize, pageNumber)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": orgs})
}

// GetOrganizationByID returns a single organization by ID (GET /api/v1/organization/:id)
func GetOrganizationByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	org, err := repository.FindOrgByID(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if org == nil {
		auth.ErrorResponse(c, http.StatusNotFound, "Organization not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": org})
}

// CreateOrganization creates a new organization (POST /api/v1/organization)
func CreateOrganization(c *gin.Context) {
	var org repository.OrgInput
	if err := c.ShouldBindJSON(&org); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	created, err := repository.CreateOrg(org)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": created})
}

// UpdateOrganization updates an organization (PUT /api/v1/organization/:id)
func UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var org repository.OrgInput
	if err := c.ShouldBindJSON(&org); err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}
	updated, err := repository.UpdateOrg(uint(id), org)
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// DeleteOrganization deletes an organization (DELETE /api/v1/organization/:id)
func DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		auth.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	err = repository.DeleteOrg(uint(id))
	if err != nil {
		auth.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}

func ListOrganizations(c *gin.Context) {
	db := database.DB
	rows, err := db.Query("SELECT ID, NAME, REGNO, LNG, LAT FROM GPS.PAY_CENTER")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var orgs []model.Org
	for rows.Next() {
		var o model.Org
		if err := rows.Scan(&o.ID, &o.Name, &o.Regno, &o.Lng, &o.Lat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orgs = append(orgs, o)
	}
	c.JSON(http.StatusOK, orgs)
}
