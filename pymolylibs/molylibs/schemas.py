from datetime import datetime
from typing import List, Optional, Any
from pydantic import BaseModel, Field
from bson import ObjectId 

class MongoObjectIdPydanticConfig():
    json_encoders = {
        ObjectId: str
    }
    arbitrary_types_allowed = True

class TreeAttributes(BaseModel):
    org_id: str = Field(alias='organization_id')
    slug: str
    label: str
    desc: str = Field(alias='description')
    created_at: datetime

class Tree(BaseModel):
    id: int
    root_id: int
    parent_id: Optional[int]
    order: int
    attributes: TreeAttributes

class HttpSuccessResponse(BaseModel):
    result: str = Field(default='success')
    data: Any

class HttpSuccessResponseWithTotalCount(HttpSuccessResponse):
    total: int

class HttpErrorResponse(BaseModel):
    result: str = Field(default='success')
    data: Optional[Any]
    error: str
    error_prompt: str