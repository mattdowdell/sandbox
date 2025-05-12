// Package drivers provides code that is neither application or business specific.
//
// Drivers are a set of generic code packages. In most cases, it can be used in other applications
// without further modification.
//
// As an example, a driver could be the UI rendering logic for a desktop application. An adapter
// would specify that a dialog should have a specific size, title and content. But the rendering
// logic for the dialog on the screen and listening for input is part of the drivers layer.
//
// In a system that promotes code re-use and some amount of uniformity across microservices, most
// drivers would exist in a separate repository and the drivers layer would be minimal or even
// empty. However, for this application which aims to be relatively self-contained, the drivers
// layer is relatively large.
package drivers
