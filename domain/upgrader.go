package domain

// Upgrader runs NetworkLoader periodically and store last updates
// It allow faster show any results via HTTP API
type Upgrader interface {

	// Run loads download/upload speed summary asynchronously
	Run()

	// GetDLSpeedMbps returns last download speed value
	GetDLSpeedMbps() float64

	// GetULSpeedMbps returns last upload speed value
	GetULSpeedMbps() float64

	// IsReady returns true if download and upload speed was set at least once
	IsReady() bool
}
