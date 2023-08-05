# ServingTumbnailsRedis
Explore redis and serve the thumbnail. Basically we need to check how the code should be serving a folder which contains thumbnail of the camera feed The change should take place from the backend, and then an event can be posted/pushed over redis with update filename The UI should listen and pickup the filename and refresh the image inside itself
