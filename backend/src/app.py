from contextlib import asynccontextmanager
import json
import logging
import os
from pathlib import Path
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
import uvicorn

from src.api.routers import auth, config, debug, quest, plan, quest_state, record, user
from src.logger_conf.helpers import set_modules_log_level
from src.logger_conf.implementations import (
    set_root_logger,
    setup_daily_json_logger,
    setup_json_logger,
    setup_rich_logging,
)
from src.config import get_config, init_config
from src.database import create_engine, get_session
from src.exceptions import AppException
from src.service.plan import PlanService

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(_app: FastAPI):

    # The code is duplicated from main.py for a reason!
    # Otherwise first logs from uvicorn will be ignored.

    # config
    config_path = os.environ.get("APP_CONFIG_PATH") or "config.toml"
    init_config(Path(config_path))

    # logging
    set_root_logger(level=get_config().logging.log_level)
    setup_rich_logging(level=logging.DEBUG)
    setup_daily_json_logger(log_dir="logs")
    setup_json_logger(log_file="logs/errors.jsonl", level=logging.ERROR)

    set_modules_log_level(get_config().logging.shushed_modules, "WARNING")

    if get_config().app.environment == "dev":
        generate_openapi_file(app)

    plan_implementations_path = Path(os.getcwd()) / "src" / "plans" / "implementations"
    with get_session() as db:
        PlanService.sync_plans_from_directory(db, plan_implementations_path)

    yield
    # Clean up


app = FastAPI(lifespan=lifespan)

# Routes
app.include_router(debug.router)
app.include_router(user.router)
app.include_router(config.router)
app.include_router(quest.router)
app.include_router(quest_state.router)
app.include_router(auth.router)
app.include_router(plan.router)
app.include_router(record.router)

# Middlewares
app.add_middleware(
    CORSMiddleware,
    allow_origins=get_config().server.allow_origins,
    allow_credentials=get_config().server.allow_credentials,
    allow_methods=get_config().server.allow_methods,
    allow_headers=get_config().server.allow_headers,
)


def run_app():

    if get_config().app.environment == "dev":
        reload = get_config().app.dev.reload_app
        reload_dirs = get_config().app.dev.reload_dirs
    else:
        reload = False
        reload_dirs = None

    create_engine()

    uvicorn.run(
        "src.app:app",
        log_config=None,
        port=get_config().server.port,
        log_level=get_config().logging.log_level,
        reload=reload,
        reload_dirs=reload_dirs,
    )


def generate_openapi_file(app: FastAPI):
    openapi_schema = app.openapi()
    root_dir = Path(os.getcwd())
    shared_tmp_dir = root_dir.parent / "shared" / "temp"
    Path(shared_tmp_dir).mkdir(parents=True, exist_ok=True)
    openapi_path = shared_tmp_dir / "openapi.json"
    with open(openapi_path, "w", encoding="utf-8") as f:
        json.dump(openapi_schema, f)


# Exception handlers (very important yes yes)
@app.exception_handler(AppException)
async def app_exception_handler(_request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"message": exc.message, "details": exc.details},
    )


@app.exception_handler(Exception)
async def generic_exception_handler(_request: Request, exc: Exception):
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(status_code=500, content={"detail": "Internal server error"})
