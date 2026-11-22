use serde::{Deserialize, Serialize};
use serde_json::Value as JsonValue;

use crate::domain::fields::schema::FieldSchema;
use crate::errors::{FieldSchemaError, FieldValueError};

/// Schema for a plain text input field.
///
/// Used when a setting expects a string within an optional length range.
#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename = "StrFieldSchema")]
pub struct StrField {
    /// Whether the user must provide a value.
    pub required: bool,
    /// Human-readable label shown next to the input.
    pub label: Option<String>,
    /// Optional helper text describing the field's purpose.
    pub description: Option<String>,
    /// Value used when the user leaves the field empty.
    pub default: Option<String>,
    /// Maximum allowed length in characters (inclusive).
    pub max_length: Option<u32>,
    /// Minimum allowed length in characters (inclusive).
    pub min_length: Option<u32>,
}

impl FieldSchema for StrField {
    type Validated = Option<String>;

    fn validate_schema(&self) -> Result<(), FieldSchemaError> {
        if let (Some(min), Some(max)) = (self.min_length, self.max_length) {
            if min > max {
                return Err(FieldSchemaError::MinGreaterThanMax {
                    min: min as i64,
                    max: max as i64,
                });
            }
        }
        Ok(())
    }

    fn validate_value(&self, value: &JsonValue) -> Result<Self::Validated, FieldValueError> {
        match value {
            JsonValue::Null => {
                if self.required {
                    return Err(FieldValueError::Required);
                }
                Ok(None)
            }
            JsonValue::String(s) => {
                let len = s.chars().count() as u32;
                if let Some(max) = self.max_length {
                    if len > max {
                        return Err(FieldValueError::TooLong { max });
                    }
                }
                if let Some(min) = self.min_length {
                    if len < min {
                        return Err(FieldValueError::TooShort { min });
                    }
                }
                Ok(Some(s.clone()))
            }
            _ => Err(FieldValueError::WrongType { expected: "string" }),
        }
    }
}