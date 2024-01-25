package buffalo

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/flect/name"
	"github.com/gorilla/handlers"
)

const (
	// AssetsAgeVarName is the ENV variable used to specify max age when ServeFiles is used.
	AssetsAgeVarName = "ASSETS_MAX_AGE"
)

// These method functions will be moved to Home structure.

// GET maps an HTTP "GET" request to the path and the specified handler.
func (a *App) GET(p string, h Handler) *RouteInfo {
	return a.addRoute("GET", p, h)
}

// POST maps an HTTP "POST" request to the path and the specified handler.
func (a *App) POST(p string, h Handler) *RouteInfo {
	return a.addRoute("POST", p, h)
}

// PUT maps an HTTP "PUT" request to the path and the specified handler.
func (a *App) PUT(p string, h Handler) *RouteInfo {
	return a.addRoute("PUT", p, h)
}

// DELETE maps an HTTP "DELETE" request to the path and the specified handler.
func (a *App) DELETE(p string, h Handler) *RouteInfo {
	return a.addRoute("DELETE", p, h)
}

// HEAD maps an HTTP "HEAD" request to the path and the specified handler.
func (a *App) HEAD(p string, h Handler) *RouteInfo {
	return a.addRoute("HEAD", p, h)
}

// OPTIONS maps an HTTP "OPTIONS" request to the path and the specified handler.
func (a *App) OPTIONS(p string, h Handler) *RouteInfo {
	return a.addRoute("OPTIONS", p, h)
}

// PATCH maps an HTTP "PATCH" request to the path and the specified handler.
func (a *App) PATCH(p string, h Handler) *RouteInfo {
	return a.addRoute("PATCH", p, h)
}

// Redirect from one URL to another URL. Only works for "GET" requests.
func (a *App) Redirect(status int, from, to string) *RouteInfo {
	return a.GET(from, func(c Context) error {
		return c.Redirect(status, to)
	})
}

// Mount mounts a http.Handler (or Buffalo app) and passes through all requests to it.
//
//	func muxer() http.Handler {
//		f := func(res http.ResponseWriter, req *http.Request) {
//			fmt.Fprintf(res, "%s - %s", req.Method, req.URL.String())
//		}
//		mux := mux.NewRouter()
//		mux.HandleFunc("/foo", f).Methods("GET")
//		mux.HandleFunc("/bar", f).Methods("POST")
//		mux.HandleFunc("/baz/baz", f).Methods("DELETE")
//		return mux
//	}
//
//	a.Mount("/admin", muxer())
//
//	$ curl -X DELETE http://localhost:3000/admin/baz/baz
func (a *App) Mount(p string, h http.Handler) {
	prefix := path.Join(a.Prefix, p)
	path := path.Join(p, "{path:.+}")
	a.ANY(path, WrapHandler(http.StripPrefix(prefix, h)))
}

// ServeFiles maps an path to a directory on disk to serve static files.
// Useful for JavaScript, images, CSS, etc...
/*
	a.ServeFiles("/assets", http.Dir("path/to/assets"))
*/
func (a *App) ServeFiles(p string, root http.FileSystem) {
	path := path.Join(a.Prefix, p)
	a.filepaths = append(a.filepaths, path)

	h := stripAsset(path, a.fileServer(root), a)
	a.router.PathPrefix(path).Handler(h)
}

func (a *App) fileServer(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	baseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			eh := a.ErrorHandlers.Get(http.StatusNotFound)
			eh(http.StatusNotFound, fmt.Errorf("could not find %s", r.URL.Path), a.newContext(RouteInfo{}, w, r))
			return
		}

		stat, _ := f.Stat()
		maxAge := envy.Get(AssetsAgeVarName, "31536000")
		w.Header().Add("ETag", fmt.Sprintf("%x", stat.ModTime().UnixNano()))
		w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%s", maxAge))
		fsh.ServeHTTP(w, r)
	})

	if a.CompressFiles {
		return handlers.CompressHandler(baseHandler)
	}

	return baseHandler
}

type newable interface {
	New(Context) error
}

type editable interface {
	Edit(Context) error
}

// Resource maps an implementation of the Resource interface
// to the appropriate RESTful mappings. Resource returns the *App
// associated with this group of mappings so you can set middleware, etc...
// on that group, just as if you had used the a.Group functionality.
//
// Resource automatically creates a URL `/resources/new` if the resource
// has a function `New()`. So it could act as a restriction for the value
// of `resource_id`. URL `/resources/new` will always show the resource
// creation page instead of showing the resource called `new`.
/*
	a.Resource("/users", &UsersResource{})

	// Is equal to this:

	ur := &UsersResource{}
	g := a.Group("/users")
	g.GET("/", ur.List) // GET /users => ur.List
	g.POST("/", ur.Create) // POST /users => ur.Create
	g.GET("/new", ur.New) // GET /users/new => ur.New
	g.GET("/{user_id}", ur.Show) // GET /users/{user_id} => ur.Show
	g.PUT("/{user_id}", ur.Update) // PUT /users/{user_id} => ur.Update
	g.DELETE("/{user_id}", ur.Destroy) // DELETE /users/{user_id} => ur.Destroy
	g.GET("/{user_id}/edit", ur.Edit) // GET /users/{user_id}/edit => ur.Edit
*/
func (a *App) Resource(p string, r Resource) *App {
	g := a.Group(p)

	if mw, ok := r.(Middler); ok {
		g.Use(mw.Use()...)
	}

	p = "/"

	rv := reflect.ValueOf(r)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	rt := rv.Type()
	resourceName := rt.Name()
	handlerName := fmt.Sprintf("%s.%s", rt.PkgPath(), resourceName) + ".%s"

	n := strings.TrimSuffix(rt.Name(), "Resource")
	paramName := name.New(n).ParamID().String()

	type paramKeyable interface {
		ParamKey() string
	}

	if pk, ok := r.(paramKeyable); ok {
		paramName = pk.ParamKey()
	}

	spath := path.Join(p, "{"+paramName+"}")

	// This order will become the order of route evaluation too.
	setFuncKey(r.List, fmt.Sprintf(handlerName, "List"))
	g.GET(p, r.List).ResourceName = resourceName

	setFuncKey(r.Create, fmt.Sprintf(handlerName, "Create"))
	g.POST(p, r.Create).ResourceName = resourceName

	// NOTE: it makes restriction that resource id cannot be 'new'.
	if n, ok := r.(newable); ok {
		setFuncKey(n.New, fmt.Sprintf(handlerName, "New"))
		g.GET(path.Join(p, "new"), n.New).ResourceName = resourceName
	}

	setFuncKey(r.Show, fmt.Sprintf(handlerName, "Show"))
	g.GET(path.Join(spath), r.Show).ResourceName = resourceName

	setFuncKey(r.Update, fmt.Sprintf(handlerName, "Update"))
	g.PUT(path.Join(spath), r.Update).ResourceName = resourceName

	setFuncKey(r.Destroy, fmt.Sprintf(handlerName, "Destroy"))
	g.DELETE(path.Join(spath), r.Destroy).ResourceName = resourceName

	if n, ok := r.(editable); ok {
		setFuncKey(n.Edit, fmt.Sprintf(handlerName, "Edit"))
		g.GET(path.Join(spath, "edit"), n.Edit).ResourceName = resourceName
	}

	g.Prefix = path.Join(g.Prefix, spath)
	g.prefix = g.Prefix

	return g
}

// ANY accepts a request across any HTTP method for the specified path
// and routes it to the specified Handler.
func (a *App) ANY(p string, h Handler) {
	a.GET(p, h)
	a.POST(p, h)
	a.PUT(p, h)
	a.PATCH(p, h)
	a.HEAD(p, h)
	a.OPTIONS(p, h)
	a.DELETE(p, h)
}

// Group creates a new `*App` that inherits from it's parent `*App`.
// This is useful for creating groups of end-points that need to share
// common functionality, like middleware.
/*
	g := a.Group("/api/v1")
	g.Use(AuthorizeAPIMiddleware)
	g.GET("/users, APIUsersHandler)
	g.GET("/users/:user_id, APIUserShowHandler)
*/
func (a *App) Group(groupPath string) *App {
	// TODO: move this function to app.go or home.go eventually.
	g := New(a.Options)
	// keep them for v0 compatibility
	g.Prefix = path.Join(a.Prefix, groupPath)
	g.Name = g.Prefix

	// for Home structure
	g.prefix = path.Join(a.prefix, groupPath)
	g.host = a.host
	g.name = g.prefix

	g.router = a.router
	g.RouteNamer = a.RouteNamer
	g.Middleware = a.Middleware.clone()
	g.ErrorHandlers = a.ErrorHandlers

	g.app = a.app  // will replace g.root
	g.root = g.app // will be deprecated

	// to be replaced with child Homes. currently, only used in grifts.
	a.children = append(a.children, g)
	return g
}

// VirtualHost creates a new `*App` that inherits from it's parent `*App`.
// All pre-configured things on the parent App such as middlewares will be
// applied, and can be modified only for this child App.
//
// This is a multi-homing feature similar to the `VirtualHost` in Apache
// or multiple `server`s in nginx. One important different behavior is that
// there is no concept of the `default` host in buffalo (at least for now)
// and the routing decision will be made with the "first match" manner.
// (e.g. if you have already set the route for '/' for the root App before
// setting up a virualhost, the route of the root App will be picked up
// even if the client makes a request to the specified domain.)
/*
	a.VirtualHost("www.example.com")
	a.VirtualHost("{subdomain}.example.com")
	a.VirtualHost("{subdomain:[a-z]+}.example.com")
*/
func (a *App) VirtualHost(h string) *App {
	g := a.Group("/")
	g.host = h
	g.router = a.router.Host(h).Subrouter()

	return g
}

// RouteHelpers returns a map of BuildPathHelper() for each route available in the app.
func (a *App) RouteHelpers() map[string]RouteHelperFunc {
	rh := map[string]RouteHelperFunc{}
	for _, route := range a.Routes() {
		cRoute := route
		rh[cRoute.PathName] = cRoute.BuildPathHelper()
	}
	return rh
}

func (e *Home) addRoute(method string, url string, h Handler) *RouteInfo {
	// NOTE: lock the root app (not this app). only the root has the affective
	// routes list.
	e.app.moot.Lock()
	defer e.app.moot.Unlock()

	url = path.Join(e.prefix, url)
	url = e.app.normalizePath(url)
	name := e.app.RouteNamer.NameRoute(url)

	hs := funcKey(h)
	r := &RouteInfo{
		Method:      method,
		Path:        url,
		HandlerName: hs,
		Handler:     h,
		App:         e.appSelf, // CHKME: to be replaced with Home
		Aliases:     []string{},
	}

	r.MuxRoute = e.router.Handle(url, r).Methods(method)
	r.Name(name)

	routes := e.app.Routes()
	routes = append(routes, r)
	// NOTE: sorting is fancy but we lose the evaluation order information
	// of routing decision. Let's keep the routes as registered order so
	// developers can easily evaluate the order with `buffalo routes` and
	// can debug any routing priority issue. (just keep the original line
	// as history reference)
	//sort.Sort(routes)

	e.app.routes = routes

	return r
}

func stripAsset(path string, h http.Handler, a *App) http.Handler {
	if path == "" {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		up := r.URL.Path
		up = strings.TrimPrefix(up, path)
		up = strings.TrimSuffix(up, "/")

		u, err := url.Parse(up)
		if err != nil {
			eh := a.ErrorHandlers.Get(http.StatusBadRequest)
			eh(http.StatusBadRequest, err, a.newContext(RouteInfo{}, w, r))
			return
		}

		r.URL = u
		h.ServeHTTP(w, r)
	})
}
