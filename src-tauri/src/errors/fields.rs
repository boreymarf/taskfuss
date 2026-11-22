#[derive(Debug, thiserror::Error)]
pub enum FieldSchemaError {
    #[error("min ({min}) is greater than max ({max})")]
    MinGreaterThanMax { min: i64, max: i64 },
}

#[derive(Debug, thiserror::Error)]
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