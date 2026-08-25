from contextlib import asynccontextmanager
import json
import logging
import os
from pathlib import Path
import traceback

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from pydantic import ValidationError
from fastapi.encoders import jsonable_encoder
import uvicorn

from src.api.routers import config, debug, plan, quest, quest_state, record, user
from src.config import get_config, init_config
from src.container import Container
from src.database import Database
from src.exception_handlers import app_exception_handler, generic_exception_handler, validation_exception_handler
from src.exceptions import AppException
from src.logging.setup import setup_logging

logger = logging.getLogger(__name__)


def create_app(config_path: str | None = None) -> FastAPI:
    """Application factory that sets up config, logging, database, and routes."""
    # Initialize configuration
    if config_path is None:
        config_path = os.environ.get("APP_CONFIG_PATH") or "config.toml"
    init_config(Path(config_path))

    # Setup logging
    setup_logging(get_config().logging)

    # Create database and store it on app.state
    database = Database.from_config()

    # Create and wire dependency injection container
    container = Container()
    container.wire()

    @asynccontextmanager
    async def lifespan(_app: FastAPI):
        # Startup
        if get_config().app.environment == "dev":
            generate_openapi_file(_app)

        # Any initialization that requires database can be done here
        # plan_implementations_path = Path(os.getcwd()) / "src" / "plans" / "implementations"
        # plan_service = container.plan_service()
        # with database.session() as db:
        #     plan_service.sync_plans_from_directory(db, plan_implementations_path)

        yield
        # Shutdown
        database.dispose()
        logger.info("Application stopped")

    app = FastAPI(
        lifespan=lifespan,
        exception_handlers={
            AppException: app_exception_handler,
            ValidationError: validation_exception_handler,
            Exception: generic_exception_handler,
        },
    )
    app.state.database = database

    # Include routers
    app.include_router(debug.router)
    app.include_router(user.router)
    app.include_router(config.router)
    app.include_router(quest.router)
    app.include_router(quest_state.router)
    app.include_router(plan.router)
    app.include_router(record.router)

    # Add CORS middleware
    app.add_middleware(
        CORSMiddleware,
        allow_origins=get_config().server.allow_origins,
        allow_credentials=get_config().server.allow_credentials,
        allow_methods=get_config().server.allow_methods,
        allow_headers=get_config().server.allow_headers,
    )

    return app


def run_app() -> None:
    """Run the application via uvicorn with import string for reload support."""
    if get_config().app.environment == "dev":
        reload = get_config().app.dev.reload_app
        reload_dirs = get_config().app.dev.reload_dirs
    else:
        reload = False
        reload_dirs = None

    uvicorn.run(
        "src.app:app",
        log_config=None,
        port=get_config().server.port,
        log_level=get_config().logging.root_level,
        reload=reload,
        reload_dirs=reload_dirs,
    )


def generate_openapi_file(app: FastAPI) -> None:
    """Write OpenAPI schema to shared temp directory in dev mode."""
    openapi_schema = app.openapi()
    root_dir = Path(os.getcwd())
    shared_tmp_dir = root_dir.parent / "shared" / "temp"
    shared_tmp_dir.mkdir(parents=True, exist_ok=True)
    openapi_path = shared_tmp_dir / "openapi.json"
    with open(openapi_path, "w", encoding="utf-8") as f:
        json.dump(openapi_schema, f)


# Global application instance (import string "src.app:app" points here)
app = create_app()
