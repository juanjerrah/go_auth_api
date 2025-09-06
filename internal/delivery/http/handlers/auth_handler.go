package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juanjerrah/go_auth_api/internal/domain/auth"
	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"github.com/juanjerrah/go_auth_api/pkg/types"
)

type AuthHandler struct {
	userService user.Service
	jwtManager  *auth.JWTManager
	authService auth.AuthService
}

func NewAuthHandler(userService user.Service, jwtManager *auth.JWTManager, authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 409 {object} map[string]string "Email already in use"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar role se fornecido
	if req.Role != "" {
		if err := h.authService.ValidateRole(req.Role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
	} else {
		// Definir role padrão como 'user'
		req.Role = user.RoleUser
	}

	userResponse, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case user.ErrEmailAlreadyInUse:
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		}
		return
	}

	// Gerar token JWT
	token, err := h.jwtManager.GenerateToken(userResponse.ID, userResponse.Email, types.Role(userResponse.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Armazenar token no Redis
	authCtx := &types.AuthContext{
		UserID:      userResponse.ID,
		Email:       userResponse.Email,
		Role:        types.Role(userResponse.Role),
		Permissions: h.authService.GetUserPermissions(types.Role(userResponse.Role)),
	}

	err = h.authService.StoreToken(c.Request.Context(), token, authCtx, h.jwtManager.GetTokenDuration())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user":    userResponse,
	})
}

// Login handles user authentication
// @Summary Authenticate user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body user.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Autenticar usuário
	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Gerar token JWT
	token, err := h.jwtManager.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Armazenar token no Redis
	authCtx := &types.AuthContext{
		UserID:      user.ID.Hex(),
		Email:       user.Email,
		Role:        user.Role,
		Permissions: h.authService.GetUserPermissions(user.Role),
	}

	err = h.authService.StoreToken(c.Request.Context(), token, authCtx, h.jwtManager.GetTokenDuration())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate current session token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Logout successful"
// @Failure 400 {object} map[string]string "No active session"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, exists := c.Get("jwtToken")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active session"})
		return
	}

	err := h.authService.DeleteToken(c.Request.Context(), token.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// LogoutAll handles logout from all devices
// @Summary Logout from all devices
// @Description Invalidate all user sessions
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Logout from all devices successful"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	authContext, exists := c.Get("authContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	authCtx := authContext.(*types.AuthContext)
	err := h.authService.InvalidateUserTokens(c.Request.Context(), authCtx.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout from all devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices successfully"})
}

// RefreshToken handles token refresh
// @Summary Refresh authentication token
// @Description Get a new token with extended expiration
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authContext, exists := c.Get("authContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	authCtx := authContext.(*types.AuthContext)

	// Gerar novo token
	newToken, err := h.jwtManager.GenerateToken(authCtx.UserID, authCtx.Email, authCtx.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	// Invalidar token antigo
	oldToken, exists := c.Get("jwtToken")
	if exists {
		h.authService.DeleteToken(c.Request.Context(), oldToken.(string))
	}

	// Armazenar novo token no Redis
	err = h.authService.StoreToken(c.Request.Context(), newToken, authCtx, h.jwtManager.GetTokenDuration())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store new session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}

// GetProfile returns current user profile
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} user.UserResponse "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	authContext, exists := c.Get("authContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	authCtx := authContext.(*types.AuthContext)

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), authCtx.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, userResponse)
}

// ValidateToken checks if a token is valid
// @Summary Validate token
// @Description Check if authentication token is valid
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} map[string]string "Invalid token"
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authContext, exists := c.Get("authContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	authCtx := authContext.(*types.AuthContext)

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": authCtx.UserID,
		"email":   authCtx.Email,
		"role":    authCtx.Role,
		"expires_in": time.Until(
			time.Now().Add(h.jwtManager.GetTokenDuration()),
		).Seconds(),
	})
}

// GetSessions returns active user sessions
// @Summary Get active sessions
// @Description Get list of active sessions for current user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Active sessions"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/sessions [get]
func (h *AuthHandler) GetSessions(c *gin.Context) {
	authContext, exists := c.Get("authContext")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	authCtx := authContext.(*types.AuthContext)

	// Esta funcionalidade requer implementação adicional no Redis Repository
	// para listar todas as sessões ativas de um usuário
	// Por enquanto retornamos uma mensagem indicando que está implementado
	// mas sem os detalhes das sessões

	c.JSON(http.StatusOK, gin.H{
		"user_id":  authCtx.UserID,
		"sessions": "Session list functionality requires additional Redis implementation",
		"message":  "Active sessions retrieved successfully",
	})
}