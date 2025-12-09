package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getCompanyIDFromContext extracts and validates company ID from context
func getCompanyIDFromContext(c *gin.Context) (uuid.UUID, error) {
	companyIDStr, exists := c.Get("company_id")
	if !exists {
		return uuid.Nil, http.ErrNoCookie // Use a standard error or generic one
	}

	companyID, err := uuid.Parse(companyIDStr.(string))
	if err != nil {
		return uuid.Nil, err
	}

	return companyID, nil
}

// getCompanyIDOrError extracts company ID and returns error response if invalid
func getCompanyIDOrError(c *gin.Context) (uuid.UUID, bool) {
	companyID, err := getCompanyIDFromContext(c)
	if err != nil || companyID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id not found or invalid"})
		return uuid.Nil, false
	}
	return companyID, true
}
