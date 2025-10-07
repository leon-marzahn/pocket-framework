package pocketframework

import (
  "github.com/pocketbase/pocketbase/apis"
  "github.com/pocketbase/pocketbase/core"
)

type ModuleRegistry struct {
  modules   []Module
  app       core.App
  apiPrefix string
}

func NewModuleRegistry(app core.App, apiPrefix string) *ModuleRegistry {
  return &ModuleRegistry{
    modules:   []Module{},
    app:       app,
    apiPrefix: apiPrefix,
  }
}

// Register registers a module in the module registry.
func (m *ModuleRegistry) Register(module Module) {
  m.modules = append(m.modules, module)
}

// Init adds the module registry to the pocketbase app.
func (m *ModuleRegistry) Init() error {
  for _, module := range m.modules {
    if err := registerModuleHooks(module, m.app); err != nil {
      return err
    }
  }

  m.app.OnServe().BindFunc(
    func(se *core.ServeEvent) error {
      baseGroup := se.Router.Group(m.apiPrefix)

      authenticatedGroup := se.Router.Group(m.apiPrefix)
      authenticatedGroup.Bind(apis.RequireAuth())

      adminGroup := se.Router.Group(m.apiPrefix)
      adminGroup.Bind(apis.RequireSuperuserAuth())

      baseGroups := RouterGroups{
        Public:        baseGroup,
        Authenticated: authenticatedGroup,
        Admin:         adminGroup,
      }

      for _, module := range m.modules {
        if err := serveModule(module, baseGroups); err != nil {
          return err
        }
      }

      return se.Next()
    },
  )

  return nil
}

func registerModuleHooks(module Module, app core.App) error {
  if err := module.RegisterHooks(app); err != nil {
    return err
  }

  if moduleWithChildren, ok := module.(ModuleWithChildren); ok {
    for _, childModule := range moduleWithChildren.Children() {
      if err := registerModuleHooks(childModule, app); err != nil {
        return err
      }
    }
  }

  return nil
}

func serveModule(module Module, baseGroups RouterGroups) error {
  groups := baseGroups.WithPrefix(module.Prefix())
  if err := module.RegisterRoutes(groups); err != nil {
    return err
  }

  if moduleWithChildren, ok := module.(ModuleWithChildren); ok {
    for _, childModule := range moduleWithChildren.Children() {
      if err := serveModule(childModule, groups); err != nil {
        return err
      }
    }
  }

  return nil
}
