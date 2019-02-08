# Pinterest clone

## Tools used

- Docker

## Instructions

- Run `docker-compose up -d`
- Run `psql pc postgres` from inside the db container
- Run `select id from users limit 1` and save the UUID in a text file
- Get the Github CLIENT-ID for the app you have registered in Github
- Point your browser to `https://github.com/login/oauth/authorize?scope=user:email&client_id=YOURGHCLIENTID&redirect_uri=http://localhost:8080/success_GH_authn_callback/YOUR-UUID-HERE`
- Your browser should be redirected around and then should get a UUID access token from this program.  Save this access-token.  You will use it to make authenticated calls to this program.

## Endpoints

### Add a new image as an authenticated user

```
curl -H "Authorization: YOUR-ACCESS-TOKEN" -XPOST -d '{"href": "http://foo.com"}' localhost:8080/images
```

### List your images

```
curl -H "Authorization: YOUR-ACCESS-TOKEN" -XGET localhost:8080/images
```

### Delete a particular image

```
curl -H "Authorization: YOUR-ACCESS-TOKEN" -XDELETE localhost:8080/images/UUID-OF-THE-IMAGE
```

### List images of a particular user

```
curl -XGET localhost:8080/images/THEIR-UUID
```
