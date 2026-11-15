use diesel::prelude::*;
use log::trace;

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
            .inspect(|plan| trace!("inserted plan {}", plan.id))
    }

    pub fn get(&self, conn: &mut DbConn, plan_id: &str) -> QueryResult<Plan> {
        plans::table
            .find(plan_id)
            .first(conn)
            .inspect(|_| trace!("found plan {}", plan_id))
            .inspect_err(|e| {
                if matches!(e, diesel::result::Error::NotFound) {
                    trace!("plan {} not found", plan_id);
                }
            })
    }

    pub fn get_all(&self, conn: &mut DbConn) -> QueryResult<Vec<Plan>> {
        plans::table
            .load(conn)
            .inspect(|rows| trace!("loaded {} plans", rows.len()))
    }

    pub fn update(&self, conn: &mut DbConn, plan_id: &str, changes: NewPlan) -> QueryResult<Plan> {
        diesel::update(plans::table.find(plan_id))
            .set(changes)
            .returning(Plan::as_returning())
            .get_result(conn)
            .inspect(|_| trace!("updated plan {}", plan_id))
    }

    pub fn delete(&self, conn: &mut DbConn, plan_id: &str) -> QueryResult<usize> {
        diesel::delete(plans::table.find(plan_id))
            .execute(conn)
            .inspect(|&n| trace!("deleted plan {} ({} rows)", plan_id, n))
    }

    pub fn delete_all(&self, conn: &mut DbConn) -> QueryResult<usize> {
        diesel::delete(plans::table)
            .execute(conn)
            .inspect(|&n| trace!("deleted all plans ({} rows)", n))
    }
}
