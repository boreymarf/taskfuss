from fastapi import APIRouter

from src.api.openapi_responses import INTERNAL_ERROR


router = APIRouter(prefix="/api/debug", tags=["debug"])


@router.get("/pong/")
async def ping():
    return {"detail": "Pong!"}


@router.get(
    "/crash",
    responses={**INTERNAL_ERROR},
)
async def crash():
    raise Exception("Something bad happened!")
