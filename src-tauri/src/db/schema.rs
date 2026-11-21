// @generated automatically by Diesel CLI.

diesel::table! {
    plans (id) {
        id -> Text,
        name -> Text,
        author -> Text,
        version -> Text,
        entry -> Text,
        language -> Text,
        permissions -> Text,
    }
}

diesel::table! {
    quest_states (uuid) {
        uuid -> Text,
        quest_uuid -> Text,
        internal -> Text,
        starts_at -> Nullable<Timestamp>,
        ends_at -> Nullable<Timestamp>,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    quests (uuid) {
        uuid -> Text,
        plan_id -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    records (uuid) {
        uuid -> Text,
        state_uuid -> Text,
        field_id -> Text,
        author -> Text,
        data -> Text,
        created_at -> Timestamp,
    }
}

diesel::joinable!(quest_states -> quests (quest_uuid));
diesel::joinable!(records -> quest_states (state_uuid));

diesel::allow_tables_to_appear_in_same_query!(plans, quest_states, quests, records,);
