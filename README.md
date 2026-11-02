# vanity

Vanity import path server for `m4o.io`, hosted on Cloud Run.

## How it works

The Go toolchain resolves a custom import path with an HTTP request. For
`go get m4o.io/pbf`, it requests `https://m4o.io/pbf?go-get=1` and reads the
`go-import` meta tag in the response. The tag names the module root, the
version control system, and the repository. This server returns that tag.

[`main.go`](main.go) is a small wrapper. It builds a TOML backend and a
vanity HTTP handler from the [`l7e.io/vanity`](https://l7e.io/vanity)
library, then serves them on the port in `PORT`. The library holds all the
request handling.

[`config.toml`](config.toml) lists the import paths. Each entry gives the
import path, the version control system, and the account URL that holds the
repository. The handler trims a request path to its first element before the
lookup, so a major version subpath such as `m4o.io/pbf/v2` resolves from the
`m4o.io/pbf` entry.

A request without `go-get=1` comes from a browser. The handler redirects it
to the package on `pkg.go.dev`.

The [`Dockerfile`](Dockerfile) builds the server with a multi-stage build and
copies `config.toml` into the image at `/var/vanity/config.toml`. The
environment variables `PORT`, `TOML_FILE`, and `TOML_TABLE` set the listen
port and the location of the import path list.

The domain `m4o.io` reaches the service through a Cloud Run domain mapping.
The apex A records point at Google, and Cloud Run terminates TLS and routes
the request to the service.

## Deployment

The server runs as the Cloud Run service `m4o-vanity` in the
`github-maguro-infrastructure` project, region `us-central1`.

Image: `gcr.io/github-maguro-infrastructure/vanity`

Build the image. Use `--platform linux/amd64` on an Apple silicon machine,
because Cloud Run does not start an arm64 image.

```bash
docker build --platform linux/amd64 -t gcr.io/github-maguro-infrastructure/vanity:latest .
```

Push the image.

```bash
docker push gcr.io/github-maguro-infrastructure/vanity:latest
```

Deploy the new image to the service.

```bash
gcloud run deploy m4o-vanity \
  --image gcr.io/github-maguro-infrastructure/vanity:latest \
  --region us-central1 \
  --project github-maguro-infrastructure
```

The image contains `config.toml`, so a change to the import paths needs a new
build and a new deployment.

Verify a deployment with a direct request.

```bash
curl "https://m4o.io/pbf?go-get=1"
```
