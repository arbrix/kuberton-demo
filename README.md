# kuberton-demo
Demo monolith system as initial task

## How to use

### Build

```
docker build -t <your-image-name> .
```

### Run

```
docker run --rm -it -p 3000:3000 <your-image-name>
```

Open: [shop](localhost:3000)

### Configuration

Key | Default Value | Description
PORT | `3000` |
BANNER_COLOR | "css property" |

## API

### Routes

`GET` | `/` | home page (product list, link to the cart)
`GET`| `/product/{id}` | product page, select quantity, add to the cart
`POST` | `/setCurrency` | change user currency preference
`GET` | `/static/*` | static files server
`GET` | `/_healthz` | container health check
