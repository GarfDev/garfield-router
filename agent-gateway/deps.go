package main

import (
	"github.com/GarfDev/garfield-router/agent-gateway/accounts"
	"github.com/GarfDev/garfield-router/agent-gateway/registry"
)

// Type aliases keep the rest of the main package readable while still
// referencing the types defined in the sub-packages.
type registryDeps = registry.Registry
type accountsDeps = accounts.Manager
