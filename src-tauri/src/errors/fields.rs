use std::collections::HashMap;

use serde::Serialize;

#[derive(specta::Type, Debug, thiserror::Error, Serialize)]
pub enum FieldSchemaError {
    #[error("min ({min}) is greater than max ({max})")]
    MinGreaterThanMax { min: i32, max: i32 },
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
    #[error("number {value} is below the minimum of {min}")]
    TooSmall { min: i32, value: i32 },
    #[error("number {value} is above the maximum of {max}")]
    TooBig { max: i32, value: i32 },
    #[error("number {value} does not fit in i32 (valid range: -2147483648 to 2147483647)")]
    I32Overflow { value: String },
    #[error("number {value} does not fit in u32 (valid range: 0 to 4294967295)")]
    U32Overflow { value: String },
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