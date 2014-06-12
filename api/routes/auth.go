package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/validation"
)

func logout(c *AppContext) error {

	if err := c.Logout(); err != nil {
		return err
	}

	return c.OK("Logged out")

}

func authenticate(c *AppContext) error {

	user, err := c.GetCurrentUser()
	if err != nil {
		return err
	}

    var s *session.Session

	if user == nil {
		s = &session.Session{LoggedIn: false}
	} else {
        s = &session.Session{LoggedIn: true,
                              ID: user.ID,
                              Name: user.Name}
    }

	return c.OK(s)
}

func login(c *AppContext) error {

	auth := &session.Authenticator{}
	if err := c.ParseJSON(auth); err != nil {
		return err
	}
	user, err := auth.Identify()
	if err != nil {
		if err == session.MissingLoginFields {
			return c.BadRequest("Missing email or password")
		}
		return err
	}
	if user == nil {
		return c.BadRequest("Invalid email or password")
	}
	if err := c.Login(user); err != nil {
		return err
	}
    s := &session.Session{LoggedIn: true,
                          ID: user.ID,
                          Name: user.Name}
	return c.OK(s)
}

func signup(c *AppContext) error {

	user := &models.User{}

	if err := c.ParseJSON(user); err != nil {
		return err
	}

	// ensure nobody tries to make themselves an admin
	user.IsAdmin = false

	validator := &validation.UserValidator{user}

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			return err
		}
		return c.BadRequest(result)
	}

	if err := userMgr.Insert(user); err != nil {
		return err
	}

	if err := c.Login(user); err != nil {
		return err
	}

	return c.OK(user)

}
