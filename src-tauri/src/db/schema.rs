// @generated automatically by Diesel CLI.

diesel::table! {
    plan_snapshots (uuid) {
        uuid -> Text,
        plan_id -> Text,
        path -> Text,
        hash -> Text,
        content -> Text,
        created_at -> Timestamp,
    }
}

diesel::table! {
    quests (uuid) {
        uuid -> Text,
        plan_snapshot_uuid -> Text,
        title -> Nullable<Text>,
        created_at -> Text,
        updated_at -> Text,
    }
}

diesel::joinable!(quests -> plan_snapshots (plan_snapshot_uuid));

diesel::allow_tables_to_appear_in_same_query!(plan_snapshots, quests,);
