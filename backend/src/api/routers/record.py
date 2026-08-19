from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.orm import Session

from src.api.dependencies.session import get_session
from src.api.dependencies.auth import get_current_user
from src.domain.record import Record, RecordCreateRequest, RecordQueryParams
from src.exceptions import NotFoundError
from src.exceptions.generic import BadRequestError
from src.repositories.quest_state import QuestStateRepository
from src.service.record import RecordService

router = APIRouter(prefix="/api/record", tags=["record"])


@router.get("/", response_model=list[Record])
def get_all_records(
    params: RecordQueryParams = Depends(),
    state_id: int | None = Query(None, description="Filter records by state (uses its time range)"),
    db: Session = Depends(get_session),
    # current_user: int = Depends(get_current_user),
):
    if state_id is not None:
        if params.created_at__gte is not None or params.created_at__lte is not None:
            raise BadRequestError("Cannot use 'state_id' together with 'created_at__gte' or 'created_at__lte'. Please choose one filtering method.")

        state = QuestStateRepository.get_by_id(db, state_id)
        if state is None:
            raise NotFoundError("state", state_id)

        params.created_at__gte = state.start_date
        params.created_at__lte = state.end_date

    return RecordService.get_all(db, params)


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
