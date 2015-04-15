# go-angular
1. [USAGE](https://github.com/mseshachalam/go-angular/blob/master/USAGE.md) 
2. [TODO](https://github.com/mseshachalam/go-angular/blob/master/TODO.md)

angular js web app(backed by golang) to show and rate cute puppies

1. create http server with golang
  * expose apis
    * fetch puppies
    * like a puppy
    * dislike a puppy :(
  
2. fetch cute puppies(all are cute) images from [Flickr Photo Search API](https://flickr https://www.flickr.com/services/api/flickr.photos.search.html "Flickr Photo Search")

  * add models to represent the photos from flickr
    * initialize votes on a new photo to 0
  * add orm to deal with databases(sqlite/mysql)

3. web app
  *  use bootstrap
  *  use angular js
  *  call fetch puppies api
  *  display fetched puppies in a grid with pagination
  *  give options to like or dislike
  *  options to sort based on most voted/by time
