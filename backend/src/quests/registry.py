import importlib
from pathlib import Path
import importlib.util

from pydantic import BaseModel, ConfigDict

from src.quests.base import BaseQuestPlan


class PlanRegistry(BaseModel):
    id: str
    name: str
    class_path: str

    def import_class(self) -> type[BaseQuestPlan]:
        """
        Imports a class from an absolute file path.
        
        Expected format: '/path/to/file.py:ClassName'
        """
        try:
            file_path_str, class_name = self.class_path.rsplit(":", 1)
            file_path = Path(file_path_str)

            if not file_path.exists():
                raise ImportError(f"File not found: {file_path}")

            module_name = f"dynamic_module_{abs(hash(file_path))}"
            spec = importlib.util.spec_from_file_location(module_name, file_path)

            if spec is None or spec.loader is None:
                raise ImportError(f"Could not load spec for {file_path}")

            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)

            return getattr(module, class_name)

        except (ValueError, ImportError, AttributeError, OSError) as e:
            raise ImportError(f"Failed to import {self.class_path}: {e}") from e

    def to_public(self) -> PlanRegistryPublic:
        return PlanRegistryPublic(**self.model_dump())

    model_config = ConfigDict(from_attributes=True)


class PlanRegistryPublic(BaseModel):
    id: str
    name: str

    model_config = ConfigDict(from_attributes=True)
