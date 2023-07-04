package http

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	exceptions "github.com/mousav1/weiser/app/exceptions"
	middleware "github.com/mousav1/weiser/app/http/middlewares"
)

// Define middleware aliases
var MiddlewareAliases = map[string]func(*fiber.Ctx) error{
	"logger": middleware.LoggerMiddleware,
}

// Define main middleware functions
var Middleware = []func(*fiber.Ctx) error{
	middleware.LoggerMiddleware,
	exceptions.ErrorHandler,
}

// Define middleware struct
type MiddlewareStruct struct {
	Handler func(*fiber.Ctx) error
	Order   int
}

// Define middleware group struct
type MiddlewareGroup struct {
	Condition   func(*fiber.Ctx) bool
	Middlewares []MiddlewareStruct
}

// Define custom middleware functions
func CustomMiddleware1(param1 string, param2 int) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// ...
		return c.Next()
	}
}

// Define middleware groups
var MiddlewareGroups = map[string]MiddlewareGroup{
	"mygroup1": {
		Condition: func(c *fiber.Ctx) bool {
			// ...
			return true
		},
		Middlewares: []MiddlewareStruct{
			{
				Handler: middleware.LoggerMiddleware,
				Order:   1,
			},
		},
	},
	// ...
}

// Apply a group of middlewares to a controller
func ApplyMiddlewareGroup(groupName string, controllerName interface{}, cxt *fiber.Ctx) {
	// Get the controller by name
	controller := reflect.ValueOf(controllerName).Elem()

	// Get the middleware group
	middlewareGroup, ok := MiddlewareGroups[groupName]
	if !ok {
		return
	}

	// Apply the middlewares to the controller
	for _, middleware := range middlewareGroup.Middlewares {
		if middlewareGroup.Condition != nil && !middlewareGroup.Condition(cxt) {
			continue
		}
		controller.Set(reflect.Append(controller, reflect.ValueOf(middleware.Handler)))
	}
}

// Apply a middleware to a controller
func ApplyMiddleware(handler func(*fiber.Ctx) error, order int, controllerName interface{}) {
	// Get the controller by name
	controller := reflect.ValueOf(controllerName).Elem()

	// Create the middleware struct
	middleware := MiddlewareStruct{
		Handler: handler,
		Order:   order,
	}

	// Get the number of existing middlewares in the controller
	numMethods := controller.NumMethod()

	// Find a place to insert the middleware
	var index int
	for i := 0; i < numMethods; i++ {
		m := controller.Method(i)
		if m.Type().String() == "func(*fiber.Ctx) error" {
			middlewareOrder := 0
			if orderMethod := controller.MethodByName("Order"); !orderMethod.IsNil() {
				orderValue := orderMethod.Call(nil)[0]
				if orderInt, ok := orderValue.Interface().(int); ok {
					middlewareOrder = orderInt
				}
			}
			if middleware.Order < middlewareOrder {
				index = i
				break
			} else {
				index = i + 1
			}
		}
	}

	// Insert the middleware into the controller
	newSlice := reflect.MakeSlice(controller.Type(), numMethods+1, numMethods+1)
	reflect.Copy(newSlice.Slice(0, index), controller.Slice(0, index))
	newSlice.Index(index).Set(reflect.ValueOf(middleware.Handler))
	reflect.Copy(newSlice.Slice(index+1, numMethods+1), controller.Slice(index, numMethods))
	controller.Set(newSlice)
}
