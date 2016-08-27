package server

import "github.com/zhengyun1112/glrest/logger"

type Controller interface {
	RegisterRoutes()
	SetServer(s *Server)
}

type BaseController struct {
	s *Server
}

// IMPORTANT: you need overwrite RegisterRoutes method to register url routes, such as
//      userController.Get("/user/:id", userController.getUserById)
func (c *BaseController) RegisterRoutes() {
	logger.Panic("Not implemented")
}

func (c *BaseController) SetServer(s *Server) {
	c.s = s
}

func (c *BaseController) HandleRaw(method, path string, handle RawHandler) {
	c.s.HandleRaw(method, path, handle)
}

// Get will register a 'GET' request handler to the router.
func (c *BaseController) Get(path string, handle Handler) {
	c.s.Get(path, handle)
}

// Post will register a 'POST' request handler to the router.
func (c *BaseController) Post(path string, handle Handler) {
	c.s.Post(path, handle)
}

// Put will register a 'PUT' request handler to the router.
func (c *BaseController) Put(path string, handle Handler) {
	c.s.Put(path, handle)
}

// Patch will register a 'PATCH' request handler to the router.
func (c *BaseController) Patch(path string, handle Handler) {
	c.s.Patch(path, handle)
}

// Head will register a 'HEAD' request handler to the router.
func (c *BaseController) Head(path string, handle Handler) {
	c.s.Head(path, handle)
}

// Delete will register a 'DELETE' request handler to the router.
func (c *BaseController) Delete(path string, handle Handler) {
	c.s.Delete(path, handle)
}
