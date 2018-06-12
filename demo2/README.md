# Demo 2

Function to show the Lambda lifecyle, which is something that needs to be considered when building for serverless.

One would expect the counter to increment by 10 every second, but this is not the case as Lambda goes into a "frozen" state. Invoke the function and then wait a few seconds before invoking again to see this behaviour.
