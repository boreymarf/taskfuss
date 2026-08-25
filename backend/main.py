import logging
import os
from pathlib import Path
import click
import sys

from src.config import get_config, init_config
from src.logging.setup import setup_logging

logger = logging.getLogger(__name__)


@click.group()
def cli():
    """CLI tool for managing configuration."""
    pass


@cli.command()
@click.option("--config", "-c", default=None, help="Custom config path")
def start_server(config: str | None):
    """Start server."""
    config_path = Path(config) if config else Path("config.toml")

    os.environ["APP_CONFIG_PATH"] = str(config_path.resolve())

    # config
    config_path = os.environ.get("APP_CONFIG_PATH") or "config.toml"
    init_config(Path(config_path))

    # logging
    setup_logging(get_config().logging)

    # So src.app has a chance of reading correct configuration file
    from src.app import run_app

    run_app()


@cli.command()
@click.option("--path", "-p", default=None, help="Custom file path to write to")
def generate_openapi_file(_path: str | None):
    pass


@cli.command()
@click.option(
    "--force", "-f", is_flag=True, help="Force create new config even if exists"
)
@click.option("--path", "-p", default=None, help="Custom config path")
def generate_config(path: str | None, force: bool = False):
    """Initialize or update configuration file. Will update config by default."""

    config_path = Path(path) if path else Path("config.toml")

    if force and config_path.exists():
        click.echo("Force deleted old configuration.")
        config_path.unlink()

    init_config(config_path=config_path)
    click.echo(f"Config ready at {config_path}")


if __name__ == "__main__":
    if len(sys.argv) == 1:
        start_server()
    else:
        cli()
