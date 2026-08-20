from datetime import datetime

from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.orm import Session

from src.api.dependencies.session import get_session
from src.api.dependencies.auth import get_current_user
from src.domain.record import Record, RecordCreateRequest
from src.service.record import RecordService

router = APIRouter(prefix="/api/record", tags=["record"])


@router.get("/", response_model=list[Record])
def get_all_records(
    quest_id: int | None = Query(None),
    field_path: str | None = Query(None),
    field_path__startswith: str | None = Query(None),
    automated: bool | None = Query(None),
    created_at__gte: datetime | None = Query(None),
    created_at__lte: datetime | None = Query(None),
    latest: bool = Query(False),
    ordering: str | None = Query(None),
    limit: int | None = Query(None, ge=1),
    offset: int | None = Query(None, ge=0),
    state_id: int | None = Query(None, description="Filter records by state (uses its time range)"),
    db: Session = Depends(get_session),
    # current_user: int = Depends(get_current_user),
):
    return RecordService.get_all(
        db,
        quest_id=quest_id,
        field_path=field_path,
        field_path__startswith=field_path__startswith,
        automated=automated,
        created_at__gte=created_at__gte,
        created_at__lte=created_at__lte,
        latest=latest,
        ordering=ordering,
        limit=limit,
        offset=offset,
        state_id=state_id,
    )


@router.get("/{record_id}", response_model=Record)
def get_record(
    record_id: int,
    db: Session = Depends(get_session),
    # current_user: int=Depends(get_current_user),
):
    record = RecordService.get(db, record_id)
    if record is None:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Record not found"
        )
    return record


@router.post("/", response_model=Record, status_code=status.HTTP_201_CREATED)
def create_record(
    request: RecordCreateRequest,
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return RecordService.create(db, request, current_user)


@router.delete("/{record_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_record(
    record_id: int,
    db: Session = Depends(get_session),
):
    RecordService.remove(db, record_id)
    return
