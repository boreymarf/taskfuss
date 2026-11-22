use std::collections::HashMap;

use serde::{Serialize, Deserialize};

use serde_json::Value as JsonValue;
use crate::errors::{FieldSchemaError, FieldValueError};
use crate::domain::{Field, FieldValue, StrField};

/// Validation contract shared by every field schema and by the [`Field`] enum.
///
/// A schema describes the rules a field enforces, and provides two entry points
/// to check them: one for the schema itself, and one for user input.
pub trait FieldSchema {
    /// The type produced by [`validate_value`](FieldSchema::validate_value).
    ///
    /// On individual field schemas this is the field's own value type (e.g.
    /// [`Option<String>`] for [`StrField`]). On the [`Field`] enum it is
    /// [`FieldValue`], so values of any field kind can be collected into a
    /// single dictionary.
    type Validated;

    /// Verifies that the schema's own configuration is coherent, e.g. that
    /// `min` is not greater than `max`.
    ///
    /// Defaults to [`Ok`] for schemas with no cross-field constraints.
    fn validate_schema(&self) -> Result<(), FieldSchemaError> {
        Ok(())
    }

    /// Takes a raw JSON value from the frontend and produces either a typed
    /// [`Validated`](FieldSchema::Validated) value or a [`FieldValueError`].
    fn validate_value(&self, value: &serde_json::Value) -> Result<Self::Validated, FieldValueError>;
}