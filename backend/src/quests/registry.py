import importlib

from pydantic import BaseModel, ConfigDict

from src.quests.base import BaseQuestPlan


class PlanRegistry(BaseModel):
    id: str
    name: str
    class_path: str

    def import_class(self) -> BaseQuestPlan:
        """Imports class from string like 'module.path:ClassName'"""
        try:
            module_path, class_name = self.class_path.split(":")
            module = importlib.import_module(module_path)
            return getattr(module, class_name)
        except (ImportError, AttributeError, ValueError) as e:
            raise ImportError(f"Failed to import {self.class_path}: {e}")

    model_config = ConfigDict(from_attributes=True)
