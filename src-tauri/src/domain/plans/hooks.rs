use serde::Serialize;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub enum HookName {
    OnQuestCreate,
    OnRecordCreate,
    GetSettingsFields,
}

impl HookName {
    pub fn as_str(self) -> &'static str {
        match self {
            Self::OnQuestCreate    => "on_quest_create",
            Self::OnRecordCreate   => "on_record_create",
            Self::GetSettingsFields => "get_settings_fields",
        }
    }
}

impl std::fmt::Display for HookName {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str(self.as_str())
    }
}