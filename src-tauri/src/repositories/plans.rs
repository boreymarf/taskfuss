use diesel::prelude::*;

use crate::db::{DbConn, plans};
use crate::domain::{NewPlan, Plan};

pub struct PlanRepository;

impl PlanRepository {
    pub fn new() -> Self {
        Self
    }

    pub fn create(&self, conn: &mut DbConn, new_plan: NewPlan) -> QueryResult<Plan> {
        diesel::insert_into(plans::table)
            .values(new_plan)
            .returning(Plan::as_returning())
            .get_result(conn)
    }

    pub fn get(&self, conn: &mut DbConn, plan_id: &str) -> QueryResult<Plan> {
        plans::table.find(plan_id).first(conn)
    }

    pub fn get_all(&self, conn: &mut DbConn) -> QueryResult<Vec<Plan>> {
        plans::table.load(conn)
    }

    pub fn update(
        &self,
        conn: &mut DbConn,
        plan_id: &str,
        changes: NewPlan,
    ) -> QueryResult<Plan> {
        diesel::update(plans::table.find(plan_id))
            .set(changes)
            .returning(Plan::as_returning())
            .get_result(conn)
    }

    pub fn delete(&self, conn: &mut DbConn, plan_id: &str) -> QueryResult<usize> {
        diesel::delete(plans::table.find(plan_id)).execute(conn)
    }

    pub fn delete_all(&self, conn: &mut DbConn) -> QueryResult<usize> {
        diesel::delete(plans::table).execute(conn)
    }
}

