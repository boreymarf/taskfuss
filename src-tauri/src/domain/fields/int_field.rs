use serde::{Deserialize, Serialize};
use serde_json::Value as JsonValue;

use crate::domain::fields::schema::FieldSchema;
use crate::errors::{FieldSchemaError, FieldValueError};

/// Schema for a plain integer input field.
///
/// Used when a setting expects a whole number within an optional range.
#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(rename = "IntFieldSchema")]
pub struct IntField {
    /// Whether the user must provide a value.
    pub required: bool,
    /// Human-readable label shown next to the input.
    pub label: Option<String>,
    /// Optional helper text describing the field's purpose.
    pub description: Option<String>,
    /// Value used when the user leaves the field empty.
    pub default: Option<i32>,
    /// Minimum allowed value (inclusive).
    pub min: Option<i32>,
    /// Maximum allowed value (inclusive).
    pub max: Option<i32>,
}

impl FieldSchema for IntField {
    type Validated = Option<i32>;

    fn validate_schema(&self) -> Result<(), FieldSchemaError> {
        if let (Some(min), Some(max)) = (self.min, self.max) {
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
            JsonValue::Number(n) => {
                let i = n.as_i64().and_then(|v| i32::try_from(v).ok()).ok_or(
                    FieldValueError::WrongType {
                        expected: "integer",
                    },
                )?;
                if let Some(min) = self.min {
                    if i < min {
                        return Err(FieldValueError::OutOfRange);
                    }
                }
                if let Some(max) = self.max {
                    if i > max {
                        return Err(FieldValueError::OutOfRange);
                    }
                }
                Ok(Some(i))
            }
            _ => Err(FieldValueError::WrongType {
                expected: "integer",
            }),
        }
    }
}