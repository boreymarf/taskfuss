
from pydantic import BaseModel, ConfigDict


class QuestData(BaseModel):
    """Info that instance returns about the quest"""


    model_config = ConfigDict(from_attributes=True)
