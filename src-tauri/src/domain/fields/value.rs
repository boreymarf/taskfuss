use serde::Serialize;

use crate::domain::FieldSchema;
use crate::managed::PlanRuntime;

/// Validated value produced by a single field.
///
/// The frontend receives field definitions as a dictionary keyed by field id,
/// where each entry is a schema describing how the field should be rendered and
/// what rules it enforces:
///
/// ```json
/// {
///   "display_name": {
///     "type": "str",
///     "required": true,
///     "label": "Display name",
///     "max_length": 64
///   },
///   "max_uses": {
///     "type": "int",
///     "required": false,
///     "label": "Max uses",
///     "min": 0
///   }
/// }
/// ```
///
/// The frontend returns user input in the same shape, keyed by the same ids:
///
/// ```json
/// {
///   "display_name": "My quest",
///   "max_uses": 3
/// }
/// ```
///
/// Each entry is then passed through the matching field's
/// [`FieldSchema::validate_value`], which either produces a value or an error
/// for the frontend. The results are collected into a dictionary of
/// `FieldValue`s and forwarded to the plan through [`PlanRuntime`] parameters,
/// preserving the same keyed format.
///
/// [`None`] inside a variant means the user left the field empty (only valid
/// when the schema is not `required`).
#[derive(Debug, Clone, Serialize)]
#[serde(untagged)]
pub enum FieldValue {
    Str(Option<String>),
    Int(Option<i32>),
}