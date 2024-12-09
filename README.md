# Authentication Service

This is a basic authentication service with register and login functionalities.

It is not intended to being used in production or any other data sensitive environment!

## Using the service

(for this to work you need a Postgres DB installed as well as an instance of Redis)

Make sure all environment variables are set correctly and point towards your local database services.

Then run the following command:

```sh
pnpm run dev
```

You can also use Docker Compose to run the entire service. Both Postgres and Redis are running inside Docker Containers.

```sh
docker-compose up --build
```

To stop the service, run

```sh
docker compose down
```

## What's inside?

This Turborepo includes the following packages/apps:

### Apps and Packages

- `backend`: a go backend
- `web`: a [Next.js](https://nextjs.org/) app with [Tailwind CSS](https://tailwindcss.com/)
- `ui`: a stub React component library with [Tailwind CSS](https://tailwindcss.com/)
- `@@authentication-service-go/eslint-config`: `eslint` configurations (includes `eslint-config-next` and `eslint-config-prettier`)
- `@@authentication-service-go/tailwind-config`: `tailwind.config.ts`s used throughout the monorepo
- `@@authentication-service-go/typescript-config`: `tsconfig.json`s used throughout the monorepo

### Building packages/ui
