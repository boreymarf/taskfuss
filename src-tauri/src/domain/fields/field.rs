use serde::{Deserialize, Serialize};
use serde_json::Value as JsonValue;
use std::collections::HashMap;

use crate::{domain::{FieldSchema, FieldValue, IntField, StrField}};
use crate::errors::{FieldSchemaError, FieldValueError};

/// Discriminated union of every supported field type.
///
/// This is the single entry point for working with fields of unknown kind:
/// it carries the per-type schema and implements [`FieldSchema`] so callers 
/// can validate without matching on variants.
#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(tag = "type")]
pub enum Field {
    #[serde(rename = "str")]
    Str(StrField),
    #[serde(rename = "int")]
    Int(IntField),
}

impl FieldSchema for Field {
    type Validated = FieldValue;

    fn validate_schema(&self) -> Result<(), FieldSchemaError> {
        match self {
            Field::Str(f) => f.validate_schema(),
            Field::Int(f) => f.validate_schema(),
        }
    }

    fn validate_value(&self, value: &serde_json::Value) -> Result<Self::Validated, FieldValueError> {
        match self {
            Field::Str(f) => f.validate_value(value).map(FieldValue::Str),
            Field::Int(f) => f.validate_value(value).map(FieldValue::Int),
        }
    }
}