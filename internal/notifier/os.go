package notifier

import "github.com/gen2brain/beeep"

// OSNotifier sends notifications using the OS's native notification system.
type OSNotifier struct{}

// Notify sends a desktop notification.
func (n *OSNotifier) Notify(title, message string) error {
	// The beeep library uses a default app icon and name.
	// For more advanced usage, you could specify an icon path.
	return beeep.Notify(title, message, "")
}
