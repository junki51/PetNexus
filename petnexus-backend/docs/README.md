# PetNexus Backend Documentation

This directory is the documentation hub for the Go/Gin backend. Sprint history,
older planning notes, and frontend notes remain available, but the documents
below are the current backend-oriented entry points.

## Backend documentation

| Document | Purpose |
| --- | --- |
| [Backend overview](./backend/backend-overview.md) | High-level product and completed backend scope |
| [Architecture](./backend/architecture.md) | Layered architecture, dependency direction, and conventions |
| [Module map](./backend/module-map.md) | Current modules, responsibilities, tables, endpoints, and access rules |
| [API reference](./backend/api-reference.md) | Current Sprint 1-9 HTTP API contracts |
| [Database schema](./backend/database-schema.md) | Implemented PostgreSQL tables, relationships, constraints, and indexes |
| [Auth and permissions](./backend/auth-and-permissions.md) | JWT, roles, public/protected routes, and ownership rules |
| [Testing guide](./backend/testing-guide.md) | Formatting, tests, local API checks, and Render smoke tests |
| [Deployment guide](./backend/deployment-guide.md) | Local Docker/PostgreSQL and Render deployment configuration |
| [Backend roadmap](./backend/roadmap.md) | Completed scope, excluded features, and next planning direction |

## Sprint summary

- [Sprint 1-9 summary](./sprints/sprint-1-to-9-summary.md)
- [Sprint 1-8 summary](./sprints/sprint-1-to-8-summary.md)
- Detailed historical logs remain in [progress](./progress/README.md).

## Existing reference material

The following files are retained because they contain useful design history or
implementation context:

- [Backend Codex rules](./backend-codex-rules.md)
- [Backend setup checklist](./backend-setup-checklist.md)
- [Original backend roadmap](./backend-roadmap.md)
- [API plan](./api-plan.md)
- [Database plan](./database-plan.md)
- [Frontend integration guide through Sprint 5](./frontend-integration-guide-sprint-1-to-5.md)

When current code and an older planning document differ, use the backend docs
in `docs/backend/` and the implementation as the source of truth.
