package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// ChangePassword updates user's password
// @Summary Change password
// @Description Change a user's password by providing old and new passwords
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body user.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]string "Password changed successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req user.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("id")
	if err := h.userService.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == user.ErrInvalidEmail {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
			return
		}
		if err == user.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// GetUserProfile returns the current user's profile
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} user.UserResponse "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Router /users/profile [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	authContext, _ := c.Get("authContext")
	authCtx := authContext.(*auth.AuthContext)

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), authCtx.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// GetUserByID returns a user by ID
// @Summary Get user by ID
// @Description Get a user's information by their ID
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.UserResponse "User data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Insufficient permissions"
// @Failure 404 {object} map[string]string "User not found"
// @Router /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// UpdateUser updates user information
// @Summary Update user
// @Description Update a user's information
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body user.UpdateUserRequest true "User data to update"
// @Success 200 {object} map[string]string "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Insufficient permissions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	authContext, _ := c.Get("authContext")
	authCtx := authContext.(*auth.AuthContext)

	// Verificar se o usuário está atualizando a si mesmo ou tem permissão
	if authCtx.UserID != userID && !auth.HasPermission(authCtx.Role, auth.PermissionUserWrite) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Insufficient permissions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	authContext, _ := c.Get("authContext")
	authCtx := authContext.(*auth.AuthContext)

	// Apenas admins podem deletar outros usuários
	if authCtx.UserID != userID && !auth.HasPermission(authCtx.Role, auth.PermissionUserDelete) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// func (h *UserHandler) ListUsers(c *gin.Context) {
// 	// Implementar listagem de usuários com paginação
// 	users, err := h.userService.ListUsers(c.Request.Context(), &user.ListUsersRequest{
// 		Page:  1,
// 		Limit: 10,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, users)
// }
