from sqlalchemy import String
from sqlalchemy.orm import Mapped, mapped_column
from src.db.base import Base


class QuestPlanDB(Base):
    __tablename__ = "quest_plan"

    id: Mapped[str] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(120))
    class_path: Mapped[str] = mapped_column(String(120))
