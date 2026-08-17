# backend

## 0.2.0

### Minor Changes

- feat: added API endpoint for getting list of quest plan registries
- feat: added conftest fixtures for managing sessions
- feat: added tuple field
- feat: created abstract base class for quest plans and structs related to it
- feat: added API endpoint that gets setup form of specific plan
- feat: added quest_state router and functionality
- feat: added a "details" property to exception that allows to send extra data besides a message
- feat: created a contextmanager for sessions in database.py
- feat: added record creation
- feat: added quest_plan table to the database
- chore: renamed "detail" to "message" in exceptions

### Patch Changes

- feat: changed field error type from dict[str, FieldError] to list[FieldError] with location inside errors
- feat: reworked defaults for fields
- fix: fixed infinite auth tokens
- chore: moved all quest files with structs to domain dir
- chore: separated exceptions to different files
- chore: moved get_session() fastapi dependency to API modules
- chore: changed all UUIDs in the db to ints

## 0.1.1

### Patch Changes

- chore: backend now spawns openapi.json when dev environment is active
- fix: fixed openapi config not creating when reloading dev server

## 0.1.0

### Minor Changes

- feat: imported API and user registration stuff from other project
