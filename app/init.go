package app

import (
	"github.com/revel/revel"
	"github.com/cbonello/revel-csrf"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		//added for CSRF protection
		csrf.CSRFFilter,              // CSRF prevention.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		AuthFilter,
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)

	revel.TemplateFuncs["animalFaceEmoji"] = func(index int) string {
		var animalFaceEmojis = [...]rune{
			'🐵', '🐶', '🐺', '🦊', '🐱',
			'🦁', '🐯', '🐴', '🦄', '🐮',
			'🐷', '🐗', '🐭', '🐹', '🐰',
			'🐻', '🐨', '🐼', '🐔', '🐤',
			'🐦', '🐧', '🐸', '🐲', '🐙',
		}

		// mod to allow wrap around if the index is out of bounds
		return string(animalFaceEmojis[index%25])
	}

	revel.TemplateFuncs["indexToCharacter"] = func(index int) string {
		character := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		// mod to allow wrap around if the index is out of bounds
		return string(character[index%26])
	}
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	if revel.RunMode == "prod" {
		c.Response.Out.Header().Add("Strict-Transport-Security","max-age=31536000; includeSubDomains; preload")
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

var AuthFilter = func(c *revel.Controller, fc []revel.Filter) {

	gc := revel.Config.StringDefault("e6e.google_client", "empty")
	c.ViewArgs["gc"] = gc

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
