package fcore

import "sync"

// Put fields related to mounting.
type mountFields struct {
	// Mounted and main apps
	appList map[string]*App
	// Prefix of app if it was mounted
	mountPath string
	// Ordered keys of apps (sorted by key length for Render)
	appListKeys []string
	// check added routes of sub-apps
	subAppsRoutesAdded sync.Once
	// check mounted sub-apps
	subAppsProcessed sync.Once
}
