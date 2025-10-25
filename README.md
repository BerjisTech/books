Berjis Books

Books is a platform where authors and publishers can work together to distribute their work primarily online.

Highlights

- Authors collaborate with publishers; publishers can scout authors.
- Readers buy ebooks, optionally share when allowed; hard copies supported.
- Book clubs with rooms and discussions.
- Public meetups (free or paid) with registration.

Structure

- Frontend (Angular): `books/frontend`
- Service (Go + Postgres): `books/service`

Dev quick start

- Compose services are defined in root `docker-compose.yml` as `books-db`, `books-service`, `books-frontend`.
- Edge proxy hosts: `books.berjis.tech` and `books-api.berjis.tech`.

