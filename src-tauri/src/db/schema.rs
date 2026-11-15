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
    quests (uuid) {
        uuid -> Text,
        plan_id -> Text,
        title -> Nullable<Text>,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::allow_tables_to_appear_in_same_query!(plans, quests,);
