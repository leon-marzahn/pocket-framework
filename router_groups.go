package pocketframework

import (
  "github.com/pocketbase/pocketbase/core"
  "github.com/pocketbase/pocketbase/tools/router"
)

type RouterGroups struct {
  Public        *router.RouterGroup[*core.RequestEvent]
  Authenticated *router.RouterGroup[*core.RequestEvent]
  Admin         *router.RouterGroup[*core.RequestEvent]
}

func (r RouterGroups) WithPrefix(prefix string) RouterGroups {
  return RouterGroups{
    Public:        r.Public.Group(prefix),
    Authenticated: r.Authenticated.Group(prefix),
    Admin:         r.Admin.Group(prefix),
  }
}
