use serde::{Deserialize, Serialize, de::Error};

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PlanLanguage {
    Luau,
    Lua,
}

impl PlanLanguage {
    pub fn as_str(&self) -> &'static str {
        match self {
            PlanLanguage::Lua => "lua",
            PlanLanguage::Luau => "luau",
        }
    }
}

impl From<&PlanLanguage> for String {
    fn from(value: &PlanLanguage) -> String {
        value.as_str().to_owned()
    }
}

impl From<PlanLanguage> for String {
    fn from(value: PlanLanguage) -> String {
        value.as_str().to_owned()
    }
}

impl TryFrom<String> for PlanLanguage {
    type Error = serde_json::Error;
    fn try_from(s: String) -> Result<Self, Self::Error> {
        serde_json::from_value(serde_json::Value::String(s))
    }
}
