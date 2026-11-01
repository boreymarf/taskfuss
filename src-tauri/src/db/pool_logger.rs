use diesel::r2d2::event::{
    AcquireEvent, CheckinEvent, CheckoutEvent, HandleEvent, ReleaseEvent, TimeoutEvent,
};

#[derive(Debug)]
pub(super) struct DbPoolLogger;

impl HandleEvent for DbPoolLogger {
    fn handle_checkout(&self, event: CheckoutEvent) {
        log::debug!(
            "pool: connection id={} checked out",
            event.connection_id()
        );
    }

    fn handle_checkin(&self, event: CheckinEvent) {
        log::debug!(
            "pool: connection id={} checked in after {:?}",
            event.connection_id(),
            event.duration()  // how long it was held
        );
    }

    fn handle_acquire(&self, event: AcquireEvent) {
        log::debug!(
            "pool: new connection id={} acquired",
            event.connection_id()
        );
    }

    fn handle_release(&self, event: ReleaseEvent) {
        log::debug!(
            "pool: connection id={} released / discarded from pool",
            event.connection_id()
        );
    }

    fn handle_timeout(&self, event: TimeoutEvent) {
        log::warn!(
            "pool: checkout timed out after {:?}",
            event.timeout()
        );
    }
}
