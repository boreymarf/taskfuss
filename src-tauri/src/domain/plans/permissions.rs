use serde::{Deserialize, Serialize};

// Permission
#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum PlanPermission {
    Notifications,
    Timers,
    Storage,
    #[serde(rename = "tasks:read")]
    TasksRead,
    #[serde(rename = "tasks:write")]
    TasksWrite,
}

// Permissions
#[derive(Debug, Clone, Default, PartialEq, Eq, Serialize, Deserialize)]
#[serde(transparent)]
pub struct PlanPermissions(pub Vec<PlanPermission>);

impl PlanPermissions {
    pub fn has(&self, p: PlanPermission) -> bool {
        self.0.contains(&p)
    }
}

impl From<PlanPermissions> for String {
    fn from(p: PlanPermissions) -> String {
        serde_json::to_string(&p.0).expect("permissions are serializable")
    }
}

impl From<&PlanPermissions> for String {
    fn from(p: &PlanPermissions) -> String {
        serde_json::to_string(&p.0).expect("permissions are serializable")
    }
}

impl TryFrom<String> for PlanPermissions {
    type Error = serde_json::Error;
    fn try_from(s: String) -> Result<Self, Self::Error> {
        serde_json::from_str(&s).map(PlanPermissions)
    }
}

impl TryFrom<Vec<String>> for PlanPermissions {
    type Error = serde_json::Error;
    fn try_from(items: Vec<String>) -> Result<Self, Self::Error> {
        serde_json::from_value(serde_json::Value::Array(
            items.into_iter().map(serde_json::Value::String).collect(),
        ))
    }
}
