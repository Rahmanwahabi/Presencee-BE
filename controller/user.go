package controller

import (
	"net/http"

	"presensee_project/model/payload"
	"presensee_project/usecase"
	"presensee_project/utils"
	"presensee_project/utils/jwt_service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService usecase.UserService
	jwtService  jwt_service.JWTService
}

func NewUserController(userService usecase.UserService, jwtService jwt_service.JWTService) *UserController {
	return &UserController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (u *UserController) SignUpUser(c echo.Context) error {
	user := new(payload.UserSignUpRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := u.userService.SignUpUser(c.Request().Context(), user)
	if err != nil {
		switch err {
		case utils.ErrUsernameAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "success creating user",
		"data":    user,
	})
}

func (u *UserController) LoginUser(c echo.Context) error {
	user := new(payload.UserLoginRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	token, err := u.userService.LogInUser(c.Request().Context(), user)
	if err != nil {
		switch err {
		case utils.ErrInvalidCredentials:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success login",
		"token":   token,
	})
}

func (u *UserController) GetSingleUser(c echo.Context) error {
	userIDParam := c.Param("user_id")
	userID64, err := strconv.ParseUint(userIDParam, 10, 0)
	userID := uint(userID64)

	user, err := u.userService.GetSingleUser(c.Request().Context(), userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	claims := u.jwtService.GetClaims(&c)
	role := claims["role"].(string)

	switch {
	case role == "pegawai":
		fallthrough
	case role == "admin":
		return c.JSON(http.StatusOK, echo.Map{
			"message": "success getting user",
			"data":    user,
		})
	default:
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}
}

func (u *UserController) GetBriefUsers(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	role := claims["role"].(string)

	if role == "pegawai" {
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}

	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrInvalidNumber.Error())
	}

	users, err := u.userService.GetBriefUsers(c.Request().Context(), int(pageInt), int(limitInt))
	if err != nil {
		if err == utils.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success get users",
		"data":    users,
		"meta": echo.Map{
			"page":  pageInt,
			"limit": limitInt,
		},
	})
}

func (u *UserController) UpdateUser(c echo.Context) error {
	claims := u.jwtService.GetClaims(&c)
	userID := uint(claims["user_id"].(float64))

	user := new(payload.UserUpdateRequest)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, utils.ErrBadRequestBody.Error())
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := u.userService.UpdateUser(c.Request().Context(), userID, user)
	if err != nil {
		switch err {
		case utils.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		case utils.ErrUsernameAlreadyExist:
			fallthrough
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	userGet, err := u.userService.GetSingleUser(c.Request().Context(), userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update user",
		"data":    userGet,
	})
}

func (d *UserController) DeleteUser(c echo.Context) error {
	claims := d.jwtService.GetClaims(&c)
	role := claims["role"].(string)
	if role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, utils.ErrDidntHavePermission.Error())
	}
	userIDParam := c.Param("user_id")
	userID64, err := strconv.ParseUint(userIDParam, 10, 0)
	userID := uint(userID64)

	err = d.userService.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		switch err {
		case utils.ErrUserNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success deleting user",
	})
}
