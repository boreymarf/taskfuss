use serde::{Serialize, Deserialize};

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(tag = "type")]
pub enum PlanField {
    #[serde(rename = "str")]
    Str(StrField),
    #[serde(rename = "int")]
    Int(IntField),
    #[serde(rename = "str_list")]
    StrList(StrListField),
    #[serde(rename = "progress")]
    Progress(ProgressBar),
    #[serde(rename = "parent")]
    Parent(ParentField),
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct StrField {
    pub required: bool,
    pub label: Option<String>,
    pub description: Option<String>,
    pub default: Option<String>,
    pub max_length: Option<usize>,
    pub min_length: Option<usize>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct IntField {
    pub required: bool,
    pub label: Option<String>,
    pub description: Option<String>,
    pub default: Option<i64>,
    pub min: Option<i64>,
    pub max: Option<i64>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct StrListField {
    pub required: bool,
    pub label: Option<String>,
    pub description: Option<String>,
    pub default: Option<Vec<String>>,
    pub default_length: Option<usize>,
    pub max_array_length: Option<usize>,
    pub min_array_length: Option<usize>,
    pub max_string_length: Option<usize>,
    pub min_string_length: Option<usize>,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct ProgressBar {
    pub required: bool,
    pub label: Option<String>,
    pub max_value: i32,
    pub default: i32,
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct ParentField {
    pub required: bool,
    pub label: Option<String>,
    pub children: Vec<PlanField>,
}