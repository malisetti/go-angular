# go-angular
angular js web app(backed by golang) to show and rate cute puppies

create http server with go
expose apis
  fetch puppies
  like a puppy
  dislike a puppy :(
  
fetch cute puppies(all are cute) images from flickr https://www.flickr.com/services/api/flickr.photos.search.html

add models to represent the photos from flickr
  initialize votes on a new photo to 0
add orm to deal with databases(sqlite/mysql)

web app
  use bootstrap
  use angular js
  call fetch puppies api
  display fetched puppies in a grid with pagination
  give options to like or dislike
  options to sort based on most voted/by time
