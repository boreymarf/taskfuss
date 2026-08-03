from sqlalchemy import ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from src.db import Base


class QuestDB(Base):
    __tablename__ = "quest"

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    owner_id: Mapped[int] = mapped_column(ForeignKey("user.id", ondelete="CASCADE"))
    plan_id: Mapped[str] = mapped_column(
        ForeignKey("plan.id", ondelete="CASCADE")
    )
