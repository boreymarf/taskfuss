from pydantic import BaseModel, ConfigDict


class PlanMetadata(BaseModel):
    """Metadata about a quest plan, exposed to the frontend."""
    id: str
    name: str
    description: str

    model_config = ConfigDict(from_attributes=True)
