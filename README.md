# handler

**Note:** **The package is not "go get" able**, just copy and paste it to your project and adjust it to your needs.

This custom handler lets you combine routes definition with your handler definition, do authentication and parameters validation in a clean expressive way, and stop thinking about handlers names :relieved:.

Table of contents
---
- [Example](#example)
- [Struct](#struct)
- [Routing](#routing)
- [Validation](#validation)
- [Authentication](#authentication)
- [Middlewares](#middlewares)

Example
---
```go
var Posts = []handler{
	handler{
		Action: "New Post Page",
		Method: "GET",
		URI:    "/posts/new",
		Auth:   true,
		Do: func(w http.ResponseWriter, r *http.Request) {
			// ...
		}},
	handler{
		Action: "Create New Post",
		Method: "POST",
		URI:    "/posts",
		Params: []string{"title:max=200"},
		Do: func(w http.ResponseWriter, r *http.Request) {
			// ...
		}},
	handler{
		Action: "Index Posts",
		Method: "GET",
		URI:    "/posts",
		Do: func(w http.ResponseWriter, r *http.Request) {
			// ...
		}},
}

func main() {
	mux := NewServeMux(Posts)
	http.ListenAndServe(":8000", mux)
}
```

Struct
---
```go
type Handler struct {
	Action string             // Optional: Description for your handler (used in API file generation).
	Method string             // Optional: Request HTTP method. If not specified, it will accept any method.
	URI    string             // Required: Request URI.
	Params []string           // Optional: Validate request parameters.
	Auth   bool               // Optional: Choose to authenticate client or not using a predefined function.
	Before []http.HandlerFunc // Optional: Before middlerwares.
	After  []http.HandlerFunc // Optional: After middlerwares.
	Do     func(w http.ResponseWriter, r *http.Request)  // Required: Standard go http.handlerFunc
}
```

Routing
---
Handler works with the standard go http.ServeMux, but allows you to provide different handlers for the same uri using different methods.  


Validation
---
Validation works with the provided **Params** string array. Which should look like this:  

```go
[]string{"name:rule=value,rule=value", "name:rule=value", ...}
// ex:
[]string{"body:max=3000,min=140", "id:numeric"} // some rules like "numeric" don't have a value
```


**Notes**

- if you didn't define the Params field, or defined empty array, Handlers won't do r.ParseForm()
- Adding a parameter without any rules validates **presence**.


**Available rules**:  

- **numeric**: ensure the parameter is a numeric value, or a string represent numeric value.
- **max**: validates the maximum length of the parameter value.
- **min**: validates the minimum length of the parameter value.
- **empty**: allow the parameter to be empty, but if it has value, it should follow the other rules.

Authentication
--- 
You should write the authentication logic in the auth function to do authentication, then you choose to authenticate the client for a certain handler by defining the **Auth** filed to be true.  


Middlewares
---
Use **Before** and **After** fields to add middlewares to the handler.


Todo
---
- Tests
- Add more validation rules
- Add validation errors to readme
- Generate API file that is compatible with popular services like apiary
- Check if generating validation code could make things faster than parsing the params strings

License
---
Do whatever you like with the code.