// @generated automatically by Diesel CLI.

diesel::table! {
    quests (id) {
        id -> Integer,
        title -> Nullable<Text>,
        created_at -> Text,
        updated_at -> Text,
    }
}
