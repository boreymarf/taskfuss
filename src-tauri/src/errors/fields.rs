use std::collections::HashMap;

use serde::Serialize;

#[derive(specta::Type, Debug, thiserror::Error, Serialize)]
pub enum FieldSchemaError {
    #[error("min ({min}) is greater than max ({max})")]
    MinGreaterThanMax { min: i64, max: i64 },
}

#[derive(specta::Type, Debug, thiserror::Error, Serialize)]
pub enum FieldValueError {
    #[error("field is required")]
    Required,
    #[error("wrong type: expected {expected}")]
    WrongType { expected: &'static str },
    #[error("value too short (min {min})")]
    TooShort { min: u32 },
    #[error("value too long (max {max})")]
    TooLong { max: u32 },
    #[error("number out of range")]
    OutOfRange,
    #[error("nested: {0}")]
    Nested(Box<FieldValueError>),
}

#[derive(specta::Type, Debug, thiserror::Error, Serialize)]
pub enum FormError {
    #[error("values must be a JSON object")]
    NotAnObject,
    #[error("unknown field: {key}")]
    UnknownField { key: String },
    #[error("invalid values")]
    Invalid(HashMap<String, FieldValueError>),
}