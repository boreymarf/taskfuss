// domain/form.rs
use serde_json::Value as JsonValue;
use std::collections::HashMap;

use crate::domain::fields::{Field, FieldSchema, FieldValue};
use crate::errors::{FieldValueError, FormError};

/// Validates raw frontend input against a set of field schemas.
///
/// Returns a dictionary of validated values on success, keyed by the same ids
/// the schemas use. Unknown keys in `values` are rejected; missing keys are
/// treated as `null` and validated by the field's own rules.
pub fn validate_form(
    fields: &HashMap<String, Field>,
    values: &JsonValue,
) -> Result<HashMap<String, FieldValue>, FormError> {
    let null = JsonValue::Null;
    let mut out = HashMap::with_capacity(fields.len()); // Ok
    let mut errors = HashMap::new(); // Err

    // Is this even an object.
    let obj = values.as_object().ok_or(FormError::NotAnObject)?;

    // Reject unknown fields. Should never happen, but who knows.
    for (key, value) in obj {
        if !fields.contains_key(key) {
            return Err(FormError::UnknownField { key: key.clone() });
        }
    }

    // Just validate.
    for (key, field) in fields {
        let raw = obj.get(key).unwrap_or(&null);
        match field.validate_value(raw) {
            Ok(v) => {
                out.insert(key.clone(), v);
            }
            Err(e) => {
                errors.insert(key.clone(), e);
            }
        }
    }

    if errors.is_empty() {
        Ok(out)
    } else {
        Err(FormError::Invalid(errors))
    }
}
